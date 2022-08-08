package goeu

import (
	"fmt"
	"math/rand"
	"time"
)

type Pool struct {
	Goeu *Goeu

	Connections map[string]*Connection
}

func (p *Pool) Init(goeu *Goeu) *Pool {
	p.Goeu = goeu
	return p
}

func (p *Pool) Start(connections int) (err error) {
	var connection *Connection

	p.Connections = make(map[string]*Connection, connections)
	for i := 0; i < connections; i++ {
		time.Sleep(time.Second)
		connection = (&Connection{}).Init(
			fmt.Sprintf("%d", i), fmt.Sprintf("wss://%s/connect", p.Goeu.Endpoint), p.Goeu.Token)
		p.Connections[connection.ID] = connection
		p.Connections[connection.ID].Start()
	}
	return
}

func (p *Pool) Stop() {
	for _, connection := range p.Connections {
		connection.Stop()
	}
}

func (p *Pool) Execute(cmd *ApiExecuteCommand, res *ApiExecuteResult, callback chan byte) bool {
	var err error

	cons := len(p.Connections)
	coni := 0
	goal := rand.Intn(cons)

	for _, con := range p.Connections {
		if coni >= goal {
			err = con.Exec(cmd, res, callback)
			if err != nil {
				err = con.Dial()
				if err != nil {
					return false
				}
			} else {
				return true
			}
		}
		coni++
	}

	return false
}