package geometry

type BoxEdge int

const TOP = BoxEdge(1)
const BOTTOM = BoxEdge(2)
const LEFT = BoxEdge(3)
const RIGHT = BoxEdge(4)

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
