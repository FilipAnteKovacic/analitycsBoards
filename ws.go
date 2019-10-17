package main

import (
	"encoding/json"
	fmt "fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/gorilla/websocket"
	"github.com/sacOO7/gowebsocket"
)

// WSMessage receive message
type WSMessage struct {
	Sender string `bson:"sender" json:"sender"`
	Action string `bson:"action" json:"action"`
	Type   string `bson:"type" json:"type"`
	Chart  Chart  `bson:"chart" json:"chart"`
	Source Source `bson:"source" json:"source"`
}

// WSListClient send list message
type WSListClient struct {
	Type    string    `bson:"type" json:"type"`
	Sources []*Source `bson:"sources" json:"sources"`
	Charts  []*Chart  `bson:"charts" json:"charts"`
}

// WSChartClient send list message
type WSChartClient struct {
	Type   string `bson:"type" json:"type"`
	Action string `bson:"action" json:"action"`
	Chart  Chart  `bson:"chart" json:"chart"`
}

// WSErrorClient send list message
type WSErrorClient struct {
	Type  string `bson:"type" json:"type"`
	Error string `bson:"error" json:"error"`
}

// Init socket conn for grpc
var wsConn gowebsocket.Socket

func wsConection() {

	wsConn = gowebsocket.New("ws://localhost:" + os.Getenv("UI") + "/ws")

	wsConn.OnConnectError = func(err error, socket gowebsocket.Socket) {
		log.Println("Received connect error - ", err)
	}
	wsConn.OnConnected = func(socket gowebsocket.Socket) {
		log.Println("Connected to server")
	}

	wsConn.Connect()

}

func sendWSMessage(wsMsg WSMessage) {

	if wsConn.IsConnected {

		msg, err := json.Marshal(wsMsg)
		if err != nil {
			fmt.Println("Error while parsing WS message", err)
			return
		}

		wsConn.SendText(string(msg))

	}

	return
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// Hub WSclients hub
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				//client.conn.Close()
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {

				if err := client.conn.WriteMessage(1, message); err != nil {
					return
				}
			}
		}
	}
}

// BrodcastMsg send msg to clients
func BrodcastMsg(msg []byte) {

	wsHub.broadcast <- msg

}

var wsHub *Hub

type wsPage struct {
	WSURL string
}

func startUIServer(uiAddress string) error {

	// Init ws path
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		uiP := wsPage{
			WSURL: "ws://localhost:" + os.Getenv("UI") + "/ws",
		}

		crudTemplate, err := template.ParseFiles(
			"ui/admin.html",
			"ui/admin/sources.html",
			"ui/admin/charts.html",
			"ui/admin/boards.html",
		)
		if err != nil {
			fmt.Println("Error occurred while parsing template", err)
			return
		}

		err = crudTemplate.Execute(w, &uiP)
		if err != nil {
			fmt.Println("Error occurred while executing the template  or writing its output", err)
			return
		}

	})

	wsHub = newHub()
	go wsHub.run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{}

		upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}

		// Init client
		client := &Client{hub: wsHub, conn: c, send: make(chan []byte, 256)}

		// Register client in hub channel
		wsHub.register <- client

		go func() {

			defer c.Close()

			for {

				// Read incoming msg
				_, message, err := c.ReadMessage()
				if err != nil {
					wsHub.unregister <- client
					break
				}

				var wsMsg WSMessage

				json.Unmarshal(message, &wsMsg)

				switch wsMsg.Sender {
				case "server":

					BrodcastMsg(wsListMsg())

					switch wsMsg.Action {
					case "create":
						fallthrough
					case "update":
						fallthrough
					case "delete":

						switch wsMsg.Type {
						case "chart":

							wsChartController(WSMessage{
								Sender: "client",
								Action: "read",
								Type:   "chart",
								Chart:  wsMsg.Chart,
							})

							break
						}

						break
					case "updateData":

						switch wsMsg.Type {
						case "chart":

							wsChartController(WSMessage{
								Sender: "client",
								Action: "updateData",
								Type:   "chart",
								Chart:  wsMsg.Chart,
							})

							break
						}

						break
					}

					break

				case "client":

					if wsMsg.Action == "list" {
						BrodcastMsg(wsListMsg())
					}

					switch wsMsg.Type {
					case "source":
						wsSourceController(wsMsg)
						break
					case "chart":
						wsChartController(wsMsg)
						break
					}

					break
				}

			}

		}()

	})

	fs := http.FileServer(http.Dir("ui/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Printf("starting HTTP/1.1 UI server on %s", uiAddress)
	err := http.ListenAndServe(uiAddress, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	return nil
}

func wsListMsg() []byte {

	// List source
	sourceList, err := listSource()
	if err != nil {
		fmt.Println("error while getting list sources", err)
	}

	// List charts
	chartsList, err := listChart()
	if err != nil {
		fmt.Println("error while getting list charts", err)
	}

	// List charts to json
	wsClientMsg, err := json.Marshal(WSListClient{
		Type:    "list",
		Sources: sourceList,
		Charts:  chartsList,
	})

	if err != nil {
		fmt.Println("cannot marshal data", err)
	}

	return wsClientMsg

}

func wsSourceController(wsMsg WSMessage) {

	switch wsMsg.Action {
	case "create":

		// Create data in DB
		err := createSource(wsMsg.Source)
		if err != nil {

			fmt.Println("error by create", err)

		}

		dataUpdateRestart()

		// Push msg to websocket clients
		sendWSMessage(WSMessage{
			Sender: "server",
			Action: "create",
		})

		break
	case "update":

		Unique := wsMsg.Source.Unique

		// Update in DB
		err := updateByUniqueSource(Unique, wsMsg.Source)
		if err != nil {

			fmt.Println("error by delete", err)
		}

		dataUpdateRestart()

		// Push msg to websocket clients
		sendWSMessage(WSMessage{
			Sender: "server",
			Action: "update",
			Source: wsMsg.Source,
		})

		break

	case "delete":

		// Delete data in DB
		err := deleteByUniqueSource(wsMsg.Source.Unique)
		if err != nil {

			fmt.Println("error by delete", err)

		}

		dataUpdateRestart()

		// Push msg to websocket clients
		sendWSMessage(WSMessage{
			Sender: "server",
			Action: "delete",
		})

		break
	}

}

func wsChartController(wsMsg WSMessage) {

	switch wsMsg.Action {
	case "read":

		// Create data in DB
		chart, err := readByUniqueChart(wsMsg.Chart.Unique)
		if err != nil {

			fmt.Println("error by create", err)

		}

		// List charts to json
		wsClientMsg, err := json.Marshal(WSChartClient{
			Type:  "chart",
			Chart: chart,
		})

		if err != nil {
			fmt.Println("cannot marshal data", err)
		}

		// Push msg to websocket clients
		BrodcastMsg(wsClientMsg)

		break

	case "updateData":

		// Create data in DB
		chart, err := readByUniqueChart(wsMsg.Chart.Unique)
		if err != nil {

			fmt.Println("error by create", err)

		}

		// List charts to json
		wsClientMsg, err := json.Marshal(WSChartClient{
			Type:   "chart",
			Action: "updateData",
			Chart:  chart,
		})

		if err != nil {
			fmt.Println("cannot marshal data", err)
		}

		// Push msg to websocket clients
		BrodcastMsg(wsClientMsg)

		break

	case "create":

		// Create data in DB
		unique, err := createChart(wsMsg.Chart)
		if err != nil {

			fmt.Println("error by create", err)

		}

		dataUpdateRestart()

		// Push msg to websocket clients
		sendWSMessage(WSMessage{
			Sender: "server",
			Action: "create",
			Type:   "chart",
			Chart: Chart{
				Unique: unique,
			},
		})

		break
	case "update":

		Unique := wsMsg.Chart.Unique

		// Update in DB
		err := updateByUniqueChart(Unique, wsMsg.Chart)
		if err != nil {

			fmt.Println("error by delete", err)
		}

		dataUpdateRestart()
		updateChartDataByUnique(Unique)

		// Push msg to websocket clients
		sendWSMessage(WSMessage{
			Sender: "server",
			Action: "update",
			Type:   "chart",
			Chart:  wsMsg.Chart,
		})

		break

	case "delete":

		// Delete data in DB
		err := deleteByUniqueChart(wsMsg.Chart.Unique)
		if err != nil {

			fmt.Println("error by delete", err)

		}

		dataUpdateRestart()

		// Push msg to websocket clients
		sendWSMessage(WSMessage{
			Sender: "server",
			Action: "delete",
			Type:   "chart",
			Chart: Chart{
				Unique: wsMsg.Chart.Unique,
			},
		})

		break
	}

}
