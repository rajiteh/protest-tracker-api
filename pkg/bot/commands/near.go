package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/reliefeffortslk/protest-tracker-api/pkg/data"
	"github.com/rusq/tbcomctl/v4"
	tb "gopkg.in/tucnak/telebot.v3"
)

type NearCommand struct {
	CommandHandler
}

var NearByRadius = 10.0

func NewNear(bot *tb.Bot) *NearCommand {
	nc := NearCommand{}
	nc.bot = bot

	nc.setupHandlers()

	return &nc
}

func (nc *NearCommand) setupHandlers() {
	log.Info("Setting up handlers for near")

	locationPinInput := tbcomctl.NewInputText("location_pin", "I need to know where you are:", processLocationInput(nc.bot), tbcomctl.IOptValueResolver(func(m *tb.Message) (string, error) {
		location, err := json.Marshal(m.Location)
		if err != nil {
			return "", err
		}
		return string(location), nil
	}))

	tvc := tbcomctl.NewStaticTVC("", nil, nil)
	tvc.TextFn = finalizeNearbyForm(nc.bot)
	finalizerInput := tbcomctl.NewMessage("finalize", tvc)

	form := tbcomctl.NewForm(locationPinInput, finalizerInput).SetOverwrite(false).SetRemoveButtons(true)

	nc.bot.Handle("/nearby", form.Handler)
	nc.ResponseHandlers = append(nc.ResponseHandlers, form.OnTextMiddleware(func(ctx tb.Context) error {
		return nil
	}))

}

// TODO: convert time to firendly value
func finalizeNearbyForm(b *tb.Bot) func(ctx context.Context, c tb.Context) (string, error) {
	return func(ctx context.Context, c tb.Context) (string, error) {
		log := log.WithContext(ctx)
		log.Info("Finalizing form...")

		// ctrl, ok := tbcomctl.ControllerFromCtx(ctx)
		// if !ok {
		// 	return "There was a problem", errors.New("something went wrong trying to process the form")
		// }

		// form := ctrl.Form()
		// formData := form.Data(c.Sender())

		coordinates := &tb.Location{
			Lat: 6.918766,
			Lng: 79.862305,
		}
		// if err := json.Unmarshal([]byte(formData["location_pin"]), coordinates); err != nil {
		// 	log.Error("Error while unmarshaling location_pin: " + err.Error())
		// 	return "", err
		// }
		loc, err := time.LoadLocation("Asia/Colombo")
		if err != nil {
			log.Error("Error while getting location: " + err.Error())
			return "", err
		}

		// now := time.Now()
		now := time.Date(2022, 4, 8, 1, 2, 3, 4, loc)
		midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
		tomorrowMidnight := midnight.AddDate(0, 0, 2)
		db := data.GetDb()
		protests, err := data.GetProtestsForCoordinate(db, float64(coordinates.Lat), float64(coordinates.Lng), NearByRadius, midnight, tomorrowMidnight)
		if err != nil {
			log.Error("Error while getting activities: " + err.Error())
			return "", err
		}
		processedProtests := getSortedProtestsWithDistance(protests, coordinates)
		for idx, protWithDist := range processedProtests {
			// only first 10
			if idx > 10 {
				break
			}
			outputString := protWithDist.GenerateOutput()
			protLocation := &tb.Location{
				Lat: float32(protWithDist.Lat),
				Lng: float32(protWithDist.Lng),
			}

			if msg, err := b.Send(c.Recipient(), outputString); err != nil {
				log.Error("Error sending: " + err.Error())
				return "", err
			} else {
				b.Send(c.Recipient(), protLocation, &tb.SendOptions{ReplyTo: msg})
			}
		}
		summary := ""
		if len(processedProtests) == 0 {
			summary = "No activities found."
		} else {
			summary = fmt.Sprintf("Found %d activities.", len(protests))
			if len(protests) > 10 {
				summary = summary + " Only showing the closest 10."
			}
		}

		return summary, nil
	}
}

type ProtestsWithDistance struct {
	data.Protest
	Distance float64
}

func (protest *ProtestsWithDistance) GenerateOutput() string {
	var sb strings.Builder

	distanceStr := strconv.FormatFloat(math.Round(protest.Distance*100)/100, 'f', -1, 64)
	fmt.Fprintf(&sb, "Location: %s\n", protest.Location)
	fmt.Fprintf(&sb, "Distance from you: %s km\n", distanceStr)
	fmt.Fprintf(&sb, "Size: %s\n", protest.Size)
	fmt.Fprintf(&sb, "Time: %s\n", protest.Date)
	fmt.Fprintf(&sb, "Notes: %s\n", protest.Notes)
	fmt.Fprintf(&sb, "Tracker URL: https://protests.projects.ukr.lk/?current=%s\n", protest.ImportID)

	if linkBytes, err := protest.Links.MarshalJSON(); err == nil {
		var links []string
		if err := json.Unmarshal(linkBytes, &links); err != nil {
			log.Warn("Error unmarshalling links: " + err.Error())
		}
		if len(links) > 0 {
			fmt.Fprint(&sb, "Links:\n")
			for _, link := range links {
				fmt.Fprintf(&sb, "- %s\n", link)
			}
		}
	}
	return sb.String()
}

func getSortedProtestsWithDistance(protests []data.Protest, current *tb.Location) []ProtestsWithDistance {

	protestsWithDistance := make([]ProtestsWithDistance, len(protests))
	for idx, protest := range protests {
		protestsWithDistance[idx] = ProtestsWithDistance{
			Protest:  protest,
			Distance: distance(protest.Lat, protest.Lng, float64(current.Lat), float64(current.Lng)),
		}
	}

	sort.Slice(protestsWithDistance, func(i, j int) bool {
		return protestsWithDistance[i].Distance < protestsWithDistance[j].Distance
	})
	return protestsWithDistance
}

func distance(lat1 float64, lng1 float64, lat2 float64, lng2 float64, unit ...string) float64 {

	radlat1 := float64(math.Pi * lat1 / 180)
	radlat2 := float64(math.Pi * lat2 / 180)

	theta := float64(lng1 - lng2)
	radtheta := float64(math.Pi * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)

	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / math.Pi
	dist = dist * 60 * 1.1515
	dist = dist * 1.609344

	return dist
}
