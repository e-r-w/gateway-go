package gateway

import (
	"encoding/json"
	"fmt"
	"net/http"

	sparta "github.com/mweagle/Sparta"
)

// Context ...
type Context struct {
	Request        *json.RawMessage
	LambdaContext  *sparta.LambdaContext
	ResponseWriter http.ResponseWriter
}

// JSON ...
func (ctx *Context) JSON(object interface{}) {
	json.NewEncoder(ctx.ResponseWriter).Encode(object)
}

// String ...
func (ctx *Context) String(object interface{}) {
	fmt.Fprint(ctx.ResponseWriter, object)
}
