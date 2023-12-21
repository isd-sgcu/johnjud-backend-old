package adopt

import (
	"github.com/isd-sgcu/johnjud-backend/src/app/model/adopt"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindAll(result *[]*adopt.Adopt) error {
	return r.db.Model(&adopt.Adopt{}).Find(result).Error
}

func (r *Repository) Create(in *adopt.Adopt) error {
	return r.db.Create(&in).Error
}

func (r *Repository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&adopt.Adopt{}).Error
}
