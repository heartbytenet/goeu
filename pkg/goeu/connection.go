package goeu

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"math/rand"
	"sync"
	"time"
)

type Connection struct {
	ID    string
	URL   string
	Token string
	Conn  *websocket.Conn

	Promises map[string]*Promise

	sync.Mutex
}

func (c *Connection) Init(ID string, url string, token string) *Connection {
	c.ID    = ID
	c.URL   = url
	c.Token = token
	c.Conn  = nil

	c.Promises = map[string]*Promise{}

	return c
}

func (c *Connection) Start() {
	_ = c.Dial()
}

func (c *Connection) Stop() {

}

func (c *Connection) Dial() (err error) {
	c.Lock()
	defer c.Unlock()

	c.Conn, _, err = websocket.DefaultDialer.Dial(c.URL, nil)
	return
}

func (c *Connection) Exec(cmd *ApiExecuteCommand, res *ApiExecuteResult, callback chan byte) (err error) {
	var data []byte
	var msgt int
	var _res ApiExecuteResult

	now := func() int64 { return time.Now().UnixMilli() }
	gid := func() (res string) { res = ""; for i := 0; i < 16; i++ { res += string(rune(int('a') + rand.Intn(26))) }; return }

	c.Lock()
	defer c.Unlock()

	if c.Conn == nil {
		return errors.New("connection is nil")
	}

	cmd.Token = c.Token
	cmd.ID = gid()
	for {
		if _, in := c.Promises[cmd.ID]; !in {
			break
		}
		cmd.ID = gid()
	}

	c.Promises[cmd.ID] = (&Promise{}).Init(res, now(), callback)

	data, err = json.Marshal(cmd)
	if err != nil {
		return
	}

	err = c.Conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		return
	}

	msgt, data, err = c.Conn.ReadMessage()
	if err != nil {
		return
	}

	if msgt != 1 {
		return
	}

	err = json.Unmarshal(data, &_res)
	if err != nil {
		return
	}

	promise, in := c.Promises[_res.ID]
	if !in {
		err = errors.New("promise not found")
		return
	}

	promise.DoneSet(true)
	*promise.Res = _res
	promise.Callback <- 42

	ks := make([]string, 0)
	for k, v := range c.Promises {
		if v.Ended(now(), 30 * 1000) {
			ks = append(ks, k)
		}
	}
	for _, v := range ks {
		delete(c.Promises, v)
	}

	return
}
