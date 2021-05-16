package mysql

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/li4n0/revsuit/internal/database"
	"github.com/li4n0/revsuit/internal/rule"
	"gorm.io/gorm/clause"
	log "unknwon.dev/clog/v2"
)

type Rule struct {
	rule.BaseRule     `yaml:",inline"`
	Files             string            `form:"files" json:"files"`
	ExploitJdbcClient bool              `gorm:"exploit_jdbc_client" form:"exploit_jdbc_client" json:"exploit_jdbc_client" yaml:"exploit_jdbc_client"`
	Payloads          database.MapField `json:"payloads" form:"payloads"`
}

func (Rule) TableName() string {
	return "mysql_rules"
}

// CreateOrUpdate creates or updates the mysql rule in database and ruleSet
func (r *Rule) CreateOrUpdate() (err error) {
	db := database.DB.Model(r)
	err = db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns(
			[]string{
				"name",
				"flag_format",
				"rank",
				"files",
				"exploit_jdbc_client",
				"payloads",
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

// Delete deletes the mysql rule in database and ruleSet
func (r *Rule) Delete() (err error) {
	db := database.DB.Model(r)
	err = db.Delete(r).Error
	if err != nil {
		return
	}
	err = GetServer().updateRules()
	return err
}

// ListRules lists all mysql rules those satisfy the filter
func ListRules(c *gin.Context) {
	var (
		mysqlRule Rule
		res       []Rule
		count     int64
		order     = c.Query("order")
		pageSize  int
	)

	if c.Query("pageSize") == "" {
		pageSize = 10
	} else if n, err := strconv.Atoi(c.Query("pageSize")); err == nil {
		if n <= 0 || n > 100 {
			pageSize = 10
		}
	}

	if err := c.ShouldBind(&mysqlRule); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"result": nil,
		})
		return
	}

	db := database.DB.Model(&mysqlRule)
	if mysqlRule.Name != "" {
		db.Where("name = ?", mysqlRule.Name)
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

// UpsertRules create or update mysql rule from user submit
func UpsertRules(c *gin.Context) {
	var (
		mysqlRule Rule
		update    bool
	)

	if err := c.ShouldBind(&mysqlRule); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"data":   nil,
		})
		return
	}

	if mysqlRule.ID != 0 {
		update = true
	}

	if err := mysqlRule.CreateOrUpdate(); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"data":   nil,
		})
		return
	}

	if update {
		log.Trace("MySQL rule[id:%d] has been updated", mysqlRule.ID)
	} else {
		log.Trace("MySQL rule[id:%d] has been created", mysqlRule.ID)
	}

	c.JSON(200, gin.H{
		"status": "succeed",
		"error":  nil,
		"result": nil,
	})
}

// DeleteRules Delete mysql rule from user submit
func DeleteRules(c *gin.Context) {
	var mysqlRule Rule

	if err := c.ShouldBind(&mysqlRule); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"data":   nil,
		})
		return
	}

	if err := mysqlRule.Delete(); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"data":   nil,
		})
		return
	}

	log.Trace("MySQL rule[id:%d] has been deleted", mysqlRule.ID)

	c.JSON(200, gin.H{
		"status": "succeed",
		"error":  nil,
		"data":   nil,
	})
}
