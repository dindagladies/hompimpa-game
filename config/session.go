package config

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var store *session.Store

func InitSessionStore() {
	store = session.New(session.Config{
		Expiration: 1 * time.Hour,
	})
}

func CreateUserSession(c *fiber.Ctx, ID int, username string) error {
	session, _ := store.Get(c)

	// if new session
	if session.Fresh() {
		session.Set("ID", ID)
		session.Set("username", username)
		if err := session.Save(); err != nil {
			return err
		}
	}

	return nil
}

func GetUserSession(c *fiber.Ctx) (map[string]interface{}, error) {
	session, err := store.Get(c)
	if err != nil {
		return nil, err
	}

	ID, _ := session.Get("ID").(int)
	username, _ := session.Get("username").(string)

	return map[string]interface{}{
		"ID":       ID,
		"username": username,
	}, nil
}
