package gateway

import (
	"encoding/json"
	"net/http"

	"github.com/Sirupsen/logrus"
	sparta "github.com/mweagle/Sparta"
)

// Gateway ...
type Gateway struct {
	Resources           []*Resource
	Stage               *sparta.Stage
	API                 *sparta.API
	Lambdas             []*sparta.LambdaAWSInfo
	APIName             string
	Description         string
	CORSEnabled         bool
	APIGatewayResources []*sparta.Resource
	APIGatewayMethods   []*sparta.Method
}

// Bootstrap ...
func (g *Gateway) Bootstrap() *Gateway {

	g.API.CORSEnabled = g.CORSEnabled

	for _, resource := range g.Resources {
		lambda := sparta.NewLambda(resource.RoleDefinition, resource.Function, resource.Options)
		if resource.Decorator != nil {
			lambda.Decorator = resource.Decorator
		}
		g.Lambdas = append(g.Lambdas, lambda)

		apiGatewayResource, _ := g.API.NewResource(resource.Route, lambda)
		g.APIGatewayResources = append(g.APIGatewayResources, apiGatewayResource)

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
		g.Lambdas,
		g.API,
		nil)
}

func (g *Gateway) DynamoPassThrough() *Resource {



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
		Route:          route,
		Method:         method,
		RoleDefinition: sparta.IAMRoleDefinition{},
		Authorization:  None,
	}

	wrapped := func(event *json.RawMessage, context *sparta.LambdaContext, w http.ResponseWriter, logger *logrus.Logger) {
		wrappedCtx := Context{
			RawEvent:       event,
			LambdaContext:  context,
			ResponseWriter: w,
		}
		for _, m := range resource.Middlewares {
			m(&wrappedCtx, logger, func(){

			})
		}
		handler(&wrappedCtx, logger)
	}

	resource.Function = wrapped

	g.Resources = append(g.Resources, &resource)

	return &resource

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
