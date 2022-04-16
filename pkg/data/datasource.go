package data

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func UpdateDataSource(db *gorm.DB, datasource *DataSource) error {
	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "slug"}},
		UpdateAll: true,
	}).Create(datasource)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
