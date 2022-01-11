package auth

import (
	"context"
	firebase "firebase.google.com/go"
	"github.com/gofiber/fiber/v2"
)

type IgnoreListBuilder struct {
	list []ignoreListItem
}
type ignoreListItem struct {
	Method string
	Path   string
}

func (l *IgnoreListBuilder) GetIgnoreList() []ignoreListItem {
	return l.list
}
func (l *IgnoreListBuilder) AddItem(method string, path string) {
	for i := range l.list {
		if l.list[i].Method == method && l.list[i].Path == path {
			return
		}
	}
	l.list = append(l.list, ignoreListItem{
		Method: method,
		Path:   path,
	})
}
func (l *IgnoreListBuilder) RemoveItem(method string, path string) {
	for i := range l.list {
		if l.list[i].Method == method && l.list[i].Path == path {
			l.list = append(l.list[:i], l.list[i+1:]...)
			return
		}
	}
}

func (l *IgnoreListBuilder) ContainsItem(method string, path string) bool {
	for i := range l.list {
		if l.list[i].Method == method && l.list[i].Path == path {
			return true
		}
	}
	return false
}

func NewIgnoreList() *IgnoreListBuilder {
	return &IgnoreListBuilder{
		list: make([]ignoreListItem, 0),
	}
}

type Config struct {
	FirebaseApp *firebase.App
	IgnoreList  *IgnoreListBuilder
}

type IdToken struct {
	Email  string
	UserId string
}

var config *Config

func GetConfig() *Config {
	return config
}

func New(cfg Config) fiber.Handler {
	config = &cfg
	return func(c *fiber.Ctx) error {

		// do firebase authentication
		idToken := c.Get(fiber.HeaderAuthorization)
		if len(idToken) == 0 {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid or expired Token")
		}

		client, err := cfg.FirebaseApp.Auth(context.Background())
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid or expired token")
		}

		token, err := client.VerifyIDTokenAndCheckRevoked(context.Background(), idToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid or expired token")
		}
		if !cfg.IgnoreList.ContainsItem(c.Method(), c.Path()) && token.Claims["email_verified"].(bool) == false {
			return c.Status(fiber.StatusUnauthorized).SendString("Email not verified")
		}
		c.Locals("idtoken", &IdToken{
			Email:  token.Claims["email"].(string),
			UserId: token.UID,
		})
		return c.Next()
	}
}
