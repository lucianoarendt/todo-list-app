package models

import (
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

type List struct {
	gorm.Model
	// ID     uint   `json:"id" gorm:"primaryKey"`
	UserID    uint     `json:"user_id"`
	IsDefault bool     `json:"default"`
	Name      string   `json:"name"`
	Symbols   []Symbol `json:"symbols" db:"symbols" gorm:"foreignKey:ListID"`
}

func (l *List) Equals(list List) bool {
	if l.ID != list.ID {
		return false
	}

	if l.UserID != list.UserID {
		return false
	}

	if l.IsDefault != list.IsDefault {
		return false
	}

	if l.Name != list.Name {
		return false
	}

	for _, s := range l.Symbols {
		if !list.Contains(s) {
			return false
		}
	}

	for _, s := range list.Symbols {
		if !l.Contains(s) {
			return false
		}
	}

	return true
}

func (l *List) Contains(s Symbol) bool {
	for _, e := range l.Symbols {
		if e == s {
			return true
		}
	}
	return false
}

const MaxListsAmount = 10
const MaxSimbolsAmount = 50

func (l *List) CreateList(db *gorm.DB) error {
	//validates max list amount
	var count int64
	db.Model(List{}).Where("user_id=?", l.UserID).Count(&count)

	if count >= MaxListsAmount {
		return errors.New("user exceeded maximum list amount")
	}
	//-----------------------
	if err := db.Create(&l).Error; err != nil {
		return err
	}

	return nil
}

func (l *List) ReadListById(db *gorm.DB, userId int, id int) error {
	if err := db.Where("user_id=? AND id=?", userId, id).Find(&l).Error; err != nil {
		return err
	}

	if l.ID == 0 {
		return errors.New("List Not Found")
	}

	if err := l.PopulateWithSymbols(db); err != nil {
		return err
	}

	return nil
}

func (l *List) PopulateWithSymbols(db *gorm.DB) error {
	if err := db.Where("list_id=?", l.ID).Find(&l.Symbols).Error; err != nil {
		return err
	}
	return nil
}

func (l *List) ReadAllLists(db *gorm.DB, userID int) ([]List, error) {
	var lists []List

	if err := db.Where("user_id=? and is_default=0", userID).Find(&lists).Error; err != nil {
		return nil, err
	}

	for i := range lists {
		lists[i].PopulateWithSymbols(db)
	}

	return lists, nil
}

func (l *List) UpdateList(db *gorm.DB, userID int, id int) (*List, error) {
	l.UserID = uint(userID)
	//Gets list and validates inexistence
	var list List
	if err := db.Where("user_id=? AND id=?", l.UserID, id).Find(&list).Error; err != nil {
		return nil, err
	}

	if list.ID == 0 {
		return nil, errors.New("List Not Found")
	}
	//----------------

	//Updates list fields
	list.IsDefault = l.IsDefault
	list.Name = l.Name
	list.UserID = l.UserID
	if l.Symbols != nil {
		var symbol Symbol
		db.Where("list_id=?", id).Delete(&symbol)

		list.Symbols = l.Symbols
	}
	db.Where("id=?", id).Save(&list)
	return &list, nil
}

func (l *List) DeleteListByID(db *gorm.DB, userID int, id int) error {
	db.Where("user_id=? AND id=?", userID, id).Find(&l)
	if l.ID == 0 {
		return errors.New("List Not Found")
	}

	if err := db.Delete(&l, id).Error; err != nil {
		return err
	}
	return nil
}

func (l *List) ReadAllDefault(db *gorm.DB) ([]List, error) {
	var lists []List

	db.Where("is_default=1").Find(&lists)

	for i := range lists {
		db.Where("list_id=?", lists[i].ID).Find(&lists[i].Symbols)
	}

	return lists, nil
}

func (l List) MarshalBinary() ([]byte, error) {
	return json.Marshal(l)
}

func (l *List) Unmarshal(jsonList string) error {
	return json.Unmarshal([]byte(jsonList), l)
}
