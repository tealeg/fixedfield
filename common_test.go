package fixedfield

import (
	. "launchpad.net/gocheck"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }


type Target struct {
	Name             string  `length:"5"`
	Age              int     `length:"12" encoding:"ascii" padding:" "`
	ShoeSize         int     `length:"2" encoding:"bigendian"`
	CollarSize       int     `length:"2" encoding:"le"`
	ElbowBreadth     uint    `length:"8" encoding:"le"`
	NoseCapacity     float64 `length:"6" encoding:"ascii"`
	Pi               float64 `length:"8" encoding:"le"`
	UpsideDownCake   float32 `length:"4" encoding:"be"`
	Enrolled         bool
	ShouldBeEnrolled bool  `encoding:"ascii"`
	Dispatched       bool  `encoding:"ascii" trueChars:"jJ"`
	Ratings          []int `length:"1" repeat:"10" encoding:"ascii"`
}

type Person struct {
	Name string `length:"5"`
	Age  int    `length:"1"`
}

type Transaction struct {
	Buyer  Person
	Seller Person
}

