package main

import (
	"context"
	"log"

	. "service-secret-santa/config"
	"service-secret-santa/docs"
	"service-secret-santa/resources/di"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			Service Secret Santa
//	@version		1.0
//	@description	Service the secret santa website.
//	@termsOfService	http://swagger.io/terms/
//	@host
//	@BasePath /secret-santa
//
// @externalDocs.description	ReadMe
func main() {
	LoadConfig()
	mongoClient := di.InitializeMongoClient()
	defer func() {
		if err := mongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	docs.SwaggerInfo.Host = Cfg.SwaggerHost

	router := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AddAllowHeaders("Authorization")
	corsConfig.AllowAllOrigins = true

	router.Use(cors.New(corsConfig))

	secretSantaGroup := router.Group("/secret-santa")

	di.InitializeDI(mongoClient)
	di.Invoke(secretSantaGroup)
	secretSantaGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := router.Run(":" + Cfg.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
