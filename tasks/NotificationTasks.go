package tasks

import (
	"fmt"
	"shed/bookservice/repos/notification"

	"github.com/jasonlvhit/gocron"
)

const JOB_RUNTIME = 30

func ScheduleNotificationTasks() {

	fmt.Println("Scheduing tasks for noitfication")

	repo := notification.NewNotificationRepo()
	gocron.Every(JOB_RUNTIME).Second().Do(func() {
		repo.SendNotificationsToAll()
	})
	<-gocron.Start()
}
