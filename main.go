package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/danielmadu/goexpose/config"
	"github.com/danielmadu/goexpose/router"
	"github.com/urfave/cli/v2"
)

var (
	token    string
	hostname string
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
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}

func startServer(cli *cli.Context) error {
	config := config.GetConfig()

	config.Token = token

	fmt.Println(config.Token)

	return router.Init()
}

func startShare(cli *cli.Context) error {
	fmt.Println("Press CTRL+C to exit")
	for {
		if err := ping(); err != nil {
			return err
		}
		time.Sleep(15 * time.Second)
	}
}

type Ping struct {
	Port string `json:"port"`
}

func ping() error {
	client := &http.Client{}
	ping := config.Ping{
		Port: "8080",
	}
	body, _ := json.Marshal(ping)
	req, _ := http.NewRequest(http.MethodPost, "http://localhost:3000/goexpose/api/ping", bytes.NewReader(body))
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("fail to try refresh server infos")
	}

	return nil
}
