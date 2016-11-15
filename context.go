package gateway

import (
	"encoding/json"
	"fmt"
	"net/http"

	sparta "github.com/mweagle/Sparta"
)

type Context struct {
	Request        *json.RawMessage
	LambdaContext  *sparta.LambdaContext
	ResponseWriter http.ResponseWriter
}

func (ctx *Context) JSON(object interface{}) {
	json.NewEncoder(ctx.ResponseWriter).Encode(object)
}

func (ctx *Context) String(object interface{}) {
	fmt.Fprint(ctx.ResponseWriter, object)
}
