package rmi

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/li4n0/revsuit/internal/database"
	"github.com/li4n0/revsuit/internal/rule"
	"gorm.io/gorm/clause"
	log "unknwon.dev/clog/v2"
)

// Rule RMI rule struct
type Rule struct {
	rule.BaseRule `yaml:",inline"`
}

func (Rule) TableName() string {
	return "rmi_rules"
}

// NewRule New rmi rule struct
func NewRule(name, flagFormat string, pushToClient, notice bool) *Rule {
	return &Rule{
		BaseRule: rule.BaseRule{
			Name:         name,
			FlagFormat:   flagFormat,
			PushToClient: pushToClient,
			Notice:       notice,
		},
	}
}

// CreateOrUpdate Create or update the rmi rule in database and ruleSet
func (r *Rule) CreateOrUpdate() (err error) {
	db := database.DB.Model(r)
	err = db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns(
			[]string{
				"name",
				"flag_format",
				"rank",
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

// Delete the rmi rule in database and ruleSet
func (r *Rule) Delete() (err error) {
	db := database.DB.Model(r)
	err = db.Delete(r).Error
	if err != nil {
		return
	}

	err = GetServer().updateRules()
	return err
}

// ListRules List all rmi rules those satisfy the filter
func ListRules(c *gin.Context) {
	var (
		rmiRule  Rule
		res      []Rule
		count    int64
		order    = c.Query("order")
		pageSize int
	)

	if c.Query("pageSize") == "" {
		pageSize = 10
	} else if n, err := strconv.Atoi(c.Query("pageSize")); err == nil {
		if n <= 0 || n > 100 {
			pageSize = n
		}
	}

	if err := c.ShouldBind(&rmiRule); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"result": nil,
		})
		return
	}

	db := database.DB.Model(&rmiRule)
	if rmiRule.Name != "" {
		db.Where("name = ?", rmiRule.Name)
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

	if err := db.Order("rank desc").Order("id" + " " + order).Count(&count).Offset((page - 1) * pageSize).Limit(pageSize).Find(&res).Error; err != nil {
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

// UpsertRules Create or update rmi rule from user submit
func UpsertRules(c *gin.Context) {
	var (
		rmiRule Rule
		update  bool
	)

	if err := c.ShouldBind(&rmiRule); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"data":   nil,
		})
		return
	}

	if rmiRule.ID != 0 {
		update = true
	}

	if err := rmiRule.CreateOrUpdate(); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"result": nil,
		})
		return
	}

	if update {
		log.Trace("RMI rule[id:%d] has been updated", rmiRule.ID)
	} else {
		log.Trace("RMI rule[id:%d] has been created", rmiRule.ID)
	}

	c.JSON(200, gin.H{
		"status": "succeed",
		"error":  nil,
		"result": nil,
	})
}

// DeleteRules deletes rmi rule from user submit
func DeleteRules(c *gin.Context) {
	var rmiRule Rule

	if err := c.ShouldBind(&rmiRule); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"data":   nil,
		})
		return
	}

	if err := rmiRule.Delete(); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"data":   nil,
		})
		return
	}

	log.Trace("RMI rule[id:%d] has been deleted", rmiRule.ID)

	c.JSON(200, gin.H{
		"status": "succeed",
		"error":  nil,
		"data":   nil,
	})
}
