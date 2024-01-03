package pet

import (
	"context"
	"errors"

	"time"

	"github.com/google/uuid"
	"github.com/isd-sgcu/johnjud-backend/src/app/model"
	"github.com/isd-sgcu/johnjud-backend/src/app/model/pet"
	petConst "github.com/isd-sgcu/johnjud-backend/src/constant/pet"
	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/pet/v1"
	image_proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/file/image/v1"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type Service struct {
	proto.UnimplementedPetServiceServer
	repository   IRepository
	imageService ImageService
}

type IRepository interface {
	FindAll(result *[]*pet.Pet) error
	FindOne(id string, result *pet.Pet) error
	Create(in *pet.Pet) error
	Update(id string, result *pet.Pet) error
	Delete(id string) error
}

type ImageService interface {
	FindByPetId(petId string) ([]*image_proto.Image, error)
}

func NewService(repository IRepository, imageService ImageService) *Service {
	return &Service{repository: repository, imageService: imageService}
}

func (s *Service) Delete(ctx context.Context, req *proto.DeletePetRequest) (*proto.DeletePetResponse, error) {
	err := s.repository.Delete(req.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "pet not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &proto.DeletePetResponse{Success: true}, nil
}

func (s *Service) Update(_ context.Context, req *proto.UpdatePetRequest) (res *proto.UpdatePetResponse, err error) {
	raw, err := DtoToRaw(req.Pet)
	if err != nil {
		return nil, status.Error(codes.Internal, "error converting dto to raw")
	}

	err = s.repository.Update(req.Pet.Id, raw)
	if err != nil {
		return nil, status.Error(codes.NotFound, "pet not found")
	}

	images, err := s.imageService.FindByPetId(req.Pet.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "error querying image service")
	}
	imageUrls := ExtractImageUrls(images)

	return &proto.UpdatePetResponse{Pet: RawToDto(raw, imageUrls)}, nil
}

func (s *Service) ChangeView(_ context.Context, req *proto.ChangeViewPetRequest) (res *proto.ChangeViewPetResponse, err error) {
	petData, err := s.FindOne(context.Background(), &proto.FindOnePetRequest{Id: req.Id})
	if err != nil {
		return nil, status.Error(codes.NotFound, "pet not found")
	}
	pet, err := DtoToRaw(petData.Pet)
	if err != nil {
		return nil, status.Error(codes.Internal, "error converting dto to raw")
	}
	pet.IsVisible = req.Visible

	err = s.repository.Update(req.Id, pet)
	if err != nil {
		return nil, status.Error(codes.NotFound, "pet not found")
	}

	return &proto.ChangeViewPetResponse{Success: true}, nil
}

func (s *Service) FindAll(_ context.Context, req *proto.FindAllPetRequest) (res *proto.FindAllPetResponse, err error) {
	var pets []*pet.Pet
	var imageUrlsList [][]string

	err = s.repository.FindAll(&pets)
	if err != nil {
		log.Error().Err(err).Str("service", "event").Str("module", "find all").Msg("Error while querying all events")
		return nil, status.Error(codes.Unavailable, "Internal error")
	}

	for _, pet := range pets {
		images, err := s.imageService.FindByPetId(pet.ID.String())
		if err != nil {
			return nil, status.Error(codes.Internal, "error querying image service")
		}
		imageUrls := ExtractImageUrls(images)
		imageUrlsList = append(imageUrlsList, imageUrls)
	}

	petWithImageUrls, err := RawToDtoList(&pets, imageUrlsList)
	if err != nil {
		return nil, status.Error(codes.Internal, "error converting raw to dto list")
	}

	return &proto.FindAllPetResponse{Pets: petWithImageUrls}, nil
}

func (s Service) FindOne(_ context.Context, req *proto.FindOnePetRequest) (res *proto.FindOnePetResponse, err error) {
	var pet pet.Pet

	err = s.repository.FindOne(req.Id, &pet)
	if err != nil {
		log.Error().Err(err).
			Str("service", "pet").Str("module", "find one").Str("id", req.Id).Msg("Not found")
		return nil, status.Error(codes.NotFound, err.Error())
	}

	images, err := s.imageService.FindByPetId(req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "error querying image service")
	}
	imageUrls := ExtractImageUrls(images)

	return &proto.FindOnePetResponse{Pet: RawToDto(&pet, imageUrls)}, err
}

func (s *Service) Create(_ context.Context, req *proto.CreatePetRequest) (res *proto.CreatePetResponse, err error) {
	raw, err := DtoToRaw(req.Pet)
	if err != nil {
		return nil, status.Error(codes.Internal, "error converting dto to raw: "+err.Error())
	}

	imgUrls := []string{}

	err = s.repository.Create(raw)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create pet")
	}

	return &proto.CreatePetResponse{Pet: RawToDto(raw, imgUrls)}, nil
}

func (s *Service) AdoptPet(ctx context.Context, req *proto.AdoptPetRequest) (res *proto.AdoptPetResponse, err error) {
	dtoPet, err := s.FindOne(context.Background(), &proto.FindOnePetRequest{Id: req.PetId})
	if err != nil {
		return nil, status.Error(codes.NotFound, "pet not found")
	}
	pet, err := DtoToRaw(dtoPet.Pet)
	if err != nil {
		return nil, status.Error(codes.Internal, "error converting dto to raw")
	}
	pet.AdoptBy = req.UserId

	err = s.repository.Update(req.PetId, pet)
	if err != nil {
		return nil, status.Error(codes.NotFound, "pet not found")
	}

	return &proto.AdoptPetResponse{Success: true}, nil
}

func RawToDtoList(in *[]*pet.Pet, imageUrls [][]string) ([]*proto.Pet, error) {
	var result []*proto.Pet
	if len(*in) != len(imageUrls) {
		return nil, errors.New("length of in and imageUrls have to be the same")
	}

	for i, e := range *in {
		result = append(result, RawToDto(e, imageUrls[i]))
	}
	return result, nil
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
		AdoptBy:      in.AdoptBy,
	}
}

func DtoToRaw(in *proto.Pet) (res *pet.Pet, err error) {
	var id uuid.UUID
	var gender petConst.Gender
	var status petConst.Status

	if in.Id != "" {
		id, err = uuid.Parse(in.Id)
		if err != nil {
			return nil, err
		}
	}

	switch in.Gender {
	case 0:
		gender = petConst.MALE
	case 1:
		gender = petConst.FEMALE
	}

	switch in.Status {
	case 0:
		status = petConst.ADOPTED
	case 1:
		status = petConst.FINDHOME
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
		AdoptBy:      in.AdoptBy,
	}, nil
}

func ExtractImageUrls(in []*image_proto.Image) []string {
	var result []string
	for _, e := range in {
		result = append(result, e.ImageUrl)
	}
	return result
}
