package geometry

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
