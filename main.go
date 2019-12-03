package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

//PingMessage return true or false
type PingMessage struct {
	Message string `json:"message"`
	Ping    bool   `json:"ping"`
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func greet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Welcome to Appcelerator Arrow Cloud! cuurent time: %s", time.Now())
}

//Hello world
func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "Welcome to Appcelerator Arrow Cloud! Hello %s!\n", ps.ByName("name"))
}

//ArrowPing respond with ok
func ArrowPing(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ping := &PingMessage{Ping: true, Message: "service is reachable"}
	js, err := json.Marshal(ping)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

//Log number of lines to console
func Log(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	finished := make(chan bool)
	lines, err := strconv.Atoi(ps.ByName("number"))
	if err != nil {
		fmt.Printf("error while converting to number: %v", err)
	}
	go func() {
		for index := 0; index < lines; index++ {
			fmt.Printf("Hello World! %s", time.Now())
		}
		finished <- true
	}()
	isFinished := <-finished
	fmt.Printf("finished: %v %v lines printing", isFinished, lines)
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("number"))
}

func main() {
	var port = getenv("PORT", "8080")
	router := httprouter.New()
	router.GET("/", greet)
	router.GET("/hello/:name", Hello)
	router.GET("/log/:number", Log)
	router.GET("/arrowPing.json", ArrowPing)
	fmt.Println("server listening on port: ", port)
	http.ListenAndServe(":"+port, router)
}
