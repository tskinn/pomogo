package rpc

import (
	"encoding/json"
	"io"
	"math/rand"
)

type clientRequest struct {
	Version string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      uint64      `json:"id"`
}

type clientResponse struct {
	Version string           `json:"jsonrpc"`
	Result  *json.RawMessage `json:"result"`
	Error   *json.RawMessage `json:"error"`
}

func EncodeClientRequest(method string, args interface{}) ([]byte, error) {
	c := &clientRequest{
		Version: "2.0",
		Method:  method,
		Params:  args,
		ID:      uint64(rand.Int63()),
	}
	return json.Marshal(c)
}

func DecodeClientResponse(r io.Reader, reply interface{}) error {
	var response clientResponse
	if err := json.NewDecoder(r).Decode(&response); err != nil {
		return err
	}

	if response.Error != nil {
		jsonErr := &Error{}
		if err := json.Unmarshal(*response.Error, jsonErr); err != nil {
			return &Error{
				Code:    ErrorCodeServer,
				Message: string(*response.Error),
			}
		}
		return jsonErr
	}

	if response.Result == nil {
		return ErrorNullResult
	}

	return json.Unmarshal(*response.Result, reply)
}
