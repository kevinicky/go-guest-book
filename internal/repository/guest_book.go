package repository

import "gorm.io/gorm"

type GuestBookRepository interface {
	Ping() error
}

type guestBookRepository struct {
	db *gorm.DB
}

func NewGuestBookRepository(db *gorm.DB) GuestBookRepository {
	return &guestBookRepository{
		db: db,
	}
}

func (gb *guestBookRepository) Ping() error {
	db, err := gb.db.DB()
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	return nil
}
