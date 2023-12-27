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
	StoreRoom(c *gin.Context, room *model.Room) error
}

type RoomRepositoryMongo struct {
	db *database.Database
}

type NotFound string

func (f NotFound) Error() string {
	return fmt.Sprintf("Not found id: %s", string(f))
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

func (repo *RoomRepositoryMongo) StoreRoom(c *gin.Context, room *model.Room) error {
	filter := filterById(room.Id)
	coll := repo.collection()
	coll.DeleteOne(c, filter)
	coll.InsertOne(c, room)
	return nil
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
			return nil, nil
		}
		return nil, fmt.Errorf("roomById %s: %v", id, err)
	}

	userRows, err := repo.db.Query(c, "select id, life_time_end, name, vote from user where poker_room_id=?", id)

	if err != nil {
		return nil, err
	}

	for userRows.Next() {
		var user model.User
		userRows.Scan(&user.Id, &user.LifeTimeEnd, &user.Name, &user.Vote)
		room.Users = append(room.Users, user)
		fmt.Println(user.Name)
	}

	userRows.Close()

	return &room, nil
}

func (repo *RoomRepositorySql) StoreRoom(c *gin.Context, room *model.Room) error {
	tx, err := repo.db.BeginTx(c)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if _, err := tx.ExecContext(c, "insert into poker_room(id, values_hidden) values(?, ?)", room.Id, room.HiddenVotes); err != nil {
		return err
	}

	for _, user := range room.Users {
		if _, err := tx.ExecContext(c, "insert into user(id, poker_room_id, name, life_time_end, vote) values(?, ?, ?, ?, ?)", user.Id, room.Id, user.Name, user.LifeTimeEnd, user.Vote); err != nil {
			return err
		}
	}

	err = tx.Commit()

	return err
}

// func (repo *RoomRepositorySql) findUsers(c *gin.Context, roomId string) (model.User, error){
// 	rows, err := repo.db.Query(c, "SELECT id, name, life_time_end, vote FROM user WHERE poker_room_id = ?", roomId)
// 	if err != nil {

// 	}
// 	defer rows.Close()
// 	for rows.Scan()() {

// 	}
// }
