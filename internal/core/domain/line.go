package domain

import (
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"github.com/twpayne/go-geom/encoding/geojson"
	"github.com/twpayne/go-geom/encoding/wkt"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Line struct {
	Type        string      `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}

func NewLine(coordinates [][]float64) Line {
	return Line{
		Type:        "Linestring",
		Coordinates: coordinates,
	}
}

func (l Line) Value() (driver.Value, error) {

	if l.Coordinates == nil {
		return "", nil
	}

	var poly geom.T

	j, err := json.Marshal(l)

	if err != nil {
		return nil, err
	}

	err = geojson.Unmarshal(j, &poly)

	if err != nil {
		return nil, err
	}

	e, err := wkt.Marshal(poly)

	if err != nil {
		return nil, err
	}

	return "SRID=4326;" + e, nil
}

func (l *Line) Scan(value interface{}) error {
	t, err := hex.DecodeString(value.(string))
	if err != nil {
		return err
	}

	gt, err := ewkb.Unmarshal(t)
	if err != nil {
		return err
	}

	var points [][]float64

	for i := 0; i < len(gt.FlatCoords()); i += 2 {
		points = append(points, []float64{
			gt.FlatCoords()[i],
			gt.FlatCoords()[i+1],
		})
	}

	*l = NewLine(points)

	return nil
}

func (Line) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return "geometry(Linestring, 4326)"
}
