# stout

Tool of convert struct. Now, it support json only.

## Install

```shell
go get github.com/komem3/stout
```

## Usage

```shell
$ stout -h
usage: stout -path $struct_path [-no-format] $struct_name
options:
  -path string
        File path of defined struct. (required)
  -no-format bool
        Not format the output json.
$ stout -path ./define_test.go SampleJson
```

## Caution

Not yet supported

- `json.Marshaler`
- [support format](./define_test.go)

## Author

komem3

## License

MIT
