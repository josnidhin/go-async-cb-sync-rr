/**
 * @author Jose Nidhin
 */
package main

import (
	"encoding/json"
)

type Message struct {
	RequestId string `json:"requestId"`
	Data      []byte `json:"data"`
}

func (m Message) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

func NewMessageFromJSON(b []byte) (Message, error) {
	m := Message{}
	err := json.Unmarshal(b, &m)
	return m, err
}
