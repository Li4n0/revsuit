package rhttp

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
	Method string `gorm:"index" form:"method" json:"method"`
	Path   string `form:"path" json:"path"`
	record.BaseRecord
	RawRequest string `json:"raw_request" notice:"-"`
	Rule       Rule   `gorm:"foreignKey:RuleName;references:Name;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" form:"-" json:"-" notice:"-"`
}

func (Record) TableName() string {
	return "http_records"
}

func (r Record) Notice() {
	notice.Notice(r)
}

func NewRecord(rule *Rule, flag, method, url, ip, area, raw string) (r *Record, err error) {
	r = &Record{
		BaseRecord: record.BaseRecord{
			Flag:        flag,
			RemoteIP:    ip,
			IpArea:      area,
			RequestTime: time.Now(),
		},
		Method:     method,
		Path:       url,
		RawRequest: raw,
		Rule:       *rule,
	}
	err = database.DB.Create(r).Error
	return r, err
}

func ListRecords(c *gin.Context) {
	var (
		httpRecord Record
		res        []Record
		count      int64
		order      = c.Query("order")
	)

	if err := c.ShouldBind(&httpRecord); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err,
			"result": nil,
		})
		return
	}

	db := database.DB.Model(&httpRecord)
	if httpRecord.Flag != "" {
		db.Where("flag = ?", httpRecord.Flag)
	}
	if httpRecord.Method != "" {
		db.Where("method = ?", httpRecord.Method)
	}
	if httpRecord.Path != "" {
		db.Where("path like ?", "%"+httpRecord.Path+"%")
	}
	if httpRecord.RemoteIP != "" {
		db.Where("remote_ip = ?", httpRecord.RemoteIP)
	}
	if httpRecord.RuleName != "" {
		db.Where("rule_name = ?", httpRecord.RuleName)
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
