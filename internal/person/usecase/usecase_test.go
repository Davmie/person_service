package usecase

import (
	personMocks "github.com/Davmie/person_service/internal/person/repository/mocks"
	"github.com/Davmie/person_service/internal/testBuilders"
	"github.com/Davmie/person_service/models"
	"github.com/bxcodec/faker"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/pkg/errors"
	"testing"
)

type PersonTestSuite struct {
	suite.Suite
	uc             PersonUseCaseI
	personRepoMock *personMocks.PersonRepositoryI
	personBuilder  *testBuilders.PersonBuilder
}

func TestPersonTestSuite(t *testing.T) {
	suite.RunSuite(t, new(PersonTestSuite))
}

func (s *PersonTestSuite) BeforeEach(t provider.T) {
	s.personRepoMock = personMocks.NewPersonRepositoryI(t)
	s.uc = New(s.personRepoMock)
	s.personBuilder = testBuilders.NewPersonBuilder()
}

func (s *PersonTestSuite) TestCreatePerson(t provider.T) {
	person := s.personBuilder.WithID(1).
		WithName("Name").
		WithAge(20).
		WithAddress("Address").
		WithWork("Work").
		Build()

	s.personRepoMock.On("Create", &person).Return(nil)
	err := s.uc.Create(&person)

	t.Assert().NoError(err)
	t.Assert().Equal(person.ID, 1)
}

func (s *PersonTestSuite) TestUpdatePerson(t provider.T) {
	var err error
	person := s.personBuilder.WithID(1).
		WithName("Name").
		WithAge(20).
		WithAddress("Address").
		WithWork("Work").
		Build()

	notFoundPerson := s.personBuilder.WithID(0).Build()

	s.personRepoMock.On("Get", person.ID).Return(&person, nil)
	s.personRepoMock.On("Update", &person).Return(nil)
	s.personRepoMock.On("Get", notFoundPerson.ID).Return(&notFoundPerson, errors.Wrap(err, "Person not found"))
	s.personRepoMock.On("Update", &notFoundPerson).Return(errors.Wrap(err, "Person not found"))

	cases := map[string]struct {
		ArgData *models.Person
		Error   error
	}{
		"success": {
			ArgData: &person,
			Error:   nil,
		},
		"Person not found": {
			ArgData: &notFoundPerson,
			Error:   errors.Wrap(err, "Person not found"),
		},
	}

	for name, test := range cases {
		t.Run(name, func(t provider.T) {
			err := s.uc.Update(test.ArgData)
			t.Assert().ErrorIs(err, test.Error)
		})
	}
}

func (s *PersonTestSuite) TestGetPerson(t provider.T) {
	person := s.personBuilder.WithID(1).
		WithName("Name").
		WithAge(20).
		WithAddress("Address").
		WithWork("Work").
		Build()

	s.personRepoMock.On("Get", person.ID).Return(&person, nil)
	result, err := s.uc.Get(person.ID)

	t.Assert().NoError(err)
	t.Assert().Equal(&person, result)
}

func (s *PersonTestSuite) TestDeletePerson(t provider.T) {
	var err error
	person := s.personBuilder.WithID(1).
		WithName("Name").
		WithAge(20).
		WithAddress("Address").
		WithWork("Work").
		Build()
	notFoundPerson := s.personBuilder.WithID(0).Build()

	s.personRepoMock.On("Get", person.ID).Return(&person, nil)
	s.personRepoMock.On("Delete", person.ID).Return(nil)
	s.personRepoMock.On("Get", notFoundPerson.ID).Return(&notFoundPerson, errors.Wrap(err, "Person not found"))
	s.personRepoMock.On("Delete", notFoundPerson.ID).Return(errors.Wrap(err, "Person not found"))

	cases := map[string]struct {
		PersonID int
		Error    error
	}{
		"success": {
			PersonID: person.ID,
			Error:    nil,
		},
		"Person not found": {
			PersonID: notFoundPerson.ID,
			Error:    errors.Wrap(err, "Person not found"),
		},
	}

	for name, test := range cases {
		t.Run(name, func(t provider.T) {
			err := s.uc.Delete(test.PersonID)
			t.Assert().ErrorIs(err, test.Error)
		})
	}
}

func (s *PersonTestSuite) TestGetAll(t provider.T) {
	persons := make([]models.Person, 0, 10)
	err := faker.FakeData(&persons)
	t.Assert().NoError(err)

	personsPtr := make([]*models.Person, len(persons))
	for i, person := range persons {
		personsPtr[i] = &person
	}

	s.personRepoMock.On("GetAll").Return(personsPtr, nil)

	cases := map[string]struct {
		Persons []models.Person
		Error   error
	}{
		"success": {
			Persons: persons,
			Error:   nil,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t provider.T) {
			resPersons, err := s.uc.GetAll()
			t.Assert().ErrorIs(err, test.Error)
			t.Assert().Equal(personsPtr, resPersons)
		})
	}
}
