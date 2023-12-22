package pet

import (
	"github.com/isd-sgcu/johnjud-backend/src/app/model/pet"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) FindAll(result *[]*pet.Pet) error {
	// unimplemented for now
	return nil
}

func (r *RepositoryMock) FindOne(id string, result *[]*pet.Pet) error {
	// unimplemented for now
	return nil
}

func (r *RepositoryMock) Create(in *pet.Pet) error {
	// unimplemented for now
	return nil
}

func (r *RepositoryMock) Update(id string, result *pet.Pet) error {
	// unimplemented for now
	return nil
}

func (r *RepositoryMock) Delete(id string) error {
	args := r.Called(id)
	return args.Error(0)
}
