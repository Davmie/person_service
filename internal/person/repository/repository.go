package repository

import "github.com/Davmie/person_service/models"

type PersonRepositoryI interface {
	Create(p *models.Person) error
	Get(id int) (*models.Person, error)
	Update(p *models.Person) error
	Delete(id int) error
	GetAll() ([]*models.Person, error)
}
