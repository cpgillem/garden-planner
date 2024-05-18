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

func (aabb *AxisAlignedBoundingBox) IsVertical() bool {
	return aabb.Size.Y >= aabb.Size.X
}

func NewBox() AxisAlignedBoundingBox {
	return AxisAlignedBoundingBox{
		Location: NewVector(0, 0, 0),
		Size:     NewVector(0, 0, 0),
	}
}
