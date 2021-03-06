package callerhandler

import (
	"fmt"
	"github.com/gorilla/schema"
	_ "io/ioutil"
	"net/http"
	"twiml"
	"utils"
	"webui"
)

type StatusCallbackRequest struct {
	CallDuration      string
	RecordingUrl      string
	RecordingSid      string
	RecordingDuration string
	CallStatus        string
	AccountSid        string
	CallSid           string
	From              string
}

type HangUpWrapper struct {
	Callerid chan twiml.Thingy
	Push     *[]chan webui.PushData
}

func (c HangUpWrapper) HangUpHandler(w http.ResponseWriter, r *http.Request) {
	//Check if it's a POST call, if not return immediately
	if r.Method != "POST" {
		return
	}
	//Error checking
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	// Unmarshal Form data
	var request StatusCallbackRequest
	decoder := schema.NewDecoder()
	decoder.Decode(&request, r.Form)

	userCount := utils.GetUserCount()

	//If the user has hung up, send n new request through the channel to dequeue the user
	if request.CallStatus == "completed" {
		c.Callerid <- twiml.Thingy{request.CallSid, "", false, request.From}
		fmt.Printf("%.6s has hung up\n", request.CallSid)
		fmt.Printf("There are %d total users\n", userCount)
	}

	pData := webui.PushData{UserCount: userCount}
	if webui.UseNumbers {
		pData.Call1Id = request.From
	} else {
		pData.Call1Id = request.CallSid
	}

	for _, j := range *c.Push {
		j <- pData
	}
}
