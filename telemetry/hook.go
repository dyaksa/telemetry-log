package telemetry

import (
	"context"
	"fmt"
	"github.com/dyaksa/telemetry-log/telemetry/mongo"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type MongoHook struct {
	Client   *mongo.Mongo
	Timeout  time.Duration
	WithHook bool
}

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
		return nil
	}
	return nil
}

func (m *MongoHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
