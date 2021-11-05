package models

type Meta struct {
	Date         string
	Population   int64
	Last_updated string
}

type Total struct {
	Confirmed   int
	Deceased    int
	Recovered   int
	Tested      int
	Vaccinated1 int
	Vaccinated2 int
}
type CovidData struct {
	State string
	Meta  Meta
	Total Total
}
