package handlers

import (
	"delivery-service/internal/core/ports"
	"delivery-service/internal/graph"
	"delivery-service/internal/graph/generated"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
)

func NewGraphQL(router *gin.Engine, service ports.DeliveryService) {
	router.POST("/query", graphqlHandler(service))
	router.GET("/", playgroundHandler())
}

func graphqlHandler(service ports.DeliveryService) gin.HandlerFunc {
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{DeliveryService: service}}))

	return func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	}
}

func playgroundHandler() gin.HandlerFunc {
	srv := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	}
}
