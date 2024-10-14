package repositories

import (
	"gotempl/database"
	"gotempl/model"
)

type EventRepository struct{}

func (r *EventRepository) Create(event *model.Event) error {
	return database.DB.Create(event).Error
}

func (r *EventRepository) GetAll() ([]model.Event, error) {
	var events []model.Event
	err := database.DB.Find(&events).Error
	return events, err
}

func (r *EventRepository) GetByID(id uint) (*model.Event, error) {
	var event model.Event
	err := database.DB.First(&event, id).Error
	return &event, err
}

func (r *EventRepository) Update(event *model.Event) error {
	return database.DB.Save(event).Error
}

func (r *EventRepository) Delete(id uint) error {
	return database.DB.Delete(&model.Event{}, id).Error
}
