package main

import (
	"kowaluk/go-scrum-poker/controller"
	"kowaluk/go-scrum-poker/database"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func disableCors(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
}

func main() {
	var err error
	db, err := database.OpenClientLite()

	if err != nil {
		panic(err)
	}
	defer func() {
		if err = db.CloseClient(); err != nil {
			panic(err)
		}
	}()

	router := gin.Default()

	roomController := controller.CreateRoomController(db)
	router.GET("/rooms/:roomId", disableCors, roomController.GetRoom)
	router.POST("/rooms/:roomId/users", disableCors, roomController.PostUser)
	router.POST("/rooms", roomController.PostRoom)
	router.POST("/rooms/:roomId/users/:userId/heartbeat", roomController.PostHearbeat)
	router.DELETE("/rooms/:roomId/users/:userId", roomController.DeleteUser)

	router.POST("/rooms/:roomId/users/:userId/votes", roomController.PostUserVote)
	router.Run("localhost:8080")

}
