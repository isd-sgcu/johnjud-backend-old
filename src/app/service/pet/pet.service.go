package pet

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/isd-sgcu/johnjud-backend/src/app/model"
	"github.com/isd-sgcu/johnjud-backend/src/app/model/pet"
	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/pet/v1"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type Service struct {
	repository IRepository
}

type IRepository interface {
	FindAll() error
	FindOne(string, *pet.Pet) error
	Create(*pet.Pet) error
}

func NewService(repository IRepository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) FindAll(_ context.Context) error {
	return nil
}

func (s Service) FindOne(_ context.Context, req *proto.FindOnePetRequest) (res *proto.FindOnePetResponse, err error) {
	var pet *pet.Pet

	err = s.repository.FindOne(req.Id, pet)
	if err != nil {
		log.Error().Err(err).
			Str("service", "like").Str("module", "find one").Str("id", req.Id).Msg("Not found")
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &proto.FindOnePetResponse{Pet: RawToDto(pet)}, err
}

func (s *Service) Create(_ context.Context, req *proto.CreatePetRequest) (res *proto.CreatePetResponse, err error) {
	raw, _ := DtoToRaw(req.Pet)
	imgUrl := []string{}

	err = s.repository.Create(raw)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create pet")
	}

	return &proto.CreatePetResponse{Pet: RawToDto(raw, imgUrl)}, nil
}

func RawToDto(in *pet.Pet, imgUrl []string) *proto.Pet {
	return &proto.Pet{
		Id:           in.ID.String(),
		Type:         in.Type,
		Species:      in.Species,
		Name:         in.Name,
		Birthdate:    in.Birthdate,
		Gender:       proto.Gender(in.Gender),
		Habit:        in.Habit,
		Caption:      in.Caption,
		Status:       proto.PetStatus(in.Status),
		ImageUrls:    imgUrl,
		IsSterile:    in.IsSterile,
		IsVaccinated: in.IsVaccinated,
		IsVisible:    in.IsVisible,
		IsClubPet:    in.IsClubPet,
		Background:   in.Background,
		Address:      in.Address,
		Contact:      in.Contact,
	}
}

func DtoToRaw(in *proto.Pet) (res *pet.Pet, err error) {
	var id uuid.UUID
	if in.Id != "" {
		id, err = uuid.Parse(in.Id)
		if err != nil {
			return nil, err
		}
	}

	id, err = uuid.Parse(in.Id)
	if err != nil {
		return nil, err
	}

	return &pet.Pet{
		Base: model.Base{
			ID:        id,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: gorm.DeletedAt{},
		},
		Type:         in.Type,
		Species:      in.Species,
		Name:         in.Name,
		Birthdate:    in.Birthdate,
		Gender:       "",
		Habit:        in.Habit,
		Caption:      in.Caption,
		Status:       "",
		IsSterile:    in.IsSterile,
		IsVaccinated: in.IsVaccinated,
		IsVisible:    in.IsVisible,
		IsClubPet:    in.IsClubPet,
		Background:   in.Background,
		Address:      in.Address,
		Contact:      in.Contact,
	}, nil
}
