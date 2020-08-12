package utils

import (
	"bytes"
	"strconv"
	"sync"
)

func ReformatMetricUniqueString(endpoint, metric string, tags map[string]string, dsType string, step int) string {
	var (
		bufferPool = sync.Pool{New: func() interface{} { return new(bytes.Buffer) }}
		ret        *bytes.Buffer
	)

	ret = bufferPool.Get().(*bytes.Buffer)
	ret.Reset()
	defer bufferPool.Put(ret)

	if tags == nil || len(tags) == 0 {
		ret.WriteString(endpoint)
		ret.WriteString("/")
		ret.WriteString(metric)
		ret.WriteString("/")
		ret.WriteString(dsType)
		ret.WriteString("/")
		ret.WriteString(strconv.Itoa(step))

		return ret.String()
	}

	ret.WriteString(endpoint)
	ret.WriteString("/")
	ret.WriteString(metric)
	ret.WriteString("/")
	ret.WriteString(SortedTags(tags))
	ret.WriteString("/")
	ret.WriteString(dsType)
	ret.WriteString("/")
	ret.WriteString(strconv.Itoa(step))

	return ret.String()
}

func GetMetricPrimaryKey(endpoint, metric string, tags map[string]string) string {
	var (
		bufferPool = sync.Pool{New: func() interface{} { return new(bytes.Buffer) }}
		ret        *bytes.Buffer
	)

	ret = bufferPool.Get().(*bytes.Buffer)
	ret.Reset()
	defer bufferPool.Put(ret)

	if tags == nil || len(tags) == 0 {
		ret.WriteString(endpoint)
		ret.WriteString("/")
		ret.WriteString(metric)

		return ret.String()
	}

	ret.WriteString(endpoint)
	ret.WriteString("/")
	ret.WriteString(metric)
	ret.WriteString("/")
	ret.WriteString(SortedTags(tags))

	return ret.String()
}
