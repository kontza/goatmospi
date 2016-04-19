package data

import (
	_ "github.com/jinzhu/gorm"
)

/*
CREATE TABLE Devices(DeviceID INTEGER PRIMARY KEY, Type TEXT, SerialID TEXT, Label TEXT);
CREATE TABLE Temperature(DeviceID INT, Timestamp INT, C REAL, F REAL);
CREATE TABLE Humidity(DeviceID INT, Timestamp INT, H REAL);
CREATE TABLE Flag(DeviceID INT, Timestamp INT, Value TEXT);
CREATE INDEX temperature_dt ON Temperature(DeviceID, Timestamp);
CREATE INDEX humidity_dt ON Humidity(DeviceID, Timestamp);
CREATE INDEX flag_dt ON Flag(DeviceID, Timestamp);
*/

type Device struct {
	DeviceID uint   `gorm:"primary_key;column:DeviceID"`
	Type     string `gorm:"column:Type"`
	SerialID string `gorm:"column:SerialID"`
	Label    string `gorm:"column:Label"`
}

func (Device) TableName() string {
	return "Devices"
}

type Temperature struct {
	DeviceID  uint    `gorm:"primary_key;index:temperature_dt;column:DeviceID"`
	Timestamp int     `gorm:"index:temperature_dt;column:Timestamp"`
	C         float64 `gorm:"column:C"`
	F         float64 `gorm:"column:F"`
}

func (Temperature) TableName() string {
	return "Temperature"
}

type Humidity struct {
	DeviceID  uint    `gorm:"primary_key;index:humidity_dt;column:DeviceID"`
	Timestamp int     `gorm:"index:humidity_dt;column:Timestamp"`
	H         float64 `gorm:"column:H"`
}

func (Humidity) TableName() string {
	return "Humidity"
}

type Flag struct {
	DeviceID  uint   `gorm:"primary_key;index:flag_dt;column:DeviceID"`
	Timestamp int    `gorm:"index:flag_dt;column:Timestamp"`
	Value     string `gorm:"column:Value"`
}

func (Flag) TableName() string {
	return "Flag"
}
