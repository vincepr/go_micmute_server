/*
*
*
 */

package main

import (
	"context"
	"log"
	"sync"
	"time"
)

type Manager struct {
	receivers    ReceiverList   // stores all open WebSocket Connections (whose Microphones get controlled)
	sync.RWMutex                // Mutex for async safety
	otps         RetentionMap   // all currently valid login tokens (just a few sec valid)
}

func NewManager(ctx context.Context) *Manager {
	manager := &Manager{
		otps:        NewRetentionMap(ctx, 3*time.Second),
		receivers:   make(ReceiverList),
	}
	return manager
}

/*
*		Methods for adding or removing Websocket connections from our maps.
*		we Lock our mutex whenever reading/writing.
*/

type ReceiverList map[*ReceiverClient]bool

// add the newly connected client to our List of all current clients
func (m *Manager) addReceiverClient(client *ReceiverClient) {
	m.Lock()
	defer m.Unlock()
	m.receivers[client] = true
}

// remove client and cleanup (example after they disconect/timeout)
func (m *Manager) removeReceiverClient(client *ReceiverClient) {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.receivers[client]; ok {
		client.conn.Close()         // gracefully close the connection
		delete(m.receivers, client) // and delete the reference to the connection from current list
	}
}

// check if a username is already in use in any of the Receivers
func (m *Manager) isUsernameInUse(username string) bool {
	m.RLock()
	defer m.RUnlock()
	for client := range m.receivers {
		if client.username == username {
			return true
		}
	}
	return false
}

// TODO: once everything runs refactor the .clients map to be map[username] = Client. Would remove the whole looping over the map.
// check if username - password provided match up
func (m *Manager) isValidUsernamePw(username string, password string) bool {
	m.RLock()
	defer m.RUnlock()
	log.Println("u:", username, "pw:", password)
	for client := range m.receivers {
		log.Println("client found with:", client.username, client.password)
		if client.username == username && client.password == password {
		log.Println("u:", username, "pw:", password, "was found!")
			return true
		}
	}
	return false
}
