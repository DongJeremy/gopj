package models

import (
	"fmt"

	"github.com/go-xorm/xorm"
)

type User struct {
	Id      int    `json:"id" form:"id" xorm:"pk autoincr"`
	Name    string `json:"name" form:"name" xorm:"name"`
	Address string `json:"address" form:"address" xorm:"address"`
	Age     int    `json:"age" form:"age"`
}

// get one
func (model *User) UserGet(db *xorm.Engine, id int) (user User, err error) {
	// find one record
	_, err = db.Table("t_user").Where("id = ?", id).Get(&user)
	return
}

// UserGetList get list
func (model *User) UserGetList(db *xorm.Engine, page int, pageSize int) (users []User, err error) {
	users = make([]User, 0)
	offset := pageSize * (page - 1)
	limit := pageSize
	err = db.Table("t_user").Limit(limit, offset).Find(&users)
	fmt.Println(users)
	return
}

// UserAdd create
func (model *User) UserAdd(db *xorm.Engine) (id int64, err error) {
	user := User{Name: model.Name, Age: model.Age, Address: model.Address}
	id, err = db.Table("t_user").Insert(&user)
	return
}

// UserUpdate update
func (model *User) UserUpdate(db *xorm.Engine, id int) (afr int64, err error) {
	user := User{Name: model.Name, Age: model.Age, Address: model.Address}
	afr, err = db.Table("t_user").Id(id).Update(&user)
	return
}

// UserDelete delete
func (model *User) UserDelete(db *xorm.Engine, id int) (afr int64, err error) {
	afr, err = db.Table("t_user").Id(id).Delete(&User{})
	return
}
