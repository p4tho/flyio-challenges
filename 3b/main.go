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
	neighbors []string
	mu sync.Mutex
	data []int // simulates a set
}

func (mon *Monitor) push(val int) bool {
	mon.mu.Lock()
	defer mon.mu.Unlock()
	
	// Only add val if not in data
	for _, ele := range mon.data {
		if ele == val {
			return false
		}
	}
	
	mon.data = append(mon.data, val)
	return true
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
		neighbors: []string{},
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
	if serv.push(req.Message) {
		for _, id := range serv.neighbors {
			req := Gossip {
				Message: req.Message,
			}
			
			err := utils.SendAsync(serv.node, "gossip", id, req)
			if err != nil {
				return BroadcastOk{}, err
			}
		}	
	}
	
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
	if neighbors, ok := req.Topology[serv.node.ID()]; ok {
        serv.neighbors = neighbors
    }
	
	res := TopologyOk{}
	return res, nil
}

func (serv *Monitor) gossip_handler(req Gossip) error {
	if serv.push(req.Message) {
		for _, id := range serv.neighbors {
			req := Gossip {
				Message: req.Message,
			}
			
			err := utils.SendAsync(serv.node, "gossip", id, req)
			if err != nil {
				return err
			}
		}	
	}
	
	return nil
}
