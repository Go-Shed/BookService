package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

var Client *auth.Client

func init() {
	Client = createFireBaseClient()
}

type RequestContextKey struct {
	Key string `json:"key"`
}

/////export GOOGLE_APPLICATION_CREDENTIALS="/home/user/Downloads/service-account-file.json"
func AuthorizeUsers(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authToken := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authToken) != 2 {
			fmt.Println("Malformed token")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Malformed Token"))
			return
		}

		ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*500)
		token, err := Client.VerifyIDToken(ctx, authToken[1])

		if err != nil {
			fmt.Println("Malformed token")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}

		requestContext := context.WithValue(r.Context(), RequestContextKey{Key: "uid"}, token)
		next.ServeHTTP(w, r.WithContext(requestContext))
	})
}

func createFireBaseClient() *auth.Client {
	app, err := firebase.NewApp(context.Background(), &firebase.Config{ProjectID: "shed-477d9"})
	if err != nil {
		log.Fatalf("error initializing firebase auth: %v\n", err)
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*500)
	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting firebase Auth client: %v\n", err)
	}

	return client
}
