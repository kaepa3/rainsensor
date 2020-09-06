package main

import (
	"fmt"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

const (
	On  = 1
	Off = 0
)

func main() {
	r := raspi.NewAdaptor()
	// 初期状態の取得
	isRain, _ := r.DigitalRead("12")
	sensor := gpio.NewPIRMotionDriver(r, "12")

	work := func() {
		// センサOnになったら
		sensor.On(gpio.MotionDetected, func(data interface{}) {
			fmt.Println(gpio.MotionDetected)

			if isRain == Off {
				fmt.Println("sensor on")
			}
			isRain = On
		})

		// センサOffになったら
		sensor.On(gpio.MotionStopped, func(data interface{}) {
			fmt.Println(gpio.MotionStopped)
			if isRain == On {
				fmt.Println("sensor off")
			}
			isRain = Off
		})

	}

	robot := gobot.NewRobot("rainbot",
		[]gobot.Connection{r},
		[]gobot.Device{sensor},
		work,
	)

	robot.Start()
}
