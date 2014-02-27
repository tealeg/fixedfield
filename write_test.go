package fixedfield

import (
	. "launchpad.net/gocheck"
	"math"
)

type WriteSuite struct{}

var _ = Suite(&WriteSuite{})

func (s *WriteSuite) TestMarshalIntegerASCIIPositive(c *C) {
	type target struct {
		Value int `encoding:"ascii", length:"1"`
	}
	var t *target
	var block []byte

	t = &target{Value: 3}
	specs, err := buildSpecs(t)
	c.Assert(err, IsNil)
	block, err = marshalInteger(specs[0])
	c.Assert(err, IsNil)
	c.Assert(string(block), Equals, "3")
}

func (s *WriteSuite) TestMarshalIntegerASCIINegative(c *C) {
	type target struct {
		Value int `encoding:"ascii", length:"2"`
	}
	var t *target
	var block []byte

	t = &target{Value: -3}
	specs, err := buildSpecs(t)
	c.Assert(err, IsNil)
	block, err = marshalInteger(specs[0])
	c.Assert(err, IsNil)
	c.Assert(string(block), Equals, "-3")
}


func (s *WriteSuite) TestMarshalIntegerASCIIPositivePad(c *C) {
	type target struct {
		Value int `encoding:"ascii", length:"10", pad:" "`
	}
	var t *target
	var block []byte

	t = &target{Value: 33}
	specs, err := buildSpecs(t)
	c.Assert(err, IsNil)
	block, err = marshalInteger(specs[0])
	c.Assert(err, IsNil)
	c.Assert(string(block), Equals, "        33")
}


func (s *WriteSuite) TestPopulateBytesFromSpecAndStruct(c *C) {
	var data []byte
	var specs []spec
	var err error

	target := &Target{
		Name: "Geoff",
		Age: 36,
		ShoeSize: 10,
		CollarSize: 16,
		ElbowBreadth: 32,
		NoseCapacity: 10000.1289,
		Pi: math.Pi,
		UpsideDownCake: float32(math.Pi),
		Enrolled: false,
		ShouldBeEnrolled: true,
		Dispatched: true,
		Ratings: []int{0,1,2,3,4,5,6,7,8,9}}


	specs, err = buildSpecs(target)
	c.Assert(err, IsNil)
	data, err = populateBytesFromSpecAndStruct(specs)
	c.Assert(err, IsNil)
	c.Assert(string(data[0:5]), Equals, "Geoff")
}


// Marshal builds Specs and dumps the target struct to a byte array.
func (s *WriteSuite) TestMarshal(c *C) {
	var data []byte
	var err error

	target := &Target{
		Name: "Geoff",
		Age: 36,
		ShoeSize: 10,
		CollarSize: 16,
		ElbowBreadth: 32,
		NoseCapacity: 10000.1289,
		Pi: math.Pi,
		UpsideDownCake: float32(math.Pi),
		Enrolled: false,
		ShouldBeEnrolled: true,
		Dispatched: true,
		Ratings: []int{0,1,2,3,4,5,6,7,8,9}}


	data, err = Marshal(target)
	c.Assert(err, IsNil)
	c.Assert(string(data[0:5]), Equals, "Geoff")
	c.Assert(string(data[5:7]), Equals, "36")
}
