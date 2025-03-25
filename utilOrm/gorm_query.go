package utilOrm

import (
	"github.com/hilaoyu/go-utils/utilHttp"
	"gorm.io/gorm"
)

type GormQuery struct {
	orm   *gorm.DB
	pager *utilHttp.Paginator
}

func (q *GormQuery) WithPager(p *utilHttp.Paginator) *GormQuery {
	//var total int64
	//q.orm.Count(&total)
	//p.SetTotal(total)
	q.orm = q.orm.Limit(p.PerPage).Offset(p.Offset())
	q.pager = p
	return q
}

func (q *GormQuery) WithRelate(query string, args ...interface{}) *GormQuery {
	q.orm = q.orm.Preload(query, args...)

	return q
}
func (q *GormQuery) IncludeDeleted() *GormQuery {
	q.orm = q.orm.Unscoped()
	return q
}

func (q *GormQuery) Where(query interface{}, args ...interface{}) *GormQuery {
	q.orm = q.orm.Where(query, args...)
	return q
}
func (q *GormQuery) Not(query interface{}, args ...interface{}) *GormQuery {
	q.orm = q.orm.Not(query, args...)
	return q
}

func (q *GormQuery) Or(query interface{}, args ...interface{}) *GormQuery {
	q.orm = q.orm.Or(query, args...)
	return q
}

func (q *GormQuery) Joins(query string, args ...interface{}) *GormQuery {
	q.orm = q.orm.Joins(query, args...)
	return q
}
func (q *GormQuery) InnerJoins(query string, args ...interface{}) *GormQuery {
	q.orm = q.orm.InnerJoins(query, args...)
	return q
}
func (q *GormQuery) Group(name string) *GormQuery {
	q.orm = q.orm.Group(name)
	return q
}
func (q *GormQuery) Having(query string, args ...interface{}) *GormQuery {
	q.orm = q.orm.Having(query, args...)
	return q
}
func (q *GormQuery) Order(value interface{}) *GormQuery {
	q.orm = q.orm.Order(value)
	return q
}
func (q *GormQuery) Limit(limit int) *GormQuery {
	q.orm = q.orm.Limit(limit)
	return q
}
func (q *GormQuery) Offset(offset int) *GormQuery {
	q.orm = q.orm.Offset(offset)
	return q
}
func (q *GormQuery) Select(query interface{}, args ...interface{}) *GormQuery {
	q.orm = q.orm.Select(query, args...)
	return q
}
func (q *GormQuery) Distinct(args ...interface{}) *GormQuery {
	q.orm = q.orm.Distinct(args...)
	return q
}

func (q *GormQuery) Count() (count int64, err error) {
	result := q.orm.Count(&count)
	err = result.Error
	return
}
func (q *GormQuery) Find(models interface{}, conds ...interface{}) (err error) {
	result := q.orm.Find(models, conds...)
	if nil != result.Error {
		err = result.Error
		return
	}
	if nil != q.pager {
		q.pager.Total, _ = q.Count()
	}
	return
}
func (q *GormQuery) FindInBatches(models interface{}, batchSize int, fc func(batch int) error) (err error) {
	result := q.orm.FindInBatches(models, batchSize, func(tx *gorm.DB, batch int) error {
		return fc(batch)
	})
	if nil != result.Error {
		err = result.Error
		return
	}
	return
}
func (q *GormQuery) Pluck(column string, model interface{}) (err error) {
	result := q.orm.Pluck(column, model)
	err = result.Error
	return
}
func (q *GormQuery) Scan(model interface{}) (err error) {
	result := q.orm.Scan(model)
	err = result.Error
	return
}
func (q *GormQuery) First(model interface{}, conds ...interface{}) (err error) {
	result := q.orm.First(model, conds...)
	err = result.Error
	return
}
func (q *GormQuery) FirstOrCreate(model interface{}, conds ...interface{}) (err error) {
	result := q.orm.FirstOrCreate(model, conds...)
	err = result.Error
	return
}

func (q *GormQuery) Create(value interface{}) (err error) {
	result := q.orm.Create(value)

	err = result.Error

	return
}
func (q *GormQuery) CreateInBatches(value interface{}, batchSize int) (err error) {

	result := q.orm.CreateInBatches(value, batchSize)

	err = result.Error

	return
}
func (q *GormQuery) Save(model interface{}) (err error) {
	result := q.orm.Save(model)

	err = result.Error

	return
}

func (q *GormQuery) Update(data interface{}) (err error) {
	//omits := append(q.orm.Statement.Omits, "created_at")
	result := q.orm.Updates(data)
	err = result.Error
	return
}

func (q *GormQuery) Delete(model interface{}, conds ...interface{}) (err error) {
	result := q.orm.Unscoped().Delete(model, conds...)
	err = result.Error
	return
}

func (q *GormQuery) SoftDelete(model interface{}, conds ...interface{}) (err error) {
	result := q.orm.Delete(model, conds...)
	err = result.Error
	return
}
