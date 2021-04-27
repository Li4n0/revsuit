package mysql

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
	Username      string           `gorm:"index" form:"username" json:"username" notice:"username"`
	ClientName    string           `gorm:"index" form:"client_name" json:"client_name" notice:"client_name"`
	ClientOS      string           `gorm:"index" form:"client_os" json:"client_os" notice:"client_os"`
	LoadLocalData bool             `gorm:"index" form:"load_local_data" json:"load_local_data" notice:"load_local_data"`
	Files         []file.MySQLFile `form:"-" json:"files" notice:"-"`
	Rule          Rule             `gorm:"foreignKey:RuleName;references:Name;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" form:"-" json:"-" notice:"-"`
}

func (Record) TableName() string {
	return "mysql_records"
}

func (r Record) Notice() {
	notice.Notice(r)
}

func newRecord(rule *Rule, flag, username, clientName, clientOS, remoteIp, ipArea string, supportLoadLocalData bool, files []file.MySQLFile) (r *Record, err error) {
	r = &Record{
		BaseRecord: record.BaseRecord{
			Flag:        flag,
			RemoteIP:    remoteIp,
			IpArea:      ipArea,
			RequestTime: time.Now(),
		},
		Username:      username,
		ClientName:    clientName,
		ClientOS:      clientOS,
		LoadLocalData: supportLoadLocalData,
		Files:         files,
		Rule:          *rule,
	}
	return r, database.DB.Create(r).Error
}

func ListRecords(c *gin.Context) {
	var (
		mysqlRecord Record
		res         []Record
		count       int64
		order       = c.Query("order")
	)

	if err := c.ShouldBind(&mysqlRecord); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err,
			"result": nil,
		})
	}

	db := database.DB.Model(&mysqlRecord)
	if mysqlRecord.Flag != "" {
		db.Where("flag = ?", mysqlRecord.Flag)
	}
	if mysqlRecord.RemoteIP != "" {
		db.Where("remote_ip = ?", mysqlRecord.RemoteIP)
	}
	if mysqlRecord.RuleName != "" {
		db.Where("rule_name = ?", mysqlRecord.RuleName)
	}
	if mysqlRecord.ClientOS != "" {
		db.Where("client_os like ?", "%"+mysqlRecord.ClientOS)
	}
	if mysqlRecord.ClientName != "" {
		db.Where("client_name like ?", "%"+mysqlRecord.ClientName)
	}
	if c.Query("load_local_data") != "" {
		if c.Query("load_local_data") == "true" {
			db.Where("load_local_data = ?", true)
		} else {
			db.Where("load_local_data = ?", false)
		}
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

	if err := db.Preload("Files").Order("id" + " " + order).Count(&count).Offset((page - 1) * 10).Limit(10).Find(&res).Error; err != nil {
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
