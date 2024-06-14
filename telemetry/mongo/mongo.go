package mongo

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
)

type OptFunc func(*Mongo) error

type Mongo struct {
	client *mongo.Client

	host     string
	port     string
	username string
	password string
}

func WithConnection(host, port, username, password string) OptFunc {
	return func(m *Mongo) (err error) {
		m.host = host
		m.port = port
		m.username = username
		m.password = password
		return
	}
}

func New(opts ...OptFunc) (*Mongo, error) {
	m := &Mongo{}
	for _, opt := range opts {
		err := opt(m)
		if err != nil {
			return nil, fmt.Errorf("fail to apply options: %w", err)
		}
	}

	serverUri := strings.Builder{}
	serverUri.WriteString("mongodb://")
	serverUri.WriteString(m.host)
	serverUri.WriteString(":")
	serverUri.WriteString(m.port)

	serverApi := options.ServerAPI(options.ServerAPIVersion1)
	authCred := options.Credential{Username: m.username, Password: m.password}

	opt := options.Client().ApplyURI(serverUri.String()).SetServerAPIOptions(serverApi).SetAuth(authCred)
	client, err := mongo.Connect(context.TODO(), opt)

	if err != nil {
		err = errors.Join(err, fmt.Errorf("fail to connect to mongo: %w", err))
	}

	m.client = client

	if m.client == nil {
		err = errors.Join(err, fmt.Errorf("fail to connect to mongo"))
	}

	if err = m.client.Database("telemetry").RunCommand(context.TODO(), map[string]string{"ping": "1"}).Err(); err != nil {
		_ = m.client.Disconnect(context.Background())
		return nil, fmt.Errorf("ping failed after connection: %w", err)
	}

	return m, err
}

func (m *Mongo) Close(ctx context.Context) (err error) {
	err = m.client.Disconnect(ctx)
	return
}

func (m *Mongo) Collection(name string) *mongo.Collection {
	return m.client.Database("telemetry").Collection(name)
}
