package goeu

import (
	"github.com/heartbytenet/go-lerpc/pkg/lerpc"
	"github.com/heartbytenet/go-lerpc/pkg/proto"
	"log"
	"os"
)

type Goeu struct {
	Client *lerpc.Client
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
	err = g.Client.Start(connections)
	if err != nil {
		return
	}
	return
}

func (g *Goeu) Stop() {
	// Todo: stop leRPC client
}

func (g *Goeu) Execute(cmd *proto.ExecuteCommand, res *proto.ExecuteResult) (err error) {
	return g.Client.Execute(cmd, res)
}

func (g *Goeu) Exec(cmd *proto.ExecuteCommand) (*proto.ExecuteResult, error) {
	var res proto.ExecuteResult
	var err error

	err = g.Execute(cmd, &res)

	return &res, err
}

func (g *Goeu) Eval(filename string, keys []string, args []string) (*proto.ExecuteResult, error) {
	var (
		res  proto.ExecuteResult
		data []byte
		err  error
	)

	data, err = os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = g.Execute(
		(&proto.ExecuteCommand{}).
			SetNamespace("sc").
			SetMethod("eval").
			SetParam("sc.dt", string(data)).
			SetParam("sc.ks", keys).
			SetParam("sc.ag", args),
		&res)

	if err != nil {
		return nil, err
	}

	return &res, nil
}