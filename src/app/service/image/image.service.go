package image

import (
	"context"
	"time"

	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/file/image/v1"
	"github.com/rs/zerolog/log"
)

type Service struct {
	client proto.ImageServiceClient
}

func NewService(client proto.ImageServiceClient) *Service {
	return &Service{client: client}
}

func (s *Service) FindByPetId(petId string) ([]*proto.Image, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := s.client.FindByPetId(ctx, &proto.FindImageByPetIdRequest{PetId: petId})
	if err != nil {
		log.Error().
			Err(err).
			Str("service", "image").
			Str("module", "find by petId").
			Msg("Error while connecting to service")
		return nil, err
	}
	return res.Images, nil

}
