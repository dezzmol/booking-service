package main

import (
	"booking-service/cmd/app"
)

func main() {
	app := app.New()
	app.Run()
}
