package pet

import (
	"context"
	"errors"
	"math/rand"
	"os/user"
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
	UpdatePet        *user.User
	PetDto           *proto.Pet
	CreatePetReqMock *proto.CreatePetRequest
	UpdatePetReqMock *proto.UpdatePetRequest
}

func TestUserService(t *testing.T) {
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
