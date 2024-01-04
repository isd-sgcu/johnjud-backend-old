package like

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/isd-sgcu/johnjud-backend/src/app/model"
	"github.com/isd-sgcu/johnjud-backend/src/app/model/like"
	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/like/v1"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type Service struct {
	proto.UnimplementedLikeServiceServer
	repository IRepository
}

type IRepository interface {
	FindByUserId(string, *[]*like.Like) error
	Create(*like.Like) error
	Delete(string) error
}

func NewService(repository IRepository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) FindByUserId(_ context.Context, req *proto.FindLikeByUserIdRequest) (res *proto.FindLikeByUserIdResponse, err error) {
	var likes []*like.Like

	err = s.repository.FindByUserId(req.UserId, &likes)
	if err != nil {

		log.Error().Err(err).
			Str("service", "like").
			Str("module", "find by userId").
			Str("useId", req.UserId).
			Msg("Not found")

		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &proto.FindLikeByUserIdResponse{Likes: RawToDtoList(&likes)}, nil
}

func (s *Service) Create(_ context.Context, req *proto.CreateLikeRequest) (res *proto.CreateLikeResponse, err error) {
	raw, _ := DtoToRaw(req.Like)

	err = s.repository.Create(raw)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create like")
	}

	return &proto.CreateLikeResponse{Like: RawToDto(raw)}, nil
}

func (s *Service) Delete(_ context.Context, req *proto.DeleteLikeRequest) (res *proto.DeleteLikeResponse, err error) {
	err = s.repository.Delete(req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "something wrong when deleting like")
	}

	return &proto.DeleteLikeResponse{Success: true}, nil
}

func DtoToRaw(in *proto.Like) (result *like.Like, err error) {
	var id uuid.UUID
	if in.Id != "" {
		id, err = uuid.Parse(in.Id)
		if err != nil {
			return nil, err
		}
	}

	petId, err := uuid.Parse(in.PetId)
	if err != nil {
		return nil, err
	}
	userId, err := uuid.Parse(in.UserId)
	if err != nil {
		return nil, err
	}

	return &like.Like{
		Base: model.Base{
			ID:        id,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: gorm.DeletedAt{},
		},
		PetID:  &petId,
		UserID: &userId,
	}, nil
}

func RawToDtoList(in *[]*like.Like) []*proto.Like {
	var result []*proto.Like
	for _, b := range *in {
		result = append(result, RawToDto(b))
	}

	return result
}

func RawToDto(in *like.Like) *proto.Like {
	return &proto.Like{
		Id:     in.ID.String(),
		PetId:  in.PetID.String(),
		UserId: in.UserID.String(),
	}
}
