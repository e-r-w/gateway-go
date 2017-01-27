package gateway

import (
	"encoding/json"
	"net/http"

	"github.com/Sirupsen/logrus"
	sparta "github.com/mweagle/Sparta"
)

// None ..
const (
	None   = "NONE"
	AwsIam = "AWS_IAM"
)

// Resource ...
type Resource struct {
	Function       func(event *json.RawMessage, context *sparta.LambdaContext, w http.ResponseWriter, logger *logrus.Logger)
	RoleDefinition sparta.IAMRoleDefinition
	Method         string
	Route          string
	Decorator      sparta.TemplateDecorator
	Options        *sparta.LambdaFunctionOptions
	Authorization  string
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

// WithAuthorization ...
func (r *Resource) WithAuthorization(authorization string) *Resource {
	r.Authorization = authorization
	return r
}
