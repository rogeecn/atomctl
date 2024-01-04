package test

import (
	"errors"
	"reflect"
	"time"

	"github.com/jinzhu/copier"
)

type UserForm struct {
	Username string `form:"username" json:"username,omitempty"` //
	Password string `form:"password" json:"password,omitempty"` //
	Age      uint64 `form:"age" json:"age,omitempty"`           //
}

type UserListQueryFilter struct {
	Username *string `query:"username" json:"username,omitempty"` //
	Password *string `query:"password" json:"password,omitempty"` //
	Age      *uint64 `query:"age" json:"age,omitempty"`           //
}

type UserItem struct {
	ID        uint64    `json:"id,omitempty"`         // ID
	CreatedAt time.Time `json:"created_at,omitempty"` // 创建时间
	UpdatedAt time.Time `json:"updated_at,omitempty"` // 更新时间
	Username  string    `json:"username,omitempty"`   //
	Password  string    `json:"password,omitempty"`   //
	Age       uint64    `json:"age,omitempty"`        //
}

func UserItemFillWith(item interface{}) *UserItem {
	m := &UserItem{}
	if err := m.Fill(item); err != nil {
		return nil
	}
	return m
}

func (m *UserItem) Fill(item interface{}) error {
	if reflect.ValueOf(item).Kind() == reflect.Ptr {
		return copier.Copy(&m, item)
	}

	return errors.New("only support pointer type var")
}
