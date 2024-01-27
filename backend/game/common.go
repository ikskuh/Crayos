package game

import "encoding/json"

type Message interface {
	GetJsonType() string
	FixNils() Message
}

type Graphics interface{}

func SerializeMessage(msg Message) ([]byte, error) {

	msg = msg.FixNils()

	temp, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	var dummy map[string]interface{}

	err = json.Unmarshal(temp, &dummy)
	if err != nil {
		return nil, err
	}

	dummy["type"] = msg.GetJsonType()

	return json.Marshal(dummy)
}
