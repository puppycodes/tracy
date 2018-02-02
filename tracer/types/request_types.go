package types

import (
	"github.com/jinzhu/gorm"
)

type Request struct {
	gorm.Model
	RawRequest    string   `json:"raw_request" gorm:"not null"`
	RequestURL    string   `json:"request_url" gorm:"not null"`
	RequestMethod string   `json:"request_method" gorm:"not null"`
	Tracers       []Tracer `json:tracers"`
}
