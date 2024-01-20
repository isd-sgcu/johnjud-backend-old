package seed

func (s Seed) PetSeed1705074984036() error {
	for _, b := range pets {
		err := s.db.Save(&b).Error

		if err != nil {
			return err
		}
	}
	return nil
}
