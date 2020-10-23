# stout

Tool of convert struct. Now, it support json only.

## Install

```shell
go get github.com/komem3/stout
```

## Usage

```shell
$ stout -h
Usage of stout:
  -no-format
    	Not format.
  -path string
    	Definetion file. (required)
  -type string
    	Target struct type. (required)

$ stout -path ./define_test.go -type SampleJson
```

## Caution

Not yet supported

- `json.Marshaler`

## Author

komem3

## License

MIT
