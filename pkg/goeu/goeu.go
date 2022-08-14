package goeu

import (
	"errors"
	"log"
	"os"
	"time"
)

type Goeu struct {
	Pool     *Pool
	Token    string
	Endpoint string
}

func (g *Goeu) Init(endpoint string, token string) *Goeu {
	g.Pool     = (&Pool{}).Init(g)
	g.Token    = token
	g.Endpoint = endpoint
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
	err = g.Pool.Start(connections)
	if err != nil {
		return
	}
	return
}

func (g *Goeu) Stop() {
	g.Pool.Stop()
}

func (g *Goeu) Execute(cmd *ApiExecuteCommand, res *ApiExecuteResult) (callback chan byte, err error) {
	callback = make(chan byte, 1)

	if cmd == nil || res == nil {
		err = errors.New("cmd or res is nil")
		return
	}

	go func() {
		for attempt := 0; attempt < 10; attempt++ {
			if g.Pool.Execute(cmd, res, callback) {
				return
			}
			time.Sleep(time.Millisecond * 500)
		}
		err = errors.New("command took too many attempts")
		*res = ApiExecuteResult{
			Success: false,
			Error:   err.Error(),
		}
		callback <- 42
		return
	}()

	return
}