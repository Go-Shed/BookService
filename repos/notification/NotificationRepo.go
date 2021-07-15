package notification

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	firebase "firebase.google.com/go"
	messaging "firebase.google.com/go/messaging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type notificationRepo struct {
	mongoClient *mongo.Client
}

type Notification struct {
	Id               primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	UserToSend       string             `json:"userToSend"`
	NotificationType string             `json:"notificationType"`
	UserBy           string             `json:"userBy"`
	PostId           string             `json:"postId"`
	FCMToken         string             `json:"fcmToken"`
}

type NotificationToSend struct {
	FCMToken     string
	LastActionBy string //// name of person who liked/commented the photo last
	Times        int
}

func NewNotificationRepo() notificationRepo {
	uri := "mongodb+srv://troll:bq$2FxWkqNT!NVD@cluster0.dygiz.mongodb.net/shed?retryWrites=true&w=majority"
	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = mongoClient.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return notificationRepo{
		mongoClient: mongoClient,
	}
}

func (repo *notificationRepo) AddNotificationTODB(userToSend, notificationType, userBy, postId, fcmToken string, createdAt time.Time) error {

	if len(fcmToken) == 0 {
		return nil
	}

	client := repo.mongoClient

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := client.Database("shed").Collection("notification")

	_, err := collection.InsertOne(ctx, bson.D{
		{Key: "userToSend", Value: userToSend},
		{Key: "notificationType", Value: notificationType},
		{Key: "userBy", Value: userBy},
		{Key: "createdAt", Value: createdAt},
		{Key: "isSent", Value: false},
		{Key: "postId", Value: postId},
		{Key: "fcmToken", Value: fcmToken},
	})

	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func (repo *notificationRepo) SendNotificationsToAll() error {
	client := repo.mongoClient

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := client.Database("shed").Collection("notification")

	filter := bson.D{{"isSent", false}}
	var results []Notification
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Print(err)
		return err
	}

	if err = cursor.All(ctx, &results); err != nil {
		log.Print(err)
		return err
	}

	likes, comments := getNotificationBatches(results)

	var wg sync.WaitGroup
	wg.Add(2)

	go func(wg *sync.WaitGroup) {
		for _, like := range likes {
			text := fmt.Sprintf("%s and %d others liked your post", like.LastActionBy, like.Times)
			err := repo.sendNotification(text, like.FCMToken)

			if err != nil {
				log.Print(err)
				break
			}
		}
		wg.Done()
	}(&wg)

	go func(wg *sync.WaitGroup) {
		for _, comment := range comments {
			text := fmt.Sprintf("%s and %d others commented on your post", comment.LastActionBy, comment.Times)
			err := repo.sendNotification(text, comment.FCMToken)

			if err != nil {
				log.Print(err)
				break
			}
		}
		wg.Done()
	}(&wg)

	wg.Wait()
	repo.updateNotification(results)
	return nil
}

func (repo *notificationRepo) updateNotification(result []Notification) error {

	client := repo.mongoClient
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := client.Database("shed").Collection("notification")

	idsToUpdate := make([]primitive.ObjectID, len(result))

	for index, id := range result {
		idsToUpdate[index] = id.Id
	}

	filter := bson.M{"_id": bson.M{"$in": idsToUpdate}}
	update := bson.M{
		"$set": bson.M{
			"isSent": true,
		},
	}

	_, err := collection.UpdateMany(ctx, filter, update)

	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func (repo *notificationRepo) sendNotification(text, token string) error {

	app := &firebase.App{}
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Print("error getting Messaging client:", err)
		return err
	}

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: "Shed",
			Body:  text,
		},
		Token: token,
	}

	response, err := client.Send(ctx, message)
	if err != nil {
		log.Print(err)
		return err
	}

	fmt.Println("Successfully sent message:", response)
	return nil
}

//// get all notifification to send to user
/// one batch for comment
/// and batch map for likes
func getNotificationBatches(results []Notification) (map[string]NotificationToSend, map[string]NotificationToSend) {

	commentNotifications := make(map[string]NotificationToSend)
	likeNotification := make(map[string]NotificationToSend)

	for _, notification := range results {

		if notification.NotificationType == "comment" {

			if val, ok := commentNotifications[notification.UserToSend]; ok {
				item := NotificationToSend{
					FCMToken:     notification.FCMToken,
					LastActionBy: notification.UserBy,
					Times:        val.Times + 1,
				}
				commentNotifications[notification.UserToSend] = item
				continue
			}

			item := NotificationToSend{
				FCMToken:     notification.FCMToken,
				LastActionBy: notification.UserBy,
				Times:        1,
			}
			commentNotifications[notification.UserToSend] = item
		} else if notification.NotificationType == "like" {

			if val, ok := likeNotification[notification.UserToSend]; ok {
				item := NotificationToSend{
					FCMToken:     notification.FCMToken,
					LastActionBy: notification.UserBy,
					Times:        val.Times + 1,
				}
				likeNotification[notification.UserToSend] = item
				continue
			}

			item := NotificationToSend{
				FCMToken:     notification.FCMToken,
				LastActionBy: notification.UserBy,
				Times:        1,
			}
			likeNotification[notification.UserToSend] = item
		}

	}

	return likeNotification, commentNotifications
}
