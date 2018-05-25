package redisgraph

import (
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Redis - client for redis
var Redis *redis.Client

func init() {
	opts := redis.Options{Addr: "localhost:32768",
		Password: "",
		DB:       0}
	Redis = redis.NewClient(&opts)
}

//------------------------------------------------------------------------------

// Example struct implementing the redisgraph.Data interface
type SampleData struct {
	key        string
	attributes map[string]interface{}
}

func NewSampleData(key string, attribs ...string) (sd *SampleData) {
	sd = new(SampleData)
	sd.key = key
	sd.attributes = make(map[string]interface{})
	if len(attribs) > 1 {
		for index := 0; index < len(attribs)-1; index++ {
			sd.attributes[attribs[index]] = attribs[index+1]
		}
	}
	return
}

func (sd *SampleData) Key() string {
	return sd.key
}

func (sd *SampleData) Attributes() map[string]interface{} {
	return sd.attributes
}

//------------------------------------------------------------------------------

func keyExists(t *testing.T, expected int64, keys ...string) {
	actual, err := Redis.Exists(keys...).Result()
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func hasMember(t *testing.T, key string, member string, expected bool) {
	actual, err := Redis.SIsMember(key, member).Result()
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

//------------------------------------------------------------------------------

func TestRedisConnection(t *testing.T) {
	str, err := Redis.Ping().Result()
	assert.Nil(t, err)
	assert.Equal(t, "PONG", str)
}

func TestRelateKeys(t *testing.T) {
	foo := NewSampleData("foo")
	bar := NewSampleData("bar")
	keys := GetKeys(foo, bar)
	RelateKeys(Redis, keys...)
	keyExists(t, int64(2), "foo.related", "bar.related")
	hasMember(t, "foo.related", "bar", true)
	hasMember(t, "bar.related", "foo", true)
	Redis.Del("foo.related", "bar.related")
}

func TestRelateSingleProperties(t *testing.T) {
	foo := NewSampleData("foo", "prop1", "true")
	bar := NewSampleData("bar", "prop1", "true")
	baz := NewSampleData("bar", "prop1", "false")
	matchAttributes := []string{"prop1"}
	RelateProperties(Redis, matchAttributes, foo, bar, baz)
	keyExists(t, int64(2), "foo.related", "bar.related", "baz.related")
	keyExists(t, int64(0), "baz.related")
	hasMember(t, "foo.related", "bar", true)
	hasMember(t, "bar.related", "foo", true)
	hasMember(t, "foo.related", "baz", false)
	hasMember(t, "bar.related", "baz", false)
	Redis.Del("foo.related", "bar.related")
}

func TestRelateMultipleProperties(t *testing.T) {
	foo := NewSampleData("foo", "prop1", "true", "prop2", "true")
	bar := NewSampleData("bar", "prop1", "true", "prop2", "false")
	baz := NewSampleData("baz", "prop1", "false", "prop2", "true")
	matchAttributes := []string{"prop1", "prop2"}
	RelateProperties(Redis, matchAttributes, foo, bar, baz)
	keyExists(t, int64(3), "foo.related", "bar.related", "baz.related")

	hasMember(t, "foo.related", "bar", true)
	hasMember(t, "foo.related", "baz", true)

	hasMember(t, "bar.related", "foo", true)
	hasMember(t, "bar.related", "baz", false)

	hasMember(t, "baz.related", "foo", true)
	hasMember(t, "baz.related", "bar", false)

	Redis.Del("foo.related", "bar.related", "baz.related")
}

func TestRelateToKey(t *testing.T) {
	foo := NewSampleData("foo", "prop1", "true", "prop2", "true")
	bar := NewSampleData("bar", "prop1", "true", "prop2", "false")
	RelateToKey(Redis, "baz", foo, bar)
	keyExists(t, int64(3), "foo.related", "bar.related", "baz.related")
	hasMember(t, "foo.related", "bar", true)
	hasMember(t, "foo.related", "baz", true)
	hasMember(t, "bar.related", "foo", true)
	hasMember(t, "bar.related", "baz", true)
	hasMember(t, "baz.related", "foo", true)
	hasMember(t, "baz.related", "bar", true)
	Redis.Del("test.related", "foo.related", "bar.related", "baz.related")
}
