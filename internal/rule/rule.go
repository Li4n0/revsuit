package rule

import (
	"regexp"
	"strings"
	"time"

	log "unknwon.dev/clog/v2"
)

type Rule interface {
	CreateOrUpdate() error
	Delete() error
}

type BaseRule struct {
	Rule        `gorm:"-" json:"-" yaml:"-"`
	ID          uint           `gorm:"primarykey" form:"id" json:"id" yaml:"-"`
	CreatedAt   time.Time      `json:"created_at" yaml:"-"`
	UpdatedAt   time.Time      `json:"updated_at" yaml:"-"`
	Name        string         `gorm:"index;unique;not null;" form:"name" json:"name"`
	FlagFormat  string         `gorm:"unique;not null;" form:"flag_format" json:"flag_format" yaml:"flag_format"`
	flagCatcher *regexp.Regexp `gorm:"-" json:"-"`
	// base_bank 解决 mysql 中关键字 rank 冲突和 pg 下不能使用 ``
	Rank         int  `gorm:"default:0;column:base_rank" json:"rank" form:"rank"`
	PushToClient bool `gorm:"default:false;not null;" form:"push_to_client" json:"push_to_client" yaml:"push_to_client"`
	Notice       bool `gorm:"default:false;not null;" form:"notice" json:"notice"`
}

func (br BaseRule) Match(s string) (flag, flagGroup string, vars map[string]string) {
	vars = make(map[string]string)

	if br.flagCatcher == nil {
		if br.FlagFormat == "*" {
			flag = "*"
			return
		} else {
			if catcher, err := regexp.Compile(br.FlagFormat); err != nil {
				log.Warn("%s[rule:%s]", err, br.Name)
				return
			} else {
				br.flagCatcher = catcher
			}
		}
	}

	matched := br.flagCatcher.FindStringSubmatch(s)
	groupNames := br.flagCatcher.SubexpNames()

	if len(matched) == 0 {
		return
	}

	flag = matched[0]
	if len(matched) > 1 && groupNames[1] == "" {
		flagGroup = matched[1]
	}

	for j, name := range groupNames {
		if j != 0 && name != "" {
			vars[name] = strings.TrimSpace(matched[j])
		}
	}
	vars["flag"] = flag

	return flag, flagGroup, vars
}
