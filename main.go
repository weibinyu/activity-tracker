package main

import (
	"fmt"
	"net/http"
	"time"

	"activity_tracker/components"
	"activity_tracker/db"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
)

var activities = []db.ActivityRecord{
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

func postActivity(newActivity db.ActivityRecord) {
	activities = append(activities, newActivity)
}

func deleteActivity(ID string) {
	i := SliceIndex(len(activities), func(i int) bool {
		return activities[i].ID == ID
	})

	if i == -1 {
		return
	}

	activities[i] = activities[len(activities)-1]
	activities = activities[:len(activities)-1]
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		components.Root(activities).Render(r.Context(), w)
	})
	r.HandleFunc("/activities/{id}", func(w http.ResponseWriter, r *http.Request) {
		deleteActivity(mux.Vars(r)["id"])
	}).Methods("DELETE")

	r.HandleFunc("/activities/", func(w http.ResponseWriter, r *http.Request) {
		var newActivity = db.ActivityRecord{
			ID: "3", Name: "Bjjing", Date: time.Now(),
			CreatorID: "user", TotalMinutes: 60,
		}
		postActivity(newActivity)
		components.ActivitiesTableRow(newActivity).Render(r.Context(), w)
	}).Methods("POST")

	fmt.Println(time.Now().Format("2006-1-02"))
	fmt.Println("Server listening on port 8080")
	http.ListenAndServe("localhost:8080", r)
}

func SliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}
