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
	sync.RWMutex                         // Mutex for async safety
	receivers    ReceiverList            // stores all open WebSocket Connections (whose Microphones get controlled)
	otps         RetentionMap            // all currently valid login tokens (just a few sec valid)
	handlers     map[string]EventHandler // map all supported Even-types to their Handler-Function
}

func NewManager(ctx context.Context) *Manager {
	manager := &Manager{
		otps:      NewRetentionMap(ctx, 3*time.Second),
		receivers: make(ReceiverList),
	}
	return manager
}

// Pool of all connected WebSocket clients. Uses username as
type ReceiverList map[string]*ReceiverClient

/*
*		Methods for adding or removing Websocket connections from our maps.
*		we Lock our mutex whenever reading/writing.
 */

// add the newly connected client to our List of all current clients
// can fail if username is already in taken
func (m *Manager) addReceiverClient(client *ReceiverClient) bool {
	if m.isUsernameTaken(client.username) {
		return false
	}
	m.Lock()
	defer m.Unlock()
	m.receivers[client.username] = client
	return true
}

// remove client and cleanup (example after they disconect/timeout)
func (m *Manager) removeReceiverClient(client *ReceiverClient) {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.receivers[client.username]; ok {
		client.conn.Close()                  // gracefully close the connection
		delete(m.receivers, client.username) // and delete the reference to the connection from current list
	}
}

// check if a username is already in use in any of the Receivers
func (m *Manager) isUsernameTaken(username string) bool {
	m.RLock()
	defer m.RUnlock()
	if _, ok := m.receivers[username]; ok {
		return true
	}
	return false
}

// TODO: once everything runs refactor the .clients map to be map[username] = Client. Would remove the whole looping over the map.
// check if username - password provided match up
func (m *Manager) isValidUsernamePw(username string, password string) (*ReceiverClient, bool) {
	m.RLock()
	defer m.RUnlock()
	if client, ok := m.receivers[username]; ok {
		if password == client.password {
			log.Println("u:", username, "pw:", password, "was found!")
			return client, true
		}
	}
	log.Println("user not found OR wrong pw!")
	return nil, false
}
