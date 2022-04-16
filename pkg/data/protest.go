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

type ProtestFilters struct {
	DatasourceSlug string
	BeforeDate     *time.Time
	AfterDate      *time.Time
}

func GetAllProtests(db *gorm.DB, filters ProtestFilters) ([]Protest, error) {
	protests := []Protest{}

	query := db.Joins("DataSource")
	if filters.DatasourceSlug != "" {
		query = query.Where("data_source_slug = ?", filters.DatasourceSlug)
	}

	if filters.BeforeDate != nil {
		query = query.Where("date <= ?", filters.BeforeDate)
	}

	if filters.AfterDate != nil {
		query = query.Where("date >= ?", filters.AfterDate)
	}
	result := query.Find(&protests)
	return protests, result.Error
}

func GetProtestById(db *gorm.DB, id uint) (Protest, error) {
	protest := Protest{}
	result := db.Joins("DataSource").First(&protest, id)
	return protest, result.Error
}
