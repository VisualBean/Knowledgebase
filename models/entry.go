package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Entry struct {
	ID               uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Version          uint16    `gorm:"not null" json:"version"`
	UUID             string    `gorm:"size:36; unique_index"`
	DocumentLocation string    `gorm:"not null" json:"documentLocation"`
	Elements         []Element `gorm:"not null" json:"elements"`
	LastUpdated      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"lastUpdated"`
}

type Element struct {
	UUID string `gorm:"size:36; foreignkey:UUID"`
	Type string `gorm:"not null" json:"type"`
	Text string `gorm:"not null" json:"text"`
}

func (e *Entry) Validate() error {

	if e.DocumentLocation == "" {
		return errors.New("documentLocation required")
	}
	if len(e.Elements) == 0 {
		return errors.New("atleast one element required")
	}

	for _, element := range e.Elements {
		if element.Text == "" {
			return errors.New("element must have text")
		}
	}

	return nil
}

func (e *Entry) GetAllEntries(db *gorm.DB) (*[]Entry, error) {
	entries := []Entry{}
	db = db.Model(&Entry{}).Find(&entries).Order("Version DESC") // Consider Limit

	if db.Error != nil {
		return nil, db.Error
	}

	return &entries, nil

}

func (e *Entry) GetAllVersions(db *gorm.DB, id uint32) (*[]Entry, error) {
	entries := []Entry{}
	db = db.Model(&Entry{}).Find(&entries).Where("id = ?", id).Order("Version DESC") // Consider Limit

	if db.Error != nil {
		return nil, db.Error
	}

	return &entries, nil
}

func (e *Entry) GetNewestVersion(db *gorm.DB, id uint32) (*Entry, error) {
	db = db.Model(&Entry{}).Where("id = ?", id).Order("Version DESC").Take(&e)

	if db.Error != nil {
		return nil, db.Error
	}

	return e, nil
}

func (e *Entry) UpdateEntry(db *gorm.DB, id uint32) (*Entry, error) {
	var err error
	var existing *Entry

	err = db.Model(&Entry{}).Where("id = ?", id).Take(&existing).Error

	existing.DocumentLocation = e.DocumentLocation
	existing.Elements = e.Elements

	e, err = existing.saveEntry(db)

	if err != nil {
		return nil, err
	}

	return e, nil
}

func (e *Entry) CreateEntry(db *gorm.DB) (*Entry, error) {
	created, err := e.saveEntry(db)

	if err != nil {
		return nil, err
	}

	return created, nil
}

func (e *Entry) DeleteVersion(db *gorm.DB, id uint32, version uint16) (int64, error) {
	db = db.Model(&Entry{}).Where("id = ? AND version = ?", id, version).Take(&Entry{}).Delete(&Entry{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (e *Entry) saveEntry(db *gorm.DB) (*Entry, error) {
	e.Version++
	e.UUID = uuid.New().String()
	e.LastUpdated = time.Now()

	db = db.Create(&e)
	if db.Error != nil {
		return nil, db.Error
	}

	return e, nil
}
