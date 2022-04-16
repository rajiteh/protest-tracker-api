package data

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetGeoSubscriptionForChatId(chartId string) (*GeoSubscription, error) {
	return nil, nil
}

func GetProtestsForCoordinate(db *gorm.DB, lat, lng float64, radiusInKm float64, oldestDate time.Time, newestDate time.Time) ([]Protest, error) {
	var protests []Protest
	minLat, minLng, maxLat, maxLng := GetBoundingBox(lat, lng, radiusInKm*1000)

	query := db.Where(
		"lat >= ? AND lat <= ? AND lng >= ? AND lng <= ? AND date BETWEEN ? AND ?",
		minLat, maxLat, minLng, maxLng, oldestDate, newestDate)

	res := query.Find(&protests)

	if res.Error != nil {
		return nil, res.Error
	}

	return protests, nil
}

func UpdateProtest(db *gorm.DB, protest *Protest) (int64, error) {
	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "import_id"}},
		UpdateAll: true,
	}).Create(protest)
	return result.RowsAffected, result.Error
}
