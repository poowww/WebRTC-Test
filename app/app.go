package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/websocket"
)

var conns []*websocket.Conn

func main() {
	conns = make([]*websocket.Conn, 0, 100)
	http.Handle("/", http.FileServer(http.Dir("./static/")))

	http.HandleFunc("/ws/", func(w http.ResponseWriter, req *http.Request) {
		s := websocket.Server{Handler: websocket.Handler(wsHandler)}
		s.ServeHTTP(w, req)
	})

	http.ListenAndServe(":8080", nil)
}

func wsHandler(conn *websocket.Conn) {
	conns = append(conns, conn)
	wsMsgHandler(conn)
}

func wsMsgHandler(conn *websocket.Conn) {
	for {
		var msg string
		err := websocket.Message.Receive(conn, &msg)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("message = %s", msg)
		broadcastRoom(msg, conn)
	}
}

func broadcastRoom(msg string, from *websocket.Conn) {
	for _, conn := range conns {
		if conn != from {
			websocket.Message.Send(conn, msg)
		}
	}
}
