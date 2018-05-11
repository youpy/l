package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"../.."
)

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}

func read(conn *net.TCPConn) {
	buf := make([]byte, 1)
	for {
		n, err := conn.Read(buf)
		if n == 0 {
			break
		}
		checkError(err)

		if string(buf[:n]) == "\u0000" {
			break
		}

		fmt.Print(string(buf[:n]))
	}
}

func write(turtle *lindenmayer.Turtle, conn *net.TCPConn) {
	var z int

	if turtle.Pen {
		z = 0
	} else {
		z = 1000
	}

	_, err := conn.Write(
		[]byte(
			"G01" +
				" X" + strconv.FormatFloat(turtle.Pos.X, 'f', 4, 64) +
				" Y" + strconv.FormatFloat(turtle.Pos.Y, 'f', 4, 64) +
				" Z" + strconv.Itoa(z) +
				"\u0000",
		),
	)
	checkError(err)

	time.Sleep(200 * time.Millisecond)
}

func main() {
	x := flag.Float64("x", 1000.0, "initial X position")
	y := flag.Float64("y", -1300.0, "initial Y position")
	ruleF := flag.String("ruleF", "F[+F-F-F]F[--F+F+F]", "rule F")
	initialState := flag.String("initial", "F", "initial state")
	iterations := flag.Int("iterations", 4, "number of iterations")
	side := flag.Float64("side", 40.0, "side")
	angle := flag.Float64("angle", 15.0, "angle")
	hostname := flag.String("hostname", "line-us.local:1337", "hostname of line-us machine")

	flag.Parse()

	tcpAddr, err := net.ResolveTCPAddr("tcp4", *hostname)
	checkError(err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	read(conn)

	turtle := lindenmayer.NewTurtle(*x, *y)
	done := make(chan struct{})

	write(turtle, conn)

	go func() {
		for _ = range turtle.Dispatcher {
			write(turtle, conn)
			read(conn)
		}

		close(done)
	}()

	lindenmayer.LSystem(turtle, *ruleF, *initialState, *iterations, *side, *angle)

	<-done
}
