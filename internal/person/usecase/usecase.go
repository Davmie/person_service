package usecase

import (
	personRep "github.com/Davmie/person_service/internal/person/repository"
	"github.com/Davmie/person_service/models"
	"github.com/pkg/errors"
)

type PersonUseCaseI interface {
	Create(p *models.Person) error
	Get(id int) (*models.Person, error)
	Update(p *models.Person) error
	Delete(id int) error
	GetAll() ([]*models.Person, error)
}

type personUseCase struct {
	personRepository personRep.PersonRepositoryI
}

func New(aRep personRep.PersonRepositoryI) PersonUseCaseI {
	return &personUseCase{
		personRepository: aRep,
	}
}

func (pUC *personUseCase) Create(p *models.Person) error {
	err := pUC.personRepository.Create(p)

	if err != nil {
		return errors.Wrap(err, "personUseCase.Create error")
	}

	return nil
}

func (pUC *personUseCase) Get(id int) (*models.Person, error) {
	resPerson, err := pUC.personRepository.Get(id)

	if err != nil {
		return nil, errors.Wrap(err, "personUseCase.Get error")
	}

	return resPerson, nil
}

func (pUC *personUseCase) Update(p *models.Person) error {
	_, err := pUC.personRepository.Get(p.ID)

	if err != nil {
		return errors.Wrap(err, "personUseCase.Update error: Person not found")
	}

	err = pUC.personRepository.Update(p)

	if err != nil {
		return errors.Wrap(err, "personUseCase.Update error: Can't update in repo")
	}

	return nil
}

func (pUC *personUseCase) Delete(id int) error {
	_, err := pUC.personRepository.Get(id)

	if err != nil {
		return errors.Wrap(err, "personUseCase.Delete error: Person not found")
	}

	err = pUC.personRepository.Delete(id)

	if err != nil {
		return errors.Wrap(err, "personUseCase.Delete error: Can't delete in repo")
	}

	return nil
}

func (pUC *personUseCase) GetAll() ([]*models.Person, error) {
	persons, err := pUC.personRepository.GetAll()
	if err != nil {
		return nil, errors.Wrap(err, "personUseCase.GetAll error")
	}

	return persons, nil
}
