package api

import (
	"context"

	"github.com/labstack/echo"
	"github.com/pangpanglabs/echoswagger"
)

type APIService struct {
}

func (bs *APIService) Serve(_ context.Context) error {
	e := initServer().Echo()
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
		SetScheme("https", "http")

	// PetController{}.Init(se.Group("pet", "/pet"))
	// StoreController{}.Init(se.Group("store", "/store"))
	// UserController{}.Init(se.Group("user", "/user"))

	return se
}
