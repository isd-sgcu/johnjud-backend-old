package pet

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/isd-sgcu/johnjud-backend/src/app/model"
	"github.com/isd-sgcu/johnjud-backend/src/app/model/pet"
	cst "github.com/isd-sgcu/johnjud-backend/src/constant/pet"
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
	FindAll(*[]*pet.Pet) error
	FindOne(string, *pet.Pet) error
	Create(*pet.Pet) error
}

func NewService(repository IRepository) *Service {
	return &Service{repository}
}

func (s *Service) FindAll(_ context.Context, req *proto.FindAllPetRequest) (res *proto.FindAllPetResponse, err error) {
	var pets []*pet.Pet

	err = s.repository.FindAll(&pets)
	if err != nil {
		log.Error().Err(err).Str("service", "event").Str("module", "find all").Msg("Error while querying all events")
		return nil, status.Error(codes.Unavailable, "Internal error")
	}
	return &proto.FindAllPetResponse{Pets: RawToDtoList(&pets)}, nil
}

func (s Service) FindOne(_ context.Context, req *proto.FindOnePetRequest) (res *proto.FindOnePetResponse, err error) {
	var pet pet.Pet

	err = s.repository.FindOne(req.Id, &pet)
	if err != nil {
		log.Error().Err(err).
			Str("service", "like").Str("module", "find one").Str("id", req.Id).Msg("Not found")
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &proto.FindOnePetResponse{Pet: RawToDto(&pet, []string{""})}, nil
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

func RawToDtoList(in *[]*pet.Pet) []*proto.Pet {
	var result []*proto.Pet
	for _, e := range *in {
		result = append(result, RawToDto(e, []string{""}))
	}
	return result
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
	var gender cst.Gender
	var status cst.Status

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

	switch in.Gender {
	case 1:
		gender = cst.MALE
	case 2:
		gender = cst.FEMALE
	}

	switch in.Status {
	case 1:
		status = cst.ADOPTED
	case 2:
		status = cst.FINDHOME
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
		Gender:       gender,
		Habit:        in.Habit,
		Caption:      in.Caption,
		Status:       status,
		IsSterile:    in.IsSterile,
		IsVaccinated: in.IsVaccinated,
		IsVisible:    in.IsVisible,
		IsClubPet:    in.IsClubPet,
		Background:   in.Background,
		Address:      in.Address,
		Contact:      in.Contact,
	}, nil
}
