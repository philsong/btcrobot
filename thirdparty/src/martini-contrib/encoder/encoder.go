package encoder

// Original code borrowed from https://github.com/PuerkitoBio/martini-api-example
// TextEncoder and XmlEncoder has been removed. If someone really needs it, let me know.

// Supported tags:
// 	 - "out" if it sets to "false", value won't be set to field
import (
	"encoding/json"
	"reflect"
)

// An Encoder implements an encoding format of values to be sent as response to
// requests on the API endpoints.
type Encoder interface {
	Encode(v ...interface{}) ([]byte, error)
}

// Because `panic`s are caught by martini's Recovery handler, it can be used
// to return server-side errors (500). Some helpful text message should probably
// be sent, although not the technical error (which is printed in the log).
func Must(data []byte, err error) []byte {
	if err != nil {
		panic(err)
	}
	return data
}

type JsonEncoder struct{}

// jsonEncoder is an Encoder that produces JSON-formatted responses.
func (_ JsonEncoder) Encode(v ...interface{}) ([]byte, error) {
	var data interface{} = v
	var result interface{}

	if v == nil {
		// So that empty results produces `[]` and not `null`
		data = []interface{}{}
	} else if len(v) == 1 {
		data = v[0]
	}

	t := reflect.TypeOf(data)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() == reflect.Struct {
		result = copyStruct(reflect.ValueOf(data), t).Interface()
	} else {
		result = data
	}

	b, err := json.Marshal(result)

	return b, err
}

func copyStruct(v reflect.Value, t reflect.Type) reflect.Value {
	result := reflect.New(t).Elem()

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		if tag := t.Field(i).Tag.Get("out"); tag == "false" {
			continue
		}

		if v.Field(i).Kind() == reflect.Struct {
			result.Field(i).Set(copyStruct(v.Field(i), t.Field(i).Type))
			continue
		}

		result.Field(i).Set(v.Field(i))
	}

	return result
}
