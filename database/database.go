package database

import (
	"context"
	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const mongoUri = "mongodb://root:mongo@localhost:27017"
const dbName = "poker"

type Database struct {
	client *mongo.Client
	dbName string
}

type DatabaseLite struct {
	client *sql.DB
	dbName string
}

func OpenClientLite() (*DatabaseLite, error) {
	db, err := sql.Open("sqlite3", "poker.db")
	if err != nil {
		return nil, err
	}
	return &DatabaseLite{
		client: db,
		dbName: dbName,
	}, nil
}

func (d *DatabaseLite) QueryRow(c *gin.Context, query string, args string) *sql.Row {
	return d.client.QueryRowContext(c, query, args)
}

func (d *DatabaseLite) Query(c *gin.Context, query string, args string) (*sql.Rows, error) {
	return d.client.QueryContext(c, query, args)
}

func (d *DatabaseLite) Exec(c *gin.Context, query string, args string) (sql.Result, error) {
	return d.client.ExecContext(c, query, args)
}

func (d *DatabaseLite) BeginTx(c *gin.Context) (*sql.Tx, error) {
	return d.client.BeginTx(c, nil)
}

func (d *DatabaseLite) CloseClient() error {
	return d.client.Close()
}

func OpenClient() (*Database, error) {
	serverApi := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoUri).SetServerAPIOptions(serverApi)
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return nil, err
	}
	db := &Database{
		client: client,
		dbName: dbName,
	}
	return db, nil
}

func (d *Database) CloseClient() error {
	return d.client.Disconnect(context.Background())
}

func (d *Database) Collection(name string) *mongo.Collection {
	return d.client.Database(d.dbName).Collection(name)
}
