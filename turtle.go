package lindenmayer

import "math"

type Turtle struct {
	Pos         position
	Orientation float64
	Pen         bool
	Dispatcher  chan *Turtle
}

type position struct {
	X, Y float64
}

func NewTurtle(x float64, y float64) *Turtle {
	return &Turtle{position{x, y}, 0.0, false, make(chan *Turtle)}
}

func (t *Turtle) Forward(dist float64) {
	t.Pos.X += dist * math.Sin(t.Orientation)
	t.Pos.Y += dist * math.Cos(t.Orientation)

	t.Dispatcher <- t
}

func (t *Turtle) Left(radians float64) {
	t.Orientation += radians
}

func (t *Turtle) Right(radians float64) {
	t.Left(-radians)
}

func (t *Turtle) PenUp() {
	t.Pen = false
}

func (t *Turtle) PenDown() {
	t.Pen = true
}

func (t *Turtle) Clone() *Turtle {
	return &Turtle{position{t.Pos.X, t.Pos.Y}, t.Orientation, t.Pen, t.Dispatcher}
}

func (t *Turtle) Restore(from *Turtle) {
	t.PenUp()

	t.Pos.X = from.Pos.X
	t.Pos.Y = from.Pos.Y
	t.Orientation = from.Orientation

	t.Dispatcher <- t

	t.PenDown()
}
