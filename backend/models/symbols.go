package models

import (
	"errors"

	"gorm.io/gorm"
)

type Symbol struct {
	gorm.Model
	Symbol string `json:"symbol"`
	ListID uint   `json:"list_id" db:"lists"`
}

func (s *Symbol) CreateSymbol(db *gorm.DB, listID int) error {
	//validates max symbols amount
	var count int64
	db.Model(&Symbol{}).Where("list_id=?", listID).Count(&count)

	if count >= MaxSimbolsAmount {
		return errors.New("list exceeded maximum symbol amount")
	}
	//-------------

	db.Model(&List{}).Where("id=?", listID).Count(&count)
	if count <= 0 {
		return errors.New("List Not Found")
	}

	//validates inexistence
	db.Where("symbol=? AND list_id=?", s.Symbol, listID).Find(&s)
	if s.ID != 0 {
		return errors.New("symbol already exists in this list")
	}
	//-------------

	s.ListID = uint(listID)
	if err := db.Create(&s).Error; err != nil {
		return err
	}

	return nil
}

func (s *Symbol) DeleteSymbol(db *gorm.DB, userID, id int) error {
	//Gets Symbol and validates existence
	//db.Where("userID=? AND id=?", userID, id).Find(&s)
	db.Where("id=?", id).Find(&s)

	if s.ID == 0 {
		return errors.New("Symbol Not Found")
	}
	//-------------------------

	if err := db.Delete(&s, id).Error; err != nil {
		return err
	}

	return nil
}
