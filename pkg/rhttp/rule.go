package rhttp

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/li4n0/revsuit/internal/database"
	"github.com/li4n0/revsuit/internal/rule"
	"gorm.io/gorm/clause"
	log "unknwon.dev/clog/v2"
)

// Rule Http rule struct
type Rule struct {
	rule.BaseRule      `yaml:",inline"`
	ResponseStatusCode string            `gorm:"index;default:200;not null" form:"response_status_code" json:"response_status_code" yaml:"response_status_code"`
	ResponseHeaders    database.MapField `form:"response_headers" json:"response_headers" yaml:"response_headers"`
	ResponseBody       string            `gorm:"type:longtext" form:"response_body" json:"response_body" yaml:"response_body"`
}

func (Rule) TableName() string {
	return "http_rules"
}

// NewRule new http rule struct
func NewRule(name, flagFormat, responseBody string, pushToClient, notice bool, responseStatus string, responseHeaders database.MapField) *Rule {
	return &Rule{
		BaseRule: rule.BaseRule{
			Name:         name,
			FlagFormat:   flagFormat,
			PushToClient: pushToClient,
			Notice:       notice,
		},
		ResponseStatusCode: responseStatus,
		ResponseHeaders:    responseHeaders,
		ResponseBody:       responseBody,
	}
}

// CreateOrUpdate creates or updates the http rule in database and ruleSet
func (r *Rule) CreateOrUpdate() (err error) {
	db := database.DB.Model(r)
	err = db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns(
			[]string{
				"name",
				"flag_format",
				"base_rank",
				"response_status_code",
				"response_headers",
				"response_body",
				"push_to_client",
				"notice",
			}),
	}).Create(r).Error
	if err != nil {
		return
	}

	err = GetServer().UpdateRules()
	return err
}

// Delete deletes the http rule in database and ruleSet
func (r *Rule) Delete() (err error) {
	db := database.DB.Model(r)
	err = db.Delete(r).Error
	if err != nil {
		return
	}

	err = GetServer().UpdateRules()
	return err
}

// ListRules lists all http rules those satisfy the filter
func ListRules(c *gin.Context) {
	var (
		httpRule Rule
		res      []Rule
		count    int64
		order    = c.Query("order")
		pageSize = 10
	)

	if c.Query("pageSize") != "" {
		if n, err := strconv.Atoi(c.Query("pageSize")); err == nil {
			if n > 0 && n < 100 {
				pageSize = n
			}
		}
	}

	if err := c.ShouldBind(&httpRule); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"result": nil,
		})
		return
	}

	db := database.DB.Model(&httpRule)
	db.Where(&httpRule).Count(&count)

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"result": nil,
		})
		return
	}

	if order != "asc" {
		order = "desc"
	}

	if err := db.Order("base_rank desc").Order("id" + " " + order).Count(&count).Offset((page - 1) * pageSize).Limit(pageSize).Find(&res).Error; err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"data":   nil,
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "succeed",
		"error":  nil,
		"result": gin.H{"count": count, "data": res},
	})
}

// UpsertRules create or update http rule from user submit
func UpsertRules(c *gin.Context) {
	var httpRule Rule

	if err := c.ShouldBind(&httpRule); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"data":   nil,
		})
		return
	}

	if err := httpRule.CreateOrUpdate(); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"result": nil,
		})
		return
	}

	if httpRule.ID != 0 {
		log.Trace("HTTP rule[id:%d] has been updated", httpRule.ID)
	} else {
		log.Trace("HTTP rule[id:%d] has been created", httpRule.ID)
	}

	c.JSON(200, gin.H{
		"status": "succeed",
		"error":  nil,
		"result": nil,
	})
}

// DeleteRules deletes http rule from user submit
func DeleteRules(c *gin.Context) {
	var httpRule Rule

	if err := c.ShouldBind(&httpRule); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"data":   nil,
		})
		return
	}

	if err := httpRule.Delete(); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"data":   nil,
		})
		return
	}

	log.Trace("HTTP rule[id:%d] has been deleted", httpRule.ID)

	c.JSON(200, gin.H{
		"status": "succeed",
		"error":  nil,
		"data":   nil,
	})
}
