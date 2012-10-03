
package parse_args


import (
	"errors"
	"reflect"
	"strconv"
	. "appengine"
	_"appengine/datastore"
	_ "appengine/user"
	_"encoding/json"
	_ "fmt"
	_"gomfk/json_api"
	"net/http"
	_ "net/url"
	"strings"
)

type parsedArg struct {
	name string
	path []int
	required bool
	kind reflect.Kind
}

var ARG_CACHE = make(map[reflect.Type][]parsedArg)

func getArgSpec(tipe reflect.Type) []parsedArg {
	if cached := ARG_CACHE[tipe]; len(cached) > 0 {
		return cached
	}

	var args []parsedArg
	traverseType(tipe, "", []int{}, &args)
	ARG_CACHE[tipe] = args
	return args
}

func traverseType(tipe reflect.Type, prefix string, pathPrefix []int,
		args *[]parsedArg) {

	for i := 0; i < tipe.NumField(); i++ {
		f := tipe.Field(i)
		tag := f.Tag.Get("parseArg")

		if tag == "" {
			continue
		}

		tagParts := strings.Split(tag, ",")

 		name := prefix + tagParts[0]
		path := make([]int, len(pathPrefix) + 1)
		copy(path, pathPrefix)
		path[len(path) - 1] = i

		if f.Type.Kind() == reflect.Struct {
			traverseType(f.Type, name, path, args)
		} else {
			var arg parsedArg
			arg.name = name
			arg.path = path
			arg.required = false
			for _, opt := range tagParts[1:] {
				if opt == "required" {
					arg.required = true
				}
			}
			arg.kind = f.Type.Kind()
			*args = append(*args, arg)
		}
	}
}

func ParseArgs(r *http.Request, into interface{}) error {
	intoType := reflect.TypeOf(into).Elem()
	intoValue := reflect.ValueOf(into)

	argSpec := getArgSpec(intoType)
	NewContext(r).Errorf("%v", argSpec)

	for _, arg := range argSpec {
		formVal := r.FormValue(arg.name)

		if formVal == "" {
			if arg.required {
				return errors.New("Required field missing: " + arg.name)
			}
			continue
		}
		fieldValue := intoValue.Elem().FieldByIndex(arg.path)
		switch fieldValue.Kind() {
		case reflect.String:
			fieldValue.SetString(formVal)
		case reflect.Int64:
			intVal, err := strconv.ParseInt(formVal, 10, 64)
			if err != nil {
				return err
			}
			fieldValue.SetInt(intVal)
		}
	}

	validator := intoValue.MethodByName("Validate")
	if validator.IsValid() {
		validatorResults := validator.Call([]reflect.Value{})
		if len(validatorResults) != 1 {
			panic("validator must return error only")
		}
		err := validatorResults[0]
		if !err.IsNil() {
			return err.Interface().(error)
		}
	}
	return nil
}
