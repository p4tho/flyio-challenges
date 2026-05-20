package utils

import (
	"context"
	"encoding/json"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func AddHandler[Req any, Res any](n *maelstrom.Node, msg_type string, handler func(Req) (Res, error)) {
	n.Handle(msg_type, func(msg maelstrom.Message) error {
		var msg_body Req
		if err := json.Unmarshal(msg.Body, &msg_body); err != nil {
			return err
		}

		// Handler extracts and manipulates message body
		res, err := handler(msg_body)
		if err != nil {
			return err
		}

		// Convert handler result to a map with json
		res_json_raw, err := json.Marshal(res)
		if err != nil {
			return err
		}
		
		var res_json map[string]any
		err = json.Unmarshal(res_json_raw, &res_json)
		if err != nil {
			return err
		}
		
		// Replies always require <msg_type>_ok
		res_json["type"] = msg_type + "_ok"
		return n.Reply(msg, res_json)
	})
}

func AddAsyncHandler[Req any](n *maelstrom.Node, msg_type string, handler func(Req) error) {
	n.Handle(msg_type, func(msg maelstrom.Message) error {
		var req Req
		if err := json.Unmarshal(msg.Body, &req); err != nil {
			return err
		}

		if err := handler(req); err != nil {
			return err
		}
		return nil
	})
}

func SendAsync[Req any](n *maelstrom.Node, msg_type string, dest string, req Req) error {
	req_json_raw, err := json.Marshal(req)
	if err != nil {
		return err
	}
	
	var req_json map[string]any
	err = json.Unmarshal(req_json_raw, &req_json)
	if err != nil {
		return err
	}
	
	req_json["type"] = msg_type

	if err := n.Send(dest, req_json); err != nil {
		return err
	}
	return nil
}
