package data

import (
	_ "github.com/jinzhu/gorm"
)

type Device struct {
	DeviceID uint   `gorm:"primary_key;index:devices_pkey;column:deviceid"`
	Type     string `gorm:"column:type"`
	SerialID string `gorm:"column:serialid"`
	Label    string `gorm:"column:label"`
}

func (Device) TableName() string {
	return "devices"
}

type Temperature struct {
	ID        uint    `gorm:"primary_key;index:temperature_pkey;column:id"`
	DeviceID  uint    `gorm:"column:deviceid"`
	Timestamp int     `gorm:"index:temperature_dt;column:timestamp"`
	C         float64 `gorm:"column:c"`
	F         float64 `gorm:"column:f"`
}

func (Temperature) TableName() string {
	return "temperature"
}

type Humidity struct {
	ID        uint    `gorm:"primary_key;index:humidity_pkey;column:id"`
	DeviceID  uint    `gorm:"column:deviceid"`
	Timestamp int     `gorm:"index:humidity_dt;column:timestamp"`
	H         float64 `gorm:"column:h"`
}

func (Humidity) TableName() string {
	return "humidity"
}

type Flag struct {
	ID        uint   `gorm:"primary_key;index:flag_pkey;column:id"`
	DeviceID  uint   `gorm:"column:deviceid"`
	Timestamp int    `gorm:"index:flag_dt;column:Timestamp"`
	Value     string `gorm:"column:value"`
}

func (Flag) TableName() string {
	return "flag"
}
