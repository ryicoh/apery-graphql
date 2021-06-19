package server

import (
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ryicoh/apery-graphql/pkg"
)

const graphqlEndpoint = "/v1/graphql"

type HasuraConfig struct {
	HasuraEndpoint    string `validate:"required"`
	HasuraAdminSecret string `validate:"required"`
}

func NewServer(port int) *http.Server {
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	router := gin.New()

	router.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	router.Use(gin.Recovery())
	router.Use(cors.Default())

	gqlSvr := handler.NewDefaultServer(pkg.NewExecutableSchema(pkg.Config{Resolvers: &Resolvers{}}))

	router.POST(graphqlEndpoint, gin.WrapH(gqlSvr))
	router.GET("/", playgroundHandler())
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}
	return srv
}

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", graphqlEndpoint)

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
