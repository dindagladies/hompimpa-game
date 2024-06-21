package config

import (
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis/v3"
)

var store *session.Store

func createStorage() *redis.Storage {
	return redis.New(redis.Config{
		Host:      "127.0.0.1",
		Port:      6379,
		Username:  "",
		Password:  "",
		Database:  0,
		Reset:     false,
		TLSConfig: nil,
		PoolSize:  10 * runtime.GOMAXPROCS(0),
	})
}

func InitSessionStore() {
	store = session.New(session.Config{
		Expiration: 1 * time.Hour,
		Storage:    createStorage(),
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

func RemoveUserSession(c *fiber.Ctx) (bool, error) {
	session, err := store.Get(c)
	if err != nil {
		return false, err
	}

	if err := session.Destroy(); err != nil {
		return false, err
	}

	return true, nil
}
