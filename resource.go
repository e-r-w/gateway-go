package gateway

import (
	"encoding/json"
	"net/http"

	"github.com/Sirupsen/logrus"
	sparta "github.com/mweagle/Sparta"
)

// Resource ...
type Resource struct {
	Function       func(event *json.RawMessage, context *sparta.LambdaContext, w http.ResponseWriter, logger *logrus.Logger)
	RoleDefinition sparta.IAMRoleDefinition
	Method         string
	Route          string
	Decorator      sparta.TemplateDecorator
	Options        *sparta.LambdaFunctionOptions
}

// WithRole ...
func (r *Resource) WithRole(roleDef sparta.IAMRoleDefinition) *Resource {
	r.RoleDefinition = roleDef
	return r
}

// WithDecorator ...
func (r *Resource) WithDecorator(template sparta.TemplateDecorator) *Resource {
	r.Decorator = template
	return r
}

// WithOptions ...
func (r *Resource) WithOptions(funcOpts *sparta.LambdaFunctionOptions) *Resource {
	r.Options = funcOpts
	return r
}
