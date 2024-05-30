package api

import "github.com/gin-gonic/gin"

type Route struct {
	Path            string
	Method          string
	PreMiddlewares  []gin.HandlerFunc
	Controller      func(*gin.Context)
	PostMiddlewares []gin.HandlerFunc
}

type ApiModule struct {
	Module     string
	ApiVersion string
	Routes     []Route
}

type HttpRoute struct {
	Module   string
	Path     string
	Method   string
	Handlers []gin.HandlerFunc
}
