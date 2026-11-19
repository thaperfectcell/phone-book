package service

import (
	"context"
	"errors"

	"phone-book/model"
	"phone-book/repository"
)

// ErrValidation возвращается при невалидных входных данных.
var ErrValidation = errors.New("validation error")

type ContactService interface {
	ListContacts(ctx context.Context) ([]model.Contact, error)
	GetContact(ctx context.Context, id int) (*model.Contact, error)
	CreateContact(ctx context.Context, contact *model.Contact) error
	UpdateContact(ctx context.Context, contact *model.Contact) error
	DeleteContact(ctx context.Context, id int) error
}

type contactService struct {
	repo repository.ContactRepository
}

func NewContactService(repo repository.ContactRepository) ContactService {
	return &contactService{repo: repo}
}

func (s *contactService) ListContacts(ctx context.Context) ([]model.Contact, error) {
	return s.repo.GetAll(ctx)
}

func (s *contactService) GetContact(ctx context.Context, id int) (*model.Contact, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *contactService) CreateContact(ctx context.Context, contact *model.Contact) error {
	if err := validate(contact); err != nil {
		return err
	}
	return s.repo.Create(ctx, contact)
}

func (s *contactService) UpdateContact(ctx context.Context, contact *model.Contact) error {
	if contact.ID <= 0 {
		return errors.New("invalid id")
	}
	if err := validate(contact); err != nil {
		return err
	}
	return s.repo.Update(ctx, contact)
}

func (s *contactService) DeleteContact(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid id")
	}
	return s.repo.Delete(ctx, id)
}

func validate(c *model.Contact) error {
	if c.Username == "" {
		return errors.New("username is required")
	}
	return nil
}