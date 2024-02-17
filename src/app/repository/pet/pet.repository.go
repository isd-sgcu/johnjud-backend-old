package pet

import (
	"errors"

	"github.com/isd-sgcu/johnjud-backend/src/app/model/pet"
	petUtils "github.com/isd-sgcu/johnjud-backend/src/app/utils/pet"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindAll(result *[]*pet.Pet) error {
	return r.db.Model(&pet.Pet{}).Find(result).Error
}

func (r *Repository) FindOne(id string, result *pet.Pet) error {
	return r.db.Model(&pet.Pet{}).First(result, "id = ?", id).Error
}

func (r *Repository) Create(in *pet.Pet) error {
	return r.db.Create(&in).Error
}

func (r *Repository) Update(id string, result *pet.Pet) error {
	updateMap := petUtils.UpdateMap(result)
	return r.db.Model(&result).Updates(updateMap).First(&result, "id = ?", id).Error
}

func (r *Repository) Delete(id string) error {
	var pet pet.Pet
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&pet).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
		return err
	}
	return r.db.Delete(&pet).Error
}
