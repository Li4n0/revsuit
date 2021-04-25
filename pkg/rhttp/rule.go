package rhttp

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/li4n0/revsuit/internal/database"
	"github.com/li4n0/revsuit/internal/rule"
	"gorm.io/gorm/clause"
	log "unknwon.dev/clog/v2"
)

// Http rule struct
type Rule struct {
	rule.BaseRule
	ResponseStatusCode string            `gorm:"index;default:200;not null" form:"response_status_code" json:"response_status_code"`
	ResponseHeaders    database.MapField `form:"response_headers" json:"response_headers"`
	ResponseBody       string            `gorm:"default:Hello RevSuit!" form:"response_body" json:"response_body"`
}

func (Rule) TableName() string {
	return "http_rules"
}

// New http rule struct
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

// Create or update the http rule in database and ruleSet
func (r *Rule) CreateOrUpdate() (err error) {
	db := database.DB.Model(r)
	err = db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns(
			[]string{
				"name",
				"flag_format",
				"rank",
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

	err = GetServer().updateRules()
	return err
}

// Delete the http rule in database and ruleSet
func (r *Rule) Delete() (err error) {
	db := database.DB.Model(r)
	err = db.Delete(r).Error
	if err != nil {
		return
	}

	err = GetServer().updateRules()
	return err
}

// List all http rules those satisfy the filter
func ListRules(c *gin.Context) {
	var (
		httpRule Rule
		res      []Rule
		count    int64
		order    = c.Query("order")
	)

	if err := c.ShouldBind(&httpRule); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err,
			"result": nil,
		})
		return
	}

	db := database.DB.Model(&httpRule)
	if httpRule.Name != "" {
		db.Where("name = ?", httpRule.Name)
	}
	db.Count(&count)

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

	if err := db.Order("rank desc").Order("id" + " " + order).Count(&count).Offset((page - 1) * 10).Limit(10).Find(&res).Error; err != nil {
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

// Create or update http rule from user submit
func UpsertRules(c *gin.Context) {
	var (
		httpRule Rule
		update   bool
	)

	if err := c.ShouldBind(&httpRule); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"data":   nil,
		})
		return
	}

	if httpRule.ID != 0 {
		update = true
	}

	if err := httpRule.CreateOrUpdate(); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"result": nil,
		})
		return
	}

	if update {
		log.Trace("HTTP rule[id%d] has been updated", httpRule.ID)
	} else {
		log.Trace("HTTP rule[id%d] has been created", httpRule.ID)
	}

	c.JSON(200, gin.H{
		"status": "succeed",
		"error":  nil,
		"result": nil,
	})
}

// Delete http rule from user submit
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

	log.Trace("HTTP rule[id%d] has been deleted", httpRule.ID)

	c.JSON(200, gin.H{
		"status": "succeed",
		"error":  nil,
		"data":   nil,
	})
}
