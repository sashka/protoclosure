protoclosure
============

Protocol Buffer interoperability between goprotobuf and closure-library's
goog.proto2 (JavaScript).

protoc --go_out=. gopkg.in/samegoal/protoclosure.v0/test.proto
protoc --go_out=. gopkg.in/samegoal/protoclosure.v0/package_test.proto

mv gopkg.in/samegoal/protoclosure.v0/test.pb.go gopkg.in/samegoal/protoclosure.v0/test.pb/
mv gopkg.in/samegoal/protoclosure.v0/package_test.pb.go gopkg.in/samegoal/protoclosure.v0/package_test.pb/