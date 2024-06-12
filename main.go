package main

import (
	"dataProcessor/cache"
	"dataProcessor/config"
	"dataProcessor/controllers"
	"dataProcessor/database"
	"dataProcessor/routes"
	"dataProcessor/services"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	router := gin.Default()

	app, err := Initialize("config.yaml")
	if err != nil {
		log.Println("Error initializing processes")
	}

	routes.SetupRoutes(router, app)
	router.Run(":8080")
}

func Initialize(configPath string) (*controllers.Application, error) {
	log.Println("Initialize Starts: ", "")

	configData := config.GetConfig(configPath)

	mysqlStorage, err := database.ConnectMySQL(configData.MySQL)
	if err != nil {
		log.Fatalln("error initializing the mysql service ", err)
		panic(err)
	}

	log.Println("MYSQL Config Data: ", configData.MySQL)

	redisStorage, err := cache.ConnectRedis(configData.Redis)
	if err != nil {
		log.Fatalln("error initializing the redis service ", err)
		panic(err)
	}

	log.Println("Redis Config Data: ", configData.Redis)

	app := controllers.Application{
		ContactController: &controllers.ContactController{
			ContactService: &services.ContactService{
				MySQLRepo: mysqlStorage,
				RedisRepo: redisStorage,
			},
		},
	}
	return &app, nil

}
