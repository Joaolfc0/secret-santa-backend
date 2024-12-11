package di

import (
	"context"
	"log"
	. "service-secret-santa/config"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/dig"

	groupHandler "service-secret-santa/handlers/group"
	groupRepository "service-secret-santa/repositories/group"
	groupRoute "service-secret-santa/routes/group"
	groupService "service-secret-santa/services/group"
)

var Container *dig.Container

func InitializeDI(client *mongo.Client) {
	Container = dig.New()

	Container.Provide(func() *mongo.Client {
		return client
	})

	Container.Provide(groupRepository.NewGroupRepository)
	Container.Provide(groupService.NewGroupService)
	Container.Provide(groupHandler.NewGroupHandler)
}

func Invoke(defaultGroup *gin.RouterGroup) {
	if errGroupRoute := Container.Invoke(func(handler groupHandler.Handler) {
		groupRoute.Routes(defaultGroup, handler)
	}); errGroupRoute != nil {
		panic(errGroupRoute)
	}
}

func InitializeMongoClient() *mongo.Client {
	uri := Cfg.MongoURI
	if uri == "" {
		log.Fatal("You must set your 'MONGO_URI' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	return client
}
