package database

import (
	//"database/sql"
	//"log"
	//"os"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type Card struct {
	ID        int64
	Name      string
	Rain      bool `gorm:"default:false"`
	Snow      bool `gorm:"default:false"`
	Night     bool `gorm:"default:false"`
	Bullet    bool `gorm:"default:false"`
	Mask      bool `gorm:"default:false"`
	Whistle   bool `gorm:"default:false"`
	Trap      bool `gorm:"default:false"`
	HardKnock bool `gorm:"default:false"`
}

type Hero struct {
	ID     int64
	Name   string
	Handle string
}

var (
	DB  *gorm.DB
	err error
)

func InitCard() {
	DB.CreateTable(&Card{})
	weather := []string{"Night", "Snow", "Rain"}
	threat := []string{"Bullet", "Mask", "Whistle"}
	number := []string{"0", "1", "Trap"}

	for _, w := range weather {
		for _, t := range threat {
			for _, n := range number {
				name := w + "_" + t + "_" + n
				DB.Create(cardCreater(w, t, n, name))
			}
		}
		name := w + "_HardKnock"
		DB.Create(cardCreater(w, "", "HardKnock", name))
	}
	for _, t := range threat {
		name := t + "_HardKnock"
		DB.Create(cardCreater("", t, "HardKnock", name))
	}

	DB.Create(&Card{Name: "All", Night: true, Snow: true, Rain: true, Bullet: true, Mask: true, Whistle: true})
	DB.Create(&Card{Name: "Weather", Night: true, Snow: true, Rain: true})
	DB.Create(&Card{Name: "Threat", Bullet: true, Mask: true, Whistle: true})
	DB.Create(&Card{Name: "Night_Bullet_2", Night: true, Bullet: true})
	DB.Create(&Card{Name: "Rain_Mask_2", Rain: true, Mask: true})
	DB.Create(&Card{Name: "Snow_Whistle_2", Snow: true, Whistle: true})
	DB.Create(&Card{Name: "Christmas"})
	DB.Create(&Card{Name: "HardKnock", HardKnock: true})
	DB.Create(&Card{Name: "HardKnock_2", HardKnock: true})
	DB.Create(&Card{Name: "Night_Snow", Night: true, Snow: true})
	DB.Create(&Card{Name: "Snow_Rain", Snow: true, Rain: true})
	DB.Create(&Card{Name: "Rain_Night", Night: true, Rain: true})
	DB.Create(&Card{Name: "Whistle_Bullet", Bullet: true, Whistle: true})
	DB.Create(&Card{Name: "Mask_Whistle", Mask: true, Whistle: true})
	DB.Create(&Card{Name: "Bullet_Mask", Bullet: true, Mask: true})

	DB.CreateTable(&Hero{})
	hero := []string{"Night", "Snow", "Rain", "Bullet", "Mask", "Whistle"}
	used := []string{"", "_used"}
	for _, u := range used {
		for _, h := range hero {
			DB.Create(heroCreater(h, u))
		}
	}

}

func heroCreater(h string, u string) *Hero {
	return &Hero{Name: "hero_" + h + u, Handle: h}
}

func cardCreater(w string, t string, h string, name string) *Card { //根據輸入建立卡片的結構並返回
	var nbool, sbool, rbool, bbool, mbool, wbool, hkbool, trapbool bool
	switch w {
	case "Night":
		nbool = true
	case "Snow":
		sbool = true
	case "Rain":
		rbool = true
	default:
	}
	switch t {
	case "Bullet":
		bbool = true
	case "Mask":
		mbool = true
	case "Whistle":
		wbool = true
	}
	if h == "HardKnock" {
		hkbool = true
	} else if h == "Trap" {
		trapbool = true
	}
	return &Card{Name: name, Night: nbool, Bullet: bbool, Snow: sbool, Rain: rbool, Mask: mbool, Whistle: wbool, HardKnock: hkbool, Trap: trapbool}
}
