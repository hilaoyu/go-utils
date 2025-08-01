package utilOrm

import (
	"github.com/hilaoyu/go-utils/utilStr"
	"github.com/hilaoyu/go-utils/utilUuid"
	"gorm.io/gorm"
	"strings"
)

type OrmModel interface {
	GetPrimaryKey() string
	GetPrimaryKeyFiledName() string
	GetPrimaryKeyFiledNameSnake() string
}

type OrmModelGormBaseOnlyId struct {
	Id string `gorm:"primaryKey;size:36" json:"id,omitempty" form:"id"`
}
type OrmModelGormBaseWithCU struct {
	OrmModelGormBaseOnlyId
	CreatedAt OrmTime `gorm:"autoCreateTime;<-:create" json:"created_at,omitempty" form:"-"`
	UpdatedAt OrmTime `gorm:"autoUpdateTime:nano" json:"updated_at,omitempty" form:"-"`
}
type OrmModelGormBase struct {
	OrmModelGormBaseWithCU
	DeletedAt *OrmTime `json:"deleted_at,omitempty" form:"-"`
}

func NewOrmModelGormBaseWithCU(id string) (ormModelGormBaseWithCU OrmModelGormBaseWithCU) {
	id = strings.TrimSpace(id)
	ormModelGormBaseWithCU = OrmModelGormBaseWithCU{}
	if "" != id {
		ormModelGormBaseWithCU.Id = id
	}
	return
}
func NewOrmModelGormBase(id string) (ormModelGormBase OrmModelGormBase) {
	id = strings.TrimSpace(id)
	ormModelGormBase = OrmModelGormBase{}
	if "" != id {
		ormModelGormBase.Id = id
	}
	return
}

func (om *OrmModelGormBaseOnlyId) BeforeCreate(tx *gorm.DB) (err error) {
	om.generatePrimaryKey()
	return nil
}
func (om *OrmModelGormBaseOnlyId) generatePrimaryKey() {
	if "" == om.Id {
		om.Id = utilUuid.UuidGenerate()
	}
}

func (om *OrmModelGormBaseOnlyId) GetOrGeneratePrimaryKey() string {
	om.generatePrimaryKey()
	return om.Id
}
func (om *OrmModelGormBaseOnlyId) GetPrimaryKey() string {
	//om.generatePrimaryKey()
	return om.Id
}
func (om *OrmModelGormBaseOnlyId) GetPrimaryKeyFiledName() string {
	return "Id"
}
func (om *OrmModelGormBaseOnlyId) GetPrimaryKeyFiledNameSnake() string {
	return utilStr.ToSnake(om.GetPrimaryKeyFiledName())
}
