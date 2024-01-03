# Funky Go - FS

**F**unky **S**treams, a package for dealing with streams of data.

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

Note since a stream returns a stream of subsequent values, we need to use that stream to pull more data out of the stream. Which for example can be done in a `for` loop.

```go
for value, stream := stream(); stream != nil; value, stream = stream() {
	log.Println(value)
}
```

Instead of using a for loop, it's also possible to use the `Each` function.

```go
fs.Each(stream, func(value int) error {
    log.Println(value)
    return nil
})
```

Which loops until an error is thrown by the provided callback function.

## Peeking a stream

Instead of pulling data out of the stream, it's also possible to inspect or to peek the first element(s).

```go
value, stream := fs.Peek(fs.FromArgs(1, 2, 3))
log.Println(value) // 1 (first element)
```

Note the `Peek` function, just like a stream itself, returns a value and a stream. The returned stream is guaranteed to contain the value which was 'peeked' (Although the value could have been pulled out of the underlying data source).

When interesting in more elements, use the `PeekN` function.

```go
first2, stream := fs.PeekN(fs.FromArgs(1, 2, 3), 2)
fs.Each(first2, func(x int) error {
    log.Println(x)
    return nil
})
```

`PeekN` returns two streams. The first is a stream containing the 'peeked' values. The second stream is a stream containing all values.

It's also possible to conditionally peek a stream using `PeekUntil`.

```go
ones, stream := fs.PeekUntil(fs.FromArgs(1, 1, 1, 2, 2), func(x int) bool {
    return x != 1
})
fs.Each(ones, func(x int) error {
    log.Println(x) // should only print 1's
    return nil
})
```

## Pulling more than one item from a stream

By using the stream itself as a function to pull items out of it, items are pulled out one by one. When there's the need to pull more, FunkyGo FS provides `TakeN` an `TakeUntil` functions. Just like `PeekN` and `PeekUntil` these functions return pulled out items as a separate stream.

```go
ignored, rest := fs.TakeUntil(fs.FromArgs('a', 'b', 'c', 'd', 'E'), unicode.IsUpper)
firstThreeIgnored, restIgnored := fs.TakeN(ignored,3)
```
