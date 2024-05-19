package geometry

import "fyne.io/fyne/v2"

// 3D vector for storing garden layout data in the abstract.
// Coordinates are top-left.
type Vector struct {
	X, Y, Z float32
}

func NewVector(x, y, z float32) Vector {
	return Vector{
		X: x,
		Y: y,
		Z: z,
	}
}

func (v *Vector) Scale(scalar float32) *Vector {
	return &Vector{
		X: v.X * scalar,
		Y: v.Y * scalar,
		Z: v.Z * scalar,
	}
}

func (v *Vector) Add(v2 *Vector) *Vector {
	return &Vector{
		X: v.X + v2.X,
		Y: v.Y + v2.Y,
		Z: v.Z + v2.Z,
	}
}

func (v *Vector) AddTo(v2 *Vector) *Vector {
	v.X += v2.X
	v.Y += v2.Y
	v.Z += v2.Z
	return v
}

func (v *Vector) Negate() *Vector {
	return v.Scale(-1)
}

func (v *Vector) ToSize() fyne.Size {
	return fyne.NewSize(v.X, v.Y)
}

func (v *Vector) ToPosition() fyne.Position {
	return fyne.NewPos(v.X, v.Y)
}

func NewVectorFromSize(size *fyne.Size) Vector {
	return Vector{
		X: size.Width,
		Y: size.Height,
		Z: 0,
	}
}
