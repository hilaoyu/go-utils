package utilMongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongodbClientConf struct {
	Host          string
	DbName        string
	User          string
	Password      string
	LoggerOptions *options.LoggerOptions
}

type MongodbClient struct {
	*mongo.Database
	utilConf *MongodbClientConf
}

func NewMongodbClient(conf *MongodbClientConf) (client *MongodbClient, err error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := []*options.ClientOptions{options.Client().SetHosts([]string{conf.Host})}
	if "" != conf.User {
		clientOpts = append(clientOpts, options.Client().SetAuth(options.Credential{
			//AuthMechanism:           "",
			//AuthMechanismProperties: nil,
			AuthSource: conf.DbName,
			Username:   conf.User,
			Password:   conf.Password,
			//PasswordSet:             false,
		}))
	}
	if nil != conf.LoggerOptions {
		clientOpts = append(clientOpts, options.Client().SetLoggerOptions(conf.LoggerOptions))
	}

	c, err := mongo.Connect(ctx, clientOpts...)
	if nil != err {
		return
	}

	client = &MongodbClient{
		Database: c.Database(conf.DbName),
		utilConf: conf,
	}

	return
}

func (mc *MongodbClient) New() (client *MongodbClient, err error) {
	client, err = NewMongodbClient(mc.utilConf)
	return
}

func (mc *MongodbClient) CollectionCreateIndexes(collectionName string, indexModels []mongo.IndexModel) (err error) {
	collection := mc.Collection(collectionName)

	indexView := collection.Indexes()
	if MongodbHasIndexes(indexView) {
		return
	}

	opts := options.CreateIndexes().SetMaxTime(2 * time.Second)
	_, err = indexView.CreateMany(context.TODO(), indexModels, opts)
	if err != nil {
		return fmt.Errorf("创建索引失败 %+v", err)
	}
	return
}

func (mc *MongodbClient) CollectionInsertOne(collectionName string, data interface{}) (id string, err error) {
	collection := mc.Collection(collectionName)
	cur, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		err = fmt.Errorf("保存失败 %+v", err)
		return
	}

	id, _ = cur.InsertedID.(string)

	return
}

func (mc *MongodbClient) CollectionInsertMany(collectionName string, data []interface{}) (ids []string, err error) {
	collection := mc.Collection(collectionName)
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

func (mc *MongodbClient) CollectionUpdate(collectionName string, filter interface{}, update interface{}) (count int64, err error) {
	collection := mc.Collection(collectionName)

	updateResult, err := collection.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		return 0, fmt.Errorf("更新失败 %+v", err)
	}

	return updateResult.MatchedCount, nil
}

func (mc *MongodbClient) CollectionFind(collectionName string, filter interface{}, sort interface{}, result interface{}) (err error) {
	collection := mc.Collection(collectionName)

	opts := []*options.FindOptions{options.Find().SetSort(sort), options.Find().SetLimit(1)}

	cur, err := collection.Find(context.Background(), filter, opts...)
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

func (mc *MongodbClient) CollectionSelect(collectionName string, filter interface{}, sort interface{}, limit int64, offset int64, results interface{}, total *int64) (err error) {
	collection := mc.Collection(collectionName)
	opts := []*options.FindOptions{options.Find().SetSkip(offset)}
	if nil != sort {
		opts = append(opts, options.Find().SetSort(sort))
	}
	if limit > 0 {
		opts = append(opts, options.Find().SetLimit(limit))
	}
	cur, err := collection.Find(context.Background(), filter, opts...)
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
		totalNum, _ := collection.CountDocuments(context.Background(), filter)
		*total = totalNum
	}

	return nil
}
func (mc *MongodbClient) CollectionDelete(collectionName string, filter interface{}) (count int64, err error) {
	collection := mc.Collection(collectionName)

	cur, err := collection.DeleteMany(context.Background(), filter)
	if err != nil {
		err = fmt.Errorf("查询失败 %+v", err)
		return
	}

	count = cur.DeletedCount

	return
}

func (mc *MongodbClient) CollectionAggregate(collectionName string, group bson.D, filter interface{}, sort interface{}, results interface{}) error {
	collection := mc.Collection(collectionName)
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

	cursor, err := collection.Aggregate(context.Background(), pipeline)

	if err != nil {
		return fmt.Errorf("查询错误 %+v", err)
	}

	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), results); nil != err {
		return fmt.Errorf("解析数据错误 %+v", err)
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
