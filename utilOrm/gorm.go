package utilOrm

import (
	"database/sql"
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

func (ug *UtilGorm) Debug(debug bool) *UtilGorm {
	if debug {
		ug.orm = ug.orm.Session(&gorm.Session{
			Logger: ug.orm.Logger.LogMode(logger.Info),
		})
	} else {
		ug.orm = ug.orm.Session(&gorm.Session{
			Logger: ug.orm.Logger.LogMode(logger.Error),
		})
	}
	return ug
}

func (ug *UtilGorm) Clauses(conds ...clause.Expression) *UtilGorm {
	ug.orm = ug.orm.Clauses(conds...)
	return ug
}

func (ug *UtilGorm) Scopes(funcs ...func(*gorm.DB) *gorm.DB) *UtilGorm {
	ug.orm = ug.orm.Scopes(funcs...)
	return ug
}

func (ug *UtilGorm) Raw(sql string, values ...interface{}) (err error) {
	result := ug.orm.Raw(sql, values...)
	err = result.Error
	return
}

func (ug *UtilGorm) Where(query interface{}, args ...interface{}) *UtilGorm {
	ug.orm = ug.orm.Where(query, args...)
	return ug
}
func (ug *UtilGorm) WithRelate(query string, args ...interface{}) *UtilGorm {
	ug.orm = ug.orm.Preload(query, args...)
	return ug
}

func (ug *UtilGorm) Begin(opts ...*sql.TxOptions) (err error) {
	ug.orm = ug.orm.Begin(opts...)
	err = ug.orm.Error
	return
}
func (ug *UtilGorm) Commit() (err error) {
	ug.orm = ug.orm.Commit()
	err = ug.orm.Error
	return
}
func (ug *UtilGorm) Rollback() (err error) {
	ug.orm = ug.orm.Rollback()
	err = ug.orm.Error
	return
}
func (ug *UtilGorm) SavePoint(name string) (err error) {
	ug.orm = ug.orm.SavePoint(name)
	err = ug.orm.Error
	return
}
func (ug *UtilGorm) RollbackTo(name string) (err error) {
	ug.orm = ug.orm.RollbackTo(name)
	err = ug.orm.Error
	return
}

func (ug *UtilGorm) TableQuery(tableName string, orderBy *[]string, args ...interface{}) *GormQuery {
	q := ug.orm.Table(tableName, args...)

	if nil != orderBy {
		for _, orderItem := range *orderBy {
			q = q.Order(orderItem)
		}

	}
	return &GormQuery{orm: q}
}
func (ug *UtilGorm) ModelQuery(model interface{}, orderBy *[]string) *GormQuery {
	q := ug.orm.Model(model)
	if nil != orderBy {
		for _, orderItem := range *orderBy {
			q = q.Order(orderItem)
		}

	}
	return &GormQuery{orm: q}
}

func (ug *UtilGorm) TableName(table string) (name string) {
	name = ug.orm.NamingStrategy.TableName(table)
	return
}
func (ug *UtilGorm) ModelTableName(model interface{}) (name string) {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	name = ug.TableName(t.Name())
	return
}

func (ug *UtilGorm) ModelRelatedLoad(model interface{}, related string, conds ...interface{}) (err error) {
	relatedValue := utils.GetInterfaceFiledValue(model, related)
	if !relatedValue.IsValid() {
		err = fmt.Errorf("关联关系错误")
		return
	}
	if !relatedValue.IsZero() {
		return
	}

	t := relatedValue.Type()
	isPtr := t.Kind() == reflect.Ptr
	if isPtr {
		t = t.Elem()
	}

	v := reflect.New(t).Interface()

	//fmt.Println("related err", related, nil == v)
	//continue
	qr := ug.ModelQuery(model, nil).orm.Association(related)
	err = qr.Find(v, conds...)

	if nil != err {
		if reflect.DeepEqual(err, gorm.ErrRecordNotFound) {
			err = nil
		}
		return
	}

	if isPtr {
		relatedValue.Set(reflect.ValueOf(v))
	} else {
		relatedValue.Set(reflect.ValueOf(v).Elem())
	}
	return
}
func (ug *UtilGorm) ModelRelatedAppend(model interface{}, related string, values ...interface{}) (err error) {
	err = ug.ModelQuery(model, nil).orm.Association(related).Append(values...)
	return
}
func (ug *UtilGorm) ModelRelatedReplace(model interface{}, related string, values ...interface{}) (err error) {
	err = ug.ModelQuery(model, nil).orm.Association(related).Replace(values...)
	return
}
func (ug *UtilGorm) ModelRelatedDelete(model interface{}, related string, values ...interface{}) (err error) {
	err = ug.ModelQuery(model, nil).orm.Association(related).Delete(values...)
	return
}
func (ug *UtilGorm) ModelRelatedClear(model interface{}, related string) (err error) {
	err = ug.ModelQuery(model, nil).orm.Association(related).Clear()
	return
}

func (ug *UtilGorm) Clone() *UtilGorm {
	tmp := &UtilGorm{
		orm: ug.orm.Scopes(),
	}
	return tmp
}
