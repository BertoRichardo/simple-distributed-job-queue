package main

import (
	"jobqueue/config"
	"jobqueue/delivery/graphql"
	_dataloader "jobqueue/delivery/graphql/dataloader"
	"jobqueue/delivery/graphql/mutation"
	"jobqueue/delivery/graphql/query"
	"jobqueue/delivery/graphql/schema"
	"jobqueue/entity"
	"jobqueue/pkg/handler"
	"jobqueue/pkg/server"
	inmemrepo "jobqueue/repository/inmem"
	"jobqueue/service"
	"time"

	_graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	workerCount  = 10
	jobQueueSize = 100
)

func main() {
	setupLogger()
	logger := logrus.New()
	logger.SetReportCaller(true)

	e := server.New(config.Data.Server)
	e.Echo.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${remote_ip} ${time_rfc3339_nano} \"${method} ${path}\" ${status} ${bytes_out} \"${referer}\" \"${user_agent}\"\n",
	}))
	e.Echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.OPTIONS},
	}))

	jobQueue := make(chan *entity.Job, jobQueueSize)
	inMemDb := make(map[string]*entity.Job)
	jobRepository := inmemrepo.
		NewJobRepository().
		SetInMemConnection(inMemDb).
		Build()
	workerPool := service.NewWorkerPool(jobRepository, jobQueue, logger)
	go workerPool.Start(workerCount)
	jobService := service.NewJobService().
		SetJobRepository(jobRepository).
		SetJobQueue(jobQueue).
		Build()

	dataloader := _dataloader.New().SetJobRepository(jobRepository).SetBatchFunction().Build()
	jobMutation := mutation.NewJobMutation(jobService, dataloader)
	jobQuery := query.NewJobQuery(jobService, dataloader)
	rootResolver := graphql.New().SetJobMutation(jobMutation).SetJobQuery(jobQuery).Build()
	opts := []_graphql.SchemaOpt{_graphql.SubscribeResolverTimeout(10 * time.Second)}

	graphqlSchema := _graphql.MustParseSchema(schema.String(), rootResolver, opts...)

	e.Echo.POST("/graphql", handler.GraphQLHandler(&relay.Handler{Schema: graphqlSchema}), dataloader.EchoMiddelware)
	e.Echo.GET("/graphql", handler.GraphQLHandler(&relay.Handler{Schema: graphqlSchema}), dataloader.EchoMiddelware)
	e.Echo.GET("/graphiql", handler.GraphiQLHandler)

	e.Echo.Logger.Fatal(e.Start())
}

func setupLogger() {
	configLogger := zap.NewDevelopmentConfig()
	configLogger.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	configLogger.DisableStacktrace = true
	logger, _ := configLogger.Build()
	zap.ReplaceGlobals(logger)
}