package service

type HTTPService interface {
	Serve() error
	Shutdown() error
	RegisterUtilityRoutes()
	RegisterUserRoutes()
	RegisterNotFoundRoutes()
}