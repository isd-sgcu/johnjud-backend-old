package pet

import (
	"errors"
	"math"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/isd-sgcu/johnjud-backend/src/app/model"
	"github.com/isd-sgcu/johnjud-backend/src/app/model/pet"
	"github.com/isd-sgcu/johnjud-backend/src/constant"
	petConst "github.com/isd-sgcu/johnjud-backend/src/constant/pet"
	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/pet/v1"
	imageProto "github.com/isd-sgcu/johnjud-go-proto/johnjud/file/image/v1"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func FilterPet(in *[]*pet.Pet, query *proto.FindAllPetRequest) error {
	if query.MaxAge == 0 {
		query.MaxAge = math.MaxInt32
	}
	log.Info().
		Str("service", "filter pet").
		Str("module", "FilterPet").Msgf("minAge: %d, maxAge: %d", query.MinAge, query.MaxAge)
	log.Info().
		Str("service", "filter pet").
		Str("module", "FilterPet").Interface("query", query).Msgf("query")
	var results []*pet.Pet
	for _, p := range *in {
		res, err := filterAge(p, query.MinAge, query.MaxAge)
		if err != nil {
			return err
		}
		if !res {
			log.Info().
				Str("service", "filter pet").
				Str("module", "FilterPet reject").Msg("age not in range")
			continue
		}
		if query.Search != "" && !strings.Contains(p.Name, query.Search) {
			log.Info().
				Str("service", "filter pet").
				Str("module", "FilterPet reject").Msg("not in search")
			continue
		}
		if query.Type != "" && p.Type != query.Type {
			log.Info().
				Str("service", "filter pet").
				Str("module", "FilterPet reject").Msg("not the type")
			continue
		}
		if query.Gender != "" && p.Gender != petConst.Gender(query.Gender) {
			log.Info().
				Str("service", "filter pet").
				Str("module", "FilterPet reject").Msg("not the gender")
			continue
		}
		if query.Color != "" && p.Color != query.Color {
			log.Info().
				Str("service", "filter pet").
				Str("module", "FilterPet reject").Msg("not the color")
			continue
		}
		if query.Origin != "" && p.Origin != query.Origin {
			log.Info().
				Str("service", "filter pet").
				Str("module", "FilterPet reject").Msg("not the origin")
			continue
		}
		log.Info().
			Str("service", "filter pet").
			Str("module", "FilterPet accept").Interface("pet", p).Msgf("pet accepted")
		results = append(results, p)
	}
	*in = results
	return nil
}

func PaginatePets(pets *[]*pet.Pet, page int32, pageSize int32, metadata *proto.FindAllPetMetaData) error {
	totalsPets := int32(len(*pets))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = totalsPets
	}
	start := (page - 1) * pageSize
	end := start + pageSize

	if start > totalsPets {
		*pets = []*pet.Pet{}
		return nil
	}
	if end > totalsPets {
		end = totalsPets
	}
	*pets = (*pets)[start:end]

	totalPages := int32(math.Ceil(float64(totalsPets) / float64(pageSize)))

	metadata.Page = page
	metadata.PageSize = pageSize
	metadata.Total = totalsPets
	metadata.TotalPages = totalPages
	return nil
}

func RawToDtoList(in *[]*pet.Pet, images map[string][]*imageProto.Image, query *proto.FindAllPetRequest) ([]*proto.Pet, error) {
	var result []*proto.Pet
	if len(*in) != len(images) {
		return nil, errors.New("length of in and imageUrls have to be the same")
	}

	for _, p := range *in {
		// TODO: create new filter image function this wont work
		result = append(result, RawToDto(p, images[p.ID.String()]))
	}
	return result, nil
}

func RawToDto(in *pet.Pet, images []*imageProto.Image) *proto.Pet {
	return &proto.Pet{
		Id:           in.ID.String(),
		Type:         in.Type,
		Name:         in.Name,
		Birthdate:    in.Birthdate,
		Gender:       string(in.Gender),
		Color:        in.Color,
		Habit:        in.Habit,
		Caption:      in.Caption,
		Status:       string(in.Status),
		Images:       images,
		IsSterile:    in.IsSterile,
		IsVaccinated: in.IsVaccinated,
		IsVisible:    in.IsVisible,
		Origin:       in.Origin,
		Address:      in.Address,
		Contact:      in.Contact,
		AdoptBy:      in.AdoptBy,
	}
}

func DtoToRaw(in *proto.Pet) (res *pet.Pet, err error) {
	var id uuid.UUID
	var gender petConst.Gender
	var status petConst.Status

	if in.Id != "" {
		id, err = uuid.Parse(in.Id)
		if err != nil {
			return nil, err
		}
	}

	switch in.Gender {
	case string(petConst.MALE):
		gender = petConst.MALE
	case string(petConst.FEMALE):
		gender = petConst.FEMALE
	}

	switch in.Status {
	case string(petConst.ADOPTED):
		status = petConst.ADOPTED
	case string(petConst.FINDHOME):
		status = petConst.FINDHOME
	}

	return &pet.Pet{
		Base: model.Base{
			ID:        id,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: gorm.DeletedAt{},
		},
		Type:         in.Type,
		Name:         in.Name,
		Birthdate:    in.Birthdate,
		Gender:       gender,
		Color:        in.Color,
		Habit:        in.Habit,
		Caption:      in.Caption,
		Status:       status,
		IsSterile:    in.IsSterile,
		IsVaccinated: in.IsVaccinated,
		IsVisible:    in.IsVisible,
		Origin:       in.Origin,
		Address:      in.Address,
		Contact:      in.Contact,
		AdoptBy:      in.AdoptBy,
	}, nil
}

func ExtractImageUrls(in []*imageProto.Image) []string {
	var result []string
	for _, e := range in {
		result = append(result, e.ImageUrl)
	}
	return result
}

func ExtractImageIDs(in []*imageProto.Image) []string {
	var result []string
	for _, e := range in {
		result = append(result, e.Id)
	}
	return result
}

func UpdateMap(in *pet.Pet) map[string]interface{} {
	updateMap := make(map[string]interface{})
	t := reflect.TypeOf(*in)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		typeName := field.Type.Name()
		value := reflect.ValueOf(*in).Field(i).Interface()
		if (typeName == "string" || typeName == "Gender" || typeName == "Status") && value != "" {
			updateMap[field.Name] = value
		}
		if typeName == "bool" {
			updateMap[field.Name] = value
		}
	}
	return updateMap
}

func parseDate(dateStr string) (time.Time, error) {
	parsedTime, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return time.Time{}, err
	}
	return parsedTime, nil
}

func filterAge(pet *pet.Pet, minAge, maxAge int32) (bool, error) {
	birthdate, err := parseDate(pet.Birthdate)
	if err != nil {
		return false, err
	}

	currYear := time.Now()
	birthYear := birthdate
	diff := currYear.Sub(birthYear).Hours() / constant.DAY / constant.YEAR

	return diff >= float64(minAge) && diff <= float64(maxAge), nil
}
