package fixedfield

import (
	. "launchpad.net/gocheck"
)


type SpecSuite struct {}

var _ = Suite(&SpecSuite{})


// buildSpecs can read a struct and it's tags to build a valid
// readSpec
func (s *SpecSuite) TestBuildReadSpecs(c *C) {
	target := &Target{}
	result, err := buildSpecs(target)
	c.Assert(err, IsNil)
	c.Assert(result, HasLen, 12)
	spec := result[0]
	c.Assert(spec.StructName, Equals, "*fixedfield.Target")
	c.Assert(spec.StructField.Name, Equals, "Name")
	c.Assert(spec.Length, Equals, 5)
	c.Assert(spec.Repeat, Equals, 1)
	spec = result[1]
	c.Assert(spec.StructField.Name, Equals, "Age")
	c.Assert(spec.Length, Equals, 2)
	c.Assert(spec.Repeat, Equals, 1)
	c.Assert(spec.Encoding, Equals, "ascii")
	spec = result[2]
	c.Assert(spec.StructField.Name, Equals, "ShoeSize")
	c.Assert(spec.Length, Equals, 2)
	c.Assert(spec.Repeat, Equals, 1)
	c.Assert(spec.Encoding, Equals, "bigendian")
	spec = result[3]
	c.Assert(spec.StructField.Name, Equals, "CollarSize")
	c.Assert(spec.Length, Equals, 2)
	c.Assert(spec.Repeat, Equals, 1)
	c.Assert(spec.Encoding, Equals, "le")
	spec = result[4]
	c.Assert(spec.StructField.Name, Equals, "ElbowBreadth")
	c.Assert(spec.Length, Equals, 8)
	c.Assert(spec.Repeat, Equals, 1)
	c.Assert(spec.Encoding, Equals, "le")
	spec = result[5]
	c.Assert(spec.StructField.Name, Equals, "NoseCapacity")
	c.Assert(spec.Length, Equals, 6)
	c.Assert(spec.Repeat, Equals, 1)
	c.Assert(spec.Encoding, Equals, "ascii")
	spec = result[6]
	c.Assert(spec.StructField.Name, Equals, "Pi")
	c.Assert(spec.Length, Equals, 8)
	c.Assert(spec.Repeat, Equals, 1)
	c.Assert(spec.Encoding, Equals, "le")
	spec = result[7]
	c.Assert(spec.StructField.Name, Equals, "UpsideDownCake")
	c.Assert(spec.Length, Equals, 4)
	c.Assert(spec.Repeat, Equals, 1)
	c.Assert(spec.Encoding, Equals, "be")
	spec = result[8]
	c.Assert(spec.StructField.Name, Equals, "Enrolled")
	c.Assert(spec.Length, Equals, 1)
	c.Assert(spec.Repeat, Equals, 1)
	c.Assert(spec.Encoding, Equals, "LE")
	spec = result[9]
	c.Assert(spec.StructField.Name, Equals, "ShouldBeEnrolled")
	c.Assert(spec.Length, Equals, 1)
	c.Assert(spec.Repeat, Equals, 1)
	c.Assert(spec.Encoding, Equals, "ascii")
	c.Assert(string(spec.TrueBytes), Equals, "Yy")
	spec = result[10]
	c.Assert(spec.StructField.Name, Equals, "Dispatched")
	c.Assert(spec.Length, Equals, 1)
	c.Assert(spec.Repeat, Equals, 1)
	c.Assert(spec.Encoding, Equals, "ascii")
	c.Assert(string(spec.TrueBytes), Equals, "jJ")
	spec = result[11]
	c.Assert(spec.StructField.Name, Equals, "Ratings")
	c.Assert(spec.Length, Equals, 1)
	c.Assert(spec.Repeat, Equals, 10)
	c.Assert(spec.Encoding, Equals, "ascii")
}

// Test that buildSpecs copes with nested structures
func (s *SpecSuite) TestBuildReadSpecsWithNestedStructs(c *C) {
	transaction := &Transaction{}
	result, err := buildSpecs(transaction)
	c.Assert(err, IsNil)
	c.Assert(result, HasLen, 2)
	spec := result[0]
	c.Assert(spec.StructName, Equals, "*fixedfield.Transaction")
	c.Assert(spec.StructField.Name, Equals, "Buyer")
	c.Assert(spec.Length, Equals, 0)
	c.Assert(spec.Repeat, Equals, 0)
	c.Assert(len(spec.Children), Equals, 2)
	childSpec := spec.Children[0]
	c.Assert(childSpec.StructName, Equals, "fixedfield.Person")
	c.Assert(childSpec.StructField.Name, Equals, "Name")
	c.Assert(childSpec.Length, Equals, 5)
	c.Assert(childSpec.Repeat, Equals, 1)
	childSpec = spec.Children[1]
	c.Assert(childSpec.StructName, Equals, "fixedfield.Person")
	c.Assert(childSpec.StructField.Name, Equals, "Age")
	c.Assert(childSpec.Length, Equals, 1)
	c.Assert(childSpec.Repeat, Equals, 1)
	spec = result[1]
	c.Assert(spec.StructName, Equals, "*fixedfield.Transaction")
	c.Assert(spec.StructField.Name, Equals, "Seller")
	c.Assert(spec.Length, Equals, 0)
	c.Assert(spec.Repeat, Equals, 0)
	c.Assert(len(spec.Children), Equals, 2)
	childSpec = spec.Children[0]
	c.Assert(childSpec.StructName, Equals, "fixedfield.Person")
	c.Assert(childSpec.StructField.Name, Equals, "Name")
	childSpec = spec.Children[1]
	c.Assert(childSpec.StructName, Equals, "fixedfield.Person")
	c.Assert(childSpec.StructField.Name, Equals, "Age")
}


