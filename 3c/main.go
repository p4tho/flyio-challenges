// https://fly.io/dist-sys/3c/
package main

import (
    "log"
    "sync"
    "time"

    maelstrom "github.com/jepsen-io/maelstrom/demo/go"
    utils "github.com/p4tho/flyio-challenges/utils"
)

type Monitor struct {
	node *maelstrom.Node
	neighbors []string
	mu sync.Mutex
	data []int // simulates a set
	delivered []GossipMsg
}

type GossipMsg struct {
	dest_node string
	message int
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

// Remove first instance of node name in delivered array
func (mon *Monitor) remove(target_msg GossipMsg) bool {
	mon.mu.Lock()
	defer mon.mu.Unlock()
	
	for idx, id := range mon.delivered {
		if id == target_msg {
			mon.delivered = append(mon.delivered[:idx], mon.delivered[idx+1:]...)
			return true
		}
	}
	
	return false
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
		delivered: []GossipMsg{},
	}
	
	go func() {
		for range time.Tick(time.Second) {
			for _, msg := range server.delivered {
				err := utils.SendAsync(server.node, "gossip", msg.dest_node, Gossip{ Message: msg.message })
				if err != nil {
					log.Printf("error async deliver to %s: %v", msg.dest_node, err)
				}
			}
		}
	}()
	
	utils.AddHandler(n, "broadcast", server.broadcast_handler)
	utils.AddHandler(n, "read", server.read_handler)
	utils.AddHandler(n, "topology", server.topology_handler)
	utils.AddHandler(n, "gossip", server.gossip_handler)
	utils.AddAsyncHandler(n, "gossip_ok", server.gossip_ok_handler)
	
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
			
			// Add to delivered array to check for confirmations
			msg := GossipMsg {
				dest_node: id,
				message: req.Message,
			}
			serv.delivered = append(serv.delivered, msg)
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

func (serv *Monitor) gossip_handler(req Gossip) (GossipOk, error) {
	if serv.push(req.Message) {
		for _, id := range serv.neighbors {
			req := Gossip {
				Message: req.Message,
			}
			
			err := utils.SendAsync(serv.node, "gossip", id, req)
			if err != nil {
				return GossipOk{}, err
			}
			
			msg := GossipMsg {
				dest_node: id,
				message: req.Message,
			}
			serv.delivered = append(serv.delivered, msg)
		}	
	}
	
	// Reply with confirmation to remove from delivered array of source node
	res := GossipOk {
		Src: serv.node.ID(),
		Message: req.Message,
	}
	
	return res, nil
}

func (serv *Monitor) gossip_ok_handler(req GossipOk) error {
	msg := GossipMsg {
		dest_node: req.Src,
		message: req.Message,
	}
	
	serv.remove(msg)
	
	return nil
}