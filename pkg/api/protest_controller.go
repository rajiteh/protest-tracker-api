package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pangpanglabs/echoswagger/v2"
	"github.com/reliefeffortslk/protest-tracker-api/pkg/data"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type ProtestController struct{}

type ProtestsParam struct {
	DatasourceSlug string `query:"datasource" swagger:"optional,desc(datasource name to filter)"`
	AfterDate      string `query:"afterDate" swagger:"optional,desc(filters results after date: YYYY-mm-dd)"`
	BeforeDate     string `query:"beforeDate" swagger:"optional,desc(filters results before date: YYYY-mm-dd)"`
}

func (c ProtestController) Init(g echoswagger.ApiGroup) {
	g.SetDescription("Get Information about Protests in the tracker").
		SetExternalDocs("Find out more", "https://protestslk.com")

	// security := map[string][]string{
	// 	"petstore_auth": {"write:pets", "read:pets"},
	// }

	// type Category struct {
	// 	Id   int64  `json:"id"`
	// 	Name string `json:"name"`
	// }
	// type Tag struct {
	// 	Id   int64  `json:"id"`
	// 	Name string `json:"name"`
	// }
	// type Pet struct {
	// 	Id        int64    `json:"id"`
	// 	Category  Category `json:"category"`
	// 	Name      string   `json:"name" swagger:"required"`
	// 	PhotoUrls []string `json:"photoUrls" xml:"photoUrl" swagger:"required"`
	// 	Tags      []Tag    `json:"tags" xml:"tag"`
	// 	Status    string   `json:"status" swagger:"enum(available|pending|sold),desc(pet status in the store)"`
	// }
	// pet := Pet{Name: "doggie"}

	sampleProtest := ProtestResponse{
		DataSource: DataSourceResponse{
			Reference:   "2",
			Name:        "datasourceName",
			Attribution: "datasource attribution",
		},
		Lat:      69.321,
		Lng:      69.123,
		Location: "Colombo",
		Date:     "31-04-2022",
		Notes:    "some notes",
		Links: []string{
			"http://twitter.com/someTweet/",
			"http://instagram.com/someInsta/",
		},
		Size: "medium",
	}

	g.GET("/", c.ListProtests).
		AddParamQueryNested(&ProtestsParam{}).
		AddResponse(http.StatusOK, "succesful operation", &[]ProtestResponse{sampleProtest}, nil).
		SetOperationId("listAllProtests").
		SetDescription("Returns all available protests in the tracker").
		SetSummary("Lists all protests")

	g.GET("/:protestId", c.GetProtestById).
		AddParamPath(0, "protestId", "ID of protest").
		AddResponse(http.StatusOK, "successful operation", &sampleProtest, nil).
		SetOperationId("getProtestById").
		SetDescription("Returns a single protest").
		SetSummary("Find protest by ID")

	// type StatusParam struct {
	// 	Status []string `query:"status" swagger:"required,desc(Status values that need to be considered for filter),default(available),enum(available|pending|sold)"`
	// }
	// g.GET("/findByStatus", c.FindByStatus).
	// 	AddParamQueryNested(&StatusParam{}).
	// 	AddResponse(http.StatusOK, "successful operation", &[]Pet{pet}, nil).
	// 	AddResponse(http.StatusBadRequest, "Invalid status value", nil, nil).
	// 	SetOperationId("findPetsByStatus").
	// 	SetDescription("Multiple status values can be provided with comma separated strings").
	// 	SetSummary("Finds Pets by status").
	// 	SetSecurityWithScope(security)

	// g.GET("/findByTags", c.FindByTags).
	// 	AddParamQuery([]string{}, "tags", "Tags to filter by", true).
	// 	AddResponse(http.StatusOK, "successful operation", &[]Pet{pet}, nil).
	// 	AddResponse(http.StatusBadRequest, "Invalid tag value", nil, nil).
	// 	SetOperationId("findPetsByTags").
	// 	SetDeprecated().
	// 	SetDescription("Multiple tags can be provided with comma separated strings. Use         tag1, tag2, tag3 for testing.").
	// 	SetSummary("Finds Pets by tags").
	// 	SetSecurityWithScope(security)

	// g.GET("/{petId}", c.GetById).
	// 	AddParamPath(0, "petId", "ID of pet to return").
	// 	AddResponse(http.StatusOK, "successful operation", &pet, nil).
	// 	AddResponse(http.StatusBadRequest, "Invalid ID supplied", nil, nil).
	// 	AddResponse(http.StatusNotFound, "Pet not found", nil, nil).
	// 	SetOperationId("getPetById").
	// 	SetDescription("Returns a single pet").
	// 	SetSummary("Find pet by ID").
	// 	SetSecurity("api_key")

	// g.POST("/{petId}", c.CreateById).
	// 	AddParamPath(0, "petId", "ID of pet that needs to be updated").
	// 	AddParamForm("", "name", "Updated name of the pet", false).
	// 	AddParamForm("", "status", "Updated status of the pet", false).
	// 	AddResponse(http.StatusMethodNotAllowed, "Invalid input", nil, nil).
	// 	SetRequestContentType("application/x-www-form-urlencoded").
	// 	SetOperationId("updatePetWithForm").
	// 	SetSummary("Updates a pet in the store with form data").
	// 	SetSecurityWithScope(security)

	// g.DELETE("/{petId}", c.DeleteById).
	// 	AddParamHeader("", "api_key", "", false).
	// 	AddParamPath(int64(0), "petId", "Pet id to delete").
	// 	AddResponse(http.StatusBadRequest, "Invalid ID supplied", nil, nil).
	// 	AddResponse(http.StatusNotFound, "Pet not found", nil, nil).
	// 	SetOperationId("deletePet").
	// 	SetSummary("Deletes a pet").
	// 	SetSecurityWithScope(security)

	// type ApiResponse struct {
	// 	Code    int32  `json:"code"`
	// 	Type    string `json:"type"`
	// 	Message string `json:"message"`
	// }
	// g.POST("/{petId}/uploadImage", c.UploadImageById).
	// 	AddParamPath("", "petId", "ID of pet to update").
	// 	AddParamForm("", "additionalMetadata", "Additional data to pass to server", false).
	// 	AddParamFile("file", "file to upload", false).
	// 	AddResponse(http.StatusOK, "successful operation", &ApiResponse{}, nil).
	// 	SetRequestContentType("multipart/form-data").
	// 	SetResponseContentType("application/json").
	// 	SetOperationId("uploadFile").
	// 	SetSummary("uploads an image").
	// 	SetSecurityWithScope(security)
}

type DataSourceResponse struct {
	Reference     string `json:"reference"`
	Name          string
	Attribution   string
	LastUpdatedAt time.Time
}

type ProtestResponse struct {
	DataSource DataSourceResponse `json:"dataSource"`
	Lat        float64            `json:"lat"`
	Lng        float64            `json:"lng"`
	Location   string             `json:"location"`
	Date       string             `json:"date"`
	Notes      string             `json:"notes"`
	Links      []string           `json:"links"`
	Size       string             `json:"size"`
}

func protestToResponseObject(protest *data.Protest) ProtestResponse {
	links := []string{}
	if linksBytes, err := protest.Links.MarshalJSON(); err != nil {
		log.Warnf("error marshalling json from protest: %v")
	} else {
		if err := json.Unmarshal(linksBytes, &links); err != nil {
			log.Warnf("error unmarshaling prest links: %v", err)
		}
	}

	response := ProtestResponse{
		DataSource: DataSourceResponse{
			Reference:     protest.ImportID,
			Name:          protest.DataSource.Slug,
			Attribution:   protest.DataSource.Description,
			LastUpdatedAt: protest.DataSource.LastImport.UTC(),
		},
		Lat:      protest.Lat,
		Lng:      protest.Lng,
		Location: protest.Location,
		Date:     protest.Date.Format("02-01-2006"),
		Notes:    protest.Notes,
		Links:    links,
		Size:     string(protest.Size),
	}

	return response
}

// parseDateString parses a string to a date with format YYYY-mm-ddd
func parseDateString(date string) (*time.Time, error) {
	loc, _ := time.LoadLocation("Asia/Colombo")
	if parsed, err := time.ParseInLocation("2006-01-02", date, loc); err != nil {
		return nil, err
	} else {
		return &parsed, nil
	}
}

func (ProtestController) ListProtests(c echo.Context) error {
	var err error
	db := data.GetDb()
	params := ProtestsParam{}
	err = echo.QueryParamsBinder(c).
		String("datasource", &params.DatasourceSlug).
		String("beforeDate", &params.BeforeDate).
		String("afterDate", &params.AfterDate).
		BindError()
	if err != nil {
		log.Errorf("couldnt parse query: %v", err)
		return errors.New("couldnt parse query parameters, please check")
	}

	filters := data.ProtestFilters{
		DatasourceSlug: params.DatasourceSlug,
	}

	if params.AfterDate != "" {
		if filters.AfterDate, err = parseDateString(params.AfterDate); err != nil {
			return err
		}
	}

	if params.BeforeDate != "" {
		if filters.BeforeDate, err = parseDateString(params.BeforeDate); err != nil {
			return err
		}
	}

	protests, err := data.GetAllProtests(db, filters)
	if err != nil {
		return err
	}
	protestReponse := make([]ProtestResponse, len(protests))
	for idx, protest := range protests {
		protestReponse[idx] = protestToResponseObject(&protest)
	}
	return c.JSON(http.StatusOK, protestReponse)
}

func (ProtestController) GetProtestById(c echo.Context) error {
	db := data.GetDb()
	id, err := strconv.ParseUint(c.Param("protestId"), 10, 32)
	if err != nil {
		return errors.New("id is not an integer")
	}

	protest, err := data.GetProtestById(db, uint(id))
	if err != nil {
		log.Errorf("Could not get the protest from db id=%d: %v", id, err)
		return errors.New("unable to get protest")
	}
	protestResponse := protestToResponseObject(&protest)

	return c.JSON(http.StatusOK, protestResponse)
}
