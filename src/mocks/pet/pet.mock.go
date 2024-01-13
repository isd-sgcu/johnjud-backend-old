package pet

import (
	"github.com/isd-sgcu/johnjud-backend/src/app/model/pet"
	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/pet/v1"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) FindOne(id string, result *pet.Pet) error {
	args := r.Called(id, result)

	if args.Get(0) != nil {
		*result = *args.Get(0).(*pet.Pet)
	}

	return args.Error(1)
}

func (r *RepositoryMock) Create(in *pet.Pet) error {
	args := r.Called(in)

	if args.Get(0) != nil {
		*in = *args.Get(0).(*pet.Pet)
	}

	return args.Error(1)
}

func (r *RepositoryMock) FindAll(result *[]*pet.Pet, query *proto.FindAllPetRequest) error {
	args := r.Called(*result, query)

	if args.Get(0) != nil {
		*result = *args.Get(0).(*[]*pet.Pet)
	}

	return args.Error(1)
}

func (r *RepositoryMock) Update(id string, result *pet.Pet) error {
	args := r.Called(id, result)

	if args.Get(0) != nil {
		*result = *args.Get(0).(*pet.Pet)
	}

	return args.Error(1)
}

func (r *RepositoryMock) Delete(id string) error {
	args := r.Called(id)
	return args.Error(0)
}
