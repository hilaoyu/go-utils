package utilMongodb

import (
	"fmt"
	"github.com/hilaoyu/go-utils/utilTime"
	"github.com/hilaoyu/go-utils/utilUuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type ModelMongodb interface {
	SetId(id string)
	GetId() string
	CollectionName() string
	Indexes() (indexes []mongo.IndexModel)
	SetClient(mc *MongodbClient)
	GetClient() (mc *MongodbClient)

	SetCreatedAt(t time.Time)
	SetUpdatedAt(t time.Time)
}
type ModelBaseMongodb struct {
	_client   *MongodbClient `json:"-" bson:"-"`
	Id        string         `json:"id,omitempty" bson:"_id"`
	CreatedAt string         `json:"created_at" bson:"created_at"`
	UpdatedAt string         `json:"updated_at" bson:"updated_at"`
	//DeletedAt gorm.DeletedAt `json:"deleted_at" form:"-"`
}

/*
	func (m *ModelBaseMongodb) CollectionName() string {
		return ""
	}
*/
func (m *ModelBaseMongodb) Indexes() (indexes []mongo.IndexModel) {

	return
}

func (m *ModelBaseMongodb) SetId(id string) {
	m.Id = id
}
func (m *ModelBaseMongodb) GetId() string {
	return m.Id
}

func (m *ModelBaseMongodb) SetClient(mc *MongodbClient) {
	m._client = mc
}

func (m *ModelBaseMongodb) GetClient() (mc *MongodbClient) {
	mc = m._client
	if nil == mc {
		mc = GetDefaultClient()
	}
	return
}

func (m *ModelBaseMongodb) SetCreatedAt(t time.Time) {
	m.CreatedAt = utilTime.TimeFormat(t)
}

func (m *ModelBaseMongodb) SetUpdatedAt(t time.Time) {
	m.UpdatedAt = utilTime.TimeFormat(t)
}

func CreateIndexes(m ModelMongodb) (err error) {

	indexes := m.Indexes()
	indexes = append(indexes, mongo.IndexModel{
		Keys: bson.D{{"_id", 1}},
	}, mongo.IndexModel{
		Keys: bson.D{{"created_at", -1}},
	}, mongo.IndexModel{
		Keys: bson.D{{"updated_at", -1}},
	})
	err = m.GetClient().CreateIndexes(m.CollectionName(), indexes)

	return
}
func Save(m ModelMongodb) (err error) {

	isUpdate := false
	if "" != m.GetId() {
		m.SetUpdatedAt(time.Now())
		isUpdate = true
	} else {
		m.SetId(utilUuid.UuidGenerate())
		m.SetCreatedAt(time.Now())
		m.SetUpdatedAt(time.Now())

	}

	if isUpdate {
		_, err = m.GetClient().Update(m.CollectionName(), bson.M{"_id": m.GetId()}, bson.M{"$set": m})
	} else {
		_, err = m.GetClient().InsertOne(m.CollectionName(), m)
	}

	return
}

func InsertMany(data []interface{}) (err error) {
	if len(data) <= 0 {
		return
	}

	tmp := data[0]
	m, ok := tmp.(ModelMongodb)
	if !ok {
		return fmt.Errorf("数据有错误")
	}
	_, err = m.GetClient().InsertMany(m.CollectionName(), data)

	return
}

func Read(m ModelMongodb) (err error) {
	if "" != m.GetId() {
		err = m.GetClient().Find(m.CollectionName(), bson.M{"_id": m.GetId()}, nil, m)
	}

	return
}
func Find(m ModelMongodb, filter interface{}, sort interface{}) (err error) {
	err = m.GetClient().Find(m.CollectionName(), filter, sort, m)
	return
}

func Delete(m ModelMongodb) (err error) {

	if "" != m.GetId() {
		err = m.GetClient().Delete(m.CollectionName(), bson.M{"_id": m.GetId()}, nil)
	}
	return
}

func Select(m ModelMongodb, results interface{}, filter interface{}, sort interface{}, limit int64, offset int64, total *int64) (err error) {
	err = m.GetClient().Select(m.CollectionName(), filter, sort, limit, offset, results, total)
	return
}
func Aggregate(m ModelMongodb, results interface{}, group bson.D, filter interface{}, sort interface{}) (err error) {

	err = m.GetClient().Aggregate(m.CollectionName(), group, filter, sort, results)
	return
}
