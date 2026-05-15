// https://fly.io/dist-sys/3b/
package main

import (
    "log"
    "sync"

    maelstrom "github.com/jepsen-io/maelstrom/demo/go"
    utils "github.com/p4tho/flyio-challenges/utils"
)

type Monitor struct {
	node *maelstrom.Node
	mu sync.Mutex
	data []int
}

func (mon *Monitor) push(val int) {
	mon.mu.Lock()
	defer mon.mu.Unlock()
	
	mon.data = append(mon.data, val)
}

func (mon *Monitor) get() []int {
	mon.mu.Lock()
	defer mon.mu.Unlock()
	
	data := make([]int, len(mon.data))
	copy(data, mon.data)	

	return data
}

func main() {
	n := maelstrom.NewNode()
	server := &Monitor{
		node: n,
		data: []int{},
	}
	
	utils.AddHandler(n, "broadcast", server.broadcast_handler)
	utils.AddHandler(n, "read", server.read_handler)
	utils.AddHandler(n, "topology", server.topology_handler)
	utils.AddAsyncHandler(n, "gossip", server.gossip_handler)
	
	if err := n.Run(); err != nil {
	    log.Fatal(err)
	}
}

func (serv *Monitor) broadcast_handler(req Broadcast) (BroadcastOk, error) {
	for _, id := range serv.node.NodeIDs() {
		gossip_req := map[string]any{}
		gossip_req["type"] = "gossip"
		gossip_req["message"] = req.Message

		err := serv.node.Send(id, gossip_req)
		if err != nil {
			return BroadcastOk{}, err
		}
	}
	
	serv.push(req.Message)
	
	res := BroadcastOk{}
	return res, nil
}

func (serv *Monitor) read_handler(req Read) (ReadOk, error) {
	messages := serv.get()
	
	res := ReadOk{
		Messages: messages,
	}
	return res, nil
}

func (serv *Monitor) topology_handler(req Topology) (TopologyOk, error) {
	res := TopologyOk{}
	return res, nil
}

func (serv *Monitor) gossip_handler(req Gossip) error {
	serv.push(req.Message)
	
	return nil
}
