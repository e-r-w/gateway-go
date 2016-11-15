package gateway

import (
	"encoding/json"
	"net/http"

	"github.com/Sirupsen/logrus"
	sparta "github.com/mweagle/Sparta"
)

type Resource struct {
	Function       func(event *json.RawMessage, context *sparta.LambdaContext, w http.ResponseWriter, logger *logrus.Logger)
	RoleDefinition sparta.IAMRoleDefinition
	Method         string
	Route          string
}

func (r *Resource) WithRole(roleDef sparta.IAMRoleDefinition) {
	r.RoleDefinition = roleDef
}
