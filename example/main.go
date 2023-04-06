package main

import "github.com/LordCasser/reception"

func main() {
	rec := reception.New()
	_ = rec.AddSwitch("cloud.lordcasser.com", "127.0.0.1:8001", false)
	_ = rec.AddSwitch("lordcasser.com", "127.0.0.1:8002", true)
	rec.Serve()
}
