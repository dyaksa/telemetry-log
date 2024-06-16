// Package telemetry provides functionality for telemetry logging.
package telemetry

import (
	"context"
	"fmt"
	"github.com/dyaksa/telemetry-log/telemetry/mongo"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

// MongoHook is a struct that holds the necessary information for a MongoDB hook.
type MongoHook struct {
	Client   *mongo.Mongo  // Client is a pointer to a Mongo instance.
	Timeout  time.Duration // Timeout is the duration before the hook times out.
	WithHook bool          // WithHook is a boolean that determines whether the hook is active.
}

// Fire is a method that logs an entry to a MongoDB collection.
// If the entry level is "error" and the hook is active, it logs the entry to the "application_trace" collection.
// Otherwise, it logs a sample entry to the "application_log" collection.
func (m *MongoHook) Fire(e *logrus.Entry) error {
	switch {
	case e.Level.String() == logrus.ErrorLevel.String() && m.WithHook:
		_, err := m.Client.Collection("application_trace").InsertOne(context.TODO(), bson.D{
			{"level", e.Level.String()},
			{"trace_date", e.Time},
			{"func", e.Data["func"]},
			{"file", fmt.Sprintf("%s:%d", e.Data["file"], e.Data["line"])},
			{"trace", e.Data["trace"]},
		})
		if err != nil {
			return err
		}
		break
	default:
		_, err := m.Client.Collection("application_log").InsertOne(context.TODO(), bson.D{
			{"trace_date", e.Time},
			{"func", e.Data["func"]},
			{"file", fmt.Sprintf("%s:%d", e.Data["file"], e.Data["line"])},
		})
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

// Levels is a method that returns all logrus levels.
func (m *MongoHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
