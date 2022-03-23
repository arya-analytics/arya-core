package timeseries

import (
	"context"
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/go-redis/redis/v8"
	"reflect"
	"strconv"
)

// |||| COMMANDS ||||

type Command string

const (
	CMDCreateSeries  Command = "TS.CREATE"
	CMDCreateSamples Command = "TS.MADD"
	CMDGet           Command = "TS.GET"
	CMDRange         Command = "TS.RANGE"
	CMDInfo          Command = "TS.INFO"
)

// |||| SAMPLE ||||

type Sample struct {
	Key       string
	Value     float64
	Timestamp telem.TimeStamp
}

func NewSampleFromRes(key string, res interface{}) (Sample, error) {
	resVal := reflect.ValueOf(res)
	if resVal.Kind() != reflect.Slice || resVal.Len() < 2 {
		return Sample{}, fmt.Errorf("response invalid: %s", res)
	}
	ts := telem.TimeStamp(resVal.Index(0).Interface().(int64))
	val, err := strconv.ParseFloat(resVal.Index(1).Interface().(string), 64)
	if err != nil {
		return Sample{}, err
	}
	return Sample{
		Key:       key,
		Timestamp: ts,
		Value:     val,
	}, err
}

func (s Sample) args() (args []interface{}) {
	return []interface{}{
		s.Key,
		int64(s.Timestamp),
		s.Value,
	}
}

// |||| CLIENT ||||

type Client struct {
	*redis.Client
}

func NewWrap(redis *redis.Client) *Client {
	return &Client{redis}
}

type CreateOptions struct {
	Retention int64
}

func (co CreateOptions) args(key string) []interface{} {
	return []interface{}{
		key,
		fmt.Sprintf("RETENTION %v", co.Retention),
	}
}

func (c *Client) exec(ctx context.Context, cmd Command, args ...interface{}) *redis.Cmd {
	args = append([]interface{}{reflect.ValueOf(cmd).String()}, args...)
	return c.Do(ctx, args...)
}

func (c *Client) TSCreateSeries(ctx context.Context, key string,
	opts CreateOptions) *redis.Cmd {
	return c.exec(ctx, CMDCreateSeries, opts.args(key)...)
}

func (c *Client) TSCreateSamples(ctx context.Context, samples ...Sample) *redis.Cmd {
	var sampleArgs []interface{}
	for _, sample := range samples {
		sampleArgs = append(sampleArgs, sample.args()...)
	}
	return c.exec(ctx, CMDCreateSamples, sampleArgs...)
}

func (c *Client) TSGet(ctx context.Context, key string) *redis.Cmd {
	return c.exec(ctx, CMDGet, key)
}

func (c *Client) TSGetAll(ctx context.Context, key string) *redis.Cmd {
	return c.TSGetRange(ctx, key, telem.AllTime())
}

func (c *Client) TSGetRange(ctx context.Context, key string, tr telem.TimeRange) *redis.Cmd {
	return c.exec(ctx, CMDRange, key, parseStartTS(tr), parseEndTS(tr))
}

const (
	timeRangeMin = "-"
	timeRangeMax = "+"
)

func parseStartTS(tr telem.TimeRange) string {
	if tr.Start() != telem.TimeStampMin {
		return parseTS(tr.Start())
	}
	return timeRangeMin
}

func parseEndTS(tr telem.TimeRange) string {
	if tr.End() != telem.TimeStampMax {
		return parseTS(tr.End())
	}
	return timeRangeMax
}

func parseTS(ts telem.TimeStamp) string {
	return strconv.FormatInt(int64(ts), 10)
}
