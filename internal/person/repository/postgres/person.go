package postgres

import (
	"github.com/Davmie/person_service/internal/person/repository"
	"github.com/Davmie/person_service/models"
	"github.com/Davmie/person_service/pkg/logger"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type pgPersonRepo struct {
	Logger logger.Logger
	DB     *gorm.DB
}

func New(logger logger.Logger, db *gorm.DB) repository.PersonRepositoryI {
	return &pgPersonRepo{
		Logger: logger,
		DB:     db,
	}
}

func (pr *pgPersonRepo) Create(p *models.Person) error {
	tx := pr.DB.Create(p)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "pgPersonRepo.Create error while inserting in repo")
	}

	return nil
}

func (pr *pgPersonRepo) Get(id int) (*models.Person, error) {
	var p models.Person
	tx := pr.DB.Where("id = ?", id).Take(&p)

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "pgPersonRepo.Get error")
	}

	return &p, nil
}

func (pr *pgPersonRepo) Update(p *models.Person) error {
	tx := pr.DB.Clauses(clause.Returning{}).Omit("id").Updates(p)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "pgPersonRepo.Update error while inserting in repo")
	}

	return nil
}

func (pr *pgPersonRepo) Delete(id int) error {
	tx := pr.DB.Delete(&models.Person{}, id)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "pgPersonRepo.Delete error")
	}

	return nil
}

func (pr *pgPersonRepo) GetAll() ([]*models.Person, error) {
	var persons []*models.Person

	tx := pr.DB.Find(&persons)

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "pgPersonRepo.GetAll error")
	}

	return persons, nil
}
