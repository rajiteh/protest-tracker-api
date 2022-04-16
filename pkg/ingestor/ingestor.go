package ingestor

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/reliefeffortslk/protest-tracker-api/pkg/data"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"gorm.io/datatypes"
)

var log = logrus.New()

var watchdogSheetID = "1yShvemHd_eNNAtC3pmxPs9B5RbGmfBUP1O6WGQ5Ycrg"
var watchdogSheetReadRange = "Protests!A2:K"

var loc, _ = time.LoadLocation("Asia/Colombo")
var coordSanityReg, _ = regexp.Compile(`[^0-9.]+`)

func IngestFromAll() (count int64, err error) {
	ctx := context.Background()
	db := data.GetDb()

	// Watchdog
	dataSource := data.DataSource{
		Slug:        "watchdogteam",
		Description: "Data provided by watchdog.team",
		LastImport:  time.Now(),
	}

	if err := data.UpdateDataSource(db, &dataSource); err != nil {
		return -1, errors.Wrap(err, "datasource update error")
	}

	count = 0

	callbackFn := func(p *data.Protest) error {
		p.DataSource = dataSource
		log.Infof("Inserting protest %v", p)
		if affectedRows, err := data.UpdateProtest(db, p); err != nil {
			return errors.Wrap(err, "couldnt create protest")
		} else {
			count = count + affectedRows
		}
		return nil
	}
	if err := IngestFromWatchdogGoogleSheets(ctx, callbackFn); err != nil {
		return -1, err
	}

	return count, nil
}

func IngestFromWatchdogGoogleSheets(ctx context.Context, callbackFn func(p *data.Protest) error) error {

	authJsonB64 := os.Getenv("GCP_AUTH_JSON_B64")
	if authJsonB64 == "" {
		return fmt.Errorf("auth info not supplied")
	}
	authJsonBytes, _ := base64.StdEncoding.DecodeString(authJsonB64)
	srv, err := sheets.NewService(ctx, option.WithCredentialsJSON(authJsonBytes))
	if err != nil {
		return err
	}

	resp, err := srv.Spreadsheets.Values.Get(watchdogSheetID, watchdogSheetReadRange).Do()
	if err != nil {
		return errors.Wrapf(err, "unable to retrieve data from sheet")
	}

	if len(resp.Values) == 0 {
		return fmt.Errorf("no data found from sheet")
	} else {
		for idx, row := range resp.Values {
			rowNum := idx + 1 // actual rows are 1 plus because we skip headers

			if len(row) >= 4 {
				numFields := len(row)
				// Schema
				// 0 - Protest_id (required)
				// 1 - Location (required)
				// 2 - LatLng (required)
				// 3 - Date (required)
				// 4 - Footage
				// 5 - Size
				// 6 - Status
				// 7 - Notes
				// 8 - Footage Notes

				if row[0] == nil {
					log.Warnf("No protest ID defined for row %d, skipped.", rowNum)
					continue
				}
				protestId, ok := row[0].(string)
				if !ok || protestId == "" {
					log.Warnf("unable to parse protestId for row %d, skipped.", rowNum)
					continue
				}

				if row[1] == nil {
					log.Warnf("No location defined for protestId %s, skipped.", protestId)
					continue
				}
				location, ok := row[1].(string)
				if !ok || location == "" {
					log.Warnf("unable to parse location for protestID %s, skipped.", protestId)
					continue
				}

				if row[2] == nil {
					log.Warnf("No latlng defined for protestID %s, skipped.", protestId)
					continue
				}
				latLng, ok := row[2].(string)
				if !ok || latLng == "" {
					log.Warnf("unable to parse latlng for protestID %s, skipped.", protestId)
					continue
				}

				var latLngSlice []float64
				for idx, coord := range strings.Split(latLng, ",") {
					if idx > 1 {
						break // malformed input data but, we are forgiving
					}

					if err != nil {
						log.Fatal(err)
					}
					coord := coordSanityReg.ReplaceAllString(coord, "")
					if parsedCoord, err := strconv.ParseFloat(strings.TrimSpace(coord), 64); err != nil {
						log.Warnf("unable to parse location coordinate %s: %v", coord, err)
						break
					} else {
						latLngSlice = append(latLngSlice, parsedCoord)
					}
				}

				if len(latLngSlice) != 2 {
					log.Warnf("unable to seperate lat lng for protestID %s with latlng %s, skipped.", protestId, latLng)
					continue
				}

				if row[3] == nil {
					log.Warnf("No date defined for protestID %s, skipped.", protestId)
					continue
				}
				dateRaw, ok := row[3].(string)
				if !ok || dateRaw == "" {
					log.Warnf("unable to parse date for protestID %s, skipped.", protestId)
					continue
				}

				dateParsed, err := time.ParseInLocation("2/1/2006", dateRaw, loc)
				if err != nil {
					log.Warnf("unable to parse the date to time.Time for protestID %s with value %s: %v", protestId, dateRaw, err)
					continue
				}

				links := []string{}

				if numFields > 4 && row[4] != nil {
					for _, linkByNL := range strings.Split(row[4].(string), "\n") {
						for _, link := range strings.Split(linkByNL, ",") {

							link = strings.TrimSpace(link)

							// skip blank
							if link == "" {
								continue
							}

							// Try to add the schema for links without one
							if !strings.HasPrefix(link, "http") {
								link = fmt.Sprintf("https://%s", link)
							}

							if _, err := url.ParseRequestURI(link); err != nil {
								log.Warnf("potentially invalid link in protestId %s link: %s, skipped.", protestId, link)
								continue
							}

							links = append(links, link)
						}
					}
				}

				normalizedLinks, err := json.Marshal(links)
				if err != nil {
					return errors.Wrapf(err, "cannot marshal the links")
				}

				size := string(data.SizeUnknown)
				if numFields > 5 && row[5] != nil {
					size = row[5].(string)
				}

				notes := ""
				if numFields > 7 && row[7] != nil {
					notes = row[7].(string)
				}

				notesFootage := ""
				if numFields > 8 && row[8] != nil {
					notesFootage = row[8].(string)
				}

				notesAggregate := strings.TrimSpace(fmt.Sprintf("%s\n%s", notes, notesFootage))

				p := data.Protest{
					ImportID: protestId,
					Lat:      latLngSlice[0],
					Lng:      latLngSlice[1],
					Location: location,
					Date:     dateParsed,
					Notes:    notesAggregate,
					Links:    datatypes.JSON(normalizedLinks),
					Size:     data.Size(size),
				}

				if err := callbackFn(&p); err != nil {
					return errors.Wrapf(err, "couldnt insert protest to db")
				}

			}
		}
	}

	return nil
}
