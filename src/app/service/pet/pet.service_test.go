package pet

import (
	"context"
	"testing"
	"time"

	mock "github.com/isd-sgcu/johnjud-backend/src/mock/pet"

	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"github.com/isd-sgcu/johnjud-backend/src/app/model"
	"github.com/isd-sgcu/johnjud-backend/src/app/model/pet"
	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/pet/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type PetServiceTest struct {
	suite.Suite
	Pet              *pet.Pet
	Pets             []*pet.Pet
	UpdatePet        *pet.Pet
	PetDto           *proto.Pet
	CreatePetReqMock *proto.CreatePetRequest
	UpdatePetReqMock *proto.UpdatePetRequest
}

func TestPetService(t *testing.T) {
	suite.Run(t, new(PetServiceTest))
}

func (t *PetServiceTest) SetupTest() {
	t.Pet = &pet.Pet{
		Base: model.Base{
			ID:        uuid.New(),
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: gorm.DeletedAt{},
		},
		Type:         faker.Word(),
		Species:      faker.Word(),
		Name:         faker.Word(),
		Birthdate:    faker.Word(),
		Gender:       1,
		Habit:        faker.Word(),
		Caption:      faker.Word(),
		Status:       1,
		IsSterile:    false,
		IsVaccinated: false,
		IsVisible:    false,
		IsClubPet:    false,
		Background:   faker.Word(),
		Address:      faker.Word(),
		Contact:      faker.E164PhoneNumber(),
	}

	t.PetDto = &proto.Pet{
		Id:           t.Pet.ID.String(),
		Type:         t.Pet.Type,
		Species:      t.Pet.Species,
		Name:         t.Pet.Name,
		Birthdate:    t.Pet.Birthdate,
		Gender:       proto.Gender(t.Pet.Gender),
		Habit:        t.Pet.Habit,
		Caption:      t.Pet.Caption,
		Status:       proto.PetStatus(t.Pet.Status),
		ImageUrls:    []string{""},
		IsSterile:    t.Pet.IsSterile,
		IsVaccinated: t.Pet.IsVaccinated,
		IsVisible:    t.Pet.IsVaccinated,
		IsClubPet:    t.Pet.IsClubPet,
		Background:   t.Pet.Background,
		Address:      t.Pet.Address,
		Contact:      t.Pet.Contact,
	}
}

func (t *PetServiceTest) TestFindOneSuccess() {

	t.PetDto.ImageUrls = []string{""}

	want := &proto.FindOnePetResponse{Pet: t.PetDto}

	repo := &mock.RepositoryMock{}
	repo.On("FindOne", t.Pet.ID.String(), &pet.Pet{}).Return(t.Pet, nil)

	srv := NewService(repo)
	actual, err := srv.FindOne(context.Background(), &proto.FindOnePetRequest{Id: t.Pet.ID.String()})

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *PetServiceTest) TestFindAllSuccess() {
	var pets []*pet.Pet

	want := &proto.FindAllPetResponse{Pets: createPetDto(t.Pets)}

	r := mock.RepositoryMock{}
	r.On("FindAll", pets).Return(&t.Pets, nil)

	srv := NewService(&r)

	actual, err := srv.FindAll(context.Background(), &proto.FindAllPetRequest{})

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func createPetDto(in []*pet.Pet) []*proto.Pet {
	var result []*proto.Pet

	for _, p := range in {
		r := &proto.Pet{
			Id:           p.ID.String(),
			Type:         p.Type,
			Species:      p.Species,
			Name:         p.Name,
			Birthdate:    p.Birthdate,
			Gender:       proto.Gender(p.Gender),
			Habit:        p.Habit,
			Caption:      p.Caption,
			Status:       0,
			ImageUrls:    []string{""},
			IsSterile:    p.IsSterile,
			IsVaccinated: p.IsVaccinated,
			IsVisible:    p.IsVisible,
			IsClubPet:    p.IsClubPet,
			Background:   p.Background,
			Address:      p.Address,
			Contact:      p.Contact,
		}

		result = append(result, r)
	}

	return result
}
