package file

import (
	"fmt"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/li4n0/revsuit/internal/database"
)

type File struct {
	ID       uint   `gorm:"primarykey" form:"id" json:"id"`
	RecordID uint   `json:"-"`
	Name     string `json:"name"`
	Content  []byte `json:"-"`
}

type MySQLFile File

func (MySQLFile) TableName() string {
	return "mysql_files"
}

type FTPFile File

func (FTPFile) TableName() string {
	return "ftp_files"
}

func GetFile(c *gin.Context) {
	var file File
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  fmt.Errorf("param id missed").Error(),
			"result": nil,
		})
		return
	}

	recordType := c.Param("record_type")
	if id == "" {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  fmt.Errorf("param record_type missed"),
			"result": nil,
		})
		return
	}
	if recordType == "mysql" {
		database.DB.Table("mysql_files").Where("id = ?", id).Find(&file)
	} else if recordType == "ftp" {
		database.DB.Table("ftp_files").Where("id = ?", id).Find(&file)
	}

	if file.Content == nil || len(file.Content) == 0 {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  fmt.Errorf("file not found").Error(),
			"result": nil,
		})
		return
	} else {
		mime := mimetype.Detect(file.Content)
		c.Header("Content-Type", mime.String())
		c.Header("Content-Disposition",
			fmt.Sprintf(
				"filename=%s_%d",
				file.Name,
				file.ID,
			),
		)
		c.String(200, string(file.Content))
	}
}
