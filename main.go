package main

import (
<<<<<<< HEAD
	//"github.com/caibinsong/wedding/config"
	"github.com/caibinsong/wedding/controllers"
	"github.com/caibinsong/wedding/utils"
=======
	"./utils"
	"github.com/Amniversary/wedding-logic-redpacket/controllers"
>>>>>>> dd12374ac95f08e4145cdb3fa4b628e5d98bd4d3
	"log"
)

func main() {
	utils.GetLogWriter().SetLogFile("./wedding.log")
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.SetOutput(utils.GetLogWriter())
	controllers.Run()
}
