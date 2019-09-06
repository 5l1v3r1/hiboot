//go:generate statik -src=./dist

package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/spec"
	"github.com/gorilla/handlers"
	_ "hidevops.io/hiboot/examples/web/swagger/http-server/statik"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/app/web/context"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/starter/actuator"
	"hidevops.io/hiboot/pkg/starter/logging"
	"net/http"
	"path"
)

type controller struct {
	at.RestController
	at.RequestMapping `value:"/api/greeting-server"`
}

func init() {
	app.Register(newController)
}

func newController() *controller {
	return &controller{}
}

func (c controller) loadDoc() (retVal []byte, err error) {

	swgSpec := &spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			ID:      "",
			Swagger: "2.0",
			Info: &spec.Info{
				InfoProps: spec.InfoProps{
					Description:    "",
					Title:          "HiBoot Demo Application - Greeting Server",
					TermsOfService: "",
					Contact:        nil,
					License:        nil,
					Version:        "1.0.0",
				},
			},
			Schemes:  []string{"http"},
			Host:     "localhost:8080",
			BasePath: "/api/greeting-server",
			Paths: &spec.Paths{
				Paths: map[string]spec.PathItem{
					"/hello": spec.PathItem{
						PathItemProps: spec.PathItemProps{
							Get: &spec.Operation{
								OperationProps: spec.OperationProps{
									ID: "getGreeting",
									Produces: []string{
										"text/plain",
									},
									Parameters: []spec.Parameter{
										{
											SimpleSchema: spec.SimpleSchema{
												Type: "string",
											},
											ParamProps: spec.ParamProps{
												Description: "defaults to World if not given",
												Name:        "name",
												In:          "query",
												Required:    false,
											},
										},
									},

									Responses: &spec.Responses{
										VendorExtensible: spec.VendorExtensible{
											Extensions: nil,
										},
										ResponsesProps: spec.ResponsesProps{
											Default: nil,
											StatusCodeResponses: map[int]spec.Response{
												200: {
													Refable: spec.Refable{},
													ResponseProps: spec.ResponseProps{
														Description: "returns a greeting",
														Schema: &spec.Schema{
															VendorExtensible: spec.VendorExtensible{},
															SchemaProps: spec.SchemaProps{
																Type:        spec.StringOrArray{"string"},
																Description: "contains the actual greeting as plain text",
															},
														},
													},
													VendorExtensible: spec.VendorExtensible{},
												},
												404: {
													Refable: spec.Refable{},
													ResponseProps: spec.ResponseProps{
														Description: "Resource is not found",
														Schema: &spec.Schema{
															VendorExtensible: spec.VendorExtensible{},
															SchemaProps: spec.SchemaProps{
																Type:        spec.StringOrArray{"string"},
																Description: "Report 'not found' error message",
															},
														},
													},
													VendorExtensible: spec.VendorExtensible{},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	retVal, err = json.MarshalIndent(swgSpec, "", "  ")
	if err != nil {
		return
	}
	return
}

func (c *controller) serve(ctx context.Context, docsPath string) {
	b, err := c.loadDoc()
	if err != nil {
		return
	}
	basePath := c.RequestMapping.Value

	handler := middleware.Redoc(middleware.RedocOpts{
		BasePath: basePath,
		SpecURL:  path.Join(basePath, "swagger.json"),
		Path:     docsPath,
	}, http.NotFoundHandler())

	visit := fmt.Sprintf("http://%s%s", ctx.Host(), ctx.Path())

	log.Debugf("visit: %v", visit)

	handler = handlers.CORS()(middleware.Spec(basePath, b, handler))

	ctx.WrapHandler(handler)
}

// UI serve static resource via context StaticResource method
func (c *controller) Swagger(at struct{ at.GetMapping `value:"/swagger.json"` }) (response string) {
	b, err := c.loadDoc()
	if err != nil {
		return
	}
	response = string(b)
	return
}

// UI serve static resource via context StaticResource method
func (c *controller) SwaggerUI(at struct{ at.GetMapping `value:"/swagger-ui"` }, ctx context.Context) {
	c.serve(ctx, at.GetMapping.Value)
	return
}

type HelloQueryParam struct{
	at.RequestParams
	Name string
}

// Hello
func (c *controller) Hello(at struct{ at.GetMapping `value:"/hello"` }, request *HelloQueryParam) (response string) {

	response = "Hello, " + request.Name

	return
}

//run http://localhost:8080/api/greeting-server/swagger-ui to open swagger ui
func main() {
	web.NewApplication(newController).
		SetProperty(app.ProfilesInclude, actuator.Profile, logging.Profile).
		Run()
}
