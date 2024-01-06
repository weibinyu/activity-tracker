package main

import (
	"fmt"
	"net/http"
	"path"
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
		ID: "2", Name: "Gaming", Date: time.Now(),
		CreatorID: "user", TotalMinutes: 60,
	},
	{
		ID: "3", Name: "Bjjing", Date: time.Now(),
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

func deleteActivity(ID string) {
	i := SliceIndex(len(activities), func(i int) bool { return activities[i].ID == ID })
	if i == -1 {
		return
	}

	activities[i] = activities[len(activities)-1]
	activities = activities[:len(activities)-1]
}

func main() {
	component := components.Root(activities)

	http.Handle("/", templ.Handler(component))
	http.HandleFunc("/activities/", func(w http.ResponseWriter, r *http.Request) {
		deleteActivity(path.Base(r.URL.Path))
		components.ActivitiesComp(activities).Render(r.Context(), w)
	})

	fmt.Print("Server listening on port 8080")
	http.ListenAndServe("localhost:8080", nil)
}

func SliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}
