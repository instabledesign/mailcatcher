package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

func Event(broker *Broker) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		conn, err := NewSSEConnHTTPHandler(writer, request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		broker.newClients <- conn.writeChan
		defer func() { broker.closingClients <- conn.writeChan }()
		conn.Listen()
		fmt.Println("client over")
	}
}

type SSEWriter interface {
	http.ResponseWriter
	http.Flusher
}

type SSEClient struct {
	writer  http.ResponseWriter
	flusher http.Flusher
}

func (c *SSEClient) Header() http.Header {
	return c.writer.Header()
}
func (c *SSEClient) WriteHeader(statusCode int) {
	c.writer.WriteHeader(statusCode)
}

func (c *SSEClient) Event(eventType string, data []byte) (int, error) {
	if eventType == "" {
		eventType = "message"
	}
	data = append(data, "\n\n"...)
	data = append([]byte("data: "), data...)
	data = append([]byte("event: "+eventType+"\n"), data...)
	return c.Write(data)
}

func (c *SSEClient) Ping() (int, error) {
	return c.Write([]byte("event: ping\ndata:{\"time\":" + time.Now().Format(time.RFC3339) + "}\n\n"))
}

func (c *SSEClient) Message(data []byte) (int, error) {
	return c.Event("message", data)
}

func (c *SSEClient) Write(data []byte) (int, error) {
	return c.writer.Write(data)
}
func (c *SSEClient) Flush() {
	c.flusher.Flush()
}

func NewSSEClient(writer http.ResponseWriter, flusher http.Flusher) *SSEClient {
	return &SSEClient{
		writer:  writer,
		flusher: flusher,
	}
}

func NewSSEClientFromResponseWriter(writer http.ResponseWriter) (*SSEClient, error) {
	flusher, ok := writer.(http.Flusher)

	if !ok {
		return nil, errors.New("streaming unsupported")
	}
	return NewSSEClient(writer, flusher), nil
}

type SSEConn struct {
	client    *SSEClient
	request   *http.Request
	writeChan chan []byte
}

func (c *SSEConn) Listen() {
	c.client.Header().Set("Content-Type", "text/event-stream")
	c.client.Header().Set("Cache-Control", "no-cache")
	c.client.Header().Set("Connection", "keep-alive")
	c.client.Header().Set("Access-Control-Allow-Origin", "*")
	c.client.Ping()
	c.client.Flush()

	end := c.request.Context().Done()
	for {
		select {
		case <-end:
			c.client.Event("close", []byte("server stopped"))
			c.client.Flush()
			return
		case data := <-c.writeChan:
			c.client.Event("mail", data)
			c.client.Flush()
		case <-time.Tick(5 * time.Second):
			c.client.Ping()
			c.client.Flush()
		}
	}
}

func NewSSEConn(client *SSEClient, request *http.Request) (*SSEConn, error) {
	if client == nil {
		return nil, errors.New("SSEClient is mandatory")
	}
	if request == nil {
		return nil, errors.New("http.Request is mandatory")
	}
	return &SSEConn{client: client, request: request, writeChan: make(chan []byte, 5)}, nil
}

func NewSSEConnHTTPHandler(writer http.ResponseWriter, request *http.Request) (*SSEConn, error) {
	c, err := NewSSEClientFromResponseWriter(writer)
	if err != nil {
		return nil, err
	}
	return NewSSEConn(c, request)
}
