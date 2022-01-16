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
	err = data.Connect(os.Getenv("SQL_DSN"))
	if err != nil {
		panic(err)
	}
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

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
