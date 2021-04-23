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
	err = database.DB.Create(r).Error
	return r, err
}

func ListRecords(c *gin.Context) {
	var (
		rmiRecord Record
		res       []Record
		count     int64
		order     = c.Query("order")
	)

	if err := c.ShouldBind(&rmiRecord); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err,
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
