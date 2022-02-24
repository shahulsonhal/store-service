package main

import "github.com/shahulsonhal/store-service/internal/app"

func main() {
	s := app.NewServer()
	s.Start()
}
