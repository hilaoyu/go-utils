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
	Id        string         `gorm:"primaryKey;size:36" json:"id,omitempty" form:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime;<-:create;index:index_created_at" json:"created_at" form:"-"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime:nano;index:index_updated_at" json:"updated_at" form:"-"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" form:"-"`
}

func (om *OrmModelGormBase) BeforeCreate(tx *gorm.DB) (err error) {
	om.generatePrimaryKey()
	return nil
}
func (om *OrmModelGormBase) generatePrimaryKey() {
	if "" == om.Id {
		om.Id = utilUuid.UuidGenerate()
	}
}

func (om *OrmModelGormBase) GetPrimaryKey() string {
	om.generatePrimaryKey()
	return om.Id
}
func (om *OrmModelGormBase) GetPrimaryKeyFiledName() string {
	return "Id"
}
func (om *OrmModelGormBase) GetPrimaryKeyFiledNameSnake() string {
	return utilStr.ToSnake(om.GetPrimaryKeyFiledName())
}
