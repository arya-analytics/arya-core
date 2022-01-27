package timeseries

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"reflect"
	"strconv"
)

type Command string

const (
	CMDCreate Command = "TS.CREATE"
	CMDExists Command = "EXISTS"
	CMDMAdd   Command = "TS.MADD"
	CMDGet    Command = "TS.GET"
	CMDRange  Command = "TS.RANGE"
	CMDInfo   Command = "TS.INFO"
)

type Sample struct {
	Key       string
	Value     float64
	Timestamp int64
}

func newSampleFromRes(key string, res interface{}) (Sample, error) {
	resVal := reflect.ValueOf(res)
	ts := resVal.Index(0).Interface().(int64)
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
		s.Timestamp,
		s.Value,
	}
}

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

func (c *Client) TSCreate(ctx context.Context, key string,
	opts CreateOptions) *redis.Cmd {
	return c.exec(ctx, CMDCreate, opts.args(key)...)
}

func (c *Client) TSAddSample(ctx context.Context, samples ...Sample) *redis.Cmd {
	var sampleArgs []interface{}
	for _, sample := range samples {
		sampleArgs = append(sampleArgs, sample.args()...)
	}
	return c.exec(ctx, CMDMAdd, sampleArgs...)
}

func (c *Client) TSGet(ctx context.Context, key string) (Sample, error) {
	res, err := c.exec(ctx, CMDGet, key).Result()
	if err != nil {
		return Sample{}, err
	}
	return newSampleFromRes(key, res)
}

func (c *Client) TSRange(ctx context.Context, key string, fromTS int,
	toTS int) ([]Sample, error) {
	fromTSArg := "-"
	if fromTS != 0 {
		fromTSArg = strconv.Itoa(fromTS)
	}
	toTSArg := "+"
	if toTS != 0 {
		toTSArg = strconv.Itoa(toTS)
	}
	res, err := c.exec(ctx, CMDRange, key, fromTSArg, toTSArg).Result()
	if err != nil {
		return nil, err
	}
	var samples []Sample
	resV := reflect.ValueOf(res)
	for i := 0; i < resV.Len(); i++ {
		s, sErr := newSampleFromRes(key, resV.Index(i).Interface())
		if sErr != nil {
			return nil, sErr
		}
		samples = append(samples, s)
	}
	return samples, err
}
