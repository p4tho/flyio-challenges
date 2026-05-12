// https://fly.io/dist-sys/1/
package main

import (
    "log"

    maelstrom "github.com/jepsen-io/maelstrom/demo/go"
    utils "github.com/p4tho/flyio-challenges/utils"
)

func main() {
	n := maelstrom.NewNode()
	
	utils.AddHandler(n, "echo", echo_handler)
	
	if err := n.Run(); err != nil {
	    log.Fatal(err)
	}
}

func echo_handler(req Echo) (EchoOk, error) {
	res := EchoOk{
		Echo: req.Echo,
	}
	return res, nil
}