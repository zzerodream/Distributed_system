package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
)

// Is this workerConnection type necessary?
type WorkerConnection struct {
	ID     int // incremental
	Addr   string
	C      chan any
	conn   net.Conn
	master *Master
}

func (c *WorkerConnection) CloseConn() {
	c.conn.Close()
}

// 2 goroutines per worker connection
func (c *WorkerConnection) Run() {
	//receive message from worker
	go c.RecvWorkers()
	// maybe we don't need this for for loop here
	for {
		for serverMessage := range c.C {
			c.SendToWorker(serverMessage)
		} 
	}
}

func (c *WorkerConnection) RecvWorkers() {
	reader := bufio.NewReader(c.conn)
	for {
		//assume that every message from master is sepetated by '\n'
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading message: %v\n", err)
			break
		}
		//decode the json and create the message object
		var message Message
		err = json.Unmarshal([]byte(line), &message)
		if err != nil {
			fmt.Printf("Error unmarshalling message: %v\n", err)
			continue
		}
		fmt.Printf("Received messages from worker %d with type %d.\n", message.From, message.Type)
		if message.Type == 10 {
			c.master.workerHeartBeat <- message
		} else {
			c.master.inCh <- message
		}
		
	}
}

func (c *WorkerConnection) SendToWorker(content any) {
	if msg, ok := content.(Message); ok {
		jsonBytes, err := json.Marshal(msg)
		if err != nil {
			fmt.Println(err)
		}
		jsonBytes = append(jsonBytes,'\n')
		c.conn.Write(jsonBytes)
		fmt.Printf("Sent message from Master to worker %d with type %d\n",msg.To, msg.Type)
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
		jsonBytes = append(jsonBytes,'\n')
		c.conn.Write(jsonBytes)
		fmt.Printf("Sent message from Master to worker %d with type %d\n",msg.To, msg.Type)
	}
}

//func (c *WorkerConnection) ProcessWorkerMessage(msg []byte) {
//
//	// 统一发给main thread处理
//	message := new(Message)
//	json.Unmarshal(msg, message)
//	c.master.inCh <- *message
//}
