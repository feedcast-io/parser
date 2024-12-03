package resources

import (
	"strings"
)

type ShippingObject struct {
	Country         string `xml:"country"`
	Price           string `xml:"price"`
	MinHandlingTime int    `xml:"min_handling_time"`
	MaxHandlingTime int    `xml:"max_handling_time"`
	MinTransitTime  int    `xml:"min_transit_time"`
	MaxTransitTime  int    `xml:"max_transit_time"`
}

type Shipping struct {
	Country         string
	Price           Price
	MinHandlingTime int
	MaxHandlingTime int
	MinTransitTime  int
	MaxTransitTime  int
	IsDefined       bool
}

func (s *Shipping) ParseString(raw string) {
	parts := strings.Split(raw, ":")
	s.IsDefined = false
	s.Country = ""
	s.Price = Price{}

	if len(parts) >= 2 {
		s.Country = parts[0]
		s.Price.ParseString(parts[len(parts)-1])
		s.IsDefined = s.Price.IsDefined
	}
}

func (s *Shipping) FromObject(obj ShippingObject) {
	s.Country = obj.Country
	s.Price = Price{}
	s.Price.ParseString(obj.Price)
	s.IsDefined = s.Price.IsDefined && len(s.Country) > 0
	s.MinHandlingTime = obj.MinHandlingTime
	s.MaxHandlingTime = obj.MaxHandlingTime
	s.MinTransitTime = obj.MinTransitTime
	s.MaxTransitTime = obj.MaxTransitTime
}
