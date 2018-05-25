package redisgraph

import (
	"github.com/go-redis/redis"
)

// Data -
type Data interface {
	Key() string
	Attributes() map[string]interface{}
}

// RelateKeys - explicitly relate all of the keys passed.
// given the keys "foo", "bar", this will produce redis sets:
// foo.related = { bar } and bar.related = { foo }
func RelateKeys(Redis *redis.Client, keys ...string) {
	for _, k := range keys {
		temp := make([]string, len(keys))
		for index, kk := range keys {
			if k != kk {
				temp[index] = kk
			}
		}
		Redis.SAdd(k+".related", temp)
	}
}

// RelateProperties - relates a set of objects if any of the attributes
// specified by matchAttributes are equal.
func RelateProperties(Redis *redis.Client, matchAttributes []string, objects ...Data) {
	// iterate through all of the objects
	for index := 0; index < len(objects); index++ {
		indexObj := objects[index]
		// iterate through all of the objects between index+1 and the end.
		for offset := index + 1; offset < len(objects); offset++ {
			offsetObj := objects[offset]
			// iterate through all of the properties we care about.
			for _, prop := range matchAttributes {
				if compareProperty(prop, indexObj, offsetObj) == true {
					RelateKeys(Redis, indexObj.Key(), offsetObj.Key())
				}
			}
		}
	}
}

// RelateToKey - relates all objects to an explicit key.
func RelateToKey(Redis *redis.Client, key string, objects ...Data) {
	keys := GetKeys(objects...)
	keys = append(keys, key)
	RelateKeys(Redis, keys...)
}

// GetKeys - returns all of the object keys.
func GetKeys(objects ...Data) []string {
	keys := make([]string, 0, len(objects))
	for _, o := range objects {
		keys = append(keys, o.Key())
	}
	return keys
}

// compareProperty - returns true if all of the object.Attribute[property] are equal
func compareProperty(property string, objects ...Data) bool {
	for index := 0; index < len(objects)-1; index++ {
		baseAttributes := objects[index].Attributes()
		offsetAttributes := objects[index+1].Attributes()
		if baseAttributes[property] != offsetAttributes[property] {
			return false
		}
	}
	return true
}
