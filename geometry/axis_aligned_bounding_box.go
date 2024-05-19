package geometry

type BoxEdge int

const TOP = BoxEdge(1)
const BOTTOM = BoxEdge(2)
const LEFT = BoxEdge(3)
const RIGHT = BoxEdge(4)

// Bounding box with a location and size.
// Locations are bottom-left origin.
type AxisAlignedBoundingBox struct {
	Location Vector
	Size     Vector
}

func NewBox() AxisAlignedBoundingBox {
	return AxisAlignedBoundingBox{
		Location: NewVector(0, 0, 0),
		Size:     NewVector(0, 0, 0),
	}
}

func NewBoxWithValues(x float32, y float32, width float32, height float32) AxisAlignedBoundingBox {
	return AxisAlignedBoundingBox{
		Location: NewVector(x, y, 0),
		Size:     NewVector(width, height, 0),
	}
}

func (aabb *AxisAlignedBoundingBox) IsVertical() bool {
	return aabb.Size.Y >= aabb.Size.X
}

// Mutates this box.
func (aabb *AxisAlignedBoundingBox) AddTo(b2 *AxisAlignedBoundingBox) *AxisAlignedBoundingBox {
	aabb.Location.AddTo(&b2.Location)
	aabb.Size.AddTo(&b2.Size)
	return aabb
}
