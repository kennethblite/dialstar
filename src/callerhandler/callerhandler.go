package callerhandler

import (
	"bytes"
	"encoding/xml"
	_ "fmt"
	"github.com/gorilla/schema"
	_ "io/ioutil"
	"net/http"
	"twiml"
)

const (
	start = `<?xml version="1.0" encoding="UTF-8"?><Response>`
	end   = `</Response>`
)

type Say struct {
	XMLName  xml.Name `xml:"Say"`
	Voice    string   `xml:"voice,attr"`
	Language string   `xml:"language,attr"`
	Loop     uint     `xml:"loop,attr"`
	Text     string   `xml:",chardata"`
}
type context struct {
	b *bytes.Buffer
	r *http.Request
}

type VoiceRequest struct {
	CallSid       string
	AccountSid    string
	From          string
	To            string
	CallStatus    string
	ApiVersion    string
	Direction     string
	ForwardedFrom string
	CallerName    string
	FromCity      string
	FromState     string
	FromZip       string
	FromCountry   string
	ToCity        string
	ToState       string
	ToZip         string
	ToCountry     string
}

type StatusCallbackRequest struct {
	CallDuration      string
	RecordingUrl      string
	RecordingSid      string
	RecordingDuration string
}

func CallerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	// Unmarshal Form data
	var request VoiceRequest
	decoder := schema.NewDecoder()
	decoder.Decode(&request, r.Form)

	cityName := r.Form["FromCity"]
	//fmt.Println(actual)
	b := bytes.NewBufferString(start)
	say_response := &Say{Voice: "female", Language: "en", Loop: 1, Text: "Colin from " + cityName[0]}

	str, err := xml.Marshal(say_response)
	if err != nil {
		panic(err)
	}
	b.Write(str)
	response := &twiml.Conference{Text: "foobar", EndConferenceOnExit: "true"}
	dial := &twiml.Dial{Conference: *response}
	str, err = xml.Marshal(dial)
	if err != nil {
		panic(err)
	}
	b.Write(str)
	b.WriteString(end)
	b.WriteTo(w)
}
