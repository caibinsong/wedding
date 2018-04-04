package main

import (
	"github.com/caibinsong/wedding/controllers"
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	controllers.Run()
}
