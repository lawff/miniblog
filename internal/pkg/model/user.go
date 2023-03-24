// Copyright 2023 lawff. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"time"

	"github.com/lawff/miniblog/pkg/auth"
	"gorm.io/gorm"
)

type UserM struct {
	CreatedAt time.Time `gorm:"column:createdAt"`      //
	Email     string    `gorm:"column:email"`          //
	ID        int64     `gorm:"column:id;primary_key"` //
	Nickname  string    `gorm:"column:nickname"`       //
	Password  string    `gorm:"column:password"`       //
	Phone     string    `gorm:"column:phone"`          //
	UpdatedAt time.Time `gorm:"column:updatedAt"`      //
	Username  string    `gorm:"column:username"`       //
}

// TableName sets the insert table name for this struct type
func (u *UserM) TableName() string {
	return "user"
}

// BeforeCreate 在创建数据库记录之前加密明文密码.
func (u *UserM) BeforeCreate(tx *gorm.DB) (err error) {
	// Encrypt the user password.
	u.Password, err = auth.Encrypt(u.Password)
	if err != nil {
		return err
	}

	return nil
}
