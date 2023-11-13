package repository

import (
	"database/sql"
	"fmt"
	"kowaluk/go-scrum-poker/database"
	"kowaluk/go-scrum-poker/model"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomRepository interface {
	FindRoom(c *gin.Context, id string) (*model.Room, error)
	StoreRoom(c *gin.Context, room *model.Room)
}

type RoomRepositoryMongo struct {
	db *database.Database
}

func Create(db *database.Database) RoomRepository {
	return &RoomRepositoryMongo{
		db: db,
	}
}

func (repo *RoomRepositoryMongo) collection() *mongo.Collection {
	return repo.db.Collection("rooms")
}

func filterById(id string) bson.D {
	return bson.D{{Key: "_id", Value: id}}
}

func (repo *RoomRepositoryMongo) FindRoom(c *gin.Context, id string) (*model.Room, error) {
	result := &model.Room{}
	filter := filterById(id)
	err := repo.collection().FindOne(c, filter).Decode(result)

	return result, err
}

func (repo *RoomRepositoryMongo) StoreRoom(c *gin.Context, room *model.Room) {
	filter := filterById(room.Id)
	coll := repo.collection()
	coll.DeleteOne(c, filter)
	coll.InsertOne(c, room)
}

type RoomRepositorySql struct {
	db *database.DatabaseLite
}

func CreateLite(db *database.DatabaseLite) RoomRepository {
	return &RoomRepositorySql{
		db: db,
	}
}

func (repo *RoomRepositorySql) FindRoom(c *gin.Context, id string) (*model.Room, error) {
	row := repo.db.QueryRow(c, "select id, values_hidden from poker_room where id = ? ", id)

	var room model.Room

	if err := row.Scan(&room.Id, &room.HiddenVotes); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("roomById %s: no such room", id)
		}
		return nil, fmt.Errorf("roomById %s: %v", id, err)
	}

	return &room, nil
}

func (repo *RoomRepositorySql) StoreRoom(c *gin.Context, room *model.Room) {

}

func (repo *RoomRepositorySql) findUsers(c *gin.Context, roomId string) (model.User[], error){
	rows, err := repo.db.Query(c, "SELECT id, name, life_time_end, vote FROM user WHERE poker_room_id = ?", roomId)
	if err != nil {
			
	}
	defer rows.Close()
	for rows.Scan()() {

	}
}
