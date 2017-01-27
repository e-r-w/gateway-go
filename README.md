# gateway-go

A wrapper around [sparta](http://gosparta.io/) to provide a [gin](https://gin-gonic.github.io/gin/) or [sinatra](http://www.sinatrarb.com/) like interface for deploying serverless apps written in golang.

## Usage

```go
//main.go
import (
	"github.com/Sirupsen/logrus"
	sparta "github.com/mweagle/Sparta"
	gateway "github.com/e-r-w/gateway-go"
)

func main() {

	app := gateway.NewGateway(		
		"testing-stage", // Stage name
		"my-new-api", // API name
		"my cool new api", // Description
	)

	app.Get("/hello-world", func (ctx *gateway.Context, logger *logrus.Logger) {
		ctx.String("Hello World!")
	})

	app.Post("/hello-world", func (ctx *gateway.Context, logger *logrus.Logger) {

		resp, err := DoSomething()
		if err != nil {
			ctx.Error(err)
		}

		ctx.JSON(resp)

	}).
		WithRole(sparta.IAMRoleDefinition{
			// add role here for accessing AWS resources like DynamoDB, S3, RDS etc
		}).
		WithOptions(&sparta.LambdaFunctionOptions{
			// Add function options such as execution time or environment variables
		}).
		WithDecorator(func(/*...*/){
			// Decorator function, see http://gosparta.io/docs/dynamic_infrastructure/#template-decorators
		}).
		WithAuthorization(
			// Accepts gateway.AwsIam or gateway.None (defaults to None)
		)

	app.CORSEnabled = true // CORS is disabled by default

	app.
		Bootstrap().
		Start()

}
```

Then run & deploy just like you would with a regular sparta app!

For comparison, a Gin app:
```go
import (
	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()

	app.GET("/hello-world", func (c *gin.Context) {
		c.String(http.StatusOK, "Hello World!")
	})

	app.POST("/hello-world", func (c *gin.Context) {

		resp, err := DoSomething()
		if err != nil {
			ctx.Error(err)
		}

		c.JSON(http.StatusOK, resp)

	}).

	app.Run(":8080")
}
```

The aim of this project is to attain some what of a parity between Gin to make it easy to migrate to serverless apps
