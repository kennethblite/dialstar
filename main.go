package main

import (
	"callerhandler"
	"container/list"
	"fmt"
	"net/http"
	"net/url"
)

var user_queue *list.List

//This is just all the handlers and shit
func main() {
	callers_waiting := make(chan string, 10)
	Conf_waiters := callerhandler.CallerWrapper{Callerid: callers_waiting}
	go PollWaiters(callers_waiting)
	http.HandleFunc("/caller/", Conf_waiters.CallerHandler)
	http.HandleFunc("/conference/", callerhandler.ConferenceHandler)
	http.ListenAndServe(":3000", nil)
}

func PollWaiters(c chan string) {
	user_queue = list.New()

	for element := range c {
		_ = user_queue.PushBack(element)
		if user_queue.Len() >= 2 {
			fmt.Println("[META] - Got two or more people.")
			first := user_queue.Front()
			user_queue.Remove(first)
			f := first.Value.(string)
			second := user_queue.Front()
			user_queue.Remove(second)
			s := second.Value.(string)
			fmt.Println("[META] Paired " + f + " with " + s)

			ConferenceId := f + s
			authToken := "79ed2712d0cf06c87aa2783eee6aaa7a"
			accountId := "AC6f0fa1837933462d780f6fc1daf57d44"

			values := make(url.Values)
			values.Set("Url", "http://twilio.axyjo.com/conference/?ConferenceId="+ConferenceId+"&OtherCity=Ottawa")
			http.PostForm("https://"+accountId+":"+authToken+"@api.twilio.com/2010-04-01/Accounts/"+accountId+"/Calls/"+f, values)
			http.PostForm("https://"+accountId+":"+authToken+"@api.twilio.com/2010-04-01/Accounts/"+accountId+"/Calls/"+s, values)
		}
	}
}
