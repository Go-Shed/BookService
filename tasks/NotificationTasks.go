package tasks

import (
	"fmt"
	"shed/bookservice/repos/notification"

	"github.com/jasonlvhit/gocron"
)

func ScheduleNotificationTasks() {

	fmt.Println("Scheduing tasks for noitfication")

	repo := notification.NewNotificationRepo()
	gocron.Every(2).Minute().Do(func() {
		repo.SendNotificationsToAll()
	})
	<-gocron.Start()
}
