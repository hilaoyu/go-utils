package utilOrm

import (
	"github.com/hilaoyu/go-utils/utilUuid"
	"gorm.io/gorm"
	"time"
)

type OrmModelBase struct {
	Id        string    `gorm:"primaryKey" json:"id,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at" form:"-"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:nano" json:"updated_at" form:"-"`
	//DeletedAt gorm.DeletedAt `json:"deleted_at" form:"-"`
}

func (om *OrmModelBase) BeforeCreate(tx *gorm.DB) (err error) {
	//fmt.Println("OrmModelBase BeforeCreate")
	if "" == om.Id {
		om.Id = utilUuid.UuidGenerate()
	}
	return nil
}
