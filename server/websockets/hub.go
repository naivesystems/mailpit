// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websockets

import (
	"encoding/json"

	"github.com/axllent/mailpit/utils/logger"
	"golang.org/x/exp/slices"
)

type Message struct {
	Data  []byte
	Users []string
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	Clients map[*Client]bool

	// Inbound messages from the clients.
	Broadcast chan Message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

// WebsocketNotification struct for responses
type WebsocketNotification struct {
	Type string
	Data interface{}
}

// NewHub returns a new hub configuration
func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
	}
}

// Run runs the listener
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.Clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.send)
			}
		case message := <-h.Broadcast:
			// logger.Log().Debugf("[broadcast] %s", message)
			for client := range h.Clients {
				if message.Users != nil && !slices.Contains(message.Users, client.user) {
					continue
				}
				select {
				case client.send <- message.Data:
				default:
					close(client.send)
					delete(h.Clients, client)
				}
			}
		}
	}
}

// Broadcast will spawn a broadcast message to all connected clients
func Broadcast(t string, msg interface{}, users []string) {
	if MessageHub == nil {
		return
	}

	w := WebsocketNotification{}
	w.Type = t
	w.Data = msg
	b, err := json.Marshal(w)

	if err != nil {
		logger.Log().Errorf("[http] broadcast received invalid data: %s", err)
	}

	m := Message{b, users}
	go func() { MessageHub.Broadcast <- m }()
}
