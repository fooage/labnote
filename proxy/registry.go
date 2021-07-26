package proxy

import (
	"math/rand"
	"sync"
	"time"
)

var (
	mutex sync.Mutex
	// Store a hash table of all available servers.
	servers = make(map[string]*survivor)
)

// Existing server structure.
type survivor struct {
	addr string
	last time.Time
}

// Register a server with the registry, or update the available time of the server.
func serverRegister(addr string) {
	mutex.Lock()
	defer mutex.Unlock()
	server := servers[addr]
	if server == nil {
		servers[addr] = &survivor{addr: addr, last: time.Now()}
	} else {
		server.last = time.Now()
	}
}

// Get a list of all available servers for the request.
func getAliveServer() []string {
	mutex.Lock()
	defer mutex.Unlock()
	var alive []string
	for addr, server := range servers {
		if Timeout == 0 || server.last.Add(Timeout).After(time.Now()) {
			alive = append(alive, addr)
		} else {
			delete(servers, addr)
		}
	}
	return alive
}

// Get a random server let the user use, if there is no server could be used it
// will return a empty string.
func getRandomServer() string {
	servers := getAliveServer()
	if servers == nil {
		return ""
	} else {
		return servers[rand.Intn(len(servers))]
	}
}
