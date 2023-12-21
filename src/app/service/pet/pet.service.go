package pet

import (
	"context"
	"errors"

	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/pet/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type Service struct {
	repo IRepository
}

type IRepository interface {
	FindAll() error
	FindOne() error
	Create() error
	Update() error
	ChangeView() error
	Delete(string) error
}

func NewService(repo IRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Delete(ctx context.Context, req *proto.DeletePetRequest) (*proto.DeletePetResponse, error) {
	err := s.repo.Delete(req.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "pet not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &proto.DeletePetResponse{Success: true}, nil
}

// ChangeView implements v1.PetServiceServer.
func (*Service) ChangeView(context.Context, *proto.ChangeViewPetRequest) (*proto.ChangeViewPetResponse, error) {
	panic("unimplemented")
}

// Create implements v1.PetServiceServer.
func (*Service) Create(context.Context, *proto.CreatePetRequest) (*proto.CreatePetResponse, error) {
	panic("unimplemented")
}

// FindAll implements v1.PetServiceServer.
func (*Service) FindAll(context.Context, *proto.FindAllPetRequest) (*proto.FindAllPetResponse, error) {
	panic("unimplemented")
}

// FindOne implements v1.PetServiceServer.
func (*Service) FindOne(context.Context, *proto.FindOnePetRequest) (*proto.FindOnePetResponse, error) {
	panic("unimplemented")
}

// Update implements v1.PetServiceServer.
func (*Service) Update(context.Context, *proto.UpdatePetRequest) (*proto.UpdatePetResponse, error) {
	panic("unimplemented")
}
