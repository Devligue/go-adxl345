# go-adxl345
Go package for ADXL345 sensor

### Basic Usage

```go
package main

import (
	"log"

	"github.com/devligue/go-adxl345"
)

func main() {
	adxl, err := adxl345.NewADXL345(
		1,                  // i2c interface index
		adxl345.AddressDef, // i2c address
	)
	if err != nil {
		log.Fatal(err)
	}
	defer adxl.Close()

	axes, err := adxl.GetAxesG()
	if err != nil {
		log.Println("Failed to get axes data [G].", err)
	}
	axes.Print()

	axes, err = adxl.GetAxesMS2()
	if err != nil {
		log.Println("Failed to get axes data [m/s^2].", err)
	}
	axes.Print()
}
```

### Documentation

[Look here!](https://godoc.org/github.com/Devligue/go-adxl345)

### Tips

##### Configuring I2C on Raspberry Pi

I2C on Raspberry Pi is not turned on by default. Use raspi-config to enable it.

* Run `sudo raspi-config`
* Use the down arrow to select *Advanced Options*.
* Arrow down to I2C.
* Select *yes* when it asks you to enable I2C.
* Also select yes if it asks about automatically loading the kernel module.
* Use the right arrow to select the *\<Finish>* button.
*  Select *yes* when asked to reboot.

##### Finding your device I2C interface index

If your device is at `/dev/i2c-1` then you should use index 1. As simple as that.

##### Finding your device I2C address

If you don't have it already, install `i2c-tools`from your package manager. Then after connecting your device run `i2cdetect -y 1`. Result should be self explanatory.

In the case of ADXL345, it usually has two addresses available. Depending on the manufacturer and/or model the default address and alternative address might be switched but one of them will most probably be **0x53** and the other **0x1D**. Look into your sensors datasheet for details.

### Compatibility

Tested with **ADXL345 GY-291**

