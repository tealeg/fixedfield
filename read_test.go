package fixedfield

import (
	"bytes"
	"encoding/binary"
	. "launchpad.net/gocheck"
	"math"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type ReadSuite struct{}

var _ = Suite(&ReadSuite{})

type Target struct {
	Name           string  `length:"5"`
	Age            int     `length:"2" encoding:"ascii"`
	ShoeSize       int     `length:"2" encoding:"bigendian"`
	CollarSize     int     `length:"2" encoding:"le"`
	ElbowBreadth   uint    `length:"8" encoding:"le"`
	NoseCapacity   float64 `length:"6" encoding:"ascii"`
	Pi             float64 `length:"8" encoding:"le"`
	UpsideDownCake float32 `length:"4" encoding:"be"`
	Ratings        []int   `length:"1" repeat:"10"`
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

// Test readBinaryInteger decodes an 8bit, Little Endian value.
func (s *ReadSuite) TestReadBinaryInteger8BitLittleEndian(c *C) {
	block := []byte("\x10")
	blockLength := 1
	byteOrder := binary.LittleEndian
	value, err := readBinaryInteger(block, blockLength, byteOrder)
	c.Assert(err, IsNil)
	c.Assert(value, Equals, int64(16))
}

// Test readBinaryInteger decodes a 16bit, Little Endian value.
func (s *ReadSuite) TestReadBinaryInteger16BitLittleEndian(c *C) {
	block := []byte("\x10\x01")
	blockLength := 2
	byteOrder := binary.LittleEndian
	value, err := readBinaryInteger(block, blockLength, byteOrder)
	c.Assert(err, IsNil)
	c.Assert(value, Equals, int64(272))
}

// Test readBinaryInteger decodes a 32bit, Little Endian value.
func (s *ReadSuite) TestReadBinaryInteger32BitLittleEndian(c *C) {
	block := []byte("\x10\x10\x01\x01")
	blockLength := 4
	byteOrder := binary.LittleEndian
	value, err := readBinaryInteger(block, blockLength, byteOrder)
	c.Assert(err, IsNil)
	c.Assert(value, Equals, int64(16846864))
}

// Test readBinaryInteger decodes a 64bit, Little Endian value.
func (s *ReadSuite) TestReadBinaryInteger64BitLittleEndian(c *C) {
	block := []byte("\x10\x01\x10\x01\x10\x01\x10\x01")
	blockLength := 8
	byteOrder := binary.LittleEndian
	value, err := readBinaryInteger(block, blockLength, byteOrder)
	c.Assert(err, IsNil)
	c.Assert(value, Equals, int64(76562361914229008))
}

// Test readBinaryInteger decodes an 8bit, Little Endian negative value.
func (s *ReadSuite) TestReadBinaryInteger8BitLittleEndianNegative(c *C) {
	block := []byte("\xf0")
	blockLength := 1
	byteOrder := binary.LittleEndian
	value, err := readBinaryInteger(block, blockLength, byteOrder)
	c.Assert(err, IsNil)
	c.Assert(value, Equals, int64(-16))
}

// Test readBinaryInteger decodes a 16bit, Little Endian negative value.
func (s *ReadSuite) TestReadBinaryInteger16BitLittleEndianNegative(c *C) {
	block := []byte("\x10\xf1")
	blockLength := 2
	byteOrder := binary.LittleEndian
	value, err := readBinaryInteger(block, blockLength, byteOrder)
	c.Assert(err, IsNil)
	c.Assert(value, Equals, int64(-3824))
}

// Test readBinaryInteger decodes a 32bit, Little Endian negative value.
func (s *ReadSuite) TestReadBinaryInteger32BitLittleEndianNegative(c *C) {
	block := []byte("\x10\x10\x01\xf1")
	blockLength := 4
	byteOrder := binary.LittleEndian
	value, err := readBinaryInteger(block, blockLength, byteOrder)
	c.Assert(err, IsNil)
	c.Assert(value, Equals, int64(-251588592))
}

// Test readBinaryInteger decodes a 64bit, Little Endian negative value.
func (s *ReadSuite) TestReadBinaryInteger64BitLittleEndianNegative(c *C) {
	block := []byte("\x10\x01\x10\x01\x10\x01\x10\xf1")
	blockLength := 8
	byteOrder := binary.LittleEndian
	value, err := readBinaryInteger(block, blockLength, byteOrder)
	c.Assert(err, IsNil)
	c.Assert(value, Equals, int64(-1076359142692617968))
}

// Test readBinaryInteger decodes an 8bit, Big Endian value.
func (s *ReadSuite) TestReadBinaryInteger8BitBigEndian(c *C) {
	block := []byte("\x10")
	blockLength := 1
	byteOrder := binary.BigEndian
	value, err := readBinaryInteger(block, blockLength, byteOrder)
	c.Assert(err, IsNil)
	c.Assert(value, Equals, int64(16))
}

// Test readBinaryInteger decodes a 16bit, Big Endian value.
func (s *ReadSuite) TestReadBinaryInteger16BitBigEndian(c *C) {
	block := []byte("\x10\x01")
	blockLength := 2
	byteOrder := binary.BigEndian
	value, err := readBinaryInteger(block, blockLength, byteOrder)
	c.Assert(err, IsNil)
	c.Assert(value, Equals, int64(4097))
}

// Test readBinaryInteger decodes a 32bit, Big Endian value.
func (s *ReadSuite) TestReadBinaryInteger32BitBigEndian(c *C) {
	block := []byte("\x10\x10\x01\x01")
	blockLength := 4
	byteOrder := binary.BigEndian
	value, err := readBinaryInteger(block, blockLength, byteOrder)
	c.Assert(err, IsNil)
	c.Assert(value, Equals, int64(269484289))
}

// Test readBinaryInteger decodes a 64bit, Big Endian value.
func (s *ReadSuite) TestReadBinaryInteger64BitBigEndian(c *C) {
	block := []byte("\x10\x01\x10\x01\x10\x01\x10\x01")
	blockLength := 8
	byteOrder := binary.BigEndian
	value, err := readBinaryInteger(block, blockLength, byteOrder)
	c.Assert(err, IsNil)
	c.Assert(value, Equals, int64(1153220576333074433))
}

// Test readBinaryInteger decodes an 8bit, Big Endian negative value.
func (s *ReadSuite) TestReadBinaryInteger8BitBigEndianNegative(c *C) {
	block := []byte("\xf0")
	blockLength := 1
	byteOrder := binary.BigEndian
	value, err := readBinaryInteger(block, blockLength, byteOrder)
	c.Assert(err, IsNil)
	c.Assert(value, Equals, int64(-16))
}

// Test readBinaryInteger decodes a 16bit, Big Endian negative value.
func (s *ReadSuite) TestReadBinaryInteger16BitBigEndianNegative(c *C) {
	block := []byte("\xf0\x01")
	blockLength := 2
	byteOrder := binary.BigEndian
	value, err := readBinaryInteger(block, blockLength, byteOrder)
	c.Assert(err, IsNil)
	c.Assert(value, Equals, int64(-4095))
}

// Test readBinaryInteger decodes a 32bit, Big Endian negative value.
func (s *ReadSuite) TestReadBinaryInteger32BitBigEndianNegative(c *C) {
	block := []byte("\xf0\x10\x01\x01")
	blockLength := 4
	byteOrder := binary.BigEndian
	value, err := readBinaryInteger(block, blockLength, byteOrder)
	c.Assert(err, IsNil)
	c.Assert(value, Equals, int64(-267386623))
}

// Test readBinaryInteger decodes a 64bit, Big Endian negative value.
func (s *ReadSuite) TestReadBinaryInteger64BitBigEndianNegative(c *C) {
	block := []byte("\xf0\x01\x10\x01\x10\x01\x10\x01")
	blockLength := 8
	byteOrder := binary.BigEndian
	value, err := readBinaryInteger(block, blockLength, byteOrder)
	c.Assert(err, IsNil)
	c.Assert(value, Equals, int64(-1152622432880619519))
}

// Test readBinaryUnsignedInteger decodes an 8bit, Little Endian value.
func (s *ReadSuite) TestReadBinaryUnsignedInteger8BitLittleEndian(c *C) {
	block := []byte("\x10")
	blockLength := 1
	byteOrder := binary.LittleEndian
	value, err := readBinaryUnsignedInteger(block, blockLength, byteOrder)
	c.Assert(err, IsNil)
	c.Assert(value, Equals, uint64(16))
}

// Test readBinaryUnsignedInteger decodes a 16bit, Little Endian value.
func (s *ReadSuite) TestReadBinaryUnsignedInteger16BitLittleEndian(c *C) {
	block := []byte("\x10\x01")
	blockLength := 2
	byteOrder := binary.LittleEndian
	value, err := readBinaryUnsignedInteger(block, blockLength, byteOrder)
	c.Assert(err, IsNil)
	c.Assert(value, Equals, uint64(272))
}

// Test readBinaryUnsignedInteger decodes a 32bit, Little Endian value.
func (s *ReadSuite) TestReadBinaryUnsignedInteger32BitLittleEndian(c *C) {
	block := []byte("\x10\x10\x01\x01")
	blockLength := 4
	byteOrder := binary.LittleEndian
	value, err := readBinaryUnsignedInteger(block, blockLength, byteOrder)
	c.Assert(err, IsNil)
	c.Assert(value, Equals, uint64(16846864))
}

// Test readBinaryUnsignedInteger decodes a 64bit, Little Endian value.
func (s *ReadSuite) TestReadBinaryUnsignedInteger64BitLittleEndian(c *C) {
	block := []byte("\x10\x01\x10\x01\x10\x01\x10\x01")
	blockLength := 8
	byteOrder := binary.LittleEndian
	value, err := readBinaryUnsignedInteger(block, blockLength, byteOrder)
	c.Assert(err, IsNil)
	c.Assert(value, Equals, uint64(76562361914229008))
}

// Test readBinaryUnsignedInteger decodes an 8bit, Big Endian value.
func (s *ReadSuite) TestReadBinaryUnsignedInteger8BitBigEndian(c *C) {
	block := []byte("\x10")
	blockLength := 1
	byteOrder := binary.BigEndian
	value, err := readBinaryUnsignedInteger(block, blockLength, byteOrder)
	c.Assert(err, IsNil)
	c.Assert(value, Equals, uint64(16))
}

// Test readBinaryUnsignedInteger decodes a 16bit, Big Endian value.
func (s *ReadSuite) TestReadBinaryUnsignedInteger16BitBigEndian(c *C) {
	block := []byte("\x10\x01")
	blockLength := 2
	byteOrder := binary.BigEndian
	value, err := readBinaryUnsignedInteger(block, blockLength, byteOrder)
	c.Assert(err, IsNil)
	c.Assert(value, Equals, uint64(4097))
}

// Test readBinaryUnsignedInteger decodes a 32bit, Big Endian value.
func (s *ReadSuite) TestReadBinaryUnsignedInteger32BitBigEndian(c *C) {
	block := []byte("\x10\x10\x01\x01")
	blockLength := 4
	byteOrder := binary.BigEndian
	value, err := readBinaryUnsignedInteger(block, blockLength, byteOrder)
	c.Assert(err, IsNil)
	c.Assert(value, Equals, uint64(269484289))
}

// Test readBinaryUnsignedInteger decodes a 64bit, Big Endian value.
func (s *ReadSuite) TestReadBinaryUnsignedInteger64BitBigEndian(c *C) {
	block := []byte("\x10\x01\x10\x01\x10\x01\x10\x01")
	blockLength := 8
	byteOrder := binary.BigEndian
	value, err := readBinaryUnsignedInteger(block, blockLength, byteOrder)
	c.Assert(err, IsNil)
	c.Assert(value, Equals, uint64(1153220576333074433))
}

// Test readASCIIInteger with positive value
func (s *ReadSuite) TestReadASCIIIntegerPositive(c *C) {
	block := []byte("4096")
	value, err := readASCIIInteger(block)
	c.Assert(err, IsNil)
	c.Assert(value, Equals, int64(4096))
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
