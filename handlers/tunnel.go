package handlers

import (
	"fmt"
	"io"

	"github.com/danielmadu/goexpose/config"
	"github.com/gin-gonic/gin"
)

func Tunnel(c *gin.Context) {

	defer c.Request.Body.Close()

	reqBody := getBody(c)

	headers := make(map[string]string)

	for k := range c.Request.Header {
		headers[k] = c.Request.Header.Get(k)
	}

	body, err := io.ReadAll(reqBody)
	if err != nil {
		fmt.Println(err)
		return
	}

	message := config.Message{
		Path:    c.Request.URL.String(),
		Headers: headers,
		Body:    string(body),
		Method:  c.Request.Method,
	}

	messageResponse := config.Message{}

	channel := config.GetChannel()

	channel <- message

	messageResponse = <-channel

	fmt.Println(messageResponse.Body)

	for k, v := range messageResponse.Headers {
		c.Header(k, v)
	}

	c.Data(messageResponse.Status, c.GetHeader("Contet-Type"), []byte(messageResponse.Body))
}

func getBody(c *gin.Context) io.ReadCloser {
	if c.Request.Method != "GET" || c.Request.Method != "DELETE" || c.Request.Method != "" {
		return c.Request.Body
	}

	return nil
}
