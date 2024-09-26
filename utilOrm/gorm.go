package utilOrm

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	std_ck "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/hilaoyu/go-utils/utils"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"reflect"
	"strings"
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
		Logger: logger.Default,
	})
	if err != nil {
		err = fmt.Errorf("连接数据库失败, error: %+v", err)
		return
	}

	db = db.Omit(clause.Associations)

	utilOrm = &UtilGorm{orm: db}
	return

}

func NewUtilGormClickHouse(host string, port int, user string, password string, dbName string, tablePrefix string, timeout time.Duration, useSsl ...bool) (utilOrm *UtilGorm, err error) {
	dbOption := &std_ck.Options{
		Addr: []string{fmt.Sprintf("%s:%d", host, port)},
		Auth: std_ck.Auth{
			Database: dbName,
			Username: user,
			Password: password,
		},
		Settings: std_ck.Settings{
			"max_execution_time": 60,
		},
		DialTimeout: timeout,
		Compression: &std_ck.Compression{
			Method: std_ck.CompressionLZ4,
			Level:  3,
		},
		//Debug: true,
	}
	if len(useSsl) > 0 && useSsl[0] {
		dbOption.TLS = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	clickhouseDb := std_ck.OpenDB(dbOption)

	//连接MYSQL, 获得DB类型实例，用于后面的数据库读写操作。
	db, err := gorm.Open(clickhouse.New(clickhouse.Config{
		Conn: clickhouseDb,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   tablePrefix,
			SingularTable: true,
			//NameReplacer:  nil,
			//NoLowerCase:   false,
		},
		Logger: logger.Default,
	})
	if err != nil {
		err = fmt.Errorf("连接数据库失败, error: %+v", err)
		return
	}

	db = db.Omit(clause.Associations)

	utilOrm = &UtilGorm{orm: db}
	return

}

func (ug *UtilGorm) Debug(debug bool) *UtilGorm {
	if debug {
		ug.orm.Logger = ug.orm.Logger.LogMode(logger.Info)
	} else {
		ug.orm.Logger = ug.orm.Logger.LogMode(logger.Error)
	}
	return ug
}

func (ug *UtilGorm) Original() *gorm.DB {
	return ug.orm
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

func (ug *UtilGorm) Session(config *gorm.Session) *UtilGorm {
	if nil == config {
		config = &gorm.Session{}
	}
	return &UtilGorm{orm: ug.orm.Session(config)}
}
func (ug *UtilGorm) Select(query interface{}, args ...interface{}) *UtilGorm {
	ug.orm = ug.orm.Select(query, args...)
	return ug
}
func (ug *UtilGorm) IncludeDeleted() *UtilGorm {
	ug.orm = ug.orm.Unscoped()
	return ug
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
func (ug *UtilGorm) ModelRead(model interface{}) error {
	return ug.ModelQuery(model, nil).First(model)
}
func (ug *UtilGorm) ModelSave(model interface{}) error {
	return ug.ModelQuery(model, nil).Save(model)
}
func (ug *UtilGorm) ModelDelete(model interface{}) error {
	return ug.ModelQuery(model, nil).Delete(model)
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

	t := relatedValue.Type()
	isPtr := t.Kind() == reflect.Ptr
	if isPtr {
		t = t.Elem()
	}

	v := reflect.New(t).Interface()

	//fmt.Println("related err", related, nil == v)
	//continue
	//qr := ug.Clone().ModelQuery(model, nil).orm.Association(related)
	//err = qr.Find(v, conds...)

	association := ug.Clone().ModelQuery(model, nil).orm.Association(related)
	var (
		queryConds = association.Relationship.ToQueryConditions(association.DB.Statement.Context, association.DB.Statement.ReflectValue)
		modelValue = reflect.New(association.Relationship.FieldSchema.ModelType).Interface()
		tx         = association.DB.Model(modelValue)
	)

	if association.Relationship.JoinTable != nil {
		if !tx.Statement.Unscoped && len(association.Relationship.JoinTable.QueryClauses) > 0 {
			joinStmt := gorm.Statement{DB: tx, Context: tx.Statement.Context, Schema: association.Relationship.JoinTable, Table: association.Relationship.JoinTable.Table, Clauses: map[string]clause.Clause{}}
			for _, queryClause := range association.Relationship.JoinTable.QueryClauses {
				joinStmt.AddClause(queryClause)
			}
			joinStmt.Build("WHERE")
			if len(joinStmt.SQL.String()) > 0 {
				tx.Clauses(clause.Expr{SQL: strings.Replace(joinStmt.SQL.String(), "WHERE ", "", 1), Vars: joinStmt.Vars})
			}
		}

		tx = tx.Session(&gorm.Session{QueryFields: true}).Clauses(clause.From{Joins: []clause.Join{{
			Table: clause.Table{Name: association.Relationship.JoinTable.Table},
			ON:    clause.Where{Exprs: queryConds},
		}}})
	} else {
		tx.Clauses(clause.Where{Exprs: queryConds})
	}
	result := tx.Find(v, conds...)

	//fmt.Println(related, reflect.Struct == t.Kind(), " : RowsAffected :", result.RowsAffected)
	if nil != result.Error || result.RowsAffected <= 0 {
		if !reflect.DeepEqual(result.Error, gorm.ErrRecordNotFound) {
			err = result.Error
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
	return ug.Session(nil)
}

func ErrorIsOrmNotFound(err error) bool {
	return reflect.DeepEqual(err, gorm.ErrRecordNotFound)
}
