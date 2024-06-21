package models

import "time"

type Item struct {
	ItemID int64     `xorm:"pk autoincr" json:"item_id"`
	Name   string    `xorm:"varchar(255)" json:"name"`
	Price  float64   `xorm:"decimal(10,2)" json:"price"`
	CreatedAt time.Time `xorm:"created" json:"-"`
	UpdatedAt time.Time `xorm:"updated" json:"-"`
}
