package testBuilders

import (
	"github.com/Davmie/person_service/models"
)

type PersonBuilder struct {
	person models.Person
}

func NewPersonBuilder() *PersonBuilder {
	return &PersonBuilder{}
}

func (b *PersonBuilder) WithID(id int) *PersonBuilder {
	b.person.ID = id
	return b
}

func (b *PersonBuilder) WithName(name string) *PersonBuilder {
	b.person.Name = name
	return b
}

func (b *PersonBuilder) WithAge(age int) *PersonBuilder {
	b.person.Age = age
	return b
}

func (b *PersonBuilder) WithAddress(address string) *PersonBuilder {
	b.person.Address = address
	return b
}

func (b *PersonBuilder) WithWork(work string) *PersonBuilder {
	b.person.Work = work
	return b
}

func (b *PersonBuilder) Build() models.Person {
	return b.person
}
