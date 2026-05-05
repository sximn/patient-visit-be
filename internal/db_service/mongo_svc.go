package db_service

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrNotFound = fmt.Errorf("document not found")
var ErrConflict = fmt.Errorf("conflict: document already exists")

type Identifiable interface {
	GetID() string
}

type DbService[T Identifiable] interface {
	CreateDocument(ctx context.Context, document T) error
	FindDocument(ctx context.Context, id string) (T, error)
	FindMany(ctx context.Context, filter bson.M) ([]T, error)
	UpdateDocument(ctx context.Context, document T) error
	DeleteDocument(ctx context.Context, id string) error
	Disconnect(ctx context.Context) error
}

type MongoServiceConfig struct {
	ServerHost string
	ServerPort int
	UserName   string
	Password   string
	DbName     string
	Collection string
	Timeout    time.Duration
}

type mongoSvc[T Identifiable] struct {
	MongoServiceConfig
	client     atomic.Pointer[mongo.Client]
	clientLock sync.Mutex
}

func NewMongoService[T Identifiable](config MongoServiceConfig) DbService[T] {
	env := func(name, def string) string {
		if v, ok := os.LookupEnv(name); ok {
			return v
		}
		return def
	}

	svc := &mongoSvc[T]{}
	svc.MongoServiceConfig = config

	if svc.ServerHost == "" {
		svc.ServerHost = env("PATIENT_VISIT_API_MONGODB_HOST", "localhost")
	}

	if svc.ServerPort == 0 {
		p, _ := strconv.Atoi(env("PATIENT_VISIT_API_MONGODB_PORT", "27017"))
		svc.ServerPort = p
	}

	if svc.UserName == "" {
		svc.UserName = env("PATIENT_VISIT_API_MONGODB_USERNAME", "")
	}

	if svc.Password == "" {
		svc.Password = env("PATIENT_VISIT_API_MONGODB_PASSWORD", "")
	}

	if svc.DbName == "" {
		svc.DbName = env("PATIENT_VISIT_API_MONGODB_DATABASE", "sukus-patient-visit")
	}

	if svc.Collection == "" {
		log.Fatal("Collection must be defined per service instance")
	}

	if svc.Timeout == 0 {
		sec, _ := strconv.Atoi(env("PATIENT_VISIT_API_MONGODB_TIMEOUT_SECONDS", "10"))
		svc.Timeout = time.Duration(sec) * time.Second
	}

	return svc
}

func (m *mongoSvc[T]) connect(ctx context.Context) (*mongo.Client, error) {
	if c := m.client.Load(); c != nil {
		return c, nil
	}

	m.clientLock.Lock()
	defer m.clientLock.Unlock()

	if c := m.client.Load(); c != nil {
		return c, nil
	}

	ctx, cancel := context.WithTimeout(ctx, m.Timeout)
	defer cancel()

	uri := fmt.Sprintf("mongodb://%s:%d", m.ServerHost, m.ServerPort)
	if m.UserName != "" {
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%d",
			m.UserName, m.Password, m.ServerHost, m.ServerPort)
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	m.client.Store(client)
	return client, nil
}

func (m *mongoSvc[T]) collection(ctx context.Context) (*mongo.Collection, error) {
	client, err := m.connect(ctx)
	if err != nil {
		return nil, err
	}
	return client.Database(m.DbName).Collection(m.Collection), nil
}

func (m *mongoSvc[T]) CreateDocument(ctx context.Context, document T) error {
	ctx, cancel := context.WithTimeout(ctx, m.Timeout)
	defer cancel()

	col, err := m.collection(ctx)
	if err != nil {
		return err
	}

	id := document.GetID()

	// check conflict
	err = col.FindOne(ctx, bson.M{"id": id}).Err()
	if err == nil {
		return ErrConflict
	}
	if err != mongo.ErrNoDocuments {
		return err
	}

	_, err = col.InsertOne(ctx, document)
	return err
}

func (m *mongoSvc[T]) FindDocument(ctx context.Context, id string) (T, error) {
	var zero T

	ctx, cancel := context.WithTimeout(ctx, m.Timeout)
	defer cancel()

	col, err := m.collection(ctx)
	if err != nil {
		return zero, err
	}

	var result T
	err = col.FindOne(ctx, bson.M{"id": id}).Decode(&result)

	if err == mongo.ErrNoDocuments {
		return zero, ErrNotFound
	}
	if err != nil {
		return zero, err
	}

	return result, nil
}

func (m *mongoSvc[T]) FindMany(ctx context.Context, filter bson.M) ([]T, error) {
	ctx, cancel := context.WithTimeout(ctx, m.Timeout)
	defer cancel()

	col, err := m.collection(ctx)
	if err != nil {
		return nil, err
	}

	cursor, err := col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	results := make([]T, 0)
	for cursor.Next(ctx) {
		var doc T
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		results = append(results, doc)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (m *mongoSvc[T]) UpdateDocument(ctx context.Context, document T) error {
	ctx, cancel := context.WithTimeout(ctx, m.Timeout)
	defer cancel()

	col, err := m.collection(ctx)
	if err != nil {
		return err
	}

	id := document.GetID()

	res, err := col.ReplaceOne(ctx, bson.M{"id": id}, document)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return ErrNotFound
	}

	return nil
}

func (m *mongoSvc[T]) DeleteDocument(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, m.Timeout)
	defer cancel()

	col, err := m.collection(ctx)
	if err != nil {
		return err
	}

	res, err := col.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return ErrNotFound
	}

	return nil
}

func (m *mongoSvc[T]) Disconnect(ctx context.Context) error {
	client := m.client.Load()
	if client == nil {
		return nil
	}
	return client.Disconnect(ctx)
}
