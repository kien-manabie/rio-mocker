package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kien-manabie/rio-mocker"
	"github.com/kien-manabie/rio-mocker/internal/config"
	"github.com/kien-manabie/rio-mocker/internal/log"
	"github.com/kien-manabie/rio-mocker/internal/setup"
	fs "github.com/kien-manabie/rio-mocker/internal/storage"
)

type AppOption func(*App)

// App defines app interface
type App struct {
	config      *config.Config
	fileStorage fs.FileStorage
	stubStore   rio.StubStore
	kit         *gin.Engine
}

// NewApp returns new app
func NewApp(ctx context.Context, config *config.Config, options ...AppOption) (*App, error) {
	stubStore, err := setup.ProvideStubStore(ctx, config)
	if err != nil {
		return nil, err
	}

	fileStorage, err := setup.ProvideFileStorage(ctx, config)
	if err != nil {
		return nil, err
	}

	app := &App{
		config:      config,
		stubStore:   stubStore,
		fileStorage: fileStorage,
		kit:         gin.New(),
	}

	for _, optionFunc := range options {
		optionFunc(app)
	}

	app.setup()
	return app, nil
}

// Start starts app
func (app *App) Start(ctx context.Context) error {
	address := app.config.ServerAddress
	srv := &http.Server{
		Addr:              address,
		Handler:           app.kit,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Info(ctx, "starting server", address)
	return srv.ListenAndServe()
}

func (app *App) setup() {
	app.kit.Use(RequestIDMiddleware())
	app.kit.Use(RequestTimeMiddleware())
	app.kit.Use(Recovery())
	app.initRoutes()
}

func (app *App) initRoutes() {
	app.kit.GET("/ping", app.handlePing)
	app.kit.DELETE("/reset", app.handleReset)
	app.kit.POST("/stub/create_many", app.handleCreate)
	app.kit.POST("/stub/upload", app.handleUpload)
	app.kit.GET("/stub/list", app.handleGetStubs)
	app.kit.POST("/proto/upload", app.handleUploadProto)
	app.kit.POST("/incoming_request/list", app.handleGetIncomingRequest)

	app.kit.Any("/echo/*path", func(ctx *gin.Context) {
		handler := rio.NewHandler(app.stubStore, app.fileStorage).WithBodyStoreThreshold(app.config.BodyStoreThreshold)
		handler.Handle(ctx.Writer, ctx.Request)
	})

	app.kit.Any("/:namespace/echo/*path", func(ctx *gin.Context) {
		namespace := ctx.Param("namespace")
		handler := rio.NewHandler(app.stubStore, app.fileStorage).
			WithBodyStoreThreshold(app.config.BodyStoreThreshold).
			WithNamespace(namespace)
		handler.Handle(ctx.Writer, ctx.Request)
	})
}

// TODO: This is opinionated solution to register UnwrapContext
func SetupContext() {
	defaultFunc := log.UnwrapContext
	log.UnwrapContext = func(ctx context.Context) context.Context {
		if c, ok := ctx.(*gin.Context); ok {
			return c.Request.Context()
		}

		return defaultFunc(ctx)
	}
}
