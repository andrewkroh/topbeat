package elasticsearch

import (
	"fmt"
	"os"
	"testing"

	"github.com/elastic/libbeat/logp"
)

func TestBulk(t *testing.T) {
	if testing.Verbose() {
		logp.LogInit(logp.LOG_DEBUG, "", false, true, []string{"elasticsearch"})
	}
	if testing.Short() {
		t.Skip("Skipping in short mode, because it requires Elasticsearch")
	}

	client := GetTestingElasticsearch()
	index := fmt.Sprintf("packetbeat-unittest-%d", os.Getpid())

	ops := []map[string]interface{}{
		map[string]interface{}{
			"index": map[string]interface{}{
				"_index": index,
				"_type":  "type1",
				"_id":    "1",
			},
		},
		map[string]interface{}{
			"field1": "value1",
		},
	}

	body := make([]interface{}, 0, 10)
	for _, op := range ops {
		body = append(body, op)
	}

	params := map[string]string{
		"refresh": "true",
	}
	_, err := client.Bulk(index, "type1", params, body)
	if err != nil {
		t.Errorf("Bulk() returned error: %s", err)
	}

	params = map[string]string{
		"q": "field1:value1",
	}
	result, err := client.SearchURI(index, "type1", params)
	if err != nil {
		t.Errorf("SearchUri() returns an error: %s", err)
	}
	if result.Hits.Total != 1 {
		t.Errorf("Wrong number of search results: %d", result.Hits.Total)
	}

	_, err = client.Delete(index, "", "", nil)
	if err != nil {
		t.Errorf("Delete() returns error: %s", err)
	}
}

func TestEmptyBulk(t *testing.T) {
	if testing.Verbose() {
		logp.LogInit(logp.LOG_DEBUG, "", false, true, []string{"elasticsearch"})
	}
	if testing.Short() {
		t.Skip("Skipping in short mode, because it requires Elasticsearch")
	}

	client := GetTestingElasticsearch()
	index := fmt.Sprintf("packetbeat-unittest-%d", os.Getpid())

	body := make([]interface{}, 0, 10)

	params := map[string]string{
		"refresh": "true",
	}
	resp, err := client.Bulk(index, "type1", params, body)
	if err != nil {
		t.Errorf("Bulk() returned error: %s", err)
	}
	if resp != nil {
		t.Errorf("Unexpected response: %s", resp)
	}
}

func TestBulkMoreOperations(t *testing.T) {
	if testing.Verbose() {
		logp.LogInit(logp.LOG_DEBUG, "", false, true, []string{"elasticsearch"})
	}
	if testing.Short() {
		t.Skip("Skipping in short mode, because it requires Elasticsearch")
	}

	client := GetTestingElasticsearch()
	index := fmt.Sprintf("packetbeat-unittest-%d", os.Getpid())

	ops := []map[string]interface{}{
		map[string]interface{}{
			"index": map[string]interface{}{
				"_index": index,
				"_type":  "type1",
				"_id":    "1",
			},
		},
		map[string]interface{}{
			"field1": "value1",
		},
		map[string]interface{}{
			"delete": map[string]interface{}{
				"_index": index,
				"_type":  "type1",
				"_id":    "2",
			},
		},
		map[string]interface{}{
			"create": map[string]interface{}{
				"_index": index,
				"_type":  "type1",
				"_id":    "3",
			},
		},
		map[string]interface{}{
			"field1": "value3",
		},
		map[string]interface{}{
			"update": map[string]interface{}{
				"_id":    "1",
				"_index": index,
				"_type":  "type1",
			},
		},
		map[string]interface{}{
			"doc": map[string]interface{}{
				"field2": "value2",
			},
		},
	}

	body := make([]interface{}, 0, 10)
	for _, op := range ops {
		body = append(body, op)
	}

	params := map[string]string{
		"refresh": "true",
	}
	resp, err := client.Bulk(index, "type1", params, body)
	if err != nil {
		t.Errorf("Bulk() returned error: %s [%s]", err, resp)
		return
	}

	params = map[string]string{
		"q": "field1:value3",
	}
	result, err := client.SearchURI(index, "type1", params)
	if err != nil {
		t.Errorf("SearchUri() returns an error: %s", err)
	}
	if result.Hits.Total != 1 {
		t.Errorf("Wrong number of search results: %d", result.Hits.Total)
	}

	params = map[string]string{
		"q": "field2:value2",
	}
	result, err = client.SearchURI(index, "type1", params)
	if err != nil {
		t.Errorf("SearchUri() returns an error: %s", err)
	}
	if result.Hits.Total != 1 {
		t.Errorf("Wrong number of search results: %d", result.Hits.Total)
	}

	_, err = client.Delete(index, "", "", nil)
	if err != nil {
		t.Errorf("Delete() returns error: %s", err)
	}
}