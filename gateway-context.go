package gateway

import (
	"encoding/json"
	"fmt"
	"net/http"

	sparta "github.com/mweagle/Sparta"
)

type GatewayContext struct {
	Request        *json.RawMessage
	LambdaContext  *sparta.LambdaContext
	ResponseWriter http.ResponseWriter
}

func (ctx GatewayContext) JSON(object interface{}) {
	json.NewEncoder(ctx.ResponseWriter).Encode(object)
}

func (ctx GatewayContext) String(object interface{}) {
	fmt.Fprint(ctx.ResponseWriter, object)
}
