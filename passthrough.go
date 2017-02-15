package gateway

import (
	"github.com/Sirupsen/logrus"
	gocf "github.com/crewjam/go-cloudformation"
	spartaCF "github.com/mweagle/Sparta/aws/cloudformation"
	"github.com/mweagle/Sparta"
	"strings"
	"fmt"
)

type PassThrough struct{
	apiName string
}


func (p *PassThrough) Apply(api *sparta.API) sparta.TemplateDecorator {
	return func(serviceName string,
	lambdaResourceName string,
	lambdaResource gocf.LambdaFunction,
	resourceMetadata map[string]interface{},
	S3Bucket string,
	S3Key string,
	buildID string,
	template *gocf.Template,
	context map[string]interface{},
	logger *logrus.Logger) error {
		apiGatewayResourceNameForPath := func(fullPath string) string {
			pathParts := strings.Split(fullPath, "/")
			return spartaCF.CloudFormationResourceName("%sResource", pathParts[0], fullPath)
		}
		apiGatewayResName := spartaCF.CloudFormationResourceName("APIGateway", p.apiName)

		// Create an API gateway entry
		apiGatewayRes := &gocf.ApiGatewayRestApi{
			Description:    gocf.String(api.Description),
			FailOnWarnings: gocf.Bool(false),
			Name:           gocf.String(p.apiName),
		}
		if "" != api.CloneFrom {
			apiGatewayRes.CloneFrom = gocf.String(api.CloneFrom)
		}
		if "" == api.Description {
			apiGatewayRes.Description = gocf.String(fmt.Sprintf("%s RestApi", serviceName))
		} else {
			apiGatewayRes.Description = gocf.String(api.Description)
		}

		template.AddResource(apiGatewayResName, apiGatewayRes)
		apiGatewayRestAPIID := gocf.Ref(apiGatewayResName)

		// List of all the method resources we're creating s.t. the
		// deployment can DependOn them
		optionsMethodPathMap := make(map[string]bool)
		var apiMethodCloudFormationResources []string
		for eachResourceMethodKey, eachResourceDef := range api.resources {
			// First walk all the user resources and create intermediate paths
			// to repreesent all the resources
			var parentResource *gocf.StringExpr
			pathParts := strings.Split(strings.TrimLeft(eachResourceDef.pathPart, "/"), "/")
			pathAccumulator := []string{"/"}
			for index, eachPathPart := range pathParts {
				pathAccumulator = append(pathAccumulator, eachPathPart)
				resourcePathName := apiGatewayResourceNameForPath(strings.Join(pathAccumulator, "/"))
				if _, exists := template.Resources[resourcePathName]; !exists {
					cfResource := &gocf.ApiGatewayResource{
						RestApiId: apiGatewayRestAPIID.String(),
						PathPart:  gocf.String(eachPathPart),
					}
					if index <= 0 {
						cfResource.ParentId = gocf.GetAtt(apiGatewayResName, "RootResourceId")
					} else {
						cfResource.ParentId = parentResource
					}
					template.AddResource(resourcePathName, cfResource)
				}
				parentResource = gocf.Ref(resourcePathName).String()
			}

			// Add the lambda permission
			apiGatewayPermissionResourceName := spartaCF.CloudFormationResourceName("APIGatewayLambdaPerm", eachResourceMethodKey)
			lambdaInvokePermission := &gocf.LambdaPermission{
				Action:       gocf.String("lambda:InvokeFunction"),
				FunctionName: gocf.GetAtt(eachResourceDef.parentLambda.logicalName(), "Arn"),
				Principal:    gocf.String(APIGatewayPrincipal),
			}
			template.AddResource(apiGatewayPermissionResourceName, lambdaInvokePermission)

			// BEGIN CORS - OPTIONS verb
			// CORS is API global, but it's possible that there are multiple different lambda functions
			// that are handling the same HTTP resource. In this case, track whether we've already created an
			// OPTIONS entry for this path and only append iff this is the first time through
			if api.CORSEnabled {
				methodResourceName := spartaCF.CloudFormationResourceName(fmt.Sprintf("%s-OPTIONS", eachResourceDef.pathPart), eachResourceDef.pathPart)
				_, resourceExists := optionsMethodPathMap[methodResourceName]
				if !resourceExists {
					template.AddResource(methodResourceName, corsOptionsGatewayMethod(apiGatewayRestAPIID, parentResource))
					apiMethodCloudFormationResources = append(apiMethodCloudFormationResources, methodResourceName)
					optionsMethodPathMap[methodResourceName] = true
				}
			}
			// END CORS - OPTIONS verb

			// BEGIN - user defined verbs
			for eachMethodName, eachMethodDef := range eachResourceDef.Methods {

				apiGatewayMethod := &gocf.ApiGatewayMethod{
					HttpMethod:        gocf.String(eachMethodName),
					AuthorizationType: gocf.String("NONE"),
					ResourceId:        parentResource.String(),
					RestApiId:         apiGatewayRestAPIID.String(),
					Integration:       eachMethodDef.CustomIntegration,
				}
				if nil == apiGatewayMethod.Integration {
					apiGatewayMethod.Integration = &gocf.APIGatewayMethodIntegration{
						IntegrationHttpMethod: gocf.String("POST"),
						Type:             gocf.String("AWS"),
						RequestTemplates: defaultRequestTemplates(),
						Uri: gocf.Join("",
							gocf.String("arn:aws:apigateway:"),
							gocf.Ref("AWS::Region"),
							gocf.String(":lambda:path/2015-03-31/functions/"),
							gocf.GetAtt(eachResourceDef.parentLambda.logicalName(), "Arn"),
							gocf.String("/invocations")),
					}
				}
				if len(eachMethodDef.Parameters) != 0 {
					requestParams := make(map[string]string, 0)
					for eachKey, eachBool := range eachMethodDef.Parameters {
						requestParams[eachKey] = fmt.Sprintf("%t", eachBool)
					}
					apiGatewayMethod.RequestParameters = requestParams
				}

				// Add the integration response RegExps
				apiGatewayMethod.Integration.IntegrationResponses = integrationResponses(eachMethodDef.Integration.Responses,
					api.CORSEnabled)

				// Add outbound method responses
				apiGatewayMethod.MethodResponses = methodResponses(eachMethodDef.Responses,
					api.CORSEnabled)

				prefix := fmt.Sprintf("%s%s", eachMethodDef.httpMethod, eachResourceMethodKey)
				methodResourceName := CloudFormationResourceName(prefix, eachResourceMethodKey, serviceName)
				res := template.AddResource(methodResourceName, apiGatewayMethod)
				res.DependsOn = append(res.DependsOn, apiGatewayPermissionResourceName)
				apiMethodCloudFormationResources = append(apiMethodCloudFormationResources, methodResourceName)
			}
		}
		return nil
	}
}