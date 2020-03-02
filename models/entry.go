package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

var (
	ErrNotFound           = errors.New("not found")
	ErrSomethingWentWrong = errors.New("something went wrong")
)

type Entry struct {
	ID               uint32    `gorm:"primary_key;auto_increment" json:"-"`
	Version          uint16    `gorm:"not null" json:"version"`
	UUID             uuid.UUID `gorm:"size:36;" json:"id"`
	DocumentLocation string    `gorm:"not null" json:"documentLocation"`
	Elements         []Element `gorm:"not null" json:"elements"`
	LastUpdated      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"lastUpdated"`
}

type Element struct {
	ID      uint32 `gorm:"primary_key;auto_increment" json:"-"`
	EntryID uint32 `gorm:"not null" json:"-"`
	Type    string `gorm:"not null" json:"type"`
	Text    string `gorm:"not null;" sql:"size:999999" json:"text"`
}

func DBMigrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(&Entry{}, &Element{})
	return db
}

func (e *Entry) Validate() error {

	if e.DocumentLocation == "" {
		return errors.New("documentLocation required")
	}
	if len(e.Elements) == 0 || e.Elements == nil {
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
	db.Preload("Elements").Where("id IN (select max(id) from entries group by uuid)").Find(&entries)
	if db.Error != nil {
		return nil, ErrSomethingWentWrong
	}

	return &entries, nil

}

func (e *Entry) GetVersions(db *gorm.DB, uuid uuid.UUID) (*[]Entry, error) {
	entries := []Entry{}
	db = db.Preload("Elements").Find(&entries).Where("uuid = ?", uuid).Order("version DESC") // Consider Limit

	if db.Error != nil {
		if db.Error == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}

		return nil, db.Error
	}

	return &entries, nil
}

func (e *Entry) GetVersion(db *gorm.DB, uuid uuid.UUID, version uint16) (*Entry, error) {

	db = db.Preload("Elements").Where("uuid = ? AND version = ?", uuid, version).Take(&e)

	if db.Error != nil {
		if db.Error == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}

		return nil, db.Error
	}

	return e, nil
}

func (e *Entry) GetNewestVersion(db *gorm.DB, uuid uuid.UUID) (*Entry, error) {
	db = db.Preload("Elements").Where("uuid = ?", uuid).Order("version DESC").Take(&e)

	if db.Error != nil {
		return nil, db.Error
	}

	return e, nil
}

func (e *Entry) UpdateEntry(db *gorm.DB, uuid uuid.UUID) (*Entry, error) {
	var err error
	existing := &Entry{}

	err = db.Preload("Elements").Where("uuid = ?", uuid).Order("version DESC").Take(&existing).Error

	if err != nil {
		if db.Error == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}

		return nil, db.Error
	}

	existing.DocumentLocation = e.DocumentLocation
	existing.Elements = e.Elements

	e, err = existing.saveEntry(db)
	if err != nil {
		return nil, err
	}

	return e, nil
}

func (e *Entry) CreateEntry(db *gorm.DB) (*Entry, error) {
	e.UUID = uuid.New()
	created, err := e.saveEntry(db)

	if err != nil {
		return nil, err
	}

	return created, nil
}

func (e *Entry) DeleteVersion(db *gorm.DB, uuid uuid.UUID, version uint16) (int64, error) {
	db = db.Model(&Entry{}).Where("uuid = ? AND version = ?", uuid, version).Take(&Entry{}).Delete(&Entry{})

	if db.Error != nil {
		if db.Error == gorm.ErrRecordNotFound {
			return 0, ErrNotFound
		}
		return 0, db.Error
	}

	return db.RowsAffected, nil
}

func (e *Entry) saveEntry(db *gorm.DB) (*Entry, error) {
	e.ID = 0
	e.Version++
	e.LastUpdated = time.Now()

	db = db.Create(&e)

	if db.Error != nil {
		return nil, db.Error
	}

	return e, nil
}
