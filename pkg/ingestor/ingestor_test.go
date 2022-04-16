package ingestor

import (
	"context"
	"testing"

	"github.com/reliefeffortslk/protest-tracker-api/pkg/configs"
	"github.com/reliefeffortslk/protest-tracker-api/pkg/data"
)

var _ = configs.LoadEnv()

func TestWatchdogGoogleSheets(t *testing.T) {
	ctx := context.Background()
	numInserted := 0
	err := IngestFromWatchdogGoogleSheets(ctx, func(p *data.Protest) error {
		if p == nil {
			t.Errorf("protest was nil")
			return nil
		}

		if p.ImportID == "" {
			t.Errorf("import id was blank for %v", p)
		}

		if p.Location == "" {
			t.Errorf("location was blank for %v", p)
		}

		if p.Lat == 0.0 {
			t.Errorf("lat was blank for %v", p)
		}

		if p.Lng == 0.0 {
			t.Errorf("lng was blank for %v", p)
		}
		numInserted = numInserted + 1

		return nil
	})

	if err != nil {
		t.Fatalf("error ingesting watchdog sheet: %v", err)
	}

	if numInserted == 0 {
		t.Errorf("nothing got inserted to the database")
	}
}
