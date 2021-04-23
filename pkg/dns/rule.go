package dns

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/li4n0/revsuit/internal/database"
	"github.com/li4n0/revsuit/internal/newdns"
	"github.com/li4n0/revsuit/internal/rule"
	"gorm.io/gorm/clause"
	log "unknwon.dev/clog/v2"
)

type Rule struct {
	rule.BaseRule
	Type  newdns.Type   `gorm:"default:1" form:"type" json:"type"`
	Value string        `form:"value" json:"value"`
	TTL   time.Duration `gorm:"ttl;default:10" form:"ttl" json:"ttl"`
}

func (Rule) TableName() string {
	return "dns_rules"
}

// New dns rule struct
func NewRule(name, flagFormat, value string, pushToClient, notice bool, _type newdns.Type, ttl time.Duration) *Rule {
	return &Rule{
		BaseRule: rule.BaseRule{
			Name:         name,
			FlagFormat:   flagFormat,
			PushToClient: pushToClient,
			Notice:       notice,
		},
		Type:  _type,
		Value: value,
		TTL:   ttl,
	}
}

// Create or update the dns rule in database and ruleSet
func (r *Rule) CreateOrUpdate() (err error) {
	db := database.DB.Model(r)
	err = db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns(
			[]string{
				"name",
				"flag_format",
				"rank",
				"type",
				"value",
				"ttl",
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

// Delete the dns rule in database and ruleSet
func (r *Rule) Delete() (err error) {
	db := database.DB.Model(r)
	err = db.Delete(r).Error
	if err != nil {
		return
	}
	err = GetServer().updateRules()
	return err
}

// List all dns rules those satisfy the filter
func ListRules(c *gin.Context) {
	var (
		dnsRule Rule
		res     []Rule
		count   int64
		order   = c.Query("order")
	)

	if err := c.ShouldBind(&dnsRule); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"result": nil,
		})
		return
	}

	db := database.DB.Model(&dnsRule)
	if dnsRule.Name != "" {
		db.Where("name = ?", dnsRule.Name)
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

// Create or update dns rule from user submit
func UpsertRules(c *gin.Context) {
	var (
		dnsRule Rule
		update  bool
	)

	if err := c.ShouldBind(&dnsRule); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"data":   nil,
		})
		return
	}

	if dnsRule.ID != 0 {
		update = true
	}

	if err := dnsRule.CreateOrUpdate(); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"data":   nil,
		})
		return
	}

	if update {
		log.Trace("DNS rule(id:%d) has been updated", dnsRule.ID)
	} else {
		log.Trace("DNS rule(id:%d) has been created", dnsRule.ID)
	}

	c.JSON(200, gin.H{
		"status": "succeed",
		"error":  nil,
		"result": nil,
	})
}

// Delete dns rule from user submit
func DeleteRules(c *gin.Context) {
	var dnsRule Rule

	if err := c.ShouldBind(&dnsRule); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"data":   nil,
		})
		return
	}

	if err := dnsRule.Delete(); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"data":   nil,
		})
		return
	}

	log.Trace("DNS rule(id:%d) has been deleted", dnsRule.ID)

	c.JSON(200, gin.H{
		"status": "succeed",
		"error":  nil,
		"data":   nil,
	})
}
