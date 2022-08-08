# Goeu

Goeu is a [Go](https://go.dev/) client for the [Noeu](https://gitlab.com/heartbytenet/noeu) database

Features
-------

* Websocket by default (Http is on the way)
* Multiple connections for reliability
* Automatic connection regeneration

Installation
-------

Install Goeu using the "go get" command:
```shell
go get -u gitlab.com/heartbytenet/goeu
```

Example
-------

```go
package main

import (
	"gitlab.com/heartbytenet/goeu/pkg/goeu"
	"log"
)

func main() {
	g := (&goeu.Goeu{}).InitEnv()
	
	err := g.Start(4)
	if err != nil {
		log.Fatalln(err)
	}
	
	var cmd goeu.ApiExecuteCommand
	var res goeu.ApiExecuteResult
	
	cmd = goeu.ApiExecuteCommand{
		Namespace: "misc",
		Method:    "ping",
		Params: map[string]interface{}{},
    }
	callback, err := g.Execute(&cmd, &res)
	if err != nil {
		log.Fatalln(err)
    }
	
	<- callback
	log.Println(res)
}
```