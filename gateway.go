package gateway

import (
	"encoding/json"
	"net/http"

	"github.com/Sirupsen/logrus"
	sparta "github.com/mweagle/Sparta"
)

// Gateway ...
type Gateway struct {
	Resources []*Resource
}

// Bootstrap ...
func (g *Gateway) Bootstrap(stageName, apiName, description string) {

	var allTheLambdas []*sparta.LambdaAWSInfo
	apiStage := sparta.NewStage(stageName)
	api := sparta.NewAPIGateway(apiName, apiStage)

	for _, resource := range g.Resources {
		lambda := sparta.NewLambda(resource.RoleDefinition, resource.Function, resource.Options)
		if resource.Decorator != nil {
			lambda.Decorator = resource.Decorator
		}
		allTheLambdas = append(allTheLambdas, lambda)
		apiGatewayResource, _ := api.NewResource(resource.Route, lambda)
		apiGatewayResource.NewMethod(resource.Method, http.StatusOK)
	}

	sparta.Main(apiName,
		description,
		allTheLambdas,
		api,
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

	wrapped := func(event *json.RawMessage, context *sparta.LambdaContext, w http.ResponseWriter, logger *logrus.Logger) {
		wrappedCtx := Context{
			RawEvent:       event,
			LambdaContext:  context,
			ResponseWriter: w,
		}
		handler(&wrappedCtx, logger)
	}

	resource := Resource{
		Route:          route,
		Method:         method,
		RoleDefinition: sparta.IAMRoleDefinition{},
		Function:       wrapped,
	}

	g.Resources = append(g.Resources, &resource)

	return &resource

}

// NewGateway ...
func NewGateway() *Gateway {
	return &Gateway{}
}
