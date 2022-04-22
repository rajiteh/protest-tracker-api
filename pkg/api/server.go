package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pangpanglabs/echoswagger/v2"
	"github.com/reliefeffortslk/protest-tracker-api/pkg/ingestor"
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

	type IngestResponse struct {
		Count int64 `json:"count"`
	}

	se.POST("/ingest", func(c echo.Context) error {
		ingestToken := strings.TrimSpace(os.Getenv("HTTP_INGEST_TOKEN_SECRET"))
		if len(ingestToken) == 0 {
			return fmt.Errorf("HTTP_INGEST_TOKEN_SECRET was not set")
		}
		count, err := ingestor.IngestFromAll()
		if err != nil {
			return err
		}
		resp := IngestResponse{
			Count: count,
		}

		return c.JSON(http.StatusOK, resp)
	})
	ProtestController{}.Init(se.Group("protest", "/protest"))
	return se
}
