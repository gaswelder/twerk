package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func runGUI() {
	src, err := ioutil.ReadFile("/home/gas/go/src/github.com/gaswelder/twerk/gui.html")
	if err != nil {
		log.Fatal(err)
	}
	tpl := template.Must(template.New("").Parse(string(src)))
	http.HandleFunc("/s", handleWebSockets)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, "ws://"+r.Host+"/echo")
	})
	log.Println("Serving GUI admin at http://localhost:8080")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func handleWebSockets(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	cfg, err := parseConfig("twerks.json")
	if err != nil {
		c.WriteJSON(err.Error())
		c.Close()
		return
	}

	for k := range cfg {
		c.WriteJSON(k)
	}

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
