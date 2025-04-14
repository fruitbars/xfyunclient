package ase

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type WSCB func(p []byte)

type ASEWebsocketClient struct {
	ASEClientBase
	conn *websocket.Conn
}

// NewASEWebsocketClient initializes a new ASEWebsocketClient instance
func NewASEWebsocketClient(serverURL, appid, apikey, apisecret, httpProto, aseAlgorithm string) *ASEWebsocketClient {
	baseClient := NewASEClientBase(serverURL, appid, apikey, apisecret, httpProto, aseAlgorithm)
	return &ASEWebsocketClient{ASEClientBase: baseClient}
}

// Connect establishes a WebSocket connection
func (c *ASEWebsocketClient) Connect() error {
	svrURL, headers, err := c.getAuthServerURL("GET")
	if err != nil {
		return err
	}

	c.conn, _, err = websocket.DefaultDialer.Dial(svrURL, *headers)
	if err != nil {
		return err
	}

	return c.conn.SetReadDeadline(time.Now().Add(30 * time.Second))
}

// SendJson sends a JSON-encoded message
func (c *ASEWebsocketClient) SendJson(data interface{}) error {
	return c.conn.WriteJSON(data)
}

// SendBinary sends a binary message
func (c *ASEWebsocketClient) SendBinary(data []byte) error {
	return c.conn.WriteMessage(websocket.TextMessage, data)
}

// Receive reads a message and returns the raw data
func (c *ASEWebsocketClient) Receive() ([]byte, error) {
	_, message, err := c.conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	return message, nil
}

// Close closes the WebSocket connection
func (c *ASEWebsocketClient) Close() {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			log.Println("Error closing connection:", err)
		}
	}
}

// ReadMessage reads a single message from the WebSocket connection
func (c *ASEWebsocketClient) ReadMessage() ([]byte, error) {
	_, message, err := c.conn.ReadMessage()
	if err != nil {
		log.Println("Error reading message:", err)
		return nil, err
	}
	return message, nil
}

// ReadAllMessage reads all messages until a terminating condition is met
func (c *ASEWebsocketClient) ReadAllMessage() ([]ASEBaseResponse, error) {
	var msgSet []ASEBaseResponse
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		var message ASEBaseResponse
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Println("Error unmarshaling JSON:", err)
			continue
		}

		msgSet = append(msgSet, message)

		if message.Code != 0 || message.Data.Status == 2 {
			log.Println("Closing connection")
			c.Close()
			break
		} else {
			log.Println("Connection status normal")
		}
	}
	return msgSet, nil
}

func (c *ASEWebsocketClient) ReadAllMessageCallBack(cb WSCB) (string, error) {
	var sid string

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			// Handle normal WebSocket close gracefully
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseNoStatusReceived) {
				log.Println("WebSocket closed normally:", err)
				return sid, nil
			}
			log.Println("Error reading WebSocket message:", err)
			return sid, err
		}

		if cb != nil {
			cb(msg)
		}

		var message ASEBaseResponse
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Println("Failed to unmarshal JSON:", err)
			continue
		}

		if sid == "" {
			sid = message.Sid
		}

		if message.Code != 0 {
			c.Close()
			return sid, fmt.Errorf("code: %d, message: %s", message.Code, message.Message)
		}

		// status: 2 means end of stream
		if message.Data.Status == 2 {
			//			log.Println("Stream completed. Closing connection.")
			c.Close()
			break
		}
	}

	return sid, nil
}

// CallASEAPI handles the complete WebSocket API call
func (c *ASEWebsocketClient) CallASEAPI(reqJsonByte []byte) ([]ASEBaseResponse, error) {
	if err := c.Connect(); err != nil {
		return nil, err
	}

	if err := c.SendBinary(reqJsonByte); err != nil {
		return nil, err
	}

	return c.ReadAllMessage()
}

func (c *ASEWebsocketClient) CallASEAPIJson(req interface{}) ([]ASEBaseResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	return c.CallASEAPI(jsonData)
}

func (c *ASEWebsocketClient) CallASEAPICallBack(req interface{}, cb WSCB) (string, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	log.Println(string(jsonData))
	if err := c.Connect(); err != nil {
		return "", err
	}

	if err := c.SendBinary(jsonData); err != nil {
		return "", err
	}

	return c.ReadAllMessageCallBack(cb)
}

func (c *ASEWebsocketClient) CallASEAPIJsonCallBack(req interface{}) ([]ASEBaseResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	return c.CallASEAPI(jsonData)
}
