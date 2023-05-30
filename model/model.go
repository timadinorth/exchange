package model

import (
	"html"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Default struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type User struct {
	Default
	Username string    `gorm:"unique;not null" json:"username"`
	Password string    `gorm:"not null" json:"-"`
	Accounts []Account `json:",omitempty"`
}

func (user *User) Save(DB *gorm.DB) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)
	user.Username = html.EscapeString(strings.TrimSpace(user.Username))

	if err := DB.Create(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (user *User) FindByUsername(DB *gorm.DB, username string) error {
	return DB.Model(User{}).Where("username = ?", username).Take(user).Error
}

type Account struct {
	Default
	UserId  uint
	Address string
	Balance uint64
	ChainID int
	Chain   Chain
}

type Chain struct {
	Default
	ExternalId      uint
	Name            string
	DepositAllowed  bool
	WithdrawAllowed bool
}

type Category struct {
	Default
	Name string `gorm:"not null" json:"name" example:"Soccer"`
	Icon string `json:"icon" example:"https://example.com/example.png"`
	Type string `json:"type" example:"sport"`
}

type Competition struct {
	Default
	Name       string `gorm:"not null" json:"name"`
	Type       string `json:"type"`
	EventCount int    `gorm:"default: 0" json:"event_count"`
}

type Event struct {
}

type Outcome struct {
}
