package image

import (
	"context"

	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/file/image/v1"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type ClientMock struct {
	mock.Mock
}

func (c *ClientMock) Upload(_ context.Context, in *proto.UploadImageRequest, _ ...grpc.CallOption) (res *proto.UploadImageResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.UploadImageResponse)
	}

	return res, args.Error(1)
}

func (c *ClientMock) FindAll(_ context.Context, in *proto.FindAllImageRequest, _ ...grpc.CallOption) (res *proto.FindAllImageResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.FindAllImageResponse)
	}

	return res, args.Error(1)
}

func (c *ClientMock) FindByPetId(_ context.Context, in *proto.FindImageByPetIdRequest, _ ...grpc.CallOption) (res *proto.FindImageByPetIdResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.FindImageByPetIdResponse)
	}

	return res, args.Error(1)
}

func (c *ClientMock) AssignPet(_ context.Context, in *proto.AssignPetRequest, _ ...grpc.CallOption) (res *proto.AssignPetResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.AssignPetResponse)
	}

	return res, args.Error(1)
}

func (c *ClientMock) Delete(_ context.Context, in *proto.DeleteImageRequest, _ ...grpc.CallOption) (res *proto.DeleteImageResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.DeleteImageResponse)
	}

	return res, args.Error(1)
}

func (c *ClientMock) DeleteByPetId(_ context.Context, in *proto.DeleteByPetIdRequest, _ ...grpc.CallOption) (res *proto.DeleteByPetIdResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.DeleteByPetIdResponse)
	}

	return res, args.Error(1)
}

type ServiceMock struct {
	mock.Mock
}

func (c *ServiceMock) FindByPetId(petId string) (res []*proto.Image, err error) {
	args := c.Called(petId)

	if args.Get(0) != nil {
		res = args.Get(0).([]*proto.Image)
	}

	return res, args.Error(1)
}

func (c *ServiceMock) AssignPet(petId string, imageIds []string) (err error) {
	args := c.Called(petId, imageIds)

	return args.Error(0)
}
