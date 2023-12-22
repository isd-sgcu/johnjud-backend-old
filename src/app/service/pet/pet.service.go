package pet

import (
	"context"
	"errors"

	"github.com/isd-sgcu/johnjud-backend/src/app/model/pet"
	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/pet/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type Service struct {
	repository IRepository
}

type IRepository interface {
	FindAll(result *[]*pet.Pet) error
	FindOne(id string, result *[]*pet.Pet) error
	Create(in *pet.Pet) error
	Update(id string, result *pet.Pet) error
	Delete(id string) error
}

func NewService(repository IRepository) *Service {
	return &Service{repository: repository}
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

func (*Service) Create(context.Context, *proto.CreatePetRequest) (*proto.CreatePetResponse, error) {
	panic("unimplemented")
}

func (*Service) FindAll(context.Context, *proto.FindAllPetRequest) (*proto.FindAllPetResponse, error) {
	panic("unimplemented")
}

func (*Service) FindOne(context.Context, *proto.FindOnePetRequest) (*proto.FindOnePetResponse, error) {
	panic("unimplemented")
}

func (*Service) Update(context.Context, *proto.UpdatePetRequest) (*proto.UpdatePetResponse, error) {
	panic("unimplemented")
}
