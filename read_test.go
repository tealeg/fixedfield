package fixedfield

import (
	"testing"
	. "launchpad.net/gocheck"
)

func Test(t *testing.T) {TestingT(t)}

type ReadSuite struct {}

var _ = Suite(&ReadSuite{})

type Target struct {
	Foreame string `length:"4"`
	Surname string `length:"4"`
	Ratings []int `length:"1" repeat:"10"`
}

func (s *ReadSuite) TestBuildReadSpecs(c *C) {
	target := &Target{}
	result, err := buildReadSpecs(target)
	c.Assert(err, IsNil)
	c.Assert(result, HasLen, 3)
	spec := result[0]
	c.Assert(spec.GetFieldName(), Equals, "Forename")
}
