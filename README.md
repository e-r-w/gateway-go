# gateway-go

A wrapper around [sparta](http://gosparta.io/) to provide a [gin](https://gin-gonic.github.io/gin/) or [sinatra](http://www.sinatrarb.com/) like interface for deploying serverless apps written in golang.

## Usage

```go
//main.go
import (
	"github.com/e-r-w/gateway-go"
  	"github.com/Sirupsen/logrus"
	sparta "github.com/mweagle/Sparta"
)

func main() {
	app := gateway.Gateway{}

	app.Get("/hello-world", func (c *gateway.Context, logger *logrus.Logger) {
		c.String("Hello World!")
	})

	app.Post("/hello-world", func (c *gateway.Context, logger *logrus.Logger) {
		c.JSON(map[string]interface{}{
			"foo": "bar",
		})
	}).
		WithRole(sparta.IAMRoleDefinition{
			// add role here for accessing AWS resources like DynamoDB, S3, RDS etc
		})

	app.Bootstrap("testing-stage", "my-new-api", "my cool new api")
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
		c.JSON(http.StatusOK, map[string]interface{}{
			"foo": "bar",
		})
	}).

	app.Run(":8080")
}
```

The aim of this project is to attain some what of a parity between Gin to make it easy to migrate to serverless apps
