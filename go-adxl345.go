// Package adxl345 is a Go library for ADXL345
package adxl345

import (
	"fmt"
	"math"

	"github.com/corrupt/go-smbus"
)

// Available I2C addresses
const (
	AddressDef = 0x53
	AddressAlt = 0x1D
)

// Earth Gravity constant in [m/s^2]
const earthGravityMS2 = 9.80665

// The typical scale factor in g/LSB
const scaleMultiplier = 0.0039

// ADXL345 Registers
const (
	dataFormat = 0x31
	bwRate     = 0x2C
	powerCTL   = 0x2D
	measure    = 0x08
)

// Device bandwidth and output data rates
const (
	Rate1600HZ = 0x0F
	Rate800HZ  = 0x0E
	Rate400HZ  = 0x0D
	Rate200HZ  = 0x0C
	Rate100HZ  = 0x0B
	Rate50HZ   = 0x0A
	Rate25HZ   = 0x09
)

// Measurement Range
const (
	Range2G  = 0x00
	Range4G  = 0x01
	Range8G  = 0x02
	Range16G = 0x03
)

// Axes Data
const (
	dataX0 = 0x32
	dataX1 = 0x33
	dataY0 = 0x34
	dataY1 = 0x35
	dataZ0 = 0x36
	dataZ1 = 0x37
)

func round(f float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return math.Floor(f*shift+.5) / shift
}

// Vector is a simple vector struct
type Vector struct {
	X float64
	Y float64
	Z float64
}

// NewVector is a factory function creating instance of Vector
func NewVector() Vector {
	return Vector{
		X: 0,
		Y: 0,
		Z: 0,
	}
}

// Print will nicely print the acceleration vector
func (v Vector) Print() {
	fmt.Printf("%+v\n", v)
}

// ADXL345 is a struct holding the device I2C address and
// I2C interface. It has associated several methods allowing
// to set up connection with ADXL345 over I2C and read
// measurement data.
type ADXL345 struct {
	bus          *smbus.SMBus
	Address      byte
	InterfaceIdx uint
}

// NewADXL345 is a factory method creating instance of ADXL345,
// setting base Bandwidth (100HZ), Range (2G) and enabling the
// measurement
func NewADXL345(interfaceIdx uint, address byte) (ADXL345, error) {
	smb, err := smbus.New(interfaceIdx, address)
	adxl345 := ADXL345{
		bus:          smb,
		Address:      address,
		InterfaceIdx: interfaceIdx,
	}
	if err != nil {
		return adxl345, err
	}

	err = adxl345.SetBandwidthRate(Rate100HZ)
	if err != nil {
		return adxl345, err
	}

	err = adxl345.SetRange(Range2G)
	if err != nil {
		return adxl345, err
	}

	err = adxl345.EnableMeasurement()
	if err != nil {
		return adxl345, err
	}

	return adxl345, err
}

// SetBandwidthRate changes the device bandwidth and output data rate.
func (a ADXL345) SetBandwidthRate(newRate byte) error {
	return a.bus.Write_byte_data(bwRate, newRate)
}

// SetRange changes the range of ADXL345. Available ranges are 2G,
// 4G, 8G and 16G.
func (a ADXL345) SetRange(newRange byte) error {
	retval, err := a.bus.Read_byte_data(dataFormat)
	if err != nil {
		return err
	}

	value := int32(retval)
	value &= ^0x0F
	value |= int32(newRange)
	value |= 0x08

	return a.bus.Write_byte_data(dataFormat, byte(value))
}

// EnableMeasurement enables measurement on ADXL345
func (a ADXL345) EnableMeasurement() error {
	return a.bus.Write_byte_data(powerCTL, measure)
}

// GetAxesG retrives axes acceleration data from ADXL345. Values
// are returned as multiplications of G
func (a ADXL345) GetAxesG() (Vector, error) {
	buf := make([]byte, 6)
	_, err := a.bus.Read_i2c_block_data(dataX0, buf)
	if err != nil {
		axes := NewVector()
		return axes, err
	}

	x := int16(buf[0]) | (int16(buf[1]) << 8)
	y := int16(buf[2]) | (int16(buf[3]) << 8)
	z := int16(buf[4]) | (int16(buf[5]) << 8)

	axes := Vector{
		X: round(float64(x)*scaleMultiplier, 4),
		Y: round(float64(y)*scaleMultiplier, 4),
		Z: round(float64(z)*scaleMultiplier, 4),
	}

	return axes, err
}

// GetAxesMS2 parses data returned by GetAxesG and returns
// them in [m/s^2]
func (a ADXL345) GetAxesMS2() (Vector, error) {
	gAxes, err := a.GetAxesG()
	if err != nil {
		axes := NewVector()
		return axes, err
	}

	axes := Vector{
		X: round(gAxes.X*earthGravityMS2, 4),
		Y: round(gAxes.Y*earthGravityMS2, 4),
		Z: round(gAxes.Z*earthGravityMS2, 4),
	}

	return axes, err
}

// Close disconnects from the device
func (a ADXL345) Close() {
	a.bus.Bus_close()
}
