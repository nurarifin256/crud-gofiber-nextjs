package appcontext

// convertAppRoutes converts types.RouteDefinition to routes.Definition
// Optimized to pre-allocate memory and reduce allocations
func ConvertAppRoutes(appRoutes []RouteDefinition) []RouteDefinition {
	if len(appRoutes) == 0 {
		return nil
	}
	// Pre-allocate slice with exact capacity
	routes := make([]RouteDefinition, len(appRoutes))
	// Use direct assignment instead of struct literal for better performance
	for i := range appRoutes {
		appRoute := &appRoutes[i]
		route := &routes[i]

		route.Method = appRoute.Method
		route.Path = appRoute.Path
		route.Action = appRoute.Action
		route.AuthRequired = appRoute.AuthRequired
		route.Middleware = appRoute.Middleware
		route.AuthMiddleware = appRoute.AuthMiddleware
	}
	return routes
}
