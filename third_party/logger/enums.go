package logger

//go:generate go run github.com/dmarkham/enumer -type=FormatEnum -linecomment -output=enums_gen.go
type FormatEnum int

const (
	Console FormatEnum = iota + 1 //console
	JSON                          //json
)
