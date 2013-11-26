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

func (s *ReadSuite) TestBuildReadSpecs(c *C) {
	target := &Target{}
	result, err := buildReadSpecs(target)
	c.Assert(err, IsNil)
	c.Assert(result, HasLen, 3)
	spec := result[0]
	c.Assert(spec.FieldType.Name, Equals, "Forename")
	c.Assert(spec.Length, Equals, 4)
	c.Assert(spec.Repeat, Equals, 1)
}
