package image

import (
	"testing"

	"github.com/bxcodec/faker/v3"
	mock "github.com/isd-sgcu/johnjud-backend/src/mocks/image"
	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/file/image/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ImageServiceTest struct {
	suite.Suite
	petId  string
	images []*proto.Image
}

func TestImageService(t *testing.T) {
	suite.Run(t, new(ImageServiceTest))
}

func (t *ImageServiceTest) SetupTest() {
	t.petId = faker.UUIDDigit()
	t.images = []*proto.Image{
		{
			Id:       faker.UUIDDigit(),
			PetId:    t.petId,
			ImageUrl: faker.URL(),
		},
		{
			Id:       faker.UUIDDigit(),
			PetId:    t.petId,
			ImageUrl: faker.URL(),
		},
	}
}

func (t *ImageServiceTest) TestFindByPetIdSuccess() {
	want := t.images

	c := mock.ClientMock{}
	c.On("FindByPetId", &proto.FindImageByPetIdRequest{PetId: t.petId}).
		Return(&proto.FindImageByPetIdResponse{Images: t.images}, nil)

	srv := NewService(&c)
	actual, err := srv.FindByPetId(t.petId)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *ImageServiceTest) TestFindByPetIdError() {
	c := mock.ClientMock{}
	c.On("FindByPetId", &proto.FindImageByPetIdRequest{PetId: t.petId}).
		Return(&proto.FindImageByPetIdResponse{Images: t.images}, status.Error(codes.Unavailable, "Connection Timeout"))

	srv := NewService(&c)
	actual, err := srv.FindByPetId(t.petId)

	st, ok := status.FromError(err)
	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.Unavailable, st.Code())
}
