package goeu

import (
	"github.com/heartbytenet/go-lerpc/pkg/lerpc"
	"log"
	"os"
)

type Goeu struct {
	Client   *lerpc.Client
}

func (g *Goeu) Init(endpoint string, token string) *Goeu {
	g.Client = (&lerpc.Client{}).Init(endpoint, token)
	return g
}

func (g *Goeu) InitEnv() *Goeu {
	endpoint := os.Getenv("NOEU_ENDPOINT")
	token    := os.Getenv("NOEU_TOKEN")

	if endpoint == "" {
		log.Fatalln("missing env var $NOEU_ENDPOINT")
	}

	if token == "" {
		log.Fatalln("missing env var $NOEU_TOKEN")
	}

	return g.Init(endpoint, token)
}

func (g *Goeu) Start(connections int) (err error) {
	// Todo: start leRPC websocket connections
	return
}

func (g *Goeu) Stop() {
	// Todo: stop leRPC client
}

// Execute currently uses HTTP-only as the websocket pool implementation is moving to go-leRPC
// https://github.com/heartbytenet/go-leRPC
// The working implementation in Goeu is available at commit 48ca169958a45c86904fb2ed9d57211a0b840852
func (g *Goeu) Execute(cmd *lerpc.ExecuteCommand, res *lerpc.ExecuteResult) (err error) {
	return g.Client.Execute(cmd, res)
}