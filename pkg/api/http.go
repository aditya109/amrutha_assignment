package api

import (
	"fmt"
	"github.com/aditya109/amrutha_assignment/pkg/middlewares"
	"github.com/aditya109/amrutha_assignment/pkg/recovery"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

func AcquireHttpServer(module Module, serviceTag string) (*http.Server, string) {
	var srv *http.Server

	PORT := viper.GetString("SERVER_PORT")
	engine := gin.New()
	engine.Use(middlewares.JSONLogMiddleware())
	engine.Use(gin.CustomRecovery(recovery.ApiPanicRecovery))
	engine.Use(gin.Recovery())
	engine.Use(middlewares.BindTraceIdToRequestHeaderMiddleware(serviceTag))

	apiRouter := engine.Group("")
	var c = construct{
		module:    module,
		apiRouter: apiRouter,
	}
	c.loadApis()

	srv = &http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: engine,
	}

	return srv, PORT
}

type construct struct {
	module    Module
	routes    []Route
	apiRouter *gin.RouterGroup
}

func (c construct) loadApis() []Route {
	var routes []Route
	module := c.module
	router := c.apiRouter.Group(fmt.Sprintf("/%s/api/%s", module.Module, module.ApiVersion))
	for _, route := range module.Routes {

		method := route.Method
		var handlers = make([]gin.HandlerFunc, 0)

		handlers = append(handlers, route.PreMiddlewares...)

		handlers = append(handlers, route.Controller)

		handlers = append(handlers, recovery.Responder)
		handlers = append(handlers, route.PostMiddlewares...)

		RegisterRoute(router, HttpRoute{
			Path:     route.Path,
			Method:   method,
			Handlers: handlers,
		})
	}

	return routes
}

func RegisterRoute(router *gin.RouterGroup, route HttpRoute) {
	switch route.Method {
	case http.MethodGet:
		{
			router.GET(route.Path, route.Handlers...)
			break
		}
	case http.MethodPost:
		{
			router.POST(route.Path, route.Handlers...)
			break
		}
	case http.MethodPut:
		{
			router.PUT(route.Path, route.Handlers...)
			break
		}
	default:
		{
			router.GET(route.Path, route.Handlers...)
			break
		}
	}
}
