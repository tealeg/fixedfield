package fixedfield

import (
	"bytes"
	"testing"
	. "launchpad.net/gocheck"
)

func Test(t *testing.T) {TestingT(t)}

type ReadSuite struct {}

var _ = Suite(&ReadSuite{})

type Target struct {
	Name string `length:"5"`
	Age int `length:"2"`
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
	c.Assert(spec.FieldType.Name, Equals, "Name")
	c.Assert(spec.Length, Equals, 5)
	c.Assert(spec.Repeat, Equals, 1)
	spec = result[1]
	c.Assert(spec.FieldType.Name, Equals, "Age")
	c.Assert(spec.Length, Equals, 2)
	c.Assert(spec.Repeat, Equals, 1)
	spec = result[2]
	c.Assert(spec.FieldType.Name, Equals, "Ratings")
	c.Assert(spec.Length, Equals, 1)
	c.Assert(spec.Repeat, Equals, 10)
}


// Test populateStructFromReadSpecAndBytes copies values from a
// ReaderSeeker into the appropriate structural elements
func (s *ReadSuite) TestPopulateStructFromReadSpecAndBytes(c *C) {
	data := bytes.NewBuffer([]byte("Geoff360123456789"))
	target := &Target{}
	readSpec, err := buildReadSpecs(target)
	c.Assert(err, IsNil)
	err = populateStructFromReadSpecAndBytes(target, readSpec, data)
	c.Assert(err, IsNil)
	c.Assert(target.Name, Equals, "Geoff")
	c.Assert(target.Age, Equals, 36)
}
