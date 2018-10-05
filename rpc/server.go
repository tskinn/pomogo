package rpc

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"sync"
	"unicode"
	"unicode/utf8"
)

var null = json.RawMessage([]byte("null"))
var version = "2.0"
var (
	// Precompute the reflect.Type of error and http.Request
	typeOfError = reflect.TypeOf((*error)(nil)).Elem()
	// typeOfRequest = reflect.TypeOf((*http.Request)(nil)).Elem()
)

type serverRequest struct {
	Version string           `json:"jsonrpc"`
	Method  string           `json:"method"`
	Params  *json.RawMessage `json:"Params"`
	ID      *json.RawMessage `json:"id"`
}

type serverResponse struct {
	Version string           `json:"jsonrpc"`
	Result  interface{}      `json:"result,omitempty"`
	Error   *Error           `json:"error,omitempty"`
	ID      *json.RawMessage `json:"id"`
}

// func NewCustomCodec(encSel rpc.EncoderSelector) *Codec {
// 	return &Codec{encSel: encSel}
// }

// func NewCode() *Code {
// 	return NewCustomCodec(rpc.DefaultEncoderSelector)
// }

// type Codec struct {
// 	encSel rpc.EncoderSelector
// }

// type CodecRequest struct {
// 	request *serverRequest
// 	err     error
// 	encoder rpc.Encoder
// }

// func (c *CodecRequest) Method() (string, error) {
// 	if c.err == nil {
// 		return c.request.Method, nil
// 	}
// 	return "", c.err
// }

// func (c *CodecRequest) ReadRequest(args interface{}) error {
// 	if c.err == nil && c.request.Params != nil {
// 		if err := json.Unmarshal(*c.request.Params, args); err != nil {
// 			params := [1]interface{}{args}
// 			if err = json.Unmarshal(*c.request.Params, &params); err != nil {
// 				c.err = &Error{
// 					Code:    ErrorCodeInvalidIDRequest,
// 					Message: err.Error(),
// 					Data:    c.request.Params,
// 				}
// 			}
// 		}
// 	}
// 	return c.err
// }

// func (c *CodecRequest) WriteServerResponse(w io.Writer, reply interface{}) {
// 	if c.request.ID != nil {
// 		encoder := json.NewEncoder(c.encoder.Encode(w))
// 		err := encoder.Encode(res)
// 		if err != nil {
// 			WriteError(w, err.Error())
// 		}
// 	}
// }

// func (c *CodecRequest) WriteError(w io.Writer, err error) {
// 	jsonErr, ok := err.(*Error)
// 	if !ok {
// 		jsonErr = &Error{
// 			Code:    ErrorCodeServer,
// 			Message: err.Error(),
// 		}
// 	}

// 	res := &serverResponse{
// 		Version: Version,
// 		Error:   jsonErr,
// 		ID:      c.request.ID,
// 	}

// 	c.WriteServerResponse(w, res)
// }

// func (c *CodecRequest) WriteResponse(w io.Writer, err error) {
// 	res := &serverResponse{
// 		Version: version,
// 		Result:  reply,
// 		ID:      c.request.ID,
// 	}

// 	c.WriteServerResponse(w, res)
// }

func WriteError(w io.Writer, msg string) {
	fmt.Fprint(w, msg)
}

type service struct {
	name     string                    // name of service
	rcvr     reflect.Value             // receiver of methods for the service
	rcvrType reflect.Type              // type of the receiver
	methods  map[string]*serviceMethod // registered methods
}

type serviceMethod struct {
	method    reflect.Method // receiver method
	argsType  reflect.Type   // type of the request argument
	replyType reflect.Type   // type of the response argument
}

type serviceMap struct {
	mutex    sync.Mutex
	services map[string]*service
}

func (m *serviceMap) register(rcvr interface{}, name string) error {
	s := &service{
		name:     name,
		rcvr:     reflect.ValueOf(rcvr),
		rcvrType: reflect.TypeOf(rcvr),
		methods:  make(map[string]*serviceMethod),
	}
	if name == "" {
		s.name = reflect.Indirect(s.rcvr).Type().Name()
		if !isExported(s.name) {
			return fmt.Errorf("rpc: type %q is not exported", s.name)
		}
	}
	if s.name == "" {
		return fmt.Errorf("rpc: no service name for type %q", s.rcvrType.String())
	}
	for i := 0; i < s.rcvrType.NumMethod(); i++ {
		method := s.rcvrType.Method(i)
		mtype := method.Type

		if method.PkgPath != "" {
			continue
		}

		if mtype.NumIn() != 3 {
			continue
		}

		args := mtype.In(1)
		if args.Kind() != reflect.Ptr || !isExportedOrBuiltin(args) {
			continue
		}

		reply := mtype.In(2)
		if reply.Kind() != reflect.Ptr || !isExportedOrBuiltin(args) {
			continue
		}

		// only return error
		if mtype.NumOut() != 1 {
			continue
		}

		if returnType := mtype.Out(0); returnType != typeOfError {
			continue
		}

		s.methods[method.Name] = &serviceMethod{
			method:    method,
			argsType:  args.Elem(),
			replyType: reply.Elem(),
		}
	}
	if len(s.methods) == 0 {
		return fmt.Errorf("rpc: %q has no exported methods of correct type", s.name)
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.services == nil {
		m.services = make(map[string]*service)
	} else if _, ok := m.services[s.name]; ok {
		return fmt.Errorf("rpc: service already registered: %q", s.name)
	}
	m.services[s.name] = s
	return nil
}

func isExported(name string) bool {
	rune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(rune)
}

func isExportedOrBuiltin(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return isExported(t.Name()) || t.PkgPath() == ""
}
