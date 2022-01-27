package users

import (
	"convention.ninja/internal/auth"
	"convention.ninja/internal/users/business"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(grp fiber.Router) {

	userGrp := grp.Group("users")
	userGrp.Get("/", auth.NewUserRequired(business.GetUsers))
	userGrp.Post("/", business.CreateUser)
	userGrp.Get("/me", auth.NewUserRequired(business.GetMe))
	userGrp.Delete("/me", auth.NewUserRequired(business.DeleteMe))
	userGrp.Get("/:userId", auth.NewUserRequired(business.GetUser))
	userGrp.Patch("/:userId", auth.NewUserRequired(business.UpdateUser))
}
