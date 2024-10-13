package repositories

import (
	"gotempl/database"
	"gotempl/models"
)

type EventRepository struct{}

func (r *EventRepository) Create(event *models.Event) error {
	return database.DB.Create(event).Error
}

func (r *EventRepository) GetAll() ([]models.Event, error) {
	var events []models.Event
	err := database.DB.Find(&events).Error
	return events, err
}

func (r *EventRepository) GetByID(id uint) (*models.Event, error) {
	var event models.Event
	err := database.DB.First(&event, id).Error
	return &event, err
}

func (r *EventRepository) Update(event *models.Event) error {
	return database.DB.Save(event).Error
}

func (r *EventRepository) Delete(id uint) error {
	return database.DB.Delete(&models.Event{}, id).Error
}
