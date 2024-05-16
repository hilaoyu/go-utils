package utilMongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongodbClient struct {
	client        *mongo.Client
	Host          string
	DbName        string
	User          string
	Password      string
	LoggerOptions *options.LoggerOptions
}

var defaultMongodbClient = &MongodbClient{}

func NewMongodbClient(host string, dbName string, user string, password string, loggerOptions ...*options.LoggerOptions) *MongodbClient {
	c := &MongodbClient{
		client:   nil,
		Host:     host,
		DbName:   dbName,
		User:     user,
		Password: password,
	}

	if len(loggerOptions) > 0 {
		c.LoggerOptions = loggerOptions[0]
	}

	return c
}
func SetDefaultConfig(host string, dbName string, user string, password string, loggerOptions ...*options.LoggerOptions) error {
	defaultMongodbClient = NewMongodbClient(host, dbName, user, password, loggerOptions...)
	return defaultMongodbClient.Connect()
}

func GetDefaultClient() *MongodbClient {
	/*if nil != defaultMongodbClient.Connect() {
		return nil
	}*/
	return defaultMongodbClient
}

func (mc *MongodbClient) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	clientOk := false
	if nil != mc.client {
		if nil == mc.client.Ping(context.TODO(), nil) {
			clientOk = true
		}
	}
	if !clientOk {
		clientOpts := []*options.ClientOptions{options.Client().SetHosts([]string{mc.Host})}
		if "" != mc.User {
			clientOpts = append(clientOpts, options.Client().SetAuth(options.Credential{
				//AuthMechanism:           "",
				//AuthMechanismProperties: nil,
				AuthSource: mc.DbName,
				Username:   mc.User,
				Password:   mc.Password,
				//PasswordSet:             false,
			}))
		}
		if nil != mc.LoggerOptions {
			clientOpts = append(clientOpts, options.Client().SetLoggerOptions(mc.LoggerOptions))
		}

		clientNew, err := mongo.Connect(ctx, clientOpts...)
		if nil != err {
			return fmt.Errorf("连接失败 %+v", err)

		}
		mc.client = clientNew
	}
	return nil
}

func (mc *MongodbClient) Ping() error {
	return mc.Connect()
}

func (mc *MongodbClient) Collection(collection string) (*mongo.Collection, error) {
	err := mc.Connect()
	if nil != err {
		return nil, err
	}

	return mc.client.Database(mc.DbName).Collection(collection), nil
}

func (mc *MongodbClient) CreateIndexes(collectionName string, indexModels []mongo.IndexModel) error {
	collection, err := mc.Collection(collectionName)
	if nil != err {
		return err
	}
	indexView := collection.Indexes()
	if MongodbHasIndexes(indexView) {
		return nil
	}

	opts := options.CreateIndexes().SetMaxTime(2 * time.Second)
	_, err = indexView.CreateMany(context.TODO(), indexModels, opts)
	if err != nil {
		return fmt.Errorf("创建索引失败 %+v", err)
	}
	return nil
}

type indexItem struct {
	Key  bson.D
	Name string
}

func MongodbHasIndexes(indexView mongo.IndexView) (found bool) {
	ctxBk := context.Background()
	cursor, err := indexView.List(ctxBk)

	if nil != err {
		return
	}
	defer cursor.Close(ctxBk)
	var index indexItem
	for cursor.Next(ctxBk) {
		err = cursor.Decode(&index)
		if nil == err && "_id_" != index.Name {
			found = true
		}
	}
	return
}

func (mc *MongodbClient) InsertOne(collectionName string, data interface{}) (id string, err error) {
	collection, err := mc.Collection(collectionName)
	if nil != err {
		return
	}
	cur, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		err = fmt.Errorf("保存失败 %+v", err)
		return
	}

	id, _ = cur.InsertedID.(string)

	return
}

func (mc *MongodbClient) InsertMany(collectionName string, data []interface{}) (ids []string, err error) {
	collection, err := mc.Collection(collectionName)
	if nil != err {
		return
	}
	cur, err := collection.InsertMany(context.TODO(), data)
	if err != nil {
		err = fmt.Errorf("保存失败 %+v", err)
		return
	}

	for _, id := range cur.InsertedIDs {
		if idStr, ok := id.(string); ok {
			ids = append(ids, idStr)
		}
	}

	return
}

func (mc *MongodbClient) Update(collectionName string, filter interface{}, update interface{}) (count int64, err error) {
	collection, err := mc.Collection(collectionName)
	if nil != err {
		return 0, err
	}

	updateResult, err := collection.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		return 0, fmt.Errorf("更新失败 %+v", err)
	}

	return updateResult.MatchedCount, nil
}

func (mc *MongodbClient) Find(collection string, filter interface{}, sort interface{}, result interface{}) (err error) {
	db, err := mc.Collection(collection)
	if nil != err {
		return err
	}

	opts := []*options.FindOptions{options.Find().SetSort(sort), options.Find().SetLimit(1)}

	cur, err := db.Find(context.Background(), filter, opts...)
	if err != nil {
		return fmt.Errorf("查询失败 %+v", err)
	}

	defer cur.Close(context.Background())

	if err = cur.Err(); err != nil {

		return fmt.Errorf("查询出错 %+v", err)
	}

	if cur.TryNext(context.TODO()) {
		err = cur.Decode(result)
		if err != nil {
			return fmt.Errorf("解析数据错误 %+v", err)
		}
	}

	return nil
}

func (mc *MongodbClient) Select(collection string, filter interface{}, sort interface{}, limit int64, offset int64, results interface{}, total *int64) (err error) {
	db, err := mc.Collection(collection)
	if nil != err {
		return err
	}

	opts := []*options.FindOptions{options.Find().SetSkip(offset)}
	if nil != sort {
		opts = append(opts, options.Find().SetSort(sort))
	}
	if limit > 0 {
		opts = append(opts, options.Find().SetLimit(limit))
	}
	cur, err := db.Find(context.Background(), filter, opts...)
	if err != nil {
		return fmt.Errorf("查询失败 %+v", err)
	}

	defer cur.Close(context.Background())

	if err = cur.Err(); err != nil {

		return fmt.Errorf("返回错误 %+v", err)
	}
	err = cur.All(context.Background(), results)
	if err != nil {
		return fmt.Errorf("解析数据错误 %+v", err)
	}

	if nil != total {
		totalNum, _ := db.CountDocuments(context.Background(), filter)
		*total = totalNum
	}

	return nil
}
func (mc *MongodbClient) Delete(collection string, filter interface{}, count *int64) (err error) {
	db, err := mc.Collection(collection)
	if nil != err {
		return err
	}

	cur, err := db.DeleteMany(context.Background(), filter)
	if err != nil {
		return fmt.Errorf("查询失败 %+v", err)
	}

	if nil != count {
		countNum := cur.DeletedCount
		count = &countNum
	}

	return nil
}

func (mc *MongodbClient) Aggregate(collection string, group bson.D, filter interface{}, sort interface{}, results interface{}) error {
	db, err := mc.Collection(collection)
	if nil != err {
		return err
	}
	pipeline := mongo.Pipeline{}

	if nil != filter {
		pipeline = append(pipeline, bson.D{bson.E{"$match", filter}})
	}
	if nil != group {
		pipeline = append(pipeline, bson.D{bson.E{"$group", group}})
	}
	if nil != sort {
		pipeline = append(pipeline, bson.D{bson.E{"$sort", sort}})
	}

	/**/

	cursor, err := db.Aggregate(context.Background(), pipeline)

	if err != nil {
		return fmt.Errorf("查询错误 %+v", err)
	}

	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), results); nil != err {
		return fmt.Errorf("解析数据错误 %+v", err)
	}

	return nil
}
