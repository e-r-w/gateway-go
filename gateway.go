package gateway

import (
	"encoding/json"
	"net/http"

	"fmt"

	"github.com/Sirupsen/logrus"
	sparta "github.com/mweagle/Sparta"
)

// Gateway ...
type Gateway struct {
	Resources           []*Resource
	Stage               *sparta.Stage
	API                 *sparta.API
	Lambda              *sparta.LambdaAWSInfo
	APIName             string
	Description         string
	CORSEnabled         bool
	APIGatewayResources []*sparta.Resource
	APIGatewayMethods   []*sparta.Method
	Options             *sparta.LambdaFunctionOptions
	RoleDefinition      sparta.IAMRoleDefinition
	Decorator           sparta.TemplateDecorator
	routeMap            map[string]*sparta.Resource
}

func (g *Gateway) createOrFindResource(route string) (*sparta.Resource, error) {
	for k, v := range g.routeMap {
		if k == route {
			return v, nil
		}
	}
	if g.routeMap == nil {
		g.routeMap = map[string]*sparta.Resource{}
	}
	apiGatewayResource, _ := g.API.NewResource(route, g.Lambda)
	g.APIGatewayResources = append(g.APIGatewayResources, apiGatewayResource)
	g.routeMap[route] = apiGatewayResource
	return apiGatewayResource, nil
}

// Bootstrap ...
func (g *Gateway) Bootstrap() *Gateway {

	g.API.CORSEnabled = g.CORSEnabled

	lambda := sparta.NewLambda(g.RoleDefinition, func(event *json.RawMessage, context *sparta.LambdaContext, w http.ResponseWriter, logger *logrus.Logger) {
		var lambdaEvent sparta.APIGatewayLambdaJSONEvent
		if err := json.Unmarshal([]byte(*event), &lambdaEvent); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// figure out the method
		method := lambdaEvent.Method
		// figure out the route
		route := lambdaEvent.Context.ResourcePath
		// invoke the right g.Resource function
		for _, resource := range g.Resources {
			if method == resource.Method && route == resource.Route {
				wrappedCtx := Context{
					Event:          &lambdaEvent,
					LambdaContext:  context,
					ResponseWriter: w,
				}
				resource.Function(&wrappedCtx, logger)
				return
			}
		}
		http.Error(w, fmt.Sprint("Unable to match route ", route, " with method ", method), http.StatusInternalServerError)
	}, g.Options)

	if g.Decorator != nil {
		lambda.Decorator = g.Decorator
	}

	g.Lambda = lambda

	for _, resource := range g.Resources {

		apiGatewayResource, _ := g.createOrFindResource(resource.Route)

		var method *sparta.Method
		if resource.Authorization == None {
			method, _ = apiGatewayResource.NewMethod(resource.Method, http.StatusOK)
		} else {
			method, _ = apiGatewayResource.NewAuthorizedMethod(resource.Method, resource.Authorization, http.StatusOK)
		}
		if resource.MethodDecorator != nil {
			method = resource.MethodDecorator(method)
		}
		g.APIGatewayMethods = append(g.APIGatewayMethods, method)
	}

	return g

}

// Start ...
func (g *Gateway) Start() {
	sparta.Main(g.APIName,
		g.Description,
		[]*sparta.LambdaAWSInfo{
			g.Lambda,
		},
		g.API,
		nil)
}

// Get ...
func (g *Gateway) Get(route string, handler func(ctx *Context, logger *logrus.Logger)) *Resource {
	return g.Route("GET", route, handler)
}

// Post ...
func (g *Gateway) Post(route string, handler func(ctx *Context, logger *logrus.Logger)) *Resource {
	return g.Route("POST", route, handler)
}

// Route ...
func (g *Gateway) Route(method string, route string, handler func(ctx *Context, logger *logrus.Logger)) *Resource {

	resource := Resource{
		Route:         route,
		Method:        method,
		Function:      handler,
		Authorization: None,
	}

	g.Resources = append(g.Resources, &resource)

	return &resource

}

// WithOptions ...
func (g *Gateway) WithOptions(funcOpts *sparta.LambdaFunctionOptions) *Gateway {
	g.Options = funcOpts
	return g
}

// WithRole ...
func (g *Gateway) WithRole(roleDef sparta.IAMRoleDefinition) *Gateway {
	g.RoleDefinition = roleDef
	return g
}

// WithDecorator ...
func (g *Gateway) WithDecorator(template sparta.TemplateDecorator) *Gateway {
	g.Decorator = template
	return g
}

// NewGateway ...
func NewGateway(stageName, apiName, description string) *Gateway {

	apiStage := sparta.NewStage(stageName)
	api := sparta.NewAPIGateway(apiName, apiStage)

	return &Gateway{
		Stage:       apiStage,
		APIName:     apiName,
		API:         api,
		Description: description,
	}

}
