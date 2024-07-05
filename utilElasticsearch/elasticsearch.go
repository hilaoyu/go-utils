package utilElasticsearch

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hilaoyu/go-utils/utils"
	"github.com/olivere/elastic/v7"
	"net/http"
	"reflect"
	"time"
)

type ElasticsearchClient struct {
	*elastic.Client
	config *ElasticsearchClientConfig
}

type esErrorLogger struct{}
type esInfoLogger struct{}

// 实现输出
func (esErrorLogger) Printf(format string, v ...interface{}) {
	fmt.Println("ELASTIC ERROR "+fmt.Sprintf(format, v...), nil)
}
func (esInfoLogger) Printf(format string, v ...interface{}) {
	fmt.Println("ELASTIC INFO " + fmt.Sprintf(format, v...))
}

func NewElasticsearchClient(conf *ElasticsearchClientConfig) (esClient *ElasticsearchClient, err error) {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	httpClient := &http.Client{
		Transport: tr,
	}

	var errorLogger elastic.Logger
	if nil != conf.ErrorLogger {
		errorLogger = conf.ErrorLogger
	} else {
		errorLogger = esErrorLogger{}
	}

	clientOptions := []elastic.ClientOptionFunc{
		elastic.SetURL(conf.Addr),
		elastic.SetHttpClient(httpClient),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(10 * time.Second),
		elastic.SetGzip(true),
		elastic.SetErrorLog(errorLogger),
	}

	if conf.Debug {
		var infoLogger elastic.Logger
		if nil != conf.InfoLogger {
			infoLogger = conf.InfoLogger
		} else {
			infoLogger = esInfoLogger{}
		}
		clientOptions = append(clientOptions, elastic.SetInfoLog(infoLogger))
	}

	if "" != conf.User && "" != conf.Password {
		clientOptions = append(clientOptions, elastic.SetBasicAuth(conf.User, conf.Password))
	}

	c, err := elastic.NewClient(clientOptions...)
	if err != nil {
		err = fmt.Errorf("连接失败: %+v", err)
		return
	}

	esClient = &ElasticsearchClient{
		Client: c,
		config: conf,
	}

	return

}

func (esClient *ElasticsearchClient) New() (newClient *ElasticsearchClient, err error) {
	newClient, err = NewElasticsearchClient(esClient.config)

	return
}
func (esClient *ElasticsearchClient) GetMap(indexName string) (mapping string, err error) {
	if "" == indexName {
		err = fmt.Errorf("indexName 不能为空")
		return
	}

	mappingResult, err := esClient.GetMapping().Index(indexName).Do(context.Background())
	if err != nil {
		err = fmt.Errorf("查询mapping失败: index: %s ,err: %+v", indexName, err)
		return
	}

	mappingTemp, ok := mappingResult[indexName]
	if !ok {
		err = fmt.Errorf("解析maping失败: index: %s ,err: %+v", indexName, errors.New("索引不存在"))
		return
	}

	mappingTempMap, ok := mappingTemp.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("解析maping失败: index: %s ,err: %+v", indexName, errors.New("索引mapping结构错误"))
		return
	}

	mappingTempMappings, ok := mappingTempMap["mappings"]
	if !ok {
		err = fmt.Errorf("解析maping失败: index: %s ,err: %+v", indexName, errors.New("mapping不存在"))
		return
	}
	jsonByte, err := json.Marshal(mappingTempMappings)
	if nil != err {
		err = fmt.Errorf("json Marshal 失败: index: %s ,err: %+v", indexName, err)
		return
	}

	mapping = string(jsonByte)

	return
}
func (esClient *ElasticsearchClient) SetMap(indexName string, mapping string) (err error) {
	if "" == indexName || "" == mapping {
		err = fmt.Errorf("indexName 或 mapping 不能为空")
		return
	}

	exists, err := esClient.IndexExists(indexName).Do(context.Background())
	if err != nil {
		err = fmt.Errorf("查询index失败: index: %s ,err: %+v", indexName, err)
		return
	}

	errFormat := "%s mapping 错误: index:%s ,err: %+v "
	if !exists {
		mapping = `{
    "settings": {
        "number_of_shards": 1,
        "number_of_replicas": 1
    },
    "mappings": ` + mapping + `
}`
		createIndex, err1 := esClient.CreateIndex(indexName).BodyString(mapping).Do(context.Background())
		if err1 != nil {
			err = fmt.Errorf(errFormat, "创建", indexName, err1)
			return
		}
		if !createIndex.Acknowledged {
			err = fmt.Errorf(errFormat, "创建", indexName, err1)
			return
		}
		return
	}

	putResp, err := esClient.PutMapping().Index(indexName).BodyString(mapping).Do(context.TODO())
	if err != nil {
		err = fmt.Errorf(errFormat, "更新", indexName, err)
		return
	}
	if putResp == nil {
		err = fmt.Errorf(errFormat, "更新", indexName, errors.New("返回为空"))
		return
	}
	if !putResp.Acknowledged {
		err = fmt.Errorf(errFormat, "更新", indexName, errors.New("!createIndex.Acknowledged"))
		return
	}

	return
}

func (esClient *ElasticsearchClient) SaveData(indexName string, data []interface{}) (err error) {

	if "" == indexName {
		err = fmt.Errorf("indexName 不能为空")
		return
	}

	bulkRequest := esClient.Bulk()

	for _, item := range data {
		bulkRequest = bulkRequest.Add(elastic.NewBulkIndexRequest().Index(indexName).Doc(item))
	}

	bulkTotal := bulkRequest.NumberOfActions()

	bulkResponse, err := bulkRequest.Do(context.TODO())
	errFormat := "保存失败: index: %s ,err: %+v "
	if err != nil {
		err = fmt.Errorf(errFormat, indexName, err)
		return
	}
	if bulkResponse == nil {
		err = fmt.Errorf(errFormat, indexName, errors.New("返回为空"))
		return
	}

	//godump.Dump(bulkResponse)

	if bulkRequest.NumberOfActions() > 0 {
		err = fmt.Errorf(errFormat, indexName, fmt.Errorf("没有全部完成,总提交: %d 个,有 %d 个没成功", bulkTotal, bulkRequest.NumberOfActions()))
		return
	}

	if bulkResponse.Errors {
		errMsg := ""
		for _, errItem := range bulkResponse.Failed() {
			errMsg += fmt.Sprintf("%s : %s ;", errItem.Error.Type, errItem.Error.Reason)
		}

		err = fmt.Errorf(errFormat, indexName, fmt.Errorf(errMsg))
		return
	}

	//fmt.Println(bulk_total,bulkRequest.NumberOfActions())
	return
}
func (esClient *ElasticsearchClient) Update(indexName string, filter *QueryFilter, script *elastic.Script) (count int64, err error) {

	if "" == indexName {
		err = fmt.Errorf("indexName 不能为空")
		return
	}

	result, err := esClient.UpdateByQuery().Index(indexName).Query(filter).Script(script).Do(context.Background())
	count = result.Updated
	return
}

func (esClient *ElasticsearchClient) Delete(indexName string, filter *QueryFilter) (count int64, err error) {

	if "" == indexName {
		err = fmt.Errorf("indexName 不能为空")
		return
	}

	result, err := esClient.DeleteByQuery().
		Index(indexName).
		Query(filter).
		Do(context.Background())

	count = result.Deleted

	return
}

func (esClient *ElasticsearchClient) Select(results interface{}, indexName string, filter *QueryFilter, sort map[string]bool, limit int64, offset int64, lastSort *QueryLastSort) (total int64, err error) {
	item, err := utils.MakeInstanceFromSlice(results)
	if nil != err {
		return
	}
	if nil == filter {
		matchAll := elastic.NewMatchAllQuery()
		matchAllSource, _ := matchAll.Source()
		filter = &QueryFilter{QuerySource: matchAllSource}
	}
	request := esClient.Search().
		TrackTotalHits(true).
		Index(indexName). // search in index "twitter"
		Query(filter).    // specify the query
		Size(int(limit)). // take documents 0-9
		Pretty(true)      // pretty print request and response JSON

	if nil != lastSort && len(*lastSort) > 0 {
		request = request.SearchAfter(*lastSort...)
	} else {
		request = request.From(int(offset))
	}

	if len(sort) > 0 {
		for sField, sAsc := range sort {
			request = request.Sort(sField, sAsc)
		}
	}

	/*f, _ := filter.Source()
	b, e1 := json.Marshal(f)
	fmt.Println("q", string(b), e1)
	*/
	searchResult, err := request.Do(context.Background()) // execute
	if err != nil {
		err = fmt.Errorf("查询失败,index: %s ,err: %+v ", indexName, err)
		return
	}

	if len(searchResult.Hits.Hits) > 0 {
		*lastSort = searchResult.Hits.Hits[len(searchResult.Hits.Hits)-1].Sort
	}

	total = searchResult.TotalHits()

	resultsTemp := searchResult.Each(reflect.TypeOf(item))

	resultsJson, err := json.Marshal(resultsTemp)
	if nil != err {
		return
	}

	err = json.Unmarshal(resultsJson, results)

	return

}
func (esClient *ElasticsearchClient) Aggregate(indexName string, group map[string]elastic.Aggregation, filter elastic.Query) (searchResult *elastic.SearchResult, err error) {

	request := esClient.Search().
		TrackTotalHits(true).
		Index(indexName). // search in index "twitter"
		Query(filter).    // specify the query
		Pretty(true)      // pretty print request and response JSON

	if len(group) > 0 {
		for gk, gv := range group {
			request = request.Aggregation(gk, gv)
		}
	}

	searchResult, err = request.Do(context.Background()) // execute

	return

}
