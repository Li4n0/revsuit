package ftp

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/li4n0/revsuit/internal/database"
	"github.com/li4n0/revsuit/internal/rule"
	"gorm.io/gorm/clause"
	log "unknwon.dev/clog/v2"
)

// FTP rule struct
type Rule struct {
	rule.BaseRule
	PasvAddress string `gorm:"pasv_address" json:"pasv_address" form:"pasv_address"`
}

func (Rule) TableName() string {
	return "ftp_rules"
}

// NewRule creates a new ftp rule struct
func NewRule(name, flagFormat, pasvAddress string, pushToClient, notice bool) *Rule {
	return &Rule{
		BaseRule: rule.BaseRule{
			Name:         name,
			FlagFormat:   flagFormat,
			PushToClient: pushToClient,
			Notice:       notice,
		},
		PasvAddress: pasvAddress,
	}
}

// Create or update the ftp rule in database and ruleSet
func (r *Rule) CreateOrUpdate() (err error) {
	db := database.DB.Model(r)
	err = db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns(
			[]string{
				"name",
				"flag_format",
				"rank",
				"pasv_address",
				"push_to_client",
				"notice",
			}),
	}).Create(r).Error
	if err != nil {
		return
	}

	return GetServer().updateRules()
}

// Delete the ftp rule in database and ruleSet
func (r *Rule) Delete() (err error) {
	db := database.DB.Model(r)
	err = db.Delete(r).Error
	if err != nil {
		return
	}

	err = GetServer().updateRules()
	return err
}

// List all ftp rules those satisfy the filter
func ListRules(c *gin.Context) {
	var (
		ftpRule Rule
		res     []Rule
		count   int64
		order   = c.Query("order")
	)

	if err := c.ShouldBind(&ftpRule); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err,
			"result": nil,
		})
		return
	}

	db := database.DB.Model(&ftpRule)
	if ftpRule.Name != "" {
		db.Where("name = ?", ftpRule.Name)
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

	if order != "desc" && order != "asc" {
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

// Create or update ftp rule from user submit
func UpsertRules(c *gin.Context) {
	var (
		ftpRule Rule
		update  bool
	)

	if err := c.ShouldBind(&ftpRule); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"data":   nil,
		})
		return
	}

	if ftpRule.ID != 0 {
		update = true
	}

	if err := ftpRule.CreateOrUpdate(); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"result": nil,
		})
		return
	}

	if update {
		log.Trace("FTP rule[id%d] has been updated", ftpRule.ID)
	} else {
		log.Trace("FTP rule[id%d] has been created", ftpRule.ID)
	}

	c.JSON(200, gin.H{
		"status": "succeed",
		"error":  nil,
		"result": nil,
	})
}

// Delete ftp rule from user submit
func DeleteRules(c *gin.Context) {
	var ftpRule Rule

	if err := c.ShouldBind(&ftpRule); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"data":   nil,
		})
		return
	}

	if err := ftpRule.Delete(); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"data":   nil,
		})
		return
	}

	log.Trace("FTP rule[id%d] has been deleted", ftpRule.ID)

	c.JSON(200, gin.H{
		"status": "succeed",
		"error":  nil,
		"data":   nil,
	})
}
