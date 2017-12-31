package adxl345

import (
	"fmt"
	"log"
	"math"

	"github.com/corrupt/go-smbus"
)

const (
	earthGravityMS2 = 9.80665
	scaleMultiplier = 0.004
	dataFormat      = 0x31
	powerCTL        = 0x2D
	measure         = 0x08
	axesData        = 0x32
	bwRate          = 0x2C
	Rate1600HZ      = 0x0F
	Rate800HZ       = 0x0E
	Rate400HZ       = 0x0D
	Rate200HZ       = 0x0C
	Rate100HZ       = 0x0B
	Rate50HZ        = 0x0A
	Rate25HZ        = 0x09
	Range2G         = 0x00
	Range4G         = 0x01
	Range8G         = 0x02
	Range16G        = 0x03
)

func round(f float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return math.Floor(f*shift+.5) / shift
}

type vector struct {
	x float64
	y float64
	z float64
}

func newVector() vector {
	return vector{
		x: 0,
		y: 0,
		z: 0,
	}
}

func (v vector) Print() {
	fmt.Printf("%+v\n", v)
}

type ADXL345 struct {
	bus     *smbus.SMBus
	Address byte
	Line    uint
}

func NewADXL345(line uint, address byte) (ADXL345, error) {
	smb, err := smbus.New(line, address)
	adxl345 := ADXL345{
		bus:     smb,
		Address: address,
		Line:    line,
	}
	if err != nil {
		log.Fatal(err)
		return adxl345, err
	}

	err = adxl345.SetBandwidthRate(Rate100HZ)
	if err != nil {
		log.Fatal(err)
		return adxl345, err
	}

	err = adxl345.SetRange(Range2G)
	if err != nil {
		log.Fatal(err)
		return adxl345, err
	}

	err = adxl345.EnableMeasurement()
	if err != nil {
		log.Fatal(err)
		return adxl345, err
	}

	return adxl345, err
}

func (a ADXL345) SetBandwidthRate(rateFlag byte) error {
	return a.bus.Write_byte_data(bwRate, rateFlag)
}

func (a ADXL345) SetRange(rangeFlag byte) error {
	value, err := a.bus.Read_byte_data(dataFormat)
	if err != nil {
		return err
	}

	value &= 0x0F
	value |= rangeFlag
	value |= 0x08

	return a.bus.Write_byte_data(dataFormat, value)
}

func (a ADXL345) EnableMeasurement() error {
	return a.bus.Write_byte_data(powerCTL, measure)
}

func (a ADXL345) GetAxesG() (vector, error) {
	buf := make([]byte, 6)
	_, err := a.bus.Read_i2c_block_data(axesData, buf)
	if err != nil {
		axes := newVector()
		return axes, err
	}

	x := buf[0] | (buf[1] << 8)
	xi := int16(x)

	y := buf[2] | (buf[3] << 8)
	yi := int16(y)

	z := buf[4] | (buf[5] << 8)
	zi := int16(z)

	axes := vector{
		x: round(float64(xi)*scaleMultiplier, 4),
		y: round(float64(yi)*scaleMultiplier, 4),
		z: round(float64(zi)*scaleMultiplier, 4),
	}

	return axes, err
}

func (a ADXL345) GetAxesAcceleration() (vector, error) {
	g_axes, err := a.GetAxesG()
	if err != nil {
		axes := newVector()
		return axes, err
	}

	axes := vector{
		x: round(g_axes.x*earthGravityMS2, 4),
		y: round(g_axes.y*earthGravityMS2, 4),
		z: round(g_axes.z*earthGravityMS2, 4),
	}

	return axes, err
}

func (a ADXL345) Close() {
	a.bus.Bus_close()
}
