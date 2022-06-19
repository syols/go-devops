package metric

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
)

type GaugeMetric float64

func (g GaugeMetric) TypeName() string {
	return "gauge"
}

func (g GaugeMetric) String() string {
	return fmt.Sprintf("%.3f", g)
}

func (g GaugeMetric) FromString(value string) (Metric, error) {
	val, _err := strconv.ParseFloat(value, 64)
	return GaugeMetric(val), _err
}

func (g GaugeMetric) FromPayload(value Payload, key string) (Metric, error) {
	mac := hmac.New(sha256.New, []byte(key))

	decodeString, err := hex.DecodeString(value.Hash)
	if err != nil {
		return nil, err
	}

	mac.Write(decodeString)
	if !hmac.Equal([]byte(value.Hash), mac.Sum(nil)) {
		return nil, errors.New("wrong hash sum")
	}
	return GaugeMetric(*value.GaugeValue), nil
}

func (g GaugeMetric) Payload(name, key string) Payload {
	value := float64(g)
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(fmt.Sprintf("%s:gauge:%f", name, g)))
	return Payload{
		Name:       name,
		MetricType: g.TypeName(),
		GaugeValue: &value,
		Hash:       hex.EncodeToString(mac.Sum(nil)),
	}
}
