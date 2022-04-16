package data

import "math"

func GetBoundingBox(lat, lng float64, distanceInMeters float64) (minLat, minLng, maxLat, maxLng float64) {
	latRadian := lat * math.Pi / 180
	degLatKm := 110.574235
	degLngKm := 110.572833 * math.Cos(latRadian)

	deltaLat := distanceInMeters / 1000.0 / degLatKm
	deltaLng := distanceInMeters / 1000.0 / degLngKm

	minLat = lat - deltaLat
	minLng = lng - deltaLng
	maxLat = lat + deltaLat
	maxLng = lng + deltaLng

	return minLat, minLng, maxLat, maxLng
}
