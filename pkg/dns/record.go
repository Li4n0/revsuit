package dns

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/li4n0/revsuit/internal/database"
	"github.com/li4n0/revsuit/internal/notice"
	"github.com/li4n0/revsuit/internal/record"
)

var _ record.Record = (*Record)(nil)

type Record struct {
	Domain string `gorm:"index" form:"domain" json:"domain"`

	record.BaseRecord
	Rule Rule `gorm:"foreignKey:RuleName;references:Name;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" form:"-" json:"-" notice:"-"`
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

func ListRecords(c *gin.Context) {
	var (
		dnsRecord Record
		res       []Record
		count     int64
		order     = c.Query("order")
	)
	if err := c.ShouldBind(&dnsRecord); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err,
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

	if err := db.Order("id" + " " + order).Count(&count).Offset((page - 1) * 10).Limit(10).Find(&res).Error; err != nil {
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
