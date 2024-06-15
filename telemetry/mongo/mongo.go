// Package mongo provides a wrapper around the mongo-driver package,
// simplifying the process of connecting to a MongoDB instance and
// performing operations on it.
package mongo

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
)

// OptFunc is a type that defines a function that modifies a Mongo instance.
type OptFunc func(*Mongo) error

// Mongo is a struct that holds the necessary information to connect to a MongoDB instance.
type Mongo struct {
	client *mongo.Client

	host     string
	port     string
	username string
	password string
}

// WithConnection is a function that returns an OptFunc which sets the connection details of a Mongo instance.
func WithConnection(host, port, username, password string) OptFunc {
	return func(m *Mongo) (err error) {
		m.host = host
		m.port = port
		m.username = username
		m.password = password
		return
	}
}

// New is a function that creates a new Mongo instance and connects to the MongoDB server.
// It applies the provided options to the Mongo instance and then attempts to connect to the server.
// If the connection is successful, it pings the server to ensure the connection is alive.
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

// Close is a method that disconnects the Mongo instance from the MongoDB server.
func (m *Mongo) Close(ctx context.Context) (err error) {
	err = m.client.Disconnect(ctx)
	return
}

// Collection is a method that returns a mongo.Collection instance for the specified collection name in the "telemetry" database.
func (m *Mongo) Collection(name string) *mongo.Collection {
	return m.client.Database("telemetry").Collection(name)
}
