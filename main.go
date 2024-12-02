package main

import (
	"github.com/anikethz/HertzDB/src/web"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	web.StartServer()
}
