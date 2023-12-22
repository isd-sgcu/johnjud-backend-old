package pet

import (
	"context"
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	mock "github.com/isd-sgcu/johnjud-backend/src/mocks/pet"
	"gorm.io/gorm"

	"github.com/isd-sgcu/johnjud-backend/src/app/model"
	"github.com/isd-sgcu/johnjud-backend/src/app/model/pet"
	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/pet/v1"

	petConst "github.com/isd-sgcu/johnjud-backend/src/constant/pet"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		Name:         faker.Name(),
		Birthdate:    faker.Word(),
		Gender:       petConst.Gender(rand.Intn(2) + 1),
		Habit:        faker.Paragraph(),
		Caption:      faker.Paragraph(),
		Status:       petConst.Status(rand.Intn(2) + 1),
		IsSterile:    true,
		IsVaccinated: true,
		IsVisible:    true,
		IsClubPet:    true,
		Background:   faker.Paragraph(),
		Address:      faker.Paragraph(),
		Contact:      faker.Paragraph(),
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
		IsSterile:    t.Pet.IsSterile,
		IsVaccinated: t.Pet.IsVaccinated,
		IsVisible:    t.Pet.IsVisible,
		IsClubPet:    t.Pet.IsClubPet,
		Background:   t.Pet.Background,
		Address:      t.Pet.Address,
		Contact:      t.Pet.Contact,
		ImageUrls:    []string{},
	}

	t.CreatePetReqMock = &proto.CreatePetRequest{
		Pet: &proto.Pet{
			Type:         t.Pet.Type,
			Species:      t.Pet.Species,
			Name:         t.Pet.Name,
			Birthdate:    t.Pet.Birthdate,
			Gender:       proto.Gender(t.Pet.Gender),
			Habit:        t.Pet.Habit,
			Caption:      t.Pet.Caption,
			Status:       proto.PetStatus(t.Pet.Status),
			ImageUrls:    []string{},
			IsSterile:    t.Pet.IsSterile,
			IsVaccinated: t.Pet.IsVaccinated,
			IsVisible:    t.Pet.IsVaccinated,
			IsClubPet:    t.Pet.IsClubPet,
			Background:   t.Pet.Background,
			Address:      t.Pet.Address,
			Contact:      t.Pet.Contact,
		},
	}

	t.UpdatePetReqMock = &proto.UpdatePetRequest{
		Pet: &proto.Pet{
			Id:           t.Pet.ID.String(),
			Type:         t.Pet.Type,
			Species:      t.Pet.Species,
			Name:         t.Pet.Name,
			Birthdate:    t.Pet.Birthdate,
			Gender:       proto.Gender(t.Pet.Gender),
			Habit:        t.Pet.Habit,
			Caption:      t.Pet.Caption,
			Status:       proto.PetStatus(t.Pet.Status),
			ImageUrls:    []string{},
			IsSterile:    t.Pet.IsSterile,
			IsVaccinated: t.Pet.IsVaccinated,
			IsVisible:    t.Pet.IsVisible,
			IsClubPet:    t.Pet.IsClubPet,
			Background:   t.Pet.Background,
			Address:      t.Pet.Address,
			Contact:      t.Pet.Contact,
		},
	}

	t.UpdatePet = &pet.Pet{
		Type:         t.Pet.Type,
		Species:      t.Pet.Species,
		Name:         t.Pet.Name,
		Birthdate:    t.Pet.Birthdate,
		Gender:       t.Pet.Gender,
		Habit:        t.Pet.Habit,
		Caption:      t.Pet.Caption,
		Status:       t.Pet.Status,
		IsSterile:    t.Pet.IsSterile,
		IsVaccinated: t.Pet.IsVaccinated,
		IsVisible:    t.Pet.IsVisible,
		IsClubPet:    t.Pet.IsClubPet,
		Background:   t.Pet.Background,
		Address:      t.Pet.Address,
		Contact:      t.Pet.Contact,
	}
}
func (t *PetServiceTest) TestDeleteSuccess() {
	want := &proto.DeletePetResponse{Success: true}

	repo := new(mock.RepositoryMock)
	repo.On("Delete", t.Pet.ID.String()).Return(nil)

	srv := NewService(repo)
	actual, err := srv.Delete(context.Background(), &proto.DeletePetRequest{Id: t.Pet.ID.String()})

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
	repo.AssertExpectations(t.T())
}

func (t *PetServiceTest) TestDeleteNotFound() {
	repo := new(mock.RepositoryMock)
	repo.On("Delete", t.Pet.ID.String()).Return(gorm.ErrRecordNotFound)

	srv := NewService(repo)
	_, err := srv.Delete(context.Background(), &proto.DeletePetRequest{Id: t.Pet.ID.String()})

	st, ok := status.FromError(err)
	assert.True(t.T(), ok)
	assert.Equal(t.T(), codes.NotFound, st.Code())
	repo.AssertExpectations(t.T())
}

func (t *PetServiceTest) TestDeleteWithDatabaseError() {
	repo := new(mock.RepositoryMock)
	repo.On("Delete", t.Pet.ID.String()).Return(errors.New("internal server error"))

	srv := NewService(repo)
	_, err := srv.Delete(context.Background(), &proto.DeletePetRequest{Id: t.Pet.ID.String()})

	st, ok := status.FromError(err)
	assert.True(t.T(), ok)
	assert.Equal(t.T(), codes.Internal, st.Code())
	repo.AssertExpectations(t.T())
}

func (t *PetServiceTest) TestDeleteWithUnexpectedError() {
	repo := new(mock.RepositoryMock)
	repo.On("Delete", t.Pet.ID.String()).Return(errors.New("unexpected error"))

	srv := NewService(repo)
	_, err := srv.Delete(context.Background(), &proto.DeletePetRequest{Id: t.Pet.ID.String()})

	assert.Error(t.T(), err)
	repo.AssertExpectations(t.T())
}

func (t *PetServiceTest) TestFindOneSuccess() {
	t.PetDto.ImageUrls = []string{}

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

	want := &proto.FindAllPetResponse{Pets: createPetsDto(t.Pets)}

	r := mock.RepositoryMock{}
	r.On("FindAll", pets).Return(&t.Pets, nil)

	srv := NewService(&r)

	actual, err := srv.FindAll(context.Background(), &proto.FindAllPetRequest{})

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *PetServiceTest) TestFindOneNotFound() {
	repo := &mock.RepositoryMock{}
	repo.On("FindOne", t.Pet.ID.String(), &pet.Pet{}).Return(nil, errors.New("Not found event"))

	srv := NewService(repo)
	actual, err := srv.FindOne(context.Background(), &proto.FindOnePetRequest{Id: t.Pet.ID.String()})

	st, ok := status.FromError(err)

	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.NotFound, st.Code())
}

func createPetsDto(in []*pet.Pet) []*proto.Pet {
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
			ImageUrls:    []string{},
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

func (t *PetServiceTest) TestCreateSuccess() {
	want := &proto.CreatePetResponse{Pet: t.PetDto}

	repo := &mock.RepositoryMock{}

	in := &pet.Pet{
		Type:         t.Pet.Type,
		Species:      t.Pet.Species,
		Name:         t.Pet.Name,
		Birthdate:    t.Pet.Birthdate,
		Gender:       t.Pet.Gender,
		Habit:        t.Pet.Habit,
		Caption:      t.Pet.Caption,
		Status:       t.Pet.Status,
		IsSterile:    t.Pet.IsSterile,
		IsVaccinated: t.Pet.IsVaccinated,
		IsVisible:    t.Pet.IsVisible,
		IsClubPet:    t.Pet.IsClubPet,
		Background:   t.Pet.Background,
		Address:      t.Pet.Address,
		Contact:      t.Pet.Contact,
	}

	repo.On("Create", in).Return(t.Pet, nil)
	srv := NewService(repo)

	actual, err := srv.Create(context.Background(), t.CreatePetReqMock)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *PetServiceTest) TestCreateInternalErr() {
	repo := &mock.RepositoryMock{}

	in := &pet.Pet{
		Type:         t.Pet.Type,
		Species:      t.Pet.Species,
		Name:         t.Pet.Name,
		Birthdate:    t.Pet.Birthdate,
		Gender:       t.Pet.Gender,
		Habit:        t.Pet.Habit,
		Caption:      t.Pet.Caption,
		Status:       t.Pet.Status,
		IsSterile:    t.Pet.IsSterile,
		IsVaccinated: t.Pet.IsVaccinated,
		IsVisible:    t.Pet.IsVisible,
		IsClubPet:    t.Pet.IsClubPet,
		Background:   t.Pet.Background,
		Address:      t.Pet.Address,
		Contact:      t.Pet.Contact,
	}

	repo.On("Create", in).Return(nil, errors.New("something wrong"))
	srv := NewService(repo)

	actual, err := srv.Create(context.Background(), t.CreatePetReqMock)

	st, ok := status.FromError(err)

	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.Internal, st.Code())
}
