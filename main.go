package main

import (
	"net/http"
	"time"

	"activity_tracker/components"
	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

type Activity_Types int64

const (
	Code Activity_Types = iota
	Rest
	BJJ
	Entertainment
)

type activity struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Date         time.Time      `json:"date"`
	TotalMinutes int            `json:"totalMinutes"`
	CreatorID    string         `json:"creatorID"`
	Type         Activity_Types `json:"type"`
}

var activities = []activity{
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
	var newActivity activity
	if err := c.BindJSON(&newActivity); err != nil {
		return
	}

	activities = append(activities, newActivity)
	c.JSON(http.StatusCreated, activities)
}

func main() {
	component := components.Root()
	http.Handle("/", templ.Handler(component))
	http.ListenAndServe(":8080", nil)
}
