package fixedfield

import (
	"bytes"
	"math"
	. "launchpad.net/gocheck"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type ReadSuite struct{}

var _ = Suite(&ReadSuite{})

type Target struct {
	Name         string  `length:"5"`
	Age          int     `length:"2" encoding:"ascii"`
	ShoeSize     int     `length:"2" encoding:"bigendian"`
	CollarSize   int     `length:"2" encoding:"le"`
	ElbowBreadth uint    `length:"8" encoding:"le"`
	NoseCapacity float64 `length:"6" encoding:"ascii"`
	Pi           float64 `length:"8" encoding:"le"`
	UpsideDownCake float32 `length:"4" encoding:"be"`
	Ratings      []int   `length:"1" repeat:"10"`
}

// buildReadSpecs can read a struct and it's tags to build a valid
// readSpec
func (s *ReadSuite) TestBuildReadSpecs(c *C) {
	target := &Target{}
	result, err := buildReadSpecs(target)
	c.Assert(err, IsNil)
	c.Assert(result, HasLen, 9)
	spec := result[0]
	c.Assert(spec.FieldType.Name, Equals, "Name")
	c.Assert(spec.Length, Equals, 5)
	c.Assert(spec.Repeat, Equals, 1)
	spec = result[1]
	c.Assert(spec.FieldType.Name, Equals, "Age")
	c.Assert(spec.Length, Equals, 2)
	c.Assert(spec.Repeat, Equals, 1)
	c.Assert(spec.Encoding, Equals, "ascii")
	spec = result[2]
	c.Assert(spec.FieldType.Name, Equals, "ShoeSize")
	c.Assert(spec.Length, Equals, 2)
	c.Assert(spec.Repeat, Equals, 1)
	c.Assert(spec.Encoding, Equals, "bigendian")
	spec = result[3]
	c.Assert(spec.FieldType.Name, Equals, "CollarSize")
	c.Assert(spec.Length, Equals, 2)
	c.Assert(spec.Repeat, Equals, 1)
	c.Assert(spec.Encoding, Equals, "le")
	spec = result[4]
	c.Assert(spec.FieldType.Name, Equals, "ElbowBreadth")
	c.Assert(spec.Length, Equals, 8)
	c.Assert(spec.Repeat, Equals, 1)
	c.Assert(spec.Encoding, Equals, "le")
	spec = result[5]
	c.Assert(spec.FieldType.Name, Equals, "NoseCapacity")
	c.Assert(spec.Length, Equals, 6)
	c.Assert(spec.Repeat, Equals, 1)
	c.Assert(spec.Encoding, Equals, "ascii")
	spec = result[6]
	c.Assert(spec.FieldType.Name, Equals, "Pi")
	c.Assert(spec.Length, Equals, 8)
	c.Assert(spec.Repeat, Equals, 1)
	c.Assert(spec.Encoding, Equals, "le")
	spec = result[7]
	c.Assert(spec.FieldType.Name, Equals, "UpsideDownCake")
	c.Assert(spec.Length, Equals, 4)
	c.Assert(spec.Repeat, Equals, 1)
	c.Assert(spec.Encoding, Equals, "be")
	spec = result[8]
	c.Assert(spec.FieldType.Name, Equals, "Ratings")
	c.Assert(spec.Length, Equals, 1)
	c.Assert(spec.Repeat, Equals, 10)
}

// Test populateStructFromReadSpecAndBytes copies values from a
// ReaderSeeker into the appropriate structural elements
func (s *ReadSuite) TestPopulateStructFromReadSpecAndBytes(c *C) {
	data := bytes.NewBuffer(
		[]byte("Geoff" +
			"36" +
			"\x00\x7f" +
			"\x7f\x00" +
			"\xff\xff\xff\xff\xff\xff\xff\xff" +
			"001.23" +
			"\x18\x2d\x44\x54\xfb\x21\x09\x40" +
			"\x40\x49\x0f\xdb" +
			"0123456789"))
	target := &Target{}
	readSpec, err := buildReadSpecs(target)
	c.Assert(err, IsNil)
	err = populateStructFromReadSpecAndBytes(target, readSpec, data)
	c.Assert(err, IsNil)
	c.Assert(target.Name, Equals, "Geoff")
	c.Assert(target.Age, Equals, 36)
	c.Assert(target.ShoeSize, Equals, 127)
	c.Assert(target.CollarSize, Equals, 127)
	c.Assert(target.ElbowBreadth, Equals, uint(18446744073709551615))
	c.Assert(target.NoseCapacity, Equals, 1.23)
	c.Assert(target.Pi, Equals, math.Pi)
	c.Assert(target.UpsideDownCake, Equals, float32(math.Pi))
}
