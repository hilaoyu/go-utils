package utilElasticsearch

import "github.com/olivere/elastic/v7"

type ElasticsearchClientConfig struct {
	Addr        string
	User        string
	Password    string
	Debug       bool
	ErrorLogger elastic.Logger
	InfoLogger  elastic.Logger
}

type QueryFilter struct {
	QuerySource interface{}
}

func (fq *QueryFilter) Source() (interface{}, error) {
	return fq.QuerySource, nil
}

type QueryLastSort []interface{}
