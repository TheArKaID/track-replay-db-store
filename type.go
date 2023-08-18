package main

import (
	"database/sql/driver"
	"errors"
)

type NullString string
type NullInt int

type Device struct {
	Id             string `ch:"id" json:"id"`
	Name           string `ch:"name" json:"name"`
	Model          string `ch:"model" json:"model"`
	Phone          string `ch:"phone" json:"phone"`
	Status         string `ch:"status" json:"status"`
	Contact        string `ch:"contact" json:"contact"`
	GroupId        int    `ch:"group_id" json:"groupId"`
	Category       string `ch:"category" json:"category"`
	Disabled       bool   `ch:"disabled" json:"disabled"`
	UniqueId       string `ch:"unique_id" json:"uniqueId"`
	Attributes     string `ch:"attributes" json:"attributes"`
	LastUpdate     string `ch:"last_update" json:"lastUpdate"`
	ExpirationTime string `ch:"expiration_time" json:"expirationTime"`
}

type Position struct {
	Id         string  `ch:"id" json:"id"`
	Speed      int     `ch:"speed" json:"speed"`
	Valid      bool    `ch:"valid" json:"valid"`
	Course     int     `ch:"course" json:"course"`
	Address    string  `ch:"address" json:"address"`
	FixTime    string  `ch:"fix_time" json:"fixTime"`
	Network    string  `ch:"network" json:"network"`
	Accuracy   int     `ch:"accuracy" json:"accuracy"`
	Altitude   int     `ch:"altitude" json:"altitude"`
	DeviceId   string  `ch:"device_id" json:"deviceId"`
	Latitude   float64 `ch:"latitude" json:"latitude"`
	Outdated   bool    `ch:"outdated" json:"outdated"`
	Protocol   string  `ch:"protocol" json:"protocol"`
	Longitude  float64 `ch:"longitude" json:"longitude"`
	Attributes string  `ch:"attributes" json:"attributes"`
	DeviceTime string  `ch:"device_time" json:"deviceTime"`
	ServerTime string  `ch:"server_time" json:"serverTime"`
}

type Data struct {
	Device   Device
	Position Position
}

func (s *NullString) Scan(value interface{}) error {
	if value == nil {
		*s = ""
		return nil
	}
	strVal, ok := value.(string)
	if !ok {
		return errors.New("column is not a string")
	}
	*s = NullString(strVal)
	return nil
}

func (s NullString) Value() (driver.Value, error) {
	if len(s) == 0 { // if nil or empty string
		return nil, nil
	}
	return string(s), nil
}

func (s *NullInt) Scan(value interface{}) error {
	if value == nil {
		*s = 0
		return nil
	}
	strVal, ok := value.(int)
	if !ok {
		return errors.New("column is not a int")
	}
	*s = NullInt(strVal)
	return nil
}

func (s NullInt) Value() (driver.Value, error) {
	if s == 0 { // if nil or empty string
		return nil, nil
	}
	return int(s), nil
}
