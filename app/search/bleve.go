package search

import (
	"fmt"
	"strings"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/lang/en"
	"github.com/blevesearch/bleve/analysis/token/keyword"
	"github.com/blevesearch/bleve/mapping"
	"github.com/revel/revel"
	//_ "github.com/blevesearch/zap/v12" //forced index v12
)

var Index bleve.Index

//CreateOrOpenIndex ...
func CreateOrOpenIndex(indexName string) {
	index, err := bleve.Open(indexName)

	if err != nil {
		revel.AppLog.Info("Main BLEVE index not found, create a new one. ", "Index Name: ", indexName)
		mapping, _ := genericBuildIndexMapping()

		Index, err = bleve.New(indexName, mapping)
	} else {
		revel.AppLog.Info("Main BLEVE index found!", "Index Name: ", indexName)
		Index = index
	}
}

//Search ...
func Search(searchText string, searchFields []string, from int, size int) *bleve.SearchResult {
	searchText = strings.TrimSpace(searchText)

	results := &bleve.SearchResult{}

	query := bleve.NewQueryStringQuery(searchText)
	search := bleve.NewSearchRequest(query)

	//Init search params for pagination
	search.From = from

	if size == 0 {
		search.Size, _ = revel.Config.Int("gorm.result_set_size")
	} else {
		search.Size = size
	}

	results = simpleSearch(searchText, searchFields, search)

	return results
}

func simpleSearch(searchText string, searchFields []string, search *bleve.SearchRequest) *bleve.SearchResult {
	search.SortBy([]string{"-_score"})

	if searchFields != nil {
		search.Fields = searchFields
	} else {
		search.Fields = []string{"*"}
	}

	searchResults, err := Index.Search(search)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(searchResults)
	return searchResults
}

//Old method for search over senteces. Not used anymore due to the use of the query system: NewQueryStringQuery
func sentenceSearch(searchText string, searchFields []string, search *bleve.SearchRequest) *bleve.SearchResult {
	search.SortBy([]string{"-_score"})

	if searchFields != nil {
		search.Fields = searchFields
	} else {
		search.Fields = []string{"*"}
	}

	searchResults, err := Index.Search(search)

	if err != nil {
		fmt.Println(err)
	}

	return searchResults
}

func genericBuildIndexMapping() (mapping.IndexMapping, error) {

	// a generic reusable mapping for english text
	englishTextFieldMapping := bleve.NewTextFieldMapping()
	englishTextFieldMapping.Analyzer = en.AnalyzerName

	// a generic reusable mapping for keyword text
	keywordFieldMapping := bleve.NewTextFieldMapping()
	keywordFieldMapping.Analyzer = keyword.Name

	entityMapping := bleve.NewDocumentMapping()
	entityMapping.AddFieldMappingsAt("name", englishTextFieldMapping)
	entityMapping.AddFieldMappingsAt("description", englishTextFieldMapping)

	nodeMapping := bleve.NewDocumentMapping()
	nodeMapping.AddFieldMappingsAt("name", englishTextFieldMapping)

	nodeInstanceMapping := bleve.NewDocumentMapping()
	nodeInstanceMapping.AddFieldMappingsAt("value", englishTextFieldMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("entity", entityMapping)
	indexMapping.AddDocumentMapping("node", nodeMapping)
	indexMapping.AddDocumentMapping("nodeInstance", nodeInstanceMapping)

	indexMapping.TypeField = "type"
	indexMapping.DefaultAnalyzer = "en"

	return indexMapping, nil
}
