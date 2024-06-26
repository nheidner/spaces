package models

import (
	"fmt"
	"spaces-p/pkg/errors"
	"spaces-p/pkg/uuid"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/mmcloughlin/geohash"
)

const (
	MaxSpaceRadiusM = 100
)

type BaseSpace struct {
	Name               string   `json:"name" binding:"required"` // does NOT have to be unique
	ThemeColorHexaCode string   `json:"themeColorHexaCode" binding:"required,hexcolor"`
	Radius             float64  `json:"radius" binding:"required,min=0,max=100"` // max MUST be same as MaxSpaceRadiusM constant
	Location           Location `json:"location" binding:"required"`
}

type Space struct {
	BaseSpace
	AdminId   UserUid   `json:"adminId" binding:"required"`
	ID        uuid.Uuid `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

type SpaceWithDistance struct {
	Space
	Distance float64 `json:"distance"`
}

type NewSpace struct {
	BaseSpace
	AdminId UserUid
}

type Location struct {
	Long float64 `json:"longitude" binding:"required,min=-180,max=180" validate:"required,min=-180,max=180"`
	Lat  float64 `json:"latitude" binding:"required,min=-90,max=90" validate:"required,min=-90,max=90"`
}

func (loc *Location) ParseString(str string) error {
	const op errors.Op = "models.Location.ParseString"

	parts := strings.Split(str, ",")
	if len(parts) != 2 {
		err := fmt.Errorf("invalid location format: %s", str)
		return errors.E(op, err)
	}

	newLoc := &Location{}
	var err error
	newLoc.Long, err = strconv.ParseFloat(parts[0], 64)
	if err != nil {
		err := fmt.Errorf("invalid longitude: %s", parts[0])
		return errors.E(op, err)
	}

	newLoc.Lat, err = strconv.ParseFloat(parts[1], 64)
	if err != nil {
		err := fmt.Errorf("invalid latitude: %s", parts[1])
		return errors.E(op, err)
	}

	validate := validator.New()
	err = validate.Struct(*newLoc)
	if err != nil {
		return errors.E(op, err)
	}

	*loc = *newLoc

	return nil
}

func (loc *Location) ParseGeoHash(geoHash string) {
	lat, lng := geohash.DecodeCenter(geoHash)

	loc = &Location{
		Long: lng,
		Lat:  lat,
	}
}

func (loc *Location) String() string {
	return fmt.Sprintf("%v,%v", loc.Long, loc.Lat)
}

func (loc *Location) GeoHash(precision int) string {
	if precision > 12 {
		precision = 12
	}

	return geohash.Encode(loc.Lat, loc.Long)[:precision]
}

type Radius float64
type Distance float64
