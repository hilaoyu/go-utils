package utilOrm

import (
	"github.com/hilaoyu/go-utils/utilHttp"
	"gorm.io/gorm"
)

type OrmQuery struct {
	orm *gorm.DB
}

func (q *OrmQuery) WithPager(p *utilHttp.Paginator) *OrmQuery {
	var total int64
	q.orm.Count(&total)
	p.SetTotal(total)
	q.orm = q.orm.Limit(p.PerPage).Offset(p.Offset())
	return q
}
func (q *OrmQuery) WithRelates(relates map[string][]interface{}) *OrmQuery {
	for relateKey, relateArgs := range relates {
		q.orm = q.orm.Preload(relateKey, relateArgs...)
	}

	return q
}

func (q *OrmQuery) Where(query interface{}, args ...interface{}) *OrmQuery {
	q.orm = q.orm.Where(query, args...)
	return q
}
func (q *OrmQuery) Not(query interface{}, args ...interface{}) *OrmQuery {
	q.orm = q.orm.Not(query, args...)
	return q
}

func (q *OrmQuery) Or(query interface{}, args ...interface{}) *OrmQuery {
	q.orm = q.orm.Or(query, args...)
	return q
}

func (q *OrmQuery) Joins(query string, args ...interface{}) *OrmQuery {
	q.orm = q.orm.Joins(query, args...)
	return q
}
func (q *OrmQuery) InnerJoins(query string, args ...interface{}) *OrmQuery {
	q.orm = q.orm.InnerJoins(query, args...)
	return q
}
func (q *OrmQuery) Group(name string) *OrmQuery {
	q.orm = q.orm.Group(name)
	return q
}
func (q *OrmQuery) Having(query string, args ...interface{}) *OrmQuery {
	q.orm = q.orm.Having(query, args...)
	return q
}
func (q *OrmQuery) Order(value interface{}) *OrmQuery {
	q.orm = q.orm.Order(value)
	return q
}
func (q *OrmQuery) Limit(limit int) *OrmQuery {
	q.orm = q.orm.Limit(limit)
	return q
}
func (q *OrmQuery) Offset(offset int) *OrmQuery {
	q.orm = q.orm.Offset(offset)
	return q
}

func (q *OrmQuery) Select(models interface{}, conds ...interface{}) (err error) {
	result := q.orm.Find(models, conds...)
	err = result.Error
	return
}
func (q *OrmQuery) First(model interface{}, conds ...interface{}) (err error) {
	result := q.orm.First(model, conds...)
	err = result.Error
	return
}
func (q *OrmQuery) FirstOrCreate(model interface{}, conds ...interface{}) (err error) {
	result := q.orm.FirstOrCreate(model, conds...)
	err = result.Error
	return
}

func (q *OrmQuery) Create(value interface{}) (err error) {
	result := q.orm.Create(value)

	err = result.Error

	return
}
func (q *OrmQuery) Save(model interface{}) (err error) {
	result := q.orm.Save(model)

	err = result.Error

	return
}

func (q *OrmQuery) Update(data interface{}) (err error) {
	result := q.orm.Updates(data)
	err = result.Error
	return
}

func (q *OrmQuery) Delete(model interface{}, conds ...interface{}) (err error) {
	result := q.orm.Delete(model, conds...)
	err = result.Error
	return
}
