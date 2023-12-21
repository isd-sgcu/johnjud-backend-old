package pet

import (
	"github.com/stretchr/testify/mock"
)

// RepositoryMock is a mock type for the IRepository interface
type RepositoryMock struct {
	mock.Mock
}

// FindAll mocks the FindAll method
func (r *RepositoryMock) FindAll() error {
	// unimplemented for now
	return nil
}

// FindOne mocks the FindOne method
func (r *RepositoryMock) FindOne() error {
	// unimplemented for now
	return nil
}

// Create mocks the Create method
func (r *RepositoryMock) Create() error {
	// unimplemented for now
	return nil
}

// Update mocks the Update method
func (r *RepositoryMock) Update() error {
	// unimplemented for now
	return nil
}

// ChangeView mocks the ChangeView method
func (r *RepositoryMock) ChangeView() error {
	// unimplemented for now
	return nil
}

// Delete mocks the Delete method
func (r *RepositoryMock) Delete(id string) error {
	args := r.Called(id)
	return args.Error(0)
}
