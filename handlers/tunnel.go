package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/danielmadu/goexpose/config"
	"github.com/gin-gonic/gin"
)

func Tunnel(c *gin.Context) {

	client := &http.Client{}

	config := config.GetConfig()

	formatted := fmt.Sprintf("%s%s", config.Shared, c.Request.URL)

	fmt.Println(formatted)

	defer c.Request.Body.Close()

	reqBody := getBody(c)

	req, _ := http.NewRequest(c.Request.Method, formatted, reqBody)

	for k := range c.Request.Header {
		c.Header(k, c.Request.Header.Get(k))
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Println(err)
		c.AbortWithError(500, err)
		return
	}

	log.Println(c.Request.Method)

	for k := range resp.Header {
		c.Header(k, resp.Header.Get(k))
	}

	c.Status(resp.StatusCode)

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	c.Data(resp.StatusCode, resp.Header.Get("Contet-Type"), body)
}

func getBody(c *gin.Context) io.ReadCloser {
	if c.Request.Method != "GET" || c.Request.Method != "DELETE" || c.Request.Method != "" {
		return c.Request.Body
	}

	return nil
}
