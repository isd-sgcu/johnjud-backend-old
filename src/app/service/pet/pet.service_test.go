package pet

import (
	"context"
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	img_mock "github.com/isd-sgcu/johnjud-backend/src/mocks/image"
	mock "github.com/isd-sgcu/johnjud-backend/src/mocks/pet"
	"gorm.io/gorm"

	"github.com/isd-sgcu/johnjud-backend/src/app/model"
	"github.com/isd-sgcu/johnjud-backend/src/app/model/pet"
	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/pet/v1"
	img_proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/file/image/v1"

	petConst "github.com/isd-sgcu/johnjud-backend/src/constant/pet"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PetServiceTest struct {
	suite.Suite
	Pet                  *pet.Pet
	UpdatePet            *pet.Pet
	ChangeViewPet        *pet.Pet
	Pets                 []*pet.Pet
	PetDto               *proto.Pet
	CreatePetReqMock     *proto.CreatePetRequest
	UpdatePetReqMock     *proto.UpdatePetRequest
	ChangeViewPetReqMock *proto.ChangeViewPetRequest
	Images               []*img_proto.Image
	ImageUrls            []string
	ChangeAdoptBy        *pet.Pet
	AdoptByReq           *proto.AdoptPetRequest
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
		Gender:       petConst.Gender(rand.Intn(1) + 1),
		Habit:        faker.Paragraph(),
		Caption:      faker.Paragraph(),
		Status:       petConst.Status(rand.Intn(1) + 1),
		IsSterile:    true,
		IsVaccinated: true,
		IsVisible:    true,
		IsClubPet:    true,
		Background:   faker.Paragraph(),
		Address:      faker.Paragraph(),
		Contact:      faker.Paragraph(),
		AdoptBy:      "",
	}

	t.Images = []*img_proto.Image{
		{
			Id:       faker.UUIDDigit(),
			PetId:    t.Pet.ID.String(),
			ImageUrl: faker.URL(),
		},
		{
			Id:       faker.UUIDDigit(),
			PetId:    t.Pet.ID.String(),
			ImageUrl: faker.URL(),
		},
	}

	t.ImageUrls = []string{
		t.Images[0].ImageUrl,
		t.Images[1].ImageUrl,
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
		ImageUrls:    t.ImageUrls,
	}

	t.UpdatePet = &pet.Pet{
		Base: model.Base{
			ID:        t.Pet.Base.ID,
			CreatedAt: t.Pet.Base.CreatedAt,
			UpdatedAt: t.Pet.Base.UpdatedAt,
			DeletedAt: t.Pet.Base.DeletedAt,
		},
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

	t.ChangeViewPet = &pet.Pet{
		Base: model.Base{
			ID:        t.Pet.Base.ID,
			CreatedAt: t.Pet.Base.CreatedAt,
			UpdatedAt: t.Pet.Base.UpdatedAt,
			DeletedAt: t.Pet.Base.DeletedAt,
		},
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
		IsVisible:    false,
		IsClubPet:    t.Pet.IsClubPet,
		Background:   t.Pet.Background,
		Address:      t.Pet.Address,
		Contact:      t.Pet.Contact,
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
			ImageUrls:    t.ImageUrls,
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
			ImageUrls:    t.ImageUrls,
			IsSterile:    t.Pet.IsSterile,
			IsVaccinated: t.Pet.IsVaccinated,
			IsVisible:    t.Pet.IsVisible,
			IsClubPet:    t.Pet.IsClubPet,
			Background:   t.Pet.Background,
			Address:      t.Pet.Address,
			Contact:      t.Pet.Contact,
		},
	}

	t.ChangeViewPetReqMock = &proto.ChangeViewPetRequest{
		Id:      t.Pet.ID.String(),
		Visible: false,
	}

	t.ChangeAdoptBy = &pet.Pet{
		Base: model.Base{
			ID:        t.Pet.Base.ID,
			CreatedAt: t.Pet.Base.CreatedAt,
			UpdatedAt: t.Pet.Base.UpdatedAt,
			DeletedAt: t.Pet.Base.DeletedAt,
		},
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
		AdoptBy:      faker.UUIDDigit(),
	}

	t.AdoptByReq = &proto.AdoptPetRequest{
		PetId:  t.ChangeAdoptBy.ID.String(),
		UserId: t.ChangeAdoptBy.AdoptBy,
	}
}
func (t *PetServiceTest) TestDeleteSuccess() {
	want := &proto.DeletePetResponse{Success: true}

	repo := new(mock.RepositoryMock)
	repo.On("Delete", t.Pet.ID.String()).Return(nil)
	imgSrv := new(img_mock.ServiceMock)

	srv := NewService(repo, imgSrv)
	actual, err := srv.Delete(context.Background(), &proto.DeletePetRequest{Id: t.Pet.ID.String()})

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
	repo.AssertExpectations(t.T())
}

func (t *PetServiceTest) TestDeleteNotFound() {
	repo := new(mock.RepositoryMock)
	repo.On("Delete", t.Pet.ID.String()).Return(gorm.ErrRecordNotFound)
	imgSrv := new(img_mock.ServiceMock)

	srv := NewService(repo, imgSrv)
	_, err := srv.Delete(context.Background(), &proto.DeletePetRequest{Id: t.Pet.ID.String()})

	st, ok := status.FromError(err)
	assert.True(t.T(), ok)
	assert.Equal(t.T(), codes.NotFound, st.Code())
	repo.AssertExpectations(t.T())
}

func (t *PetServiceTest) TestDeleteWithDatabaseError() {
	repo := new(mock.RepositoryMock)
	repo.On("Delete", t.Pet.ID.String()).Return(errors.New("internal server error"))
	imgSrv := new(img_mock.ServiceMock)

	srv := NewService(repo, imgSrv)
	_, err := srv.Delete(context.Background(), &proto.DeletePetRequest{Id: t.Pet.ID.String()})

	st, ok := status.FromError(err)
	assert.True(t.T(), ok)
	assert.Equal(t.T(), codes.Internal, st.Code())
	repo.AssertExpectations(t.T())
}

func (t *PetServiceTest) TestDeleteWithUnexpectedError() {
	repo := new(mock.RepositoryMock)
	repo.On("Delete", t.Pet.ID.String()).Return(errors.New("unexpected error"))
	imgSrv := new(img_mock.ServiceMock)

	srv := NewService(repo, imgSrv)
	_, err := srv.Delete(context.Background(), &proto.DeletePetRequest{Id: t.Pet.ID.String()})

	assert.Error(t.T(), err)
	repo.AssertExpectations(t.T())
}

func (t *PetServiceTest) TestFindOneSuccess() {
	t.PetDto.ImageUrls = []string{}

	want := &proto.FindOnePetResponse{Pet: t.PetDto}

	repo := &mock.RepositoryMock{}
	repo.On("FindOne", t.Pet.ID.String(), &pet.Pet{}).Return(t.Pet, nil)
	imgSrv := new(img_mock.ServiceMock)
	imgSrv.On("FindByPetId", t.Pet.ID.String()).Return(t.Images, nil)

	srv := NewService(repo, imgSrv)
	actual, err := srv.FindOne(context.Background(), &proto.FindOnePetRequest{Id: t.Pet.ID.String()})

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *PetServiceTest) TestFindAllSuccess() {
	var pets []*pet.Pet

	want := &proto.FindAllPetResponse{Pets: createPetsDto(t.Pets)}

	repo := &mock.RepositoryMock{}
	repo.On("FindAll", pets).Return(&t.Pets, nil)
	imgSrv := new(img_mock.ServiceMock)
	imgSrv.On("FindByPetId", t.Pet.ID.String()).Return(t.Images, nil)

	srv := NewService(repo, imgSrv)

	actual, err := srv.FindAll(context.Background(), &proto.FindAllPetRequest{})

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *PetServiceTest) TestFindOneNotFound() {
	repo := &mock.RepositoryMock{}
	repo.On("FindOne", t.Pet.ID.String(), &pet.Pet{}).Return(nil, errors.New("Not found pet"))
	imgSrv := new(img_mock.ServiceMock)
	imgSrv.On("FindByPetId", t.Pet.ID.String()).Return(nil, nil)

	srv := NewService(repo, imgSrv)
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
	want.Pet.ImageUrls = []string{} // when pet is first created, it has no images

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
	imgSrv := new(img_mock.ServiceMock)

	srv := NewService(repo, imgSrv)

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
	imgSrv := new(img_mock.ServiceMock)

	srv := NewService(repo, imgSrv)

	actual, err := srv.Create(context.Background(), t.CreatePetReqMock)

	st, ok := status.FromError(err)

	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.Internal, st.Code())
}

func (t *PetServiceTest) TestUpdateSuccess() {
	want := &proto.UpdatePetResponse{Pet: t.PetDto}

	repo := &mock.RepositoryMock{}
	repo.On("Update", t.Pet.ID.String(), t.UpdatePet).Return(t.Pet, nil)
	imgSrv := new(img_mock.ServiceMock)
	imgSrv.On("FindByPetId", t.Pet.ID.String()).Return(t.Images, nil)

	srv := NewService(repo, imgSrv)
	actual, err := srv.Update(context.Background(), t.UpdatePetReqMock)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *PetServiceTest) TestUpdateNotFound() {
	repo := &mock.RepositoryMock{}
	repo.On("Update", t.Pet.ID.String(), t.UpdatePet).Return(nil, errors.New("Not found pet"))
	imgSrv := new(img_mock.ServiceMock)
	imgSrv.On("FindByPetId", t.Pet.ID.String()).Return(t.Images, nil)

	srv := NewService(repo, imgSrv)
	actual, err := srv.Update(context.Background(), t.UpdatePetReqMock)

	st, ok := status.FromError(err)

	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.NotFound, st.Code())
}

func (t *PetServiceTest) TestChangeViewSuccess() {
	want := &proto.ChangeViewPetResponse{Success: true}

	repo := &mock.RepositoryMock{}
	repo.On("FindOne", t.Pet.ID.String(), &pet.Pet{}).Return(t.Pet, nil)
	repo.On("Update", t.Pet.ID.String(), t.ChangeViewPet).Return(t.ChangeViewPet, nil)
	imgSrv := new(img_mock.ServiceMock)

	srv := NewService(repo, imgSrv)
	actual, err := srv.ChangeView(context.Background(), t.ChangeViewPetReqMock)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *PetServiceTest) TestChangeViewNotFound() {
	repo := &mock.RepositoryMock{}
	repo.On("FindOne", t.Pet.ID.String(), &pet.Pet{}).Return(nil, errors.New("Not found pet"))
	repo.On("Update", t.Pet.ID.String(), t.UpdatePet).Return(nil, errors.New("Not found pet"))
	imgSrv := new(img_mock.ServiceMock)

	srv := NewService(repo, imgSrv)
	actual, err := srv.ChangeView(context.Background(), t.ChangeViewPetReqMock)

	st, ok := status.FromError(err)

	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.NotFound, st.Code())
}

func (t *PetServiceTest) TestAdoptBySuccess() {
	want := &proto.AdoptPetResponse{Success: true}
	repo := &mock.RepositoryMock{}

	repo.On("FindOne", t.AdoptByReq.PetId, &pet.Pet{}).Return(t.Pet, nil)
	repo.On("Update", t.AdoptByReq.PetId, t.ChangeAdoptBy).Return(t.ChangeAdoptBy, nil)

	imgSrv := new(img_mock.ServiceMock)
	srv := NewService(repo, imgSrv)

	actual, err := srv.AdoptPet(context.Background(), t.AdoptByReq)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *PetServiceTest) TestAdoptByPetNotFound() {
	wantError := status.Error(codes.NotFound, "pet not found")
	repo := &mock.RepositoryMock{}

	repo.On("FindOne", t.AdoptByReq.PetId, &pet.Pet{}).Return(nil, wantError)

	imgSrv := new(img_mock.ServiceMock)
	srv := NewService(repo, imgSrv)

	actual, err := srv.AdoptPet(context.Background(), t.AdoptByReq)

	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), wantError, err)
	assert.Nil(t.T(), actual)

	repo.AssertNotCalled(t.T(), "Update", t.AdoptByReq.PetId, t.ChangeAdoptBy)
}

func (t *PetServiceTest) TestAdoptByUpdateError() {
	wantError := status.Error(codes.NotFound, "pet not found")
	repo := &mock.RepositoryMock{}

	repo.On("FindOne", t.AdoptByReq.PetId, &pet.Pet{}).Return(t.Pet, nil)
	repo.On("Update", t.AdoptByReq.PetId, t.ChangeAdoptBy).Return(nil, errors.New("update error"))

	imgSrv := new(img_mock.ServiceMock)
	srv := NewService(repo, imgSrv)

	actual, err := srv.AdoptPet(context.Background(), t.AdoptByReq)

	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), wantError, err)
	assert.Nil(t.T(), actual)
}
