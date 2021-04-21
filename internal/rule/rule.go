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

func compileCatcher(flagFormat string) (reg *regexp.Regexp, err error) {
	if reg, err := regexp.Compile(flagFormat); err == nil {
		return reg, nil
	} else {
		// * meaning record all connections.
		if flagFormat != "*" {
			return nil, err
		} else {
			return nil, nil
		}
	}
}

func (br BaseRule) Match(s string) (flag, flagGroup string) {
	if br.flagCatcher == nil && br.FlagFormat != "*" {
		//	compile rule flags
		catcher, err := compileCatcher(br.FlagFormat)
		if err != nil {
			log.Error("%s(rule:%s)", err.Error(), br.Name)
		}
		br.flagCatcher = catcher
	}

	if br.flagCatcher == nil {
		// capture all connection.
		flag = "*"
	} else {
		matched := br.flagCatcher.FindStringSubmatch(s)
		if len(matched) == 0 {
			return
		}
		flag = matched[0]
		if len(matched) > 1 {
			flagGroup = matched[1]
		}
	}
	return
}
