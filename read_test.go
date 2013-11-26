package fixedfield

import (
	"testing"
	. "launchpad.net/gocheck"
)

func Test(t *testing.T) {TestingT(t)}

type ReadSuite struct {}

var _ = Suite(&ReadSuite{})

type Target struct {
	Forename string `length:"4"`
	Surname string `length:"4"`
	Ratings []int `length:"1" repeat:"10"`
}

// buildReadSpecs can read a struct and it's tags to build a valid
// readSpec
func (s *ReadSuite) TestBuildReadSpecs(c *C) {
	target := &Target{}
	result, err := buildReadSpecs(target)
	c.Assert(err, IsNil)
	c.Assert(result, HasLen, 3)
	spec := result[0]
	c.Assert(spec.FieldType.Name, Equals, "Forename")
	c.Assert(spec.Length, Equals, 4)
	c.Assert(spec.Repeat, Equals, 1)
	spec = result[1]
	c.Assert(spec.FieldType.Name, Equals, "Surname")
	c.Assert(spec.Length, Equals, 4)
	c.Assert(spec.Repeat, Equals, 1)
	spec = result[2]
	c.Assert(spec.FieldType.Name, Equals, "Ratings")
	c.Assert(spec.Length, Equals, 1)
	c.Assert(spec.Repeat, Equals, 10)
}
