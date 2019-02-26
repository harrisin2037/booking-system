package asset

type Boat struct {
	SizeInMeter   float64 `json:"size,omitempty" gorm:"not null;column:size"`
	SizeUnit      string  `json:"size_unit,omitempty" gorm:"not null;column:size_unit"`
	BuildYear     int     `json:"build_year,omitempty" gorm:"not null;column:build_year"`
	RefitYear     int     `json:"refit_year,omitempty" gorm:"not null;column:refit_year"`
	Shipyard      string  `json:"shipyard,omitempty" gorm:"not null;column:shipyard"`
	Speed         float64 `json:"speed,omitempty" gorm:"not null;column:speed"`
	SpeedUnit     string  `json:"speed_unit,omitempty" gorm:"not null;column:speed_unit"`
	Color         string  `json:"color,omitempty" gorm:"not null;column:color"`
	NumberOfCrew  int     `json:"number_of_crew,omitempty" gorm:"not null;column:number_of_crew"`
	NumberOfGuest int     `json:"number_of_guest,omitempty" gorm:"not null;column:number_of_guest"`
}
