package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/danielmadu/goexpose/config"
	"github.com/danielmadu/goexpose/router"
	"github.com/urfave/cli/v2"
	"golang.org/x/net/websocket"
)

var (
	token     string
	serverUrl string
)

func main() {

	config.Init()

	app := &cli.App{
		Name:  "goexpose",
		Usage: "Create public URLs for local applications",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "token",
				Value:       "",
				Aliases:     []string{"t"},
				Usage:       "Define a authentication token (optional)",
				Destination: &token,
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "serve",
				Usage:  "Start the GoExpose server",
				Action: startServer,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "token",
						Value:       "",
						Aliases:     []string{"t"},
						Usage:       "Define a authentication token (optional)",
						Destination: &token,
					},
				},
			},
			{
				Name:   "share",
				Usage:  "Share a local url with a remote GoExpose server",
				Action: startShare,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "server",
						Value:       "",
						Usage:       "Define the server to connect (required)",
						Destination: &serverUrl,
						Required:    true,
					},
					&cli.StringFlag{
						Name:        "token",
						Value:       "",
						Aliases:     []string{"t"},
						Usage:       "Pass the authentication token (optional)",
						Destination: &token,
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {

	}

}

func startServer(cli *cli.Context) error {
	config := config.GetConfig()

	config.Token = token

	fmt.Println("GoExpose listening on port :3000")

	return router.Init()
}

func startShare(cli *cli.Context) error {
	if cli.NArg() == 0 {
		fmt.Println("You must pass the local URL that you want to share")
		return nil
	}

	if cli.NArg() > 1 {
		fmt.Println("Too many arguments")
		return nil
	}

	sharedHost := cli.Args().Get(0)

	localConfig := config.GetLocalConfig()

	localConfig.SharedHostname = sharedHost

	fmt.Println("Press CTRL+C to exit")

	origin := "http://localhost/"
	url := fmt.Sprintf("ws://%s/goexpose/ws?token=%s", serverUrl, token)
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		fmt.Println(err)
	}

	message := config.GetLocalMessage()

	var msg = make([]byte, 1048576)
	client := &http.Client{}
	for {

		var n int
		if n, err = ws.Read(msg); err != nil {
			continue
		}

		resp := fmt.Sprintf("%s", msg[:n])
		if resp == "ping" {
			continue
		}

		if err := message.Decode(msg[:n]); err != nil {
			continue
		}
		if message.Path != "" {
			req, _ := http.NewRequest(message.Method, localConfig.SharedHostname+message.Path, bytes.NewBufferString(message.Body))

			for k, v := range message.Headers {
				req.Header.Set(k, v)
			}

			response, err := client.Do(req)
			if err != nil {
				fmt.Println(err)
				continue
			}

			respBody, _ := io.ReadAll(response.Body)

			headers := make(map[string]string)

			for k := range response.Header {
				headers[k] = response.Header.Get(k)
			}

			messageResponse := config.Message{
				Body:    string(respBody),
				Headers: headers,
				Status:  response.StatusCode,
			}

			encoded, _ := messageResponse.Encode()

			if _, err := ws.Write(encoded); err != nil {
				fmt.Println(err)
			}
		}

	}
}
