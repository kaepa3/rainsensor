package main

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/ashwanthkumar/slack-go-webhook"

	"fmt"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

const (
	On  = 1
	Off = 0
)

var conf AppConfig

type AppConfig struct {
	WebhookURL string
}

func main() {
	r := raspi.NewAdaptor()
	// 初期状態の取得
	isNotRain, _ := r.DigitalRead("12")
	sensor := gpio.NewPIRMotionDriver(r, "12")

	work := func() {
		// センサOn
		sensor.On(gpio.MotionDetected, func(data interface{}) {
			fmt.Println(gpio.MotionDetected)

			if isNotRain == Off {
				fmt.Println("sensor on")
			}
			isNotRain = On
		})

		// センサOff
		sensor.On(gpio.MotionStopped, func(data interface{}) {
			fmt.Println(gpio.MotionStopped)
			if isNotRain == On {
				fmt.Println("sensor off")
			}
			isNotRain = Off
		})

	}

	robot := gobot.NewRobot("rainbot",
		[]gobot.Connection{r},
		[]gobot.Device{sensor},
		work,
	)

	robot.Start()
}
func readConfig() bool {
	path := "config.toml"
	if !Exists(path) {
		return false
	}

	toml.DecodeFile("config.toml", &conf)
	return true
}
func slacker(text string) {
	if readConfig() {
		payload := slack.Payload{
			Text:      text,
			Username:  "robot",
			Channel:   "#general",
			IconEmoji: ":monkey_face:",
		}
		err := slack.Send(conf.WebhookURL, "", payload)
		if len(err) > 0 {
			fmt.Printf("error: %s\n", err)
		}
	}

}

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
