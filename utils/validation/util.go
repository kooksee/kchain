package validation

import (
	"encoding/json"
	"reflect"
	"regexp"
	"strings"
	"github.com/pkg/errors"
)

var funcMap = make(map[string]interface{})

func init() {
	// most stolen from github.com/asaskevich/govalidator
	funcMap["null"] = isNull
	funcMap["uuid"] = isUUID
	funcMap["alpha"] = isAlpha
	funcMap["alphanum"] = isAlphanumeric
	funcMap["base64"] = isBase64
	funcMap["float"] = isFloat
	funcMap["numeric"] = isNumeric
}

func isNull(str string) bool {
	return len(str) == 0
}

func isUUID(str string) bool {
	return regexp.MustCompile("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$").MatchString(str)
}

func isAlpha(str string) bool {
	if isNull(str) {
		return true
	}
	return regexp.MustCompile("^[a-zA-Z]+$").MatchString(str)
}

func isAlphanumeric(str string) bool {
	if isNull(str) {
		return true
	}
	return regexp.MustCompile("^[a-zA-Z0-9]+$").MatchString(str)
}

func isBase64(str string) bool {
	return regexp.MustCompile("^(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\\/]{3}=|[A-Za-z0-9+\\/]{4})$").MatchString(str)
}

func isFloat(str string) bool {
	return str != "" && regexp.MustCompile("^(?:[-+]?(?:[0-9]+))?(?:\\.[0-9]*)?(?:[eE][\\+\\-]?(?:[0-9]+))?$").MatchString(str)
}

func isNumeric(str string) bool {
	if isNull(str) {
		return true
	}
	return regexp.MustCompile("^[-+]?[0-9]+$").MatchString(str)
}

// TODO length and matches Ex:"length(min|max)": ByteLength, "matches(pattern)": StringMatches,

// abc_xyz to AbcXyz
func underscoreToCamelCase(s string) string {
	return strings.Replace(strings.Title(strings.Replace(strings.ToLower(s), "_", " ", -1)), " ", "", -1)
}

// usage:
// req := &pb.Version{}
// err:=misc.Request2StructWithValidate(body, req, "id:uuid", "build_version:alphanum")
// judge err
func Request2StructWithValidate(jsonStr string, req interface{}, constraints ...string) error {
	err := json.Unmarshal([]byte(jsonStr), req)
	if err != nil {
		return err
	}

	//check
	t := reflect.ValueOf(req)
	for _, v := range constraints {
		x := strings.SplitN(v, ":", 2)
		fieldName := underscoreToCamelCase(x[0])
		value := t.Elem().FieldByName(fieldName)
		if value.IsValid() {
			if len(x) == 2 {
				if !constraintValidate(value, funcMap[x[1]]) {
					return errors.New(`validate failed for "` + x[0] + `" with constraint "` + x[1] + `"`)
				}
			} else {
				if !normalValidate(value) {
					return errors.New(`validate failed for "` + x[0] + `" with missing value"`)
				}
			}
		}
		//TODO need warning if no such fileName?
	}

	return nil
}

func constraintValidate(v reflect.Value, constraintFunc interface{}) bool {
	fv := reflect.ValueOf(constraintFunc)
	if fv.Kind() == reflect.Func {
		return fv.Call([]reflect.Value{v})[0].Bool()
	}
	// TODO need warning if not validFunc?
	return true
}

func normalValidate(v reflect.Value) bool {
	k := v.Kind()
	switch k {
	case reflect.String:
		if v.String() == "" {
			return false
		}
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int, reflect.Int64:
		if v.Int() == 0 {
			return false
		}
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint, reflect.Uint64:
		if v.Uint() == 0 {
			return false
		}
	case reflect.Float32, reflect.Float64:
		if v.Float() == 0 {
			return false
		}
	}

	return true
}
