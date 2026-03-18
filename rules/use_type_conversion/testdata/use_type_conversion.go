package fixtures

type Point2D struct {
	X, Y int
}

type Coordinate struct {
	X, Y int
}

type Point3D struct {
	X, Y, Z int
}

type Vec3 struct {
	X, Y, Z int
}

type Different struct {
	X int
	Y string
}

type AlsoDifferent struct {
	X int
	Z int
}

type Single struct {
	X int
}

type AlsoSingle struct {
	X int
}

// Invalid: field-by-field copy of two compatible struct types (2 fields).
func convertTwoFields(p Point2D) Coordinate {
	var c Coordinate
	c.X = p.X // MATCH /use type conversion Coordinate(p) instead of copying struct fields one by one from Point2D to Coordinate/
	c.Y = p.Y
	return c
}

// Invalid: field-by-field copy of two compatible struct types (3 fields).
func convertThreeFields(v Point3D) Vec3 {
	var out Vec3
	out.X = v.X // MATCH /use type conversion Vec3(v) instead of copying struct fields one by one from Point3D to Vec3/
	out.Y = v.Y
	out.Z = v.Z
	return out
}

// Valid: types have different number of fields.
func differentFieldCount(p Point2D) Point3D {
	var out Point3D
	out.X = p.X
	out.Y = p.Y
	return out
}

// Valid: types have different field types.
func differentFieldTypes(p Point2D) Different {
	var d Different
	d.X = p.X
	return d
}

// Valid: types have different field names.
func differentFieldNames(p Point2D) AlsoDifferent {
	var d AlsoDifferent
	d.X = p.X
	return d
}

// Valid: only one field is copied out of two.
func partialCopy(p Point2D) Coordinate {
	var c Coordinate
	c.X = p.X
	return c
}

// Valid: different fields are used on left and right sides.
func mismatchedFields(p Point2D) Coordinate {
	var c Coordinate
	c.X = p.Y
	c.Y = p.X
	return c
}

// Valid: using a type conversion already.
func alreadyConverted(p Point2D) Coordinate {
	return Coordinate(p)
}

// Valid: non-consecutive assignments (other statement in between).
func nonConsecutive(p Point2D) Coordinate {
	var c Coordinate
	c.X = p.X
	_ = c.X
	c.Y = p.Y
	return c
}

// Valid: assignment to self (same variable).
func selfAssignment(c Coordinate) {
	c.X = c.X
	c.Y = c.Y
}

// Valid: mixed src variables.
func mixedSources(p1, p2 Point2D) Coordinate {
	var c Coordinate
	c.X = p1.X
	c.Y = p2.Y
	return c
}
