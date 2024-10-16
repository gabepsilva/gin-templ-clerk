package repository

import (
	"gotempl/model"

	"gorm.io/gorm"
)

type EventRepository struct {
	DB *gorm.DB
}

func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{DB: db}
}

func (r *EventRepository) Create(event *model.Event) error {
	return r.DB.Create(event).Error
}

func (r *EventRepository) GetAll() ([]model.Event, error) {
	var events []model.Event
	err := r.DB.Find(&events).Error
	return events, err
}

func (r *EventRepository) GetByID(id uint64) (*model.Event, error) {
	var event model.Event
	result := r.DB.First(&event, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &event, nil
}

func (r *EventRepository) Update(event *model.Event) error {
	return r.DB.Save(event).Error
}

func (r *EventRepository) Delete(id uint64) error {
	return r.DB.Delete(&model.Event{}, "id = ?", id).Error
}
