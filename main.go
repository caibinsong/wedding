package main

import (
	"./utils"
	"github.com/Amniversary/wedding-logic-redpacket/controllers"
	"log"
)

func main() {
	utils.GetLogWriter().SetLogFile("./wedding.log")
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.SetOutput(utils.GetLogWriter())
	controllers.Run()
}
