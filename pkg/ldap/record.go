package ldap

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
	Path string `form:"path" json:"path"`
	record.BaseRecord
	Rule Rule `gorm:"foreignKey:RuleName;references:Name;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" form:"-" json:"-" notice:"-"`
}

func (Record) TableName() string {
	return "ldap_records"
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

	// sqlite db-level lock to prevent too much write operation lead to error of `database is locked` #54
	if database.Driver == database.Sqlite {
		database.Locker.Lock()
		defer database.Locker.Unlock()
	}

	return r, database.DB.Create(r).Error
}

func Records(c *gin.Context) {
	var (
		ldapRecord Record
		res        []Record
		count      int64
		order      = c.Query("order")
		pageSize   = 10
	)

	if c.Query("pageSize") != "" {
		if n, err := strconv.Atoi(c.Query("pageSize")); err == nil {
			if n > 0 && n < 100 {
				pageSize = n
			}
		}
	}

	if err := c.ShouldBind(&ldapRecord); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"result": nil,
		})
		return
	}

	db := database.DB.Model(&ldapRecord)
	if ldapRecord.Flag != "" {
		db.Where("flag = ?", ldapRecord.Flag)
	}
	if ldapRecord.Path != "" {
		db.Where("path like ?", "%"+ldapRecord.Path+"%")
	}
	if ldapRecord.RemoteIP != "" {
		db.Where("remote_ip = ?", ldapRecord.RemoteIP)
	}
	if ldapRecord.RuleName != "" {
		db.Where("rule_name = ?", ldapRecord.RuleName)
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

		log.Info("%d ldap records deleted by %s", db.RowsAffected, c.Request.RemoteAddr)
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
