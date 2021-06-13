package main

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"

	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	//file:test.db
	dsn := "file:db/router.db"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Warn, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // enable color
		},
	)

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   newLogger,
	})
	if err != nil {
		log.Fatal(err)
	}

	DB = db
	db.AutoMigrate(Device{}, DeviceActivity{})
}
func StoreRow(d *Device) {
	da := DeviceActivity{}
	var daCount int64
	DB.Where("mac_address = ?", d.MacAddress).First(&da).Count(&daCount)
	if daCount > 0 {
		//update status flag
		DB.Model(&d).Where("mac_address = ?", d.MacAddress).Update("status", d.Status)
		//Calculate time since last activity for status on devices
		tt := da.TotalTime
		if d.Status == "on" {
			tt = tt + 1
		}
		DB.Model(&da).Where("mac_address = ?", d.MacAddress).Update("total_time", tt).Update("last_activity", d.DeviceActivities[0].LastActivity)

	} else {
		DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&d)
	}

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
	MacAddress   string `gorm:"unique"`
	LastActivity time.Time
	TotalTime    int
}
