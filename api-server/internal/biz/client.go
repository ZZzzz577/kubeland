package biz

import (
	"api-server/api/v1/cluster"
	"api-server/internal/data"
	"api-server/internal/data/generated"
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sync"
	"time"
)

var (
	ErrClientNotFound = errors.New("kubernetes client not found")
)

type Client struct {
	*kubernetes.Clientset
	UpdateAt time.Time
}

type ClientManager struct {
	mu              sync.RWMutex
	db              *data.Data
	refreshInterval time.Duration
	ticker          *time.Ticker
	clients         map[uint64]*Client
}

func NewClientManager(db *data.Data) (*ClientManager, func(), error) {
	connections, err := db.ClusterConnection.Query().
		All(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("list clusters error")
		return nil, nil, err
	}
	clients := make(map[uint64]*Client, len(connections))

	cm := &ClientManager{
		db:              db,
		refreshInterval: 30 * time.Second,
		clients:         clients,
	}

	if err = cm.Refresh(); err != nil {
		log.Error().Err(err).Msg("refresh client error")
		return nil, nil, err
	}

	cm.ticker = time.NewTicker(cm.refreshInterval)
	go cm.RefreshLoop()
	return cm, func() {
		cm.ticker.Stop()
	}, nil
}

func (c *ClientManager) RefreshLoop() {
	for range c.ticker.C {
		if err := c.Refresh(); err != nil {
			log.Error().Err(err).Msg("refresh client error")
		}
	}
}

func (c *ClientManager) Refresh() error {
	connections, err := c.db.ClusterConnection.Query().
		All(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("list clusters error")
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	// add or update client
	dbConnIds := make(map[uint64]struct{}, len(connections))
	for _, conn := range connections {
		id := conn.ID
		dbConnIds[id] = struct{}{}
		if c.clients[id] == nil || conn.UpdatedAt.After(c.clients[id].UpdateAt) {
			client, err := GetClientByConfig(conn)
			if err != nil {
				log.Error().
					Err(err).
					Str("address", conn.Address).
					Msg("get client error")
				continue
			}
			c.clients[id] = client
		}
	}
	// delete client
	for id := range c.clients {
		if _, ok := dbConnIds[id]; !ok {
			delete(c.clients, id)
		}
	}
	return nil
}

func GetClientByConfig(conn *generated.ClusterConnection) (*Client, error) {
	config := &rest.Config{
		Host: conn.Address,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: []byte(conn.Ca),
		},
	}
	switch conn.Type {
	case uint8(cluster.Connection_TLS_CERT):
		config.CertData = []byte(conn.Cert)
		config.KeyData = []byte(conn.Key)
	case uint8(cluster.Connection_TLS_TOKEN):
		config.BearerToken = conn.Token
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	if _, err = clientSet.ServerVersion(); err != nil {
		return nil, err
	}
	return &Client{
		clientSet,
		conn.UpdatedAt,
	}, nil
}

func (c *ClientManager) Get(clusterId uint64) (*Client, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if client := c.clients[clusterId]; client == nil {
		return nil, ErrClientNotFound
	} else {
		return client, nil
	}
}
