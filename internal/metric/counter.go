package metric

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
)

type CounterMetric uint64

func (c CounterMetric) TypeName() string {
	return "counter"
}

func (c CounterMetric) String() string {
	return strconv.FormatUint(uint64(c), 10)
}

func (c CounterMetric) FromString(value string) (Metric, error) {
	val, _err := strconv.ParseUint(value, 10, 64)
	return CounterMetric(uint64(c) + val), _err
}

func (c CounterMetric) FromPayload(value Payload, key string) (Metric, error) {
	mac := hmac.New(sha256.New, []byte(key))

	decodeString, err := hex.DecodeString(value.Hash)
	if err != nil {
		return nil, err
	}
	mac.Write(decodeString)
	if !hmac.Equal([]byte(value.Hash), mac.Sum(nil)) {
		return nil, errors.New("wrong hash sum")
	}

	return CounterMetric(uint64(c) + *value.CounterValue), nil
}

func (c CounterMetric) Payload(name, key string) Payload {
	value := uint64(c)
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(fmt.Sprintf("%s:counter:%d", name, c)))

	return Payload{
		Name:         name,
		MetricType:   c.TypeName(),
		CounterValue: &value,
		Hash:         string(mac.Sum(nil)),
	}
}
