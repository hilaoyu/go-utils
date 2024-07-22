package utilOrm

import (
	"github.com/hilaoyu/go-utils/utilStr"
	"github.com/hilaoyu/go-utils/utilUuid"
	"gorm.io/gorm"
	"time"
)

type OrmModel interface {
	GetPrimaryKey() string
	GetPrimaryKeyFiledName() string
	GetPrimaryKeyFiledNameSnake() string
}

type OrmModelGormBase struct {
	Id        string    `gorm:"primaryKey" json:"id,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at" form:"-"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:nano" json:"updated_at" form:"-"`
	//DeletedAt gorm.DeletedAt `json:"deleted_at" form:"-"`
}

func (om *OrmModelGormBase) BeforeCreate(tx *gorm.DB) (err error) {
	//fmt.Println("OrmModelBase BeforeCreate")
	if "" == om.Id {
		om.Id = utilUuid.UuidGenerate()
	}
	return nil
}

func (om *OrmModelGormBase) GetPrimaryKey() string {
	return om.Id
}
func (om *OrmModelGormBase) GetPrimaryKeyFiledName() string {
	return "Id"
}
func (om *OrmModelGormBase) GetPrimaryKeyFiledNameSnake() string {
	return utilStr.ToSnake(om.GetPrimaryKeyFiledName())
}
