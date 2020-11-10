package common

import (
	"errors"
	"github.com/26597925/EastCloud/pkg/logger"
	"sync"
)

type ConnManager struct {
	connections map[uint32]*Connection
	connLock    sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]*Connection),
	}
}

func (connMgr *ConnManager) Add(conn *Connection) {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	connMgr.connections[conn.GetConnID()] = conn

	logger.Info("connection add to ConnManager successfully: conn num = ", connMgr.Len())
}

func (connMgr *ConnManager) Remove(conn *Connection) {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	delete(connMgr.connections, conn.GetConnID())

	logger.Info("connection Remove ConnID=", conn.GetConnID(), " successfully: conn num = ", connMgr.Len())
}

func (connMgr *ConnManager) Get(connID uint32) (*Connection, error) {
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

func (connMgr *ConnManager) ClearConn() {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	for connID, conn := range connMgr.connections {
		err := conn.Close()
		if err != nil {
			logger.Error(err)
		}

		delete(connMgr.connections, connID)
	}

	logger.Info("Clear All Connections successfully: conn num = ", connMgr.Len())
}

func (connMgr *ConnManager) ShowConn() {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	for connID, _ := range connMgr.connections {
		logger.Info(connID)
	}
}