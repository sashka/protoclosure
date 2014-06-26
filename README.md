protoclosure
============

[Protocol Buffer](https://code.google.com/p/protobuf/) interoperability
between [Go](http://golang.org)'s
[goprotobuf](https://code.google.com/p/goprotobuf/) and JavaScript using
[closure-library](https://developers.google.com/closure/library/)'s
[goog.proto2](https://github.com/google/closure-library/tree/master/closure/goog/proto2)).

JS Usage
--------

Go Usage
--------

PBLite format
-------------

Example message:

```protobuf
message Person {
  optional int32 id = 1;
  optional string name = 2;
  optional string email = 3;
}
```

Example encoding:

```json
[null,1,null,"user@example.com"]
```

Example encoding (zero-index):

```json
[1,null,"user@example.com"]
```

PBObject format
---------------

Example message:

```protobuf
message Person {
  optional int32 id = 1;
  optional string name = 2;
  optional string email = 3;
}
```

Example encoding (tag name):

```json
{"id":1,"email":"user@example.com"}
```

Example encoding (tag number):

```json
{"1":1,"3":"user@example.com"}
```

protoclosure development
-------------------------

```
$ # setup environment (GOPATH, etc)
$ go get gopkg.in/samegoal/protoclosure.v0
$ cd gopkg.in/samegoal/protoclosure.v0
$ # modify the source
$ go test -race
$ make  # run vet/fmt/lint, prior to sending Pull Request
```

To regenerate unit test protobuf files:

```
protoc --go_out=. gopkg.in/samegoal/protoclosure.v0/test.proto
protoc --go_out=. gopkg.in/samegoal/protoclosure.v0/package_test.proto

mv gopkg.in/samegoal/protoclosure.v0/test.pb.go gopkg.in/samegoal/protoclosure.v0/test.pb/
mv gopkg.in/samegoal/protoclosure.v0/package_test.pb.go gopkg.in/samegoal/protoclosure.v0/package_test.pb/
```

[goprotobuf](https://code.google.com/p/goprotobuf/) limitations:

  * [Import dependencies](https://code.google.com/p/goprotobuf/issues/detail?id=32)
  * [Custom options](https://code.google.com/p/goprotobuf/issues/detail?id=34)
