package appcontext

import (
	"fmt"
	"sync"

	middleware "simple-crud/api/middlewares"

	"github.com/gofiber/fiber/v2"
)

// routeRegistry holds all registered route initializers
var (
	routeRegistry = make(map[string]RouteInitializer)
	registryMutex sync.RWMutex

	// Cache for initialized routes to avoid repeated initialization
	routeCache       []RouteDefinition
	routeCacheMutex  sync.RWMutex
	cacheInitialized bool
)

// RegisterRouteInitializer registers a route initializer with a given name
func RegisterRouteInitializer(name string, initializer RouteInitializer) {
	registryMutex.Lock()
	defer registryMutex.Unlock()
	routeRegistry[name] = initializer

	// Invalidate cache when new routes are registered
	routeCacheMutex.Lock()
	cacheInitialized = false
	routeCacheMutex.Unlock()
}

// registerRoutesToApp registers all collected routes to the Fiber application
func RegisterRoutesToApp(app *fiber.App, url string) {
	api := app.Group(url)
	routes := RouteRegistries
	mv := middleware.InitMiddleWare()
	fmt.Println("Total routes V1:", len(routes))
	for _, def := range routes {
		if def.AuthRequired {
			// fmt.Println(def.Method, def.Path, def.AuthRequired)
			handlers := append([]fiber.Handler{mv.AuthMiddleware.Auth}, def.Middleware...)
			handlers = append(handlers, def.Action)
			api.Add(def.Method, def.Path, handlers...)
		} else {
			// fmt.Println(def.Method, def.Path, def.AuthRequired)
			handlers := append(def.Middleware, def.Action)
			api.Add(def.Method, def.Path, handlers...)
		}
	}
}

// GetAllRouteInitializers returns all registered route initializers
func GetAllRouteInitializers() map[string]RouteInitializer {
	registryMutex.RLock()
	defer registryMutex.RUnlock()

	// Create a copy to avoid race conditions
	result := make(map[string]RouteInitializer, len(routeRegistry))
	for k, v := range routeRegistry {
		result[k] = v
	}
	return result
}

// InitializeAllRoutes calls all registered route initializers and returns all routes
// This function is optimized with caching to avoid repeated initialization
func InitializeAllRoutes() []RouteDefinition {
	// Check cache first
	routeCacheMutex.RLock()
	if cacheInitialized {
		defer routeCacheMutex.RUnlock()
		// Return a copy to avoid external modifications
		result := make([]RouteDefinition, len(routeCache))
		copy(result, routeCache)
		return result
	}
	routeCacheMutex.RUnlock()

	// Initialize routes and cache them
	routeCacheMutex.Lock()
	defer routeCacheMutex.Unlock()

	// Double-check pattern to avoid race conditions
	if cacheInitialized {
		result := make([]RouteDefinition, len(routeCache))
		copy(result, routeCache)
		return result
	}

	// Pre-allocate slice with estimated capacity
	initializers := GetAllRouteInitializers()
	estimatedCapacity := len(initializers) * 5 // Assume average 5 routes per module
	allRoutes := make([]RouteDefinition, 0, estimatedCapacity)

	for _, initializer := range initializers {
		routes := initializer.Initialize()
		allRoutes = append(allRoutes, routes...)
	}

	// Cache the result
	routeCache = make([]RouteDefinition, len(allRoutes))
	copy(routeCache, allRoutes)
	cacheInitialized = true

	return allRoutes
}

// ClearCache clears the route cache (useful for testing or dynamic reloading)
func ClearCache() {
	routeCacheMutex.Lock()
	defer routeCacheMutex.Unlock()
	cacheInitialized = false
	routeCache = nil
}
