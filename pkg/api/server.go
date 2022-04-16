package api

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pangpanglabs/echoswagger/v2"
)

type APIService struct {
}

func (bs *APIService) Serve(_ context.Context) error {
	e := initServer().Echo()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.HideBanner = true

	return e.Start(":1323")
}

func initServer() echoswagger.ApiRoot {
	e := echo.New()

	se := echoswagger.New(e, "docs/", &echoswagger.Info{
		Title:       "SL-Protest-Tracker-API",
		Description: "This API will provide an interface for latest protests tracked by various datasources",
		Version:     "1.0.0",
	})

	se.AddSecurityOAuth2("petstore_auth", "", echoswagger.OAuth2FlowImplicit,
		"http://petstore.swagger.io/oauth/dialog", "", map[string]string{
			"write:pets": "modify pets in your account",
			"read:pets":  "read your pets",
		},
	).AddSecurityAPIKey("api_key", "", echoswagger.SecurityInHeader)

	se.SetExternalDocs("Find out more about Swagger", "http://swagger.io").
		SetResponseContentType("application/xml", "application/json").
		SetUI(echoswagger.UISetting{DetachSpec: true, HideTop: true}).
		SetScheme("http")

	// Datasources
	// Protests

	ProtestController{}.Init(se.Group("protest", "/protest"))
	return se
}
