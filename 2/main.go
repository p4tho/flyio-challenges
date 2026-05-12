// https://fly.io/dist-sys/2/
package main

import (
	"crypto/rand"
	"encoding/binary"
    "log"

    maelstrom "github.com/jepsen-io/maelstrom/demo/go"
    utils "github.com/p4tho/flyio-challenges/utils"
)

func generate_uid() (int, error) {
	var b [4]byte
	if _, err := rand.Read(b[:]); err != nil {
		return 0, err
	}

	id := int(binary.LittleEndian.Uint32(b[:]) & 0x7fffffff)

	return id, nil
}

func main() {
	n := maelstrom.NewNode()
	
	utils.AddHandler(n, "generate", generate_handler)
	
	if err := n.Run(); err != nil {
	    log.Fatal(err)
	}
}

func generate_handler(req Generate) (GenerateOk, error) {
	uid, err := generate_uid()
	if err != nil {
		return GenerateOk{}, err
	}
	
	res := GenerateOk{
		Id: uid,
	}
	return res, nil
}