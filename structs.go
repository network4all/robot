package main

// Message object...
type Message struct {
	MessageID   string `json:"messageid"`   // timestamp+node
	MessageType int    `json:"messagetype"` // ping, sendobject, ...
	Source      string `json:"source"`      // node, serial mac
	Destination string `json:"destination"` // broadcast, serial mac
	Message     string `json:"message"`     // json object data
	Data        string `json:"data"`        // json object data
	Ack         bool   `json:"ack"`         // read ack (tcp/udp) (true/false)
}
