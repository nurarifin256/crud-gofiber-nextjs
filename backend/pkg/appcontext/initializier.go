package appcontext

// var RouteToRegistry []RouteDefinition

var RouteRegistries []RouteDefinition

// RouteInitializer handles the initialization of all route modules
// This pattern provides a clean separation of concerns for route management
type RouteInitializers struct{}

// NewRouteInitializer creates a new route initializer instance
func NewRouteInitializer() *RouteInitializers {
	return &RouteInitializers{}
}

// RegisterRouteV1 registers multiple route definitions to the registry
func RegisterRoute(defs []RouteDefinition) {
	RouteRegistries = append(RouteRegistries, defs...)
}

// // InitializeAllRoutes initializes all route modules in a centralized manner
// // This method automatically detects and registers all routes from app package
// func (ri *RouteInitializers) InitializeAllRoutes(registerFunc func([]RouteDefinition)) {
// 	// Automatically collect and register all routes from app package
// 	// All route modules are auto-registered via init() functions
// 	allRoutes := InitializeAllRoutes()
// 	registerFunc(ConvertAppRoutes(allRoutes))
// }

// InitializeAllRoutes initializes all route modules in a centralized manner
// This method automatically detects and registers all routes from app package
func (ri *RouteInitializers) InitializeAllRoutes() {
	// Automatically collect and register all routes from app package
	// All route modules are auto-registered via init() functions
	allRoutes := InitializeAllRoutes()
	RegisterRoute(ConvertAppRoutes(allRoutes))
}
