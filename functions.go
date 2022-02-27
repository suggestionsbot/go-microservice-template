package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/keyauth/v2"
	"github.com/pelletier/go-toml"
	"log"
	"os"
	"time"
)

var config *toml.Tree

func handleServer() {
	app := fiber.New(fiber.Config{
		ErrorHandler: formErrorMessage,
	})

	app.Use(logger.New(logger.Config{
		Format:     config.Get("api.logger.format").(string),
		TimeFormat: config.Get("api.logger.time_format").(string),
		TimeZone:   config.Get("api.logger.timezone").(string),
	}))
	app.Use(keyauth.New(keyauth.Config{
		KeyLookup:    config.Get("api.auth.header_key").(string),
		AuthScheme:   config.Get("api.auth.header_prefix").(string),
		ErrorHandler: formErrorMessage,
		Validator:    validateAuthToken,
	}))
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: config.GetArray("api.cors.allow_origins").(string),
		AllowHeaders: config.GetArray("api.cors.allow_headers").(string),
	}))

	app.Get("/", getRootRoute)

	api := app.Group("/api")

	v1 := api.Group("/v1")

	v1.Get("/", getTestRoute)

	port := os.Getenv("API_PORT")
	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}

func handleConfig() {
	doc, err := toml.LoadFile("config.toml")
	if err != nil {
		log.Fatal(err)
	}

	config = doc

	fmt.Println("Config file loaded!")
}

func formJsonBody(data interface{}, success bool) fiber.Map {
	return fiber.Map{
		"data":    data,
		"success": success,
		"nonce":   time.Now().UnixMilli(),
	}
}

func formErrorMessage(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "A server side error has occurred."

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	} else if err.Error() != "" {
		message = err.Error()
	}

	return ctx.Status(code).JSON(formJsonBody(
		fiber.Map{
			"code":    code,
			"message": message,
		},
		false,
	))
}

func validateAuthToken(_ *fiber.Ctx, token string) (bool, error) {
	tk := os.Getenv("API_TOKEN")

	if token != tk {
		return false, nil
	}

	return true, nil
}
