package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bouroo/neversitup-backend-test/internal/order/handlers"
	"github.com/bouroo/neversitup-backend-test/pkg/config"
	"github.com/bouroo/neversitup-backend-test/pkg/database"
	"github.com/bouroo/neversitup-backend-test/pkg/present"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/segmentio/encoding/json"
	"github.com/spf13/viper"
)

func init() {
	config.LoadConfig("order")
	viper.Set("production", bool(viper.GetString("go_env") == "production"))
}

func main() {
	// Connect to the database
	dsn, err := database.BuildConnStr(viper.GetString("db_driver"), viper.GetString("db_name"), viper.GetString("db_host"), viper.GetInt("db_port"), viper.GetString("db_user"), viper.GetString("db_password"))
	if err != nil {
		log.Fatal(err)
	}
	db, err := database.Connect(viper.GetString("db_driver"), dsn)
	if err != nil {
		log.Fatal(err)
	}
	if !viper.GetBool("production") {
		db = db.Debug()
	}

	// Create a new Fiber instance
	app := fiber.New(fiber.Config{
		ErrorHandler: present.CustomErrorHandler,
		Prefork:      viper.GetBool("production"),
		// allow 8k header for JWT
		ReadBufferSize: 8 * 1024,
		// https://docs.gofiber.io/guide/faster-fiber#custom-json-encoderdecoder
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		// allow behide reverse proxy
		EnableTrustedProxyCheck: true,
	})

	// Register middleware
	app.Use(recover.New(recover.Config{EnableStackTrace: !viper.GetBool("production")}))
	app.Use(requestid.New())
	app.Use(logger.New())
	app.Use(cors.New())

	// Register handlers
	handlers.SetupRoutes(app, db)

	// Start the server with graceful shutdown
	go func() {
		if err := app.Listen(":" + viper.GetString("app_port")); err != nil {
			log.Panic(err)
		}
	}()

	// Create channel to signify a signal being sent
	ch := make(chan os.Signal, 1)
	// When an interrupt or termination signal is sent, notify the channel
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	// This blocks the main thread until an interrupt is received
	<-ch
	fmt.Println("Gracefully shutting down...")
	_ = app.ShutdownWithTimeout(60 * time.Second)

	fmt.Println("Running cleanup tasks...")
	// Your cleanup tasks go here
	// if sqlDB, err := db.DB(); err == nil {
	// 	sqlDB.Close()
	// }
}
