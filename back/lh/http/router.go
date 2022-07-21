package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Route struct {
	// Name is the name of this Route.
	Name string
	// Method is the string for the HTTP method. ex) GET, POST etc..
	Method string
	// Pattern is the pattern of the URI.
	Pattern string
	// HandlerFunc is the handler function of this route.
	HandlerFunc gin.HandlerFunc
}

// Routes is the list of the generated Route.
type Routes []Route

func AddService(engine *gin.Engine) *gin.RouterGroup {
	group := engine.Group("/lh/")

	for _, route := range routes {
		switch route.Method {
		case "GET":
			group.GET(route.Pattern, route.HandlerFunc)
		case "POST":
			group.POST(route.Pattern, route.HandlerFunc)
		case "PUT":
			group.PUT(route.Pattern, route.HandlerFunc)
		}
	}
	return group
}

// Index is the index handler.
func Index(c *gin.Context) {
	c.String(http.StatusOK, "Hello World!")
}

var routes = Routes{
	{
		"Index",
		http.MethodGet,
		"/",
		Index,
	},

	{
		"add user right/contact",
		http.MethodPost,
		"/init/:name/:id",
		initial,
	},
	{
		"right doc validation",
		http.MethodGet,
		"/get/right/:name/:id",
		getRight,
	},

	{
		"add document",
		http.MethodPost,
		"/right/doc/add/:name/:id/:docname",
		uploadRightdoc,
	},
	{
		"right doc validation",
		http.MethodPut,
		"/right/state/change/:name/:id/:statevalue",
		changeRightstate,
	},
}
