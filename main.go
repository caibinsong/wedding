package main

import (
	//"github.com/caibinsong/wedding/config"
	"github.com/caibinsong/wedding/controllers"
	"github.com/caibinsong/wedding/utils"
	"log"
)

func main() {
	utils.GetLogWriter().SetLogFile("./wedding.log")
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.SetOutput(utils.GetLogWriter())
	controllers.Run()
}
