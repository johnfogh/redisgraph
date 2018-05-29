## redisgraph

This package alows you to use redis as a graph database.

### Installation

To install the package, use the following:

      go get -u github.com/johnfogh/redisgraph/...

### Examples

The tests provide more concrete examples of how this package works.

### General usage

Essentially, this package works by creating a redis set named `key.related`
where `key` represents the original key. The members within this set are all of
the related keys.

#### RelateKeys

Forces a relation between between the provided keys.

#### RelateProperties

Examines all of the provided objects, and if any two objects have the same
attribute value then the keys will be related. The `redisgraph.Data` interface
expects the objects to implement an `Attributes()` method that returns a
`map[string]interface{}`.

#### RelateObjects

Relates all provided objects to the specified key.

#### CommonRelatives

Examines all of the provided objects, and returns a `[]string` containing the
keys shared by all of the objects.
