package mysql

import (
	"fmt"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/li4n0/revsuit/internal/database"
)

const FILE_SPEARATOR = ";"

type File struct {
	ID       uint   `gorm:"primarykey" form:"id" json:"id"`
	RecordID uint   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Name     string `json:"name"`
	Content  []byte `json:"-"`
}

//func (File) TableName() string {
//	return "mysql_files"
//}

func GetFile(c *gin.Context) {
	var (
		file        File
		mysqlRecord Record
	)
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  fmt.Errorf("param id missed").Error(),
			"result": nil,
		})
		return
	}

	database.DB.Model(&file).Where("id = ?", id).Find(&file)
	database.DB.Model(&mysqlRecord).Where("id = ?", file.RecordID).Find(&mysqlRecord)
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
				"filename=%s_%s_%d",
				strings.Replace(mysqlRecord.RemoteIP, ".", "_", -1),
				file.Name,
				file.ID,
			),
		)
		c.String(200, string(file.Content))
	}
}
