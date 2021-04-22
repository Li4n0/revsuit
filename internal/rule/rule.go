package rule

import (
	"regexp"
	"time"

	log "unknwon.dev/clog/v2"
)

type Rule interface {
	CreateOrUpdate() error
	Delete() error
}

type BaseRule struct {
	Rule         `gorm:"-" json:"-"`
	ID           uint           `gorm:"primarykey" form:"id" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	Name         string         `gorm:"index;unique;not null;" form:"name" json:"name"`
	FlagFormat   string         `gorm:"unique;not null;" form:"flag_format" json:"flag_format"`
	flagCatcher  *regexp.Regexp `gorm:"-" json:"-"`
	Rank         int            `gorm:"default:0" json:"rank" form:"rank"`
	PushToClient bool           `gorm:"default:false;not null;" form:"push_to_client" json:"push_to_client"`
	Notice       bool           `gorm:"default:false;not null;" form:"notice" json:"notice"`
}

func (br BaseRule) Match(s string) (flag, flagGroup string) {
	if br.flagCatcher == nil {
		if br.FlagFormat == "*" {
			flag = "*"
			return
		} else {
			if catcher, err := regexp.Compile(br.FlagFormat); err != nil {
				log.Error("%s(rule:%s)", err.Error(), br.Name)
				return
			} else {
				br.flagCatcher = catcher
			}
		}
	}

	matched := br.flagCatcher.FindStringSubmatch(s)
	if len(matched) == 0 {
		return
	}

	flag = matched[0]
	if len(matched) > 1 {
		flagGroup = matched[1]
	}
	return flag, flagGroup
}
