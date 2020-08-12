package utils

import (
	"bytes"
	"sort"
	"strings"
	"sync"
)

func SortedTags(tags map[string]string) string {
	var (
		bufferPool = sync.Pool{New: func() interface{} { return new(bytes.Buffer) }}
		ret        *bytes.Buffer
		size       int
		keys       []string
	)

	if tags == nil {
		return ""
	}

	size = len(tags)

	if size == 0 {
		return ""
	}

	ret = bufferPool.Get().(*bytes.Buffer)
	ret.Reset()
	defer bufferPool.Put(ret)

	if size == 1 {
		for k, v := range tags {
			ret.WriteString(k)
			ret.WriteString("=")
			ret.WriteString(v)
		}
		return ret.String()
	}

	keys = make([]string, size)
	i := 0
	for k := range tags {
		keys[i] = k
		i++
	}

	sort.Strings(keys)

	for j, key := range keys {
		ret.WriteString(key)
		ret.WriteString("=")
		ret.WriteString(tags[key])
		if j != size-1 {
			ret.WriteString(",")
		}
	}

	return ret.String()
}

func ReformatTag(str string) map[string]string {
	var (
		tagMap map[string]string
		tags   []string
	)

	if str == "" {
		return map[string]string{}
	}

	if strings.ContainsRune(str, ' ') {
		str = strings.Replace(str, " ", "", -1)
	}

	tagMap = make(map[string]string)

	tags = strings.Split(str, ",")
	for _, tag := range tags {
		idx := strings.IndexRune(tag, '=')
		if idx != -1 {
			tagMap[tag[:idx]] = tag[idx+1:]
		}
	}
	return tagMap
}
