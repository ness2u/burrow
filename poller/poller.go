package main

import "net/http"
import "io/ioutil"
import "encoding/json"
import "fmt"

type ClusterList struct {
	Error    bool     `json:"error"`
	Message  string   `json:"message"`
	Clusters []string `json:"clusters"`
}

type ConsumerList struct {
	Error     bool     `json:"error"`
	Message   string   `json:"message"`
	Consumers []string `json:"consumers"`
}

type ConsumerLag struct {
	Error   bool           `json:"error"`
	Message string         `json:"message"`
	Status  ConsumerStatus `json:"status"`
}

type ConsumerStatus struct {
	Status         string              `json:"status"`
	TotalLag       int64               `json:"totallag"`
	Cluster        string              `json:"cluster"`
	Complete       bool                `json:"complete"`
	Group          string              `json:"group"`
	MaxLag         int                 `json:"maxlag"`
	PartitionCount int                 `json:"partition_count"`
	Partitions     []ConsumerPartition `json:"partitions"`
}

type ConsumerPartition struct {
	Topic     string    `json:"topic"`
	Status    string    `json:"status"`
	Partition int       `json:"partition"`
	Start     LagWindow `json:"start"`
	End       LagWindow `json:"end"`
}

type LagWindow struct {
	Lag       int   `json:"lag"`
	MaxOffset int64 `json:"max_offset"`
	Offset    int64 `json:"offset"`
	Timestamp int64 `json:"timestamp"`
}

type GroupLagOutput struct {
	Group          string `json:"group"`
	TotalLag       int64  `json:"totallag"`
	MaxLag         int    `json:"maxlag"`
	PartitionCount int    `json:"partition_count"`
}

type PartitionLagOutput struct {
	Timestamp int64  `json:"timestamp"`
	Group     string `json:"group"`
	Topic     string `json:"topic"`
	Partition int    `json:"partition"`
	Lag       int    `json:"lag"`
}

func main() {

	baseUri := "http://localhost:8000"

	clusters := getClusters(baseUri)
	for _, cluster := range clusters {

		consumers := getConsumers(baseUri, cluster)
		for _, consumer := range consumers {

			status := getConsumerStatus(baseUri, cluster, consumer)
			generateOutput(status)
		}

	}
}

func createGroupOutput(status ConsumerStatus) string {

	groupLagOutput := GroupLagOutput{status.Group, status.TotalLag, status.MaxLag, status.PartitionCount}
	jsonB, err := json.Marshal(groupLagOutput)
	if err != nil {
		panic("unable to encode json")
	}

	return string(jsonB)
}

func createPartitionOutputs(status ConsumerStatus) []string {

	outputs := []string{}
	for _, partition := range status.Partitions {

		partitionOutput := PartitionLagOutput{
			partition.End.Timestamp,
			status.Group,
			partition.Topic,
			partition.Partition,
			partition.End.Lag}

		jsonB, err := json.Marshal(partitionOutput)
		if err != nil {
			panic("unable to encode json")
		}
		json := string(jsonB)

		outputs = append(outputs, json)
	}

	return outputs
}

func generateOutput(status ConsumerStatus) {
	groupOutput := createGroupOutput(status)
	partitionOutputs := createPartitionOutputs(status)

	println(groupOutput)
	for _, partitionOutput := range partitionOutputs {
		println(partitionOutput)
	}
}

func getConsumerStatus(baseUri string, cluster string, consumer string) ConsumerStatus {

	uri := fmt.Sprintf("%s/v2/kafka/%s/consumer/%s/lag", baseUri, cluster, consumer)

	bytes := httpGet(uri)

	var cl ConsumerLag
	err := json.Unmarshal(bytes, &cl)
	if err != nil {
		panic("oh noes, couldn't read consumer lag")
	}

	return cl.Status

}

func getConsumers(baseUri string, cluster string) []string {

	uri := fmt.Sprintf("%s/v2/kafka/%s/consumer", baseUri, cluster)

	bytes := httpGet(uri)

	var cl ConsumerList
	err := json.Unmarshal(bytes, &cl)
	if err != nil {
		panic("oh noes, couldn't read consumer list")
	}

	return cl.Consumers
}

func getClusters(baseUri string) []string {

	uri := fmt.Sprintf("%s/v2/kafka", baseUri)

	bytes := httpGet(uri)

	var cl ClusterList
	err2 := json.Unmarshal(bytes, &cl)
	if err2 != nil {
		panic("oh noes")
	}

	return cl.Clusters
}

func httpGet(uri string) []byte {

	resp, err := http.Get(uri)
	if err != nil {
		panic(fmt.Sprintf("couldn't get %s", uri))
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic("couldn't read body")
	}

	bodyStr := string(body)
	// fmt.Println(bodyStr)

	bytes := []byte(bodyStr)
	return bytes
}
