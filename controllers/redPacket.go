package controllers

import (
	"net/http"
	"github.com/Amniversary/wedding-logic-redpacket/config"
	"encoding/json"
	"log"
	"github.com/Amniversary/wedding-logic-redpacket/business"
)

func GenRedPacket(w http.ResponseWriter, r *http.Request) {
	Response := &config.Response{Code: config.RESPONSE_ERROR}
	defer func() {
		EchoJson(w, http.StatusOK, Response)
	}()
	req := &config.GenRedPacket{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		log.Printf("genRedPacket json decode err: %v", err)
		return
	}
	redFlash, ok, err := business.GenRedPacket(req.RedPacketNum, req.RedPacketMoney)
	if !ok {
		log.Printf("genRedPacket err: %s", err)
		return
	}

	log.Printf("%v", redFlash)
}
