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
	rule.BaseRule
	Files             string            `form:"files" json:"files"`
	ExploitJdbcClient bool              `gorm:"exploit_jdbc_client" form:"exploit_jdbc_client" json:"exploit_jdbc_client"`
	Payloads          database.MapField `json:"payloads" form:"payloads"`
}

func (Rule) TableName() string {
	return "mysql_rules"
}

// Create or update the mysql rule in database and ruleSet
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

// Delete the mysql rule in database and ruleSet
func (r *Rule) Delete() (err error) {
	db := database.DB.Model(r)
	err = db.Delete(r).Error
	if err != nil {
		return
	}
	err = GetServer().updateRules()
	return err
}

// List all mysql rules those satisfy the filter
func ListRules(c *gin.Context) {
	var (
		mysqlRule Rule
		res       []Rule
		count     int64
		order     = c.Query("order")
	)

	if err := c.ShouldBind(&mysqlRule); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err,
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
			"error":  err,
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
			"error":  err,
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

// Create or update mysql rule from user submit
func UpsertRules(c *gin.Context) {
	var (
		mysqlRule Rule
		update    bool
	)

	if err := c.ShouldBind(&mysqlRule); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err,
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
			"error":  err,
			"data":   nil,
		})
		return
	}

	if update {
		log.Trace("MySQL rule[id%d] has been updated", mysqlRule.ID)
	} else {
		log.Trace("MySQL rule[id%d] has been created", mysqlRule.ID)
	}

	c.JSON(200, gin.H{
		"status": "succeed",
		"error":  nil,
		"result": nil,
	})
}

// Delete mysql rule from user submit
func DeleteRules(c *gin.Context) {
	var mysqlRule Rule

	if err := c.ShouldBind(&mysqlRule); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err,
			"data":   nil,
		})
		return
	}

	if err := mysqlRule.Delete(); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err,
			"data":   nil,
		})
		return
	}

	log.Trace("MySQL rule[id%d] has been deleted", mysqlRule.ID)

	c.JSON(200, gin.H{
		"status": "succeed",
		"error":  nil,
		"data":   nil,
	})
}
