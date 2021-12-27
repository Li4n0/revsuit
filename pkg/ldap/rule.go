package ldap

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/li4n0/revsuit/internal/database"
	"github.com/li4n0/revsuit/internal/rule"
	"gorm.io/gorm/clause"
	log "unknwon.dev/clog/v2"
)

// Rule LDAP rule struct
type Rule struct {
	rule.BaseRule `yaml:",inline"`
}

func (Rule) TableName() string {
	return "ldap_rules"
}

// NewRule new ldap rule struct
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

// CreateOrUpdate creates or updates the ldap rule in database and ruleSet
func (r *Rule) CreateOrUpdate() (err error) {
	db := database.DB.Model(r)
	err = db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns(
			[]string{
				"name",
				"flag_format",
				"base_rank",
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

// Delete deletes the ldap rule in database and ruleSet
func (r *Rule) Delete() (err error) {
	db := database.DB.Model(r)
	err = db.Delete(r).Error
	if err != nil {
		return
	}

	err = GetServer().UpdateRules()
	return err
}

// ListRules lists all ldap rules those satisfy the filter
func ListRules(c *gin.Context) {
	var (
		ldapRule Rule
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

	if err := c.ShouldBind(&ldapRule); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"result": nil,
		})
		return
	}

	db := database.DB.Model(&ldapRule)
	if ldapRule.Name != "" {
		db.Where("name = ?", ldapRule.Name)
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

// UpsertRules creates or updates ldap rule from user submit
func UpsertRules(c *gin.Context) {
	var (
		ldapRule Rule
		update   bool
	)

	if err := c.ShouldBind(&ldapRule); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"data":   nil,
		})
		return
	}

	if ldapRule.ID != 0 {
		update = true
	}

	if err := ldapRule.CreateOrUpdate(); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"result": nil,
		})
		return
	}

	if update {
		log.Trace("LDAP rule[id:%d] has been updated", ldapRule.ID)
	} else {
		log.Trace("LDAP rule[id:%d] has been created", ldapRule.ID)
	}

	c.JSON(200, gin.H{
		"status": "succeed",
		"error":  nil,
		"result": nil,
	})
}

// DeleteRules deletes ldap rule from user submit
func DeleteRules(c *gin.Context) {
	var ldapRule Rule

	if err := c.ShouldBind(&ldapRule); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"data":   nil,
		})
		return
	}

	if err := ldapRule.Delete(); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"data":   nil,
		})
		return
	}

	log.Trace("LDAP rule[id:%d] has been deleted", ldapRule.ID)

	c.JSON(200, gin.H{
		"status": "succeed",
		"error":  nil,
		"data":   nil,
	})
}
