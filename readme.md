## redisgraph

This package alows you to use redis as a graph database.

### Installation

To install the package, use the following:

	go get -u github.com/johnfogh/redisgraph/...


### RelateKeys

A redis set will be created (or updated) for each object provided `key.related`
(where key represents the object's original key), and all of the other objects keys will be
stored as members.

  example:

    ...
    Obj1.key = "foo"
    Obj2.key = "bar"
    ...
    RelateKeys( Obj1 , Obj2 )

will create two sets in redis:

    foo.related { bar }
    bar.related { foo }

### RelateProperties

Examines all of the provided objects, and if any two objects have the same
attribute value then the keys will be related.

### RelateToKey

Relates all objects to the specified key.
