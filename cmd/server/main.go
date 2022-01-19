package main

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"context"
	"convention.ninja/internal/auth"
	"convention.ninja/internal/data"
	"convention.ninja/internal/events"
	"convention.ninja/internal/inventory"
	"convention.ninja/internal/organizations"
	"convention.ninja/internal/users"
	usersMiddleware "convention.ninja/internal/users/middleware"
	firebase "firebase.google.com/go"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"log"
	"os"
)

var firebaseApp *firebase.App

func getSecret(secretKey string) (string, error) {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretKey + "/versions/latest",
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		return "", fmt.Errorf("failed to access secret version: %v", err)
	}
	return string(result.Payload.Data), nil
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	fbApp, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	firebaseApp = fbApp
	dsn := os.Getenv("SQL_DSN")
	if dsn == "" {
		dsn, err = getSecret(os.Getenv("SQL_DSN_KEY"))
		if err != nil {
			panic(err)
		}
	}
	err = data.Connect(dsn)
	if err != nil {
		panic(err)
	}

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		MaxAge:           86400,
	}))

	app.Use(auth.New(auth.Config{
		FirebaseApp: firebaseApp,
	}))

	app.Use(usersMiddleware.New())

	apiGrp := app.Group("api/")
	organizations.SetupRoutes(apiGrp)
	inventory.SetupRoutes(apiGrp)
	events.SetupRoutes(apiGrp)
	users.SetupRoutes(apiGrp)

	log.Fatal(app.Listen(":" + port))
}
