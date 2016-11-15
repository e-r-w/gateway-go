package gateway

import (
	"encoding/json"
	"net/http"

	"github.com/Sirupsen/logrus"
	sparta "github.com/mweagle/Sparta"
)

type Gateway struct {
	Resources []*Resource
}

func (g Gateway) Bootstrap(stage, apiName string) {

	var allTheLambdas []*sparta.LambdaAWSInfo

	for _, resource := range g.Resources {
		lambda := sparta.NewLambda(resource.RoleDefinition, resource.Function, nil)
		allTheLambdas = append(allTheLambdas, lambda)
		stage := sparta.NewStage(stage)
		api := sparta.NewAPIGateway(apiName, stage)
		apiGatewayResource, _ := api.NewResource(resource.Route, lambda)
		apiGatewayResource.NewMethod(resource.Method, http.StatusOK)
	}

	sparta.Main(apiName,
		"Simple Sparta application that demonstrates core functionality",
		allTheLambdas,
		nil,
		nil)
}

func (g Gateway) Get(route string, handler func(ctx *Context, logger *logrus.Logger)) *Resource {
	return g.Route("GET", route, handler)
}

func (g Gateway) Post(route string, handler func(ctx *Context, logger *logrus.Logger)) *Resource {
	return g.Route("POST", route, handler)
}

func (g Gateway) Route(method string, route string, handler func(ctx *Context, logger *logrus.Logger)) *Resource {

	wrapped := func(event *json.RawMessage, context *sparta.LambdaContext, w http.ResponseWriter, logger *logrus.Logger) {
		wrappedCtx := Context{
			Request:        event,
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
