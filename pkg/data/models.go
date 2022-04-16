package data

import (
	"fmt"
	"time"

	"crypto/sha256"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Size string

var (
	SizeSmall   Size = "small"
	SizeMedium  Size = "medium"
	SizeLarge   Size = "large"
	SizeUnknown Size = "unknown"
)

type DataSource struct {
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Slug        string         `gorm:"primarykey"`
	Description string
	LastImport  time.Time
}

type Protest struct {
	gorm.Model
	ImportID       string     `gorm:"uniqueIndex"`
	DataSource     DataSource `gorm:"foreignkey:DataSourceSlug"`
	DataSourceSlug string
	Lat            float64
	Lng            float64
	Location       string
	Date           time.Time
	Notes          string
	Links          datatypes.JSON
	Size           Size
	EntryHash      string
}

func (p *Protest) BeforeSave(tx *gorm.DB) error {
	hashSeed := fmt.Sprintf("%v%v%v%v%v%v%v%v", p.ImportID, p.Lat, p.Lng, p.Location, p.Date, p.Notes, p.Links, p.Size)
	h := sha256.New()
	h.Write([]byte(hashSeed))
	p.EntryHash = fmt.Sprintf("%x", h.Sum(nil))

	return nil
}

type GeoSubscription struct {
	gorm.Model
	Lat          float64
	Lng          float64
	Radius       float64
	ChatID       int64
	FailureCount int
}

type ProtestNotification struct {
	gorm.Model
	ChatID    int64
	Protest   Protest `gorm:"foreignkey:ProtestID"`
	ProtestID int
	EntryHash string
}
