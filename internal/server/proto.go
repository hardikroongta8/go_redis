package server

import (
	"fmt"
)

type Command interface{}

type PutCommand struct {
	key, val string
}
type GetCommand struct {
	key string
}

func parseMessage(rawMsg []byte) (Command, error) {

	cmd := deserializeString(string(rawMsg))
	if len(cmd) == 3 && cmd[0] == "PUT" {
		return PutCommand{
			key: cmd[1],
			val: cmd[2],
		}, nil
	}
	if len(cmd) == 2 && cmd[0] == "GET" {
		return GetCommand{key: cmd[1]}, nil
	}
	return nil, fmt.Errorf("INVALID MESSAGE: %s", string(rawMsg))
}

func deserializeString(msg string) []string {
	cmd := make([]string, 0)
	lf := 0
	for i := 0; i < len(msg)-1; i++ {
		if msg[i] == '\r' && msg[i+1] == '\n' {
			cmd = append(cmd, msg[lf:i])
			lf = i + 2
			i++
		}
	}
	return cmd
}
