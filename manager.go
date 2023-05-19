/*
*
*
 */

package main

import (
	"context"
	"sync"
	"time"
)

type Manager struct {
	controllers  ListController // stores all open WebSocket Connections (who act as controllers)
	receivers    ReceiverList   // stores all open WebSocket Connections (whose Microphones get controlled)
	sync.RWMutex                // Mutex for async safety
	otps         RetentionMap   // all currently valid login tokens (just a few sec valid)
}

func NewManager(ctx context.Context) *Manager {
	manager := &Manager{
		controllers: make(ListController),
		otps:        NewRetentionMap(ctx, 3*time.Second),
		receivers:   make(ReceiverList),
	}
	return manager
}

/*
*		Methods for adding or removing Websocket connections from our maps.
*		we Lock our mutex whenever reading/writing.
 */

type ListController map[*ControllerClient]bool

// add the newly connected client to our List of all current clients
func (m *Manager) addControllerClient(client *ControllerClient) {
	m.Lock()
	defer m.Unlock()
	m.controllers[client] = true
}

// remove client and cleanup (example after they disconect/timeout)
func (m *Manager) removeControllerClient(client *ControllerClient) {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.controllers[client]; ok {
		client.conn.Close()           // gracefully close the connection
		delete(m.controllers, client) // and delete the reference to the connection from current list
	}
}

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
