package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/danielmadu/goexpose/config"
	"github.com/danielmadu/goexpose/router"
	"github.com/urfave/cli/v2"
	"golang.org/x/net/websocket"
)

var (
	token       string
	serverUrl   string
	localConfig *config.LocalConfig
)

func main() {

	config.Init()

	conf := config.GetConfig()
	localConfig = config.GetLocalConfig()

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
					&cli.StringFlag{
						Name:        "certFile",
						Value:       "",
						Aliases:     []string{"c"},
						Usage:       "Certificate file",
						Destination: &conf.CertFile,
					},
					&cli.StringFlag{
						Name:        "keyFile",
						Value:       "",
						Aliases:     []string{"k"},
						Usage:       "Key file",
						Destination: &conf.KeyFile,
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
					&cli.StringFlag{
						Name:        "basicAuth",
						Value:       "",
						Usage:       "Pass the user:password to basic authentication  (optional)",
						Destination: &localConfig.BasicAuth,
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
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

	if splited := strings.Split(localConfig.BasicAuth, ":"); len(splited) == 2 {
		localConfig.EnabledBasicAuth = true
	}

	sharedHost := cli.Args().Get(0)

	localConfig.SharedHostname = sharedHost
	messageResponse := config.Message{}

	fmt.Println("Press CTRL+C to exit")

	val, err := url.Parse(serverUrl)
	if err != nil {
		fmt.Println(err)
		return err
	}

	wsprotocol := "ws"
	if val.Scheme == "https" {
		wsprotocol = "wss"
	}

	if val.Scheme == "" || (val.Scheme != "http" && val.Scheme != "https") {
		return fmt.Errorf("Invalid server url, you must define the protocol (http, https)")
	}

	serverUrl = val.Host

	origin := "http://localhost/"
	url := fmt.Sprintf("%s://%s/goexpose/ws?token=%s", wsprotocol, serverUrl, token)
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

			if localConfig.EnabledBasicAuth {

				user, passwd, _ := req.BasicAuth()

				basicAuth := strings.Split(localConfig.BasicAuth, ":")

				if user == basicAuth[0] && passwd == basicAuth[1] {
					execute(client, req, ws)
					continue
				}

				messageResponse = config.Message{
					Body:    "Unauthorized",
					Status:  401,
					Headers: make(map[string]string),
				}

				messageResponse.Headers["WWW-Authenticate"] = `Basic realm="restricted", charset="UTF-8"`
				sendResponse(messageResponse, ws)
				continue
			}

			execute(client, req, ws)

		}

	}
}

func sendResponse(messageResponse config.Message, ws *websocket.Conn) {
	encoded, _ := messageResponse.Encode()

	if _, err := ws.Write(encoded); err != nil {
		fmt.Println(err)
	}
}

func execute(client *http.Client, req *http.Request, ws *websocket.Conn) {
	response, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
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

	sendResponse(messageResponse, ws)
}
