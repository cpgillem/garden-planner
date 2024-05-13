package geometry

// Bounding box with a location and size.
// Locations are bottom-left origin.
type AxisAlignedBoundingBox struct {
	Location Vector
	Size     Vector
}
