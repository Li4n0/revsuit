package ftp

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/li4n0/revsuit/internal/database"
	"github.com/li4n0/revsuit/internal/file"
	"github.com/li4n0/revsuit/internal/notice"
	"github.com/li4n0/revsuit/internal/record"
)

var _ record.Record = (*Record)(nil)

type Record struct {
	record.BaseRecord
	User     string       `form:"user" json:"user"`
	Password string       `form:"password" json:"password"`
	Path     string       `form:"path" json:"path"`
	Method   Method       `form:"method" json:"method"`
	Status   Status       `form:"status" json:"status"`
	File     file.FTPFile `form:"file" json:"file" notice:"-"`
	Rule     Rule         `gorm:"foreignKey:RuleName;references:Name;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" form:"-" json:"-" notice:"-"`
}

func (Record) TableName() string {
	return "ftp_records"
}

func (r Record) Notice() {
	notice.Notice(r)
}

func NewRecord(rule *Rule, flag, user, password, method, path, ip, area string, file file.FTPFile, status Status) (r *Record, err error) {
	r = &Record{
		BaseRecord: record.BaseRecord{
			Flag:        flag,
			RemoteIP:    ip,
			IpArea:      area,
			RequestTime: time.Now(),
		},
		Path:     path,
		Method:   method,
		User:     user,
		Password: password,
		Status:   status,
		File:     file,
		Rule:     *rule,
	}
	return r, database.DB.Create(r).Error
}

func ListRecords(c *gin.Context) {
	var (
		ftpRecord Record
		res       []Record
		count     int64
		order     = c.Query("order")
	)

	if err := c.ShouldBind(&ftpRecord); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err,
			"result": nil,
		})
		return
	}

	db := database.DB.Model(&ftpRecord)
	if ftpRecord.Flag != "" {
		db.Where("flag = ?", ftpRecord.Flag)
	}
	if ftpRecord.User != "" {
		db.Where("user like ?", "%"+ftpRecord.User+"%")
	}
	if ftpRecord.Password != "" {
		db.Where("password like ?", "%"+ftpRecord.Password+"%")
	}
	if ftpRecord.Path != "" {
		db.Where("path like ?", "%"+ftpRecord.Path+"%")
	}
	if ftpRecord.Method != "" {
		db.Where("method = ?", ftpRecord.Method)
	}
	if ftpRecord.Status != "" {
		db.Where("status = ?", ftpRecord.Status)
	}
	if ftpRecord.RemoteIP != "" {
		db.Where("remote_ip = ?", ftpRecord.RemoteIP)
	}
	if ftpRecord.RuleName != "" {
		db.Where("rule_name = ?", ftpRecord.RuleName)
	}

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

	if err := db.Preload("File").Order("id" + " " + order).Count(&count).Offset((page - 1) * 10).Limit(10).Find(&res).Error; err != nil {
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
