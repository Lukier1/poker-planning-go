package controller

import (
	"fmt"
	"kowaluk/go-scrum-poker/database"
	"kowaluk/go-scrum-poker/model"
	"kowaluk/go-scrum-poker/repository"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomController struct {
	repo repository.RoomRepository
}

func CreateRoomController(db *database.DatabaseLite) *RoomController {
	return &RoomController{
		repo: repository.CreateLite(db),
	}
}

func LifetimeFromNow() string {
	bytes, _ := time.Now().Add(1 * time.Hour).MarshalText()
	return string(bytes)
}

func (controller *RoomController) PostUser(c *gin.Context) {
	roomId := c.Params.ByName("roomId")
	// userId := c.Params.ByName("userId")

	room, err := controller.repo.FindRoom(c, roomId)
	if room == nil {
		c.AbortWithStatus(404)
		return
	}

	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	var newUser model.User
	if err := c.BindJSON(&newUser); err != nil {
		c.AbortWithError(500, err)
		return
	}
	if newUser.Name == "" {
		c.AbortWithStatus(400)
		return
	}

	newUser.LifeTimeEnd = LifetimeFromNow()
	newUser.Vote = "0"
	newUser.Id = uuid.NewString()
	room.Users = append(room.Users, newUser)
	if err := controller.repo.StoreRoom(c, room); err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.IndentedJSON(http.StatusCreated, newUser)
}

func (controller *RoomController) DeleteUser(c *gin.Context) {
	roomId, userId, ok := getRoomAndUserId(c)
	if !ok {
		return
	}

	room, err := controller.repo.FindRoom(c, roomId)

	if err != nil {
		return
	}

	indexOfUser := -1
	for ind, user := range room.Users {
		if user.Id == userId {
			indexOfUser = ind
		}
	}
	if indexOfUser > -1 {
		room.Users = append(room.Users[:indexOfUser], room.Users[indexOfUser+1:]...)
	} else {
		c.Status(http.StatusNotFound)
	}
	controller.repo.StoreRoom(c, room)

}

func (controller *RoomController) PostRoom(c *gin.Context) {
	genId := uuid.New().String()

	room := &model.Room{
		Id: genId,
	}

	controller.repo.StoreRoom(c, room)

	c.IndentedJSON(http.StatusCreated, room)
}

func (controller *RoomController) GetRoom(c *gin.Context) {
	id := c.Params.ByName("roomId")
	result, err := controller.repo.FindRoom(c, id)
	if result == nil {
		c.AbortWithStatus(404)
	} else if err == nil {
		c.IndentedJSON(http.StatusOK, result)
	} else if err == mongo.ErrNoDocuments {
		c.IndentedJSON(http.StatusNotFound, nil)
	} else {
		fmt.Print(err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (controller *RoomController) PostHearbeat(c *gin.Context) {
	roomId, userId, ok := getRoomAndUserId(c)
	if !ok {
		return
	}

	room, err := controller.repo.FindRoom(c, roomId)

	if err != nil {
		return
	}

	for ind, user := range room.Users {
		if user.Id == userId {
			room.Users[ind].LifeTimeEnd = LifetimeFromNow()
		}
	}
	if err := controller.repo.StoreRoom(c, room); err != nil {
		c.Status(500)
		log.Fatal(err)
	}
}

func getRoomAndUserId(c *gin.Context) (string, string, bool) {
	roomId, ok := c.Params.Get("roomId")
	if !ok {
		return "", "", false
	}
	userId, ok := c.Params.Get("userId")
	if !ok {
		return "", "", false
	}
	return roomId, userId, true
}

func (controller *RoomController) PostUserVote(c *gin.Context) {
	roomId, userId, ok := getRoomAndUserId(c)
	if !ok {
		c.Status(http.StatusNotFound)
		return
	}
	room, err := controller.repo.FindRoom(c, roomId)

	if err != nil {
		c.Status(500)
	}

	var voteDTO model.VoteDTO
	c.BindJSON(&voteDTO)
	found := false
	for ind, user := range room.Users {
		if user.Id == userId {
			room.Users[ind].Vote = voteDTO.Vote
			found = true
		}
	}
	if !found {
		c.Status(http.StatusNotFound)
		return
	}
	controller.repo.StoreRoom(c, room)
}
