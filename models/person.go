package models

import (
	"github.com/mmcloughlin/geohash"
	"gopkg.in/underarmour/dynago.v2"
)

type Person struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	CarType string `json:"carType"`
	Driver  bool   `json:"driver"`

	// Location info
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
	GeoHash uint64  `json:"geohash"`
}

// Save saves a domain; use this to do updates or to create items
func (d *Person) Save() error {
	hash := geohash.EncodeInt(float64(d.Lat), float64(d.Lng))
	d.GeoHash = hash
	_, err := BasePut("people", d.ToDocument()).Execute()
	return err
}

func (d *Person) GetByEmail(email string) *Person {
	res, err := BaseGet("people", dynago.HashKey("email", email)).Execute()
	if err != nil {
		panic(err)
	}
	d.FromDocument(res.Item)
	return d
}

// ToDocument returns the document form of a domain
func (d *Person) ToDocument() dynago.Document {
	return dynago.Document{
		"name":     d.Name,
		"email":    d.Email,
		"car_type": d.CarType,
		"drive":    d.Driver,
		"lat":      d.Lat,
		"lng":      d.Lng,
		"geohash":  d.GeoHash,
	}
}

func (d *Person) FromDocument(doc dynago.Document) *Person {
	d.Name = doc.GetString("name")
	d.Email = doc.GetString("email")
	d.CarType = doc.GetString("car_type")
	d.Driver = doc.GetBool("driver")

	lat, err := doc.GetNumber("lat").FloatVal()
	if err != nil {
		panic(err)
	}
	lng, err := doc.GetNumber("lng").FloatVal()
	if err != nil {
		panic(err)
	}
	hash, err := doc.GetNumber("geohash").Uint64Val()
	if err != nil {
		panic(err)
	}

	d.Lat = lat
	d.Lng = lng
	d.GeoHash = hash
	return d
}

func (d *Person) Query() *dynago.Query {
	return BaseQuery("people")
}
