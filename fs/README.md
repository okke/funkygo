# Funky Go - FS

A package for dealing with streams of data.

## Importing Funky Go - FS

```go
import "github.com/okke/funkygo/fs"
```

## Basic concept

A stream is defined as a function which pulls the the first value and returns a stream to subsequent values.

```go
type Stream[T any] func() (T, Stream[T])
```

A stream always has an underlying source of data which is used to create a stream.

```go
stream := fs.FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
```

To take (or pull) the first value out of the stream, simply call the stream.

```go
value, stream := stream()
```

Note since a stream returns a stream of subsequent values, we need to use that stream to pull more data out of the stream.
