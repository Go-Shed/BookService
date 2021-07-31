package notification

import (
	"context"
	"fmt"
	"log"
	"shed/bookservice/common/constants"
	"sync"
	"time"

	"github.com/NaySoftware/go-fcm"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NotificationRepo struct {
	mongoClient *mongo.Client
}

//// mongo object in db
//// this is how notification is stored on mongodb
type Notification struct {
	Id               primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	UserToSend       MongoUser          `json:"userToSend"`
	NotificationType string             `json:"notificationType"`
	UserBy           MongoUser          `json:"userBy"`
	SourceId         string             `json:"sourceId"`
	FCMToken         string             `json:"fcmToken"`
	CommentId        string             `json:"commentId"`
}

type MongoUser struct {
	UserName string `json:"userName"`
	UserId   string `json:"userId"`
}

type NotificationToSend struct {
	FCMToken     string
	LastActionBy string //// name of person who liked/commented the photo last
	Times        int
	PostId       string
}

func NewNotificationRepo() NotificationRepo {
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

	return NotificationRepo{
		mongoClient: mongoClient,
	}
}

func (repo *NotificationRepo) AddNotificationTODB(notification Notification) error {

	if len(notification.FCMToken) == 0 {
		return nil
	}

	client := repo.mongoClient

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := client.Database("shed").Collection("notification")

	document := bson.D{
		{Key: "userToSend", Value: notification.UserToSend},
		{Key: "notificationType", Value: notification.NotificationType},
		{Key: "userBy", Value: notification.UserBy},
		{Key: "createdAt", Value: time.Now()},
		{Key: "isSent", Value: false},
		{Key: "fcmToken", Value: notification.FCMToken},
	}
	if notification.NotificationType == constants.NOTIFICATION_TYPE_COMMENT { ////// To support unqiue indexing over sourceId
		document = append(document, bson.E{Key: "commentId", Value: notification.SourceId})
	} else {
		document = append(document, bson.E{Key: "sourceId", Value: notification.SourceId})
	}

	_, err := collection.InsertOne(ctx, document)

	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func (repo *NotificationRepo) SendNotificationsToAll() error {
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

	if len(results) == 0 {
		return nil
	}

	likes, comments, _, follow := getNotificationBatches(results)

	var wg sync.WaitGroup
	wg.Add(2)

	go func(wg *sync.WaitGroup) {
		for _, item := range likes {

			var text string
			if item.Times > 1 {
				text = fmt.Sprintf("%s and %d others liked on your post", item.LastActionBy, item.Times-1)
			} else {
				text = fmt.Sprintf("%s liked your post", item.LastActionBy)
			}
			err := repo.sendNotification(text, item.FCMToken, item.PostId)

			if err != nil {
				log.Print(err)
				break
			}
		}

		// for _, like := range commentLikes { ///// TODO
		// 	text := fmt.Sprintf("%s and %d others liked your comment", like.LastActionBy, like.Times)
		// 	err := repo.sendNotification(text, like.FCMToken)

		// 	if err != nil {
		// 		log.Print(err)
		// 		break
		// 	}
		// }

		wg.Done()
	}(&wg)

	go func(wg *sync.WaitGroup) {
		for _, comment := range comments {
			var text string
			if comment.Times > 1 {
				text = fmt.Sprintf("%s and %d others commented on your post", comment.LastActionBy, comment.Times-1)
			} else {
				text = fmt.Sprintf("%s commented on your post", comment.LastActionBy)
			}
			err := repo.sendNotification(text, comment.FCMToken, comment.PostId)

			if err != nil {
				log.Print(err)
				break
			}
		}

		for _, item := range follow {
			var text string
			if item.Times > 1 {
				text = fmt.Sprintf("%s and %d others started followed you", item.LastActionBy, item.Times-1)
			} else {
				text = fmt.Sprintf("%s started following you", item.LastActionBy)
			}

			err := repo.sendNotification(text, item.FCMToken, item.PostId)

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

func (repo *NotificationRepo) updateNotification(result []Notification) error {

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

func (repo *NotificationRepo) sendNotification(text, token, postId string) error {

	const serverKey = "AAAAckyK83s:APA91bGzbVzVPm0SA9VZ7nitDeS1uFqoiIrkkyDonMWPZd4vMP5-IPt-hs0g6MzLw13-EBP53vpCKQQhsa2E_30Vlx59N3gEwI6sLf-U-LDPqgE0nQOfFnFqoMXCD6yJCRtpDCXYKBxx"

	var NP fcm.NotificationPayload
	NP.Title = "SHED"
	NP.Body = text
	NP.ClickAction = "OPEN_SHED_ACTIVITY"

	data := make(map[string]string)

	if len(postId) != 0 {
		data["deepLink"] = fmt.Sprintf("shed://postDetailScreen/%s", postId)
	}

	ids := []string{
		token,
	}

	c := fcm.NewFcmClient(serverKey)
	c.NewFcmRegIdsMsg(ids, data)
	c.SetNotificationPayload(&NP)

	status, err := c.Send()
	if err == nil {
		status.PrintResults()
	} else {
		fmt.Println(err)
	}

	return nil
}

/// TODO ---- Make this part a little better (DRY)
//// get all notifification to send to user
/// one batch for comment
/// and batch map for likes
func getNotificationBatches(results []Notification) (map[string]NotificationToSend, map[string]NotificationToSend, map[string]NotificationToSend, map[string]NotificationToSend) {

	commentNotifications := make(map[string]NotificationToSend)
	likeNotification := make(map[string]NotificationToSend)
	commentLikeNotification := make(map[string]NotificationToSend)
	followNotification := make(map[string]NotificationToSend)

	for _, notification := range results {

		if notification.NotificationType == constants.NOTIFICATION_TYPE_COMMENT {

			if val, ok := commentNotifications[notification.UserToSend.UserName]; ok {
				item := NotificationToSend{
					FCMToken:     notification.FCMToken,
					LastActionBy: notification.UserBy.UserName,
					Times:        val.Times + 1,
					PostId:       notification.SourceId,
				}
				commentNotifications[notification.UserToSend.UserName] = item
				continue
			}

			item := NotificationToSend{
				FCMToken:     notification.FCMToken,
				LastActionBy: notification.UserBy.UserName,
				Times:        1,
				PostId:       notification.SourceId,
			}
			commentNotifications[notification.UserToSend.UserName] = item
		} else if notification.NotificationType == constants.NOTIFICATION_TYPE_LIKE {

			if val, ok := likeNotification[notification.UserToSend.UserName]; ok {
				item := NotificationToSend{
					FCMToken:     notification.FCMToken,
					LastActionBy: notification.UserBy.UserName,
					Times:        val.Times + 1,
					PostId:       notification.SourceId,
				}
				likeNotification[notification.UserToSend.UserName] = item
				continue
			}

			item := NotificationToSend{
				FCMToken:     notification.FCMToken,
				LastActionBy: notification.UserBy.UserName,
				Times:        1,
				PostId:       notification.SourceId,
			}
			likeNotification[notification.UserToSend.UserName] = item
		} else if notification.NotificationType == constants.NOTIFICATION_TYPE_COMMENT_LIKE {

			if val, ok := commentLikeNotification[notification.UserToSend.UserName]; ok {
				item := NotificationToSend{
					FCMToken:     notification.FCMToken,
					LastActionBy: notification.UserBy.UserName,
					Times:        val.Times + 1,
				}
				commentLikeNotification[notification.UserToSend.UserName] = item
				continue
			}

			item := NotificationToSend{
				FCMToken:     notification.FCMToken,
				LastActionBy: notification.UserBy.UserName,
				Times:        1,
			}
			commentLikeNotification[notification.UserToSend.UserName] = item
		} else if notification.NotificationType == constants.NOTIFICATION_TYPE_FOLLOW {

			if val, ok := followNotification[notification.UserToSend.UserName]; ok {
				item := NotificationToSend{
					FCMToken:     notification.FCMToken,
					LastActionBy: notification.UserBy.UserName,
					Times:        val.Times + 1,
				}
				followNotification[notification.UserToSend.UserName] = item
				continue
			}

			item := NotificationToSend{
				FCMToken:     notification.FCMToken,
				LastActionBy: notification.UserBy.UserName,
				Times:        1,
			}
			followNotification[notification.UserToSend.UserName] = item
		}

	}

	return likeNotification, commentNotifications, commentLikeNotification, followNotification
}
