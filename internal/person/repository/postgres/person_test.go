package postgres

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	personRep "github.com/Davmie/person_service/internal/person/repository"
	"github.com/Davmie/person_service/internal/testBuilders"
	"github.com/Davmie/person_service/models"
	"github.com/Davmie/person_service/pkg/logger"
	"github.com/bxcodec/faker"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"regexp"
	"testing"
)

type PersonRepoTestSuite struct {
	suite.Suite
	db            *sql.DB
	gormDB        *gorm.DB
	mock          sqlmock.Sqlmock
	repo          personRep.PersonRepositoryI
	personBuilder *testBuilders.PersonBuilder
}

func TestPersonRepoSuite(t *testing.T) {
	suite.RunSuite(t, new(PersonRepoTestSuite))
}

func (s *PersonRepoTestSuite) BeforeEach(t provider.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("error while creating sql mock")
	}

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatal("error gorm open")
	}

	var logger logger.Logger

	s.db = db
	s.gormDB = gormDB
	s.mock = mock

	s.repo = New(logger, gormDB)
	s.personBuilder = testBuilders.NewPersonBuilder()
}

func (s *PersonRepoTestSuite) AfterEach(t provider.T) {
	err := s.mock.ExpectationsWereMet()
	t.Assert().NoError(err)
	s.db.Close()
}

func (s *PersonRepoTestSuite) TestCreatePerson(t provider.T) {
	person := s.personBuilder.
		WithID(1).
		WithName("Name").
		WithAge(20).
		WithAddress("Address").
		WithWork("Work").
		Build()

	s.mock.ExpectBegin()

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "people" ("name","age","address","work","id") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WithArgs(person.Name, person.Age, person.Address, person.Work, person.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	s.mock.ExpectCommit()

	err := s.repo.Create(&person)
	t.Assert().NoError(err)
	t.Assert().Equal(1, person.ID)
}

func (s *PersonRepoTestSuite) TestGetPerson(t provider.T) {
	person := s.personBuilder.
		WithID(1).
		WithName("Name").
		WithAge(20).
		WithAddress("Address").
		WithWork("Work").
		Build()

	rows := sqlmock.NewRows([]string{"id", "name", "age", "address", "work"}).
		AddRow(
			person.ID,
			person.Name,
			person.Age,
			person.Address,
			person.Work,
		)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "people" WHERE id = $1 LIMIT $2`)).
		WithArgs(person.ID, 1).
		WillReturnRows(rows)

	resPerson, err := s.repo.Get(person.ID)
	t.Assert().NoError(err)
	t.Assert().Equal(person, *resPerson)
}

func (s *PersonRepoTestSuite) TestUpdatePerson(t provider.T) {
	person := s.personBuilder.
		WithID(1).
		WithName("Name").
		WithAge(20).
		WithAddress("Address").
		WithWork("Work").
		Build()

	rows := sqlmock.NewRows([]string{"id", "name", "age", "address", "work"}).
		AddRow(
			person.ID,
			person.Name,
			person.Age,
			person.Address,
			person.Work,
		)

	s.mock.ExpectBegin()

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`UPDATE "people" SET "name"=$1,"age"=$2,"address"=$3,"work"=$4 WHERE "id" = $5 RETURNING *`)).
		WithArgs(person.Name, person.Age, person.Address, person.Work, person.ID).WillReturnRows(rows)

	s.mock.ExpectCommit()

	err := s.repo.Update(&person)
	t.Assert().NoError(err)
}

func (s *PersonRepoTestSuite) TestDeletePerson(t provider.T) {
	person := s.personBuilder.
		WithID(1).
		WithName("Name").
		WithAge(20).
		WithAddress("Address").
		WithWork("Work").
		Build()

	s.mock.ExpectBegin()

	s.mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "people" WHERE "people"."id" = $1`)).
		WithArgs(person.ID).WillReturnResult(sqlmock.NewResult(int64(person.ID), 1))

	s.mock.ExpectCommit()

	err := s.repo.Delete(person.ID)
	t.Assert().NoError(err)
}

func (s *PersonRepoTestSuite) TestGetAll(t provider.T) {
	persons := make([]models.Person, 10)
	for _, person := range persons {
		err := faker.FakeData(&person)
		t.Assert().NoError(err)
	}

	personsPtr := make([]*models.Person, len(persons))
	for i, person := range persons {
		personsPtr[i] = &person
	}

	rowsPersons := sqlmock.NewRows([]string{"id", "name", "age", "address", "work"})

	for i := range persons {
		rowsPersons.AddRow(persons[i].ID, persons[i].Name, persons[i].Age, persons[i].Address, persons[i].Work)
	}

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "people"`)).
		WillReturnRows(rowsPersons)

	resPersons, err := s.repo.GetAll()
	t.Assert().NoError(err)
	t.Assert().Equal(personsPtr, resPersons)
}
