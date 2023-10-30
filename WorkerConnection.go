package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
)

type WorkerConnection struct {
	ID   int // incremental
	Addr string
	C    chan any
	conn net.Conn

	master *Master
}

// 2 goroutines for 1 worker connection
func (c *WorkerConnection) Run() {
	//receive message from worker
	go c.ListenWorkers()
	for {
		for serverMessage := range c.C {
			c.SendToWorker(serverMessage)
		}
	}
}

func (c *WorkerConnection) ListenWorkers() {
	buf := make([]byte, 8192)
	for {
		n, err := c.conn.Read(buf)

		if err != nil && err != io.EOF {
			fmt.Println("read error: ", err)
			return
		}

		c.ProcessWorkerMessage(buf[:n])

	}
}

func (c *WorkerConnection) SendToWorker(content any) {
	if parsedContent, ok := content.(string); ok {
		c.conn.Write([]byte(parsedContent))
	} else if parsedContent, ok := content.(Node); ok {

		// TODO: how to set the fields
		msg := Message{
			From:  0,
			To:    c.ID,
			Value: parsedContent,
			Type:  ASSIGN_VERTEX,
		}
		jsonBytes, err := json.Marshal(msg)
		if err != nil {
			fmt.Println(err)
		}
		c.conn.Write(jsonBytes)
	}
}

func (c *WorkerConnection) ProcessWorkerMessage(msg []byte) {
	message := new(Message)
	json.Unmarshal(msg, message)
	switch message.Type {

	}
}
