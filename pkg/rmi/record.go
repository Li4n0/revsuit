package rmi

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
	Path string `form:"path" json:"path"`
	record.BaseRecord
	Rule Rule `gorm:"foreignKey:RuleName;references:Name;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" form:"-" json:"-" notice:"-"`
}

func (Record) TableName() string {
	return "rmi_records"
}

func (r Record) Notice() {
	notice.Notice(r)
}

func NewRecord(rule *Rule, flag, path, ip, area string) (r *Record, err error) {
	r = &Record{
		BaseRecord: record.BaseRecord{
			Flag:        flag,
			RemoteIP:    ip,
			IpArea:      area,
			RequestTime: time.Now(),
		},
		Path: path,
		Rule: *rule,
	}
	return r, database.DB.Create(r).Error
}

func ListRecords(c *gin.Context) {
	var (
		rmiRecord Record
		res       []Record
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

	if err := c.ShouldBind(&rmiRecord); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"result": nil,
		})
		return
	}

	db := database.DB.Model(&rmiRecord)
	if rmiRecord.Flag != "" {
		db.Where("flag = ?", rmiRecord.Flag)
	}
	if rmiRecord.Path != "" {
		db.Where("path like ?", "%"+rmiRecord.Path+"%")
	}
	if rmiRecord.RemoteIP != "" {
		db.Where("remote_ip = ?", rmiRecord.RemoteIP)
	}
	if rmiRecord.RuleName != "" {
		db.Where("rule_name = ?", rmiRecord.RuleName)
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
