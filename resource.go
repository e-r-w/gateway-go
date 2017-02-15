package gateway

import (
	"github.com/Sirupsen/logrus"
	sparta "github.com/mweagle/Sparta"
)

// None ..
const (
	None   = "NONE"
	AwsIam = "AWS_IAM"
)

// MethodDecorator ...
type MethodDecorator func(*sparta.Method) *sparta.Method

// Resource ...
type Resource struct {
	Function        func(context *Context, logger *logrus.Logger)
	Method          string
	Route           string
	Authorization   string
	MethodDecorator MethodDecorator
}

// WithAuthorization ...
func (r *Resource) WithAuthorization(authorization string) *Resource {
	r.Authorization = authorization
	return r
}

// WithMethodDecorator ...
func (r *Resource) WithMethodDecorator(decorator MethodDecorator) *Resource {
	r.MethodDecorator = decorator
	return r
}
