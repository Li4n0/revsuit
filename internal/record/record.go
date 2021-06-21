package record

import (
	"time"
)

var recordChan = make(chan BaseRecord, 1000)

type Record interface {
	GetFlag() string
	PushToClient()
	Notice()
}

type BaseRecord struct {
	Record      `gorm:"-" json:"-" notice:"-"`
	ID          uint      `gorm:"primarykey" form:"id" json:"id" notice:"-"`
	RuleName    string    `gorm:"index" form:"rule_name" json:"rule_name" notice:"rule"`
	Flag        string    `gorm:"index" form:"flag" json:"flag" `
	RemoteIP    string    `gorm:"index" form:"remote_ip" json:"remote_ip" notice:"remote_ip"`
	IpArea      string    `form:"ip_area" json:"ip_area" notice:"ip_area"`
	RequestTime time.Time `gorm:"index" form:"request_time" json:"request_time" notice:"-"`
}

func (b BaseRecord) GetFlag() string {
	return b.Flag
}

func (b BaseRecord) PushToClient() {
	recordChan <- b
}

func (b BaseRecord) Notice() {
}

func Channel() chan BaseRecord {
	return recordChan
}
