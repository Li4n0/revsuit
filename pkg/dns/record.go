package dns

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/li4n0/revsuit/internal/database"
	"github.com/li4n0/revsuit/internal/notice"
	"github.com/li4n0/revsuit/internal/record"
	"gorm.io/gorm"
	log "unknwon.dev/clog/v2"
)

var _ record.Record = (*Record)(nil)

type Record struct {
	Domain string `gorm:"index" form:"domain" json:"domain"`

	record.BaseRecord
	Rule Rule `gorm:"foreignKey:RuleName;references:Name;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" form:"-" json:"-" notice:"-"`
}

type Domains struct {
	Domains []string `json:"domains"`
}

func (Record) TableName() string {
	return "dns_records"
}

func (r Record) Notice() {
	notice.Notice(r)
}

func newRecord(rule *Rule, flag, domain, remoteIp, ipArea string) (r *Record, err error) {
	r = &Record{
		BaseRecord: record.BaseRecord{
			Flag:        flag,
			RemoteIP:    remoteIp,
			IpArea:      ipArea,
			RequestTime: time.Now(),
		},
		Domain: domain,
		Rule:   *rule,
	}
	return r, database.DB.Create(r).Error
}

func FindDomains(c *gin.Context) {
	var (
		dnsRecord Record
		domains   Domains
		count     int64
	)
	res := []string{}
	if err := c.ShouldBind(&domains); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"result": nil,
		})
		return
	}
	for i := 0; i < len(domains.Domains); i++ {
		db := database.DB.Model(&dnsRecord)
		domain := domains.Domains[i]
		if err := db.Where("domain like ?", "%"+domain+"%").Count(&count); err != nil {
			if count > 0 {
				res = append(res, domain)
			}
		}

	}
	c.JSON(200, gin.H{
		"found": res,
	})
}

func Records(c *gin.Context) {
	var (
		dnsRecord Record
		res       []Record
		count     int64
		order     = c.Query("order")
		pageSize  = 10
	)

	if c.Query("pageSize") != "" {
		if n, err := strconv.Atoi(c.Query("pageSize")); err == nil {
			if n > 0 && n < 100 {
				pageSize = n
			}
		}
	}

	if err := c.ShouldBind(&dnsRecord); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"result": nil,
		})
	}

	db := database.DB.Model(&dnsRecord)
	if dnsRecord.Flag != "" {
		db.Where("flag = ?", dnsRecord.Flag)
	}
	if dnsRecord.Domain != "" {
		db.Where("domain like ?", "%"+dnsRecord.Domain+"%")
	}
	if dnsRecord.RemoteIP != "" {
		db.Where("remote_ip = ?", dnsRecord.RemoteIP)
	}
	if dnsRecord.RuleName != "" {
		db.Where("rule_name = ?", dnsRecord.RuleName)
	}

	//Delete records
	if c.Request.Method == http.MethodDelete {
		if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&res).Error; err != nil {
			c.JSON(400, gin.H{
				"status": "failed",
				"error":  err.Error(),
				"data":   nil,
			})
			return
		}

		if database.Driver == database.Sqlite {
			db.Exec("VACUUM")
		}

		c.JSON(200, gin.H{
			"status": "succeed",
			"error":  nil,
		})

		log.Info("%d dns records deleted by %s", db.RowsAffected, c.Request.RemoteAddr)
		return
	}

	//List records
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

	if err := db.Order("id" + " " + order).Count(&count).Offset((page - 1) * pageSize).Limit(pageSize).Find(&res).Error; err != nil {
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
