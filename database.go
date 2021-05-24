package main

import (
	"log"
	"time"

	"gorm.io/driver/sqlite"

	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	//file:test.db
	dsn := "file:db/router.db"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	DB = db
	db.AutoMigrate(Device{}, DeviceActivity{})
}
func StoreRow(d *Device) {
	DB.Create(&d)
}
func LastActivity() *Device {
	d := Device{}
	DB.Where("mac_address=?", d.MacAddress).First(&d)
	return &d
}

type Devices struct {
	Device []Device
}
type Device struct {
	MacAddress       string `gorm:"primaryKey"`
	Name             string
	Status           string
	IP               string
	DeviceActivities []DeviceActivity `gorm:"foreignKey:MacAddress;references:MacAddress"`
}
type DeviceActivity struct {
	//uuid.UUID    `gorm:"type:uuid;primary_key;"`
	MacAddress   string
	LastActivity time.Time
}

// func (DeviceActivity *DeviceActivity) BeforeCreate(db *gorm.DB) error {
// 	DeviceActivity.UUID = uuid.NewV4()
// 	return nil
// }
