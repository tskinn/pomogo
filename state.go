package pomogo

import (
	"sync"
)

type State struct {
	mutex *sync.Mutex
	users map[string]User
}

func (state *State) StartSession() {
  state.mutex.Lock()
  defer state.mutex.Unlock()

}
