package main

import (
	"context"
	"convention.ninja/internal/auth"
	"convention.ninja/internal/data"
	"convention.ninja/internal/events"
	"convention.ninja/internal/inventory"
	"convention.ninja/internal/organizations"
	"convention.ninja/internal/users"
	usersMiddleware "convention.ninja/internal/users/middleware"
	firebase "firebase.google.com/go"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
)

var firebaseApp *firebase.App

func init() {
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	firebaseApp = app
}

func main() {

	err := data.Connect("...dsn here")
	if err != nil {
		panic(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	app := fiber.New()

	ignoreList := auth.NewIgnoreList()

	app.Use(auth.New(auth.Config{
		FirebaseApp: firebaseApp,
		IgnoreList:  ignoreList,
	}))

	app.Use(usersMiddleware.New())

	apiGrp := app.Group("api/")
	organizations.SetupRoutes(apiGrp, ignoreList)
	inventory.SetupRoutes(apiGrp, ignoreList)
	events.SetupRoutes(apiGrp, ignoreList)
	users.SetupRoutes(apiGrp, ignoreList)

	log.Fatal(app.Listen(":" + port))
}
