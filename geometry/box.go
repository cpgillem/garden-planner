package geometry

type BoxEdge int

const TOP = BoxEdge(1)
const BOTTOM = BoxEdge(2)
const LEFT = BoxEdge(3)
const RIGHT = BoxEdge(4)

type BoxCorner int

const TOP_LEFT = BoxCorner(1)
const TOP_RIGHT = BoxCorner(2)
const BOTTOM_RIGHT = BoxCorner(3)
const BOTTOM_LEFT = BoxCorner(4)

// Bounding box with a location and size.
// Locations are bottom-left origin.
type Box struct {
	Location Vector
	Size     Vector
}

func NewBoxZero() Box {
	return Box{
		Location: NewVector(0, 0, 0),
		Size:     NewVector(0, 0, 0),
	}
}

func NewBox(x float32, y float32, width float32, height float32) Box {
	return Box{
		Location: NewVector(x, y, 0),
		Size:     NewVector(width, height, 0),
	}
}

func (aabb *Box) IsVertical() bool {
	return aabb.Size.Y >= aabb.Size.X
}

// Mutates this box.
func (aabb *Box) AddTo(b2 *Box) *Box {
	aabb.Location.AddTo(&b2.Location)
	aabb.Size.AddTo(&b2.Size)
	return aabb
}

func (aabb *Box) Copy() Box {
	return Box{
		aabb.Location.Copy(),
		aabb.Size.Copy(),
	}
}

func (box *Box) GetX() float32 {
	return box.Location.X
}

func (box *Box) GetY() float32 {
	return box.Location.Y
}

func (box *Box) GetWidth() float32 {
	return box.Size.X
}

func (box *Box) GetHeight() float32 {
	return box.Size.Y
}

func (box *Box) SetX(v float32) {
	box.Location.X = v
}

func (box *Box) SetY(v float32) {
	box.Location.Y = v
}

func (box *Box) SetWidth(v float32) {
	box.Size.X = v
}

func (box *Box) SetHeight(v float32) {
	box.Size.Y = v
}
