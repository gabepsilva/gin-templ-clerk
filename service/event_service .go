package service

import (
	"errors"
	"gotempl/model"
	"gotempl/repository"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

type EventService struct {
	repo     *repository.EventRepository
	validate *validator.Validate
}

func NewEventService(repo *repository.EventRepository) *EventService {
	return &EventService{
		repo:     repo,
		validate: validator.New(),
	}
}

func (s *EventService) CreateEvent(event *model.Event) error {
	if err := s.validate.Struct(event); err != nil {
		return err
	}

	if event.Title == "" {
		return errors.New("id and eventname are required")
	}

	return s.repo.Create(event)
}

func (s *EventService) GetEvent(id uint64) (*model.Event, error) {
	return s.repo.GetByID(id)
}

func (s *EventService) GetAllEvents() ([]model.Event, error) {
	return s.repo.GetAll()
}

func (s *EventService) GetEventByID(id uint64) (*model.Event, error) {
	return s.repo.GetByID(id)
}

func (s *EventService) UpdateEvent(event *model.Event) error {
	if err := s.validate.Struct(event); err != nil {
		log.Error("Error:", err)
		//fmt.Println("Error:", err)
		return err
	}
	return s.repo.Update(event)
}

func (s *EventService) DeleteEvent(id uint64) error {
	return s.repo.Delete(id)
}

// Additional method to match the handler
func (s *EventService) GetAllEvent() ([]model.Event, error) {
	return s.repo.GetAll()
}
