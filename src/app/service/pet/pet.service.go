package pet

import (
	"context"
	"errors"
	"fmt"

	"github.com/isd-sgcu/johnjud-backend/src/app/model/pet"
	petUtils "github.com/isd-sgcu/johnjud-backend/src/app/utils/pet"
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
	FindAll(*[]*pet.Pet, bool) error
	FindOne(string, *pet.Pet) error
	Create(*pet.Pet) error
	Update(string, *pet.Pet) error
	Delete(string) error
}

type ImageService interface {
	FindByPetId(petId string) ([]*image_proto.Image, error)
	AssignPet(petId string, imageIds []string) error
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
	raw, err := petUtils.DtoToRaw(req.Pet)
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

	return &proto.UpdatePetResponse{Pet: petUtils.RawToDto(raw, images)}, nil
}

func (s *Service) ChangeView(_ context.Context, req *proto.ChangeViewPetRequest) (res *proto.ChangeViewPetResponse, err error) {
	petData, err := s.FindOne(context.Background(), &proto.FindOnePetRequest{Id: req.Id})
	if err != nil {
		return nil, status.Error(codes.NotFound, "pet not found")
	}
	pet, err := petUtils.DtoToRaw(petData.Pet)
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
	var imagesList [][]*image_proto.Image
	metaData := proto.FindAllPetMetaData{}

	err = s.repository.FindAll(&pets, req.IsAdmin)
	if err != nil {
		log.Error().Err(err).Str("service", "event").Str("module", "find all").Msg("Error while querying all events")
		return nil, status.Error(codes.Unavailable, "Internal error")
	}

	petUtils.FilterPet(&pets, req)
	petUtils.PaginatePets(&pets, req.Page, req.PageSize, &metaData)

	for _, pet := range pets {
		images, err := s.imageService.FindByPetId(pet.ID.String())
		if err != nil {
			return nil, status.Error(codes.Internal, "error querying image service")
		}
		imagesList = append(imagesList, images)
	}
	petWithImages, err := petUtils.RawToDtoList(&pets, imagesList, req)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("error converting raw to dto list: %v", err))
	}
	return &proto.FindAllPetResponse{Pets: petWithImages, Metadata: &metaData}, nil
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

	return &proto.FindOnePetResponse{Pet: petUtils.RawToDto(&pet, images)}, err
}

func (s *Service) Create(_ context.Context, req *proto.CreatePetRequest) (res *proto.CreatePetResponse, err error) {
	raw, err := petUtils.DtoToRaw(req.Pet)
	if err != nil {
		return nil, status.Error(codes.Internal, "error converting dto to raw: "+err.Error())
	}

	err = s.repository.Create(raw)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create pet")
	}

	imageIds := petUtils.ExtractImageIDs(req.Pet.Images)
	err = s.imageService.AssignPet(raw.ID.String(), imageIds)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to assign pet to images")
	}

	images, err := s.imageService.FindByPetId(raw.ID.String())
	if err != nil {
		return nil, status.Error(codes.Internal, "error finding images by pet id")
	}

	return &proto.CreatePetResponse{Pet: petUtils.RawToDto(raw, images)}, nil
}

func (s *Service) AdoptPet(ctx context.Context, req *proto.AdoptPetRequest) (res *proto.AdoptPetResponse, err error) {
	dtoPet, err := s.FindOne(context.Background(), &proto.FindOnePetRequest{Id: req.PetId})
	if err != nil {
		return nil, status.Error(codes.NotFound, "pet not found")
	}
	pet, err := petUtils.DtoToRaw(dtoPet.Pet)
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
