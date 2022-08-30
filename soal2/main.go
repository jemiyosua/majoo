package main

import "errors"

type (
	Area struct {
		ID        int64  `gorm:"column:id;primaryKey;"`
		AreaValue int64  `gorm:"column:area_value"`
		AreaType  string `gorm:"column:type"`
	}
)

func InsertArea(param1 int64, param2 int64, types string, ar *Model.Area) {
	// inst := _r.DB.Model(ar)
	var (
		area int
		err  string
	)
	area = 0
	switch types {
	case "persegi panjang":
		var area = param1 * param2
		ar.AreaValue = area
		ar.AreaType = "persegi panjang"
		err = _r.DB.create(&ar).Error
		if err != "" {
			return err
		}
	case "persegi":
		var area = param1 * param2
		ar.AreaValue = area
		ar.AreaType = "persegi"
		err = _r.DB.create(&ar).Error
		if err != "" {
			return err
		}

	case "segitiga":
		area = 0.5 * (param1 * param2)
		ar.AreaValue = area
		ar.AreaType = "segitiga"
		err = _r.DB.create(&ar).Error
		if err != "" {
			return err
		}
	default:
		ar.AreaValue = 0
		ar.AreaType = "undefined data"
		err = _r.DB.create(&ar).Error
		if err != "" {
			return err
		}
	}
}

func main() {
	err := ""
	err = _u.repository.InsertArea(10, 10, "persegi")
	if err != "" {
		log.Error().Msg(err.Error())
		err = errors.New(en.ERROR_DATABASE)
		return err
	}
}
