package gateway

import (
	"encoding/json"
	"fmt"
	"net/http"

	sparta "github.com/mweagle/Sparta"
)

// Context ...
type Context struct {
	RawEvent       *json.RawMessage
	LambdaContext  *sparta.LambdaContext
	ResponseWriter http.ResponseWriter
}

// JSON ...
func (ctx *Context) JSON(object interface{}) {
	ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(ctx.ResponseWriter).Encode(object); err != nil {
		http.Error(ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
	}
}

// String ...
func (ctx *Context) String(object interface{}) {
	fmt.Fprint(ctx.ResponseWriter, object)
}

// Event ...
func (ctx *Context) Event() (sparta.APIGatewayLambdaJSONEvent, error) {
	var lambdaEvent sparta.APIGatewayLambdaJSONEvent
	if err := json.Unmarshal([]byte(*ctx.RawEvent), &lambdaEvent); err != nil {
		return lambdaEvent, err
	}
	return lambdaEvent, nil
}
