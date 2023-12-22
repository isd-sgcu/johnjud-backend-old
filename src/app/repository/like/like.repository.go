package like

import (
	"github.com/isd-sgcu/johnjud-backend/src/app/model/like"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindByUserId(userId string, result *[]*like.Like) error {
	return r.db.Model(&like.Like{}).Find(result, "user_id = ?", userId).Error
}

func (r *Repository) Create(in *like.Like) error {
	return r.db.Create(&in).Error
}

func (r *Repository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&like.Like{}).Error
}
