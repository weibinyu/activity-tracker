package main

import (
	"net/http"
	"time"

	"activity_tracker/components"
	"activity_tracker/db"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

var activities = []db.Activity{
	{
		ID: "1", Name: "Programming", Date: time.Now(),
		CreatorID: "user", TotalMinutes: 60,
	},
	{
		ID: "2", Name: "Programming", Date: time.Now(),
		CreatorID: "user", TotalMinutes: 60,
	},
	{
		ID: "3", Name: "Programming", Date: time.Now(),
		CreatorID: "user", TotalMinutes: 60,
	},
}

func getActivities(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", activities)
}

func postActivity(c *gin.Context) {
	var newActivity db.Activity
	if err := c.BindJSON(&newActivity); err != nil {
		return
	}

	activities = append(activities, newActivity)
	c.JSON(http.StatusCreated, activities)
}

func main() {
	component := components.Root(activities)
	http.Handle("/", templ.Handler(component))

	http.ListenAndServe(":8080", nil)
}
