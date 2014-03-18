package fixedfield

import (
	"fmt"
	. "launchpad.net/gocheck"
	"math"
)

type WriteSuite struct{}

var _ = Suite(&WriteSuite{})

func (s *WriteSuite) TestMarshalIntegerASCIIPositive(c *C) {
	type target1 struct {
		Value int `encoding:"ascii" length:"1"`
	}
	var t *target1
	var block []byte

	t = &target1{Value: 3}
	specs, err := buildSpecs(t)
	c.Assert(err, IsNil)
	block, err = marshalInteger(specs[0])
	c.Assert(err, IsNil)
	c.Assert(string(block), Equals, "3")
}

func (s *WriteSuite) TestMarshalIntegerASCIINegative(c *C) {
	type target2 struct {
		Value int `encoding:"ascii" length:"2"`
	}
	var t *target2
	var block []byte

	t = &target2{Value: -3}
	specs, err := buildSpecs(t)
	c.Assert(err, IsNil)
	block, err = marshalInteger(specs[0])
	c.Assert(err, IsNil)
	c.Assert(string(block), Equals, "-3")
}


func (s *WriteSuite) TestMarshalIntegerASCIIPositivePad(c *C) {
	type target3 struct {
		Value int `encoding:"ascii" length:"10"`
	}
	var t *target3
	var block []byte

	t = &target3{Value: 33}
	specs, err := buildSpecs(t)
	fmt.Printf("Spec.length %d\n", specs[0].Length)
	c.Assert(err, IsNil)
	block, err = marshalInteger(specs[0])
	c.Assert(err, IsNil)
	c.Assert(string(block), Equals, "        33")
}

func (s *WriteSuite) TestMarshalIntegerASCIINegativePad(c *C) {
	type target4 struct {
		Value int `encoding:"ascii" length:"10"`
	}

	var t * target4
	var block []byte

	t = &target4{Value: -44}
	specs, err := buildSpecs(t)
	fmt.Printf("Specs.length %d\n", specs[0].Length)
	c.Assert(err, IsNil)
	block, err = marshalInteger(specs[0])
	c.Assert(err, IsNil)
	c.Assert(string(block), Equals, "       -44")
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
	c.Assert(string(data[5:17]), Equals, "          36")
}
