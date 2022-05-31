package internal

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

const SchemeName = "http"
const ActionType = "update"
const SkipMessage = "Skip field: %s"
const RequestErrorMessage = "Request error: %s"

type Client struct {
	client       http.Client
	scheme       string
	address      string
	gaugeMetrics GaugeMetricsValues
	count        uint64
}

func NewHttpClient(settings Settings) Client {
	var metrics = make(GaugeMetricsValues)
	for _, m := range settings.Metrics.RuntimeMetrics {
		metrics[m] = 0
	}
	metrics[RandomValueMetricName] = 0
	client := http.Client{}

	return Client{
		client:       client,
		scheme:       SchemeName,
		address:      settings.GetAddress(),
		gaugeMetrics: metrics,
		count:        0,
	}
}

func (c *Client) SendMetrics() {
	for metricName, metricValue := range c.gaugeMetrics {
		c.post(metricName, fmt.Sprintf("%f", metricValue), GaugeMetric)
	}
	c.post(PollCountMetricName, strconv.Itoa(int(c.count)), CounterMetric)
	c.post(RandomValueMetricName, fmt.Sprintf("%f", rand.Float64()), GaugeMetric)
}

func (c *Client) CollectMetrics() {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	v := reflect.ValueOf(stats)
	vsType := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := vsType.Field(i).Name
		if _, ok := (*c).gaugeMetrics[field]; ok {
			value := v.Field(i).Interface()
			switch value.(type) {
			case uint64:
				(*c).gaugeMetrics[field] = float64(value.(uint64))
			case uint32:
				(*c).gaugeMetrics[field] = float64(value.(uint32))
			case float64:
				(*c).gaugeMetrics[field] = value.(float64)
			default:
				fmt.Printf(SkipMessage, field)
			}
		}
	}
	c.count++
}

func (c *Client) post(metricName string, metricValue string, metricAlias string) {
	endpoint := c.endpoint(metricAlias, metricName, metricValue)
	payload := c.payload(metricAlias, metricName, metricValue)
	request, err := http.NewRequest(http.MethodPost, endpoint.String(), bytes.NewBufferString(payload))
	if err != nil {
		log.Fatalf(RequestErrorMessage, err)
	}
	request.Header.Add("Content-Type", "text/plain")
	request.Header.Add("Content-Length", strconv.Itoa(len(payload)))
	if response, respErr := c.client.Do(request); respErr != nil {
		log.Printf(RequestErrorMessage, respErr)
	} else {
		log.Printf(metricName, response.Status)
	}
}

func (c *Client) endpoint(alias string, metricName string, metricValue string) url.URL {
	return url.URL{
		Scheme: c.scheme,
		Host:   c.address,
		Path:   strings.Join([]string{ActionType, alias, metricName, metricValue}, "/"),
	}
}

func (c *Client) payload(alias string, metricName string, metricValue string) string {
	return strings.Join([]string{ActionType, alias, metricName, metricValue}, "::")
}
