package main

import (
<<<<<<< HEAD
	//"github.com/caibinsong/wedding/config"
	"github.com/caibinsong/wedding/controllers"
	"github.com/caibinsong/wedding/utils"
=======
<<<<<<< HEAD
	//"github.com/caibinsong/wedding/config"
	"github.com/caibinsong/wedding/controllers"
	"github.com/caibinsong/wedding/utils"
=======
	"./utils"
	"github.com/Amniversary/wedding-logic-redpacket/controllers"
>>>>>>> dd12374ac95f08e4145cdb3fa4b628e5d98bd4d3
>>>>>>> a64d7c5df01427534bebc1ec23b5463de6ce4777
	"log"
)

func main() {
	utils.GetLogWriter().SetLogFile("./wedding.log")
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.SetOutput(utils.GetLogWriter())
	controllers.Run()
}
