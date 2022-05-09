module go-toy-comp

go 1.17

require (
	gee v0.0.0
	geecache v0.0.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)

replace (
	gee => ./gee
	geecache => ./geecache
)
