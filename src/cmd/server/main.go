package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/samuelhgf/golang-study-restful-api/src/config"
	"github.com/samuelhgf/golang-study-restful-api/src/controllers"
	"github.com/samuelhgf/golang-study-restful-api/src/routes"
	"github.com/samuelhgf/golang-study-restful-api/src/services"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
)

var (
	server      *gin.Engine
	ctx         context.Context
	mongoclient *mongo.Client

	postService         services.PostService
	PostController      controllers.PostController
	postCollection      *mongo.Collection
	PostRouteController routes.PostRouteController
)

func init() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	ctx = context.TODO()

	// Connect to MongoDB
	mongoconn := options.Client().ApplyURI(config.DBUri)
	mongoclient, err = mongo.Connect(ctx, mongoconn)
	if err != nil {
		panic(err)
	}

	if err = mongoclient.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("MongoDB successfully connected...")

	postCollection = mongoclient.Database("golang_mongodb").Collection("posts")
	postService = services.NewPostService(postCollection, ctx)
	PostController = controllers.NewPostController(postService)
	PostRouteController = routes.NewPostControllerRoute(PostController)

	server = gin.Default()
}

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load config", err)
	}

	defer mongoclient.Disconnect(ctx)

	startGinServer(config)
}

func startGinServer(config config.Config) {
	router := server.Group("/")
	router.GET("healthcheck", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "All systems running"})
	})

	PostRouteController.PostRoute(router)
	log.Fatal(server.Run(":" + config.Port))
}
