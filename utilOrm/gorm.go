package utilOrm

import (
	"fmt"
	"github.com/hilaoyu/go-utils/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"reflect"
	"time"
)

type UtilGorm struct {
	orm *gorm.DB
}

func NewUtilGormMysql(host string, port int, user string, password string, dbName string, tablePrefix string, timeout time.Duration) (utilOrm *UtilGorm, err error) {
	//
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s", user, password, host, port, dbName, timeout)
	//fmt.Println(dsn)
	//连接MYSQL, 获得DB类型实例，用于后面的数据库读写操作。
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   tablePrefix,
			SingularTable: true,
			//NameReplacer:  nil,
			//NoLowerCase:   false,
		},
	})
	if err != nil {
		err = fmt.Errorf("连接数据库失败, error: %+v", err)
		return
	}

	db = db.Omit(clause.Associations)
	db = db.Session(&gorm.Session{
		Logger: db.Logger.LogMode(logger.Error),
	})

	utilOrm = &UtilGorm{orm: db}
	return

}

func (o *UtilGorm) Debug(debug bool) *UtilGorm {
	if debug {
		o.orm = o.orm.Session(&gorm.Session{
			Logger: o.orm.Logger.LogMode(logger.Info),
		})
	} else {
		o.orm = o.orm.Session(&gorm.Session{
			Logger: o.orm.Logger.LogMode(logger.Error),
		})
	}
	return o
}
func (o *UtilGorm) Raw(sql string, values ...interface{}) (err error) {
	result := o.orm.Raw(sql, values...)
	err = result.Error
	return
}

func (o *UtilGorm) TableQuery(tableName string, orderBy *[]string, args ...interface{}) *OrmQuery {
	q := o.orm.Table(tableName, args...)

	if nil != orderBy {
		for _, orderItem := range *orderBy {
			q = q.Order(orderItem)
		}

	}
	return &OrmQuery{orm: q}
}
func (o *UtilGorm) ModelQuery(model interface{}, orderBy *[]string) *OrmQuery {
	q := o.orm.Model(model)

	if nil != orderBy {
		for _, orderItem := range *orderBy {
			q = q.Order(orderItem)
		}

	}
	return &OrmQuery{orm: q}
}

func (o *UtilGorm) ModelRelatedLoad(model interface{}, relates ...string) (err error) {
	for _, related := range relates {
		relatedValue := utils.GetInterfaceFiledValue(model, related)
		if !relatedValue.IsValid() {
			continue
		}

		t := relatedValue.Type()
		isPtr := t.Kind() == reflect.Ptr
		if isPtr {
			t = t.Elem()
		}

		v := reflect.New(t).Interface()

		//fmt.Println("related err", related, nil == v)
		//continue
		qr := o.ModelQuery(model, nil).orm.Association(related)
		err = qr.Find(v)

		if nil != err {
			if reflect.DeepEqual(err, gorm.ErrRecordNotFound) {
				err = nil
				continue
			}
			return
		}

		if qr.DB.RowsAffected <= 0 {
			continue
		}

		if isPtr {
			relatedValue.Set(reflect.ValueOf(v))
		} else {
			relatedValue.Set(reflect.ValueOf(v).Elem())
		}

	}
	return
}
func (o *UtilGorm) ModelRelatedAppend(model interface{}, related string, values ...interface{}) (err error) {
	err = o.ModelQuery(model, nil).orm.Association(related).Append(values...)
	return
}
func (o *UtilGorm) ModelRelatedReplace(model interface{}, related string, values ...interface{}) (err error) {
	err = o.ModelQuery(model, nil).orm.Association(related).Replace(values...)
	return
}
func (o *UtilGorm) ModelRelatedDelete(model interface{}, related string, values ...interface{}) (err error) {
	err = o.ModelQuery(model, nil).orm.Association(related).Delete(values...)
	return
}
func (o *UtilGorm) ModelRelatedClear(model interface{}, related string) (err error) {
	err = o.ModelQuery(model, nil).orm.Association(related).Clear()
	return
}
