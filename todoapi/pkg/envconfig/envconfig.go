package envcfg

import (
	"encoding"
	"errors"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	structTag         = "envcfg"
	structTagOption   = "envcfgOptions"
	structTagDefault  = "envcfgDefault"
	tagOptionKeep     = "keep"
	tagOptionOptional = "optional"
)

var textUnmarshalerType = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()

func isTextUnmarshaler(t reflect.Type) bool {
	return t.Implements(textUnmarshalerType) || reflect.PtrTo(t).Implements(textUnmarshalerType)
}

// Unmarshal will read your environment variables and try to unmarshal them
// to the passed struct. It will return an error, if it recieves an unsupported
// non-struct type, if types of the fields are not supported or if it can't
// parse value from an environment variable, thus taking care of validation of
// environment variables values. All field are required to be defined in the
// environment by default and will result in an error if they are not. Meta
// can be marked as optional by seting the envcfgOptions:"optional" field tag
// which will cause envcfg to not return an error if the field is missing from
// the environment.
// Unmarshal supports int, string, bool and []int, []string, []bool, and
// additionally any types that support the encoding.TextUnmarshaler interface.
func Unmarshal(v interface{}) error {
	structType, err := makeSureTypeIsSupported(v)
	if err != nil {
		return err
	}
	if err := makeSureStructFieldTypesAreSupported(structType); err != nil {
		return err
	}
	makeSureValueIsInitialized(v)

	env, err := newEnviron()
	if err != nil {
		return err
	}

	structVal := getStructValue(v)

	if err := unmarshalAllStructFields(structVal, env); err != nil {
		return err
	}

	return nil
}

// ClearEnvVars will clear all environment variables based on the struct
// field names or struct field tags. It will keep all those with
// envcfgOptions:"keep" struct field tag. It will return an error,
// if it recieves an unsupported non-struct type, if types of the
// fields are not supported
func ClearEnvVars(v interface{}) error {
	structType, err := makeSureTypeIsSupported(v)
	if err != nil {
		return err
	}
	if err := makeSureStructFieldTypesAreSupported(structType); err != nil {
		return err
	}

	unsetEnvVars(structType)
	return nil
}

func unsetEnvVarFromSingleField(structField reflect.StructField) {
	if tag := structField.Tag.Get(structTagOption); strings.Contains(tag, tagOptionKeep) {
		return
	}
	envKey := getEnvKey(structField)
	os.Setenv(envKey, "") // we're using Setenv instead of Unsetenv to ensure go1.3 compatibility
}

func unsetEnvVars(structType reflect.Type) {
	for i := 0; i < structType.NumField(); i++ {
		unsetEnvVarFromSingleField(structType.Field(i))
	}
}

func getEnvKey(structField reflect.StructField) string {
	if tag := structField.Tag.Get(structTag); tag != "" {
		return tag
	}
	return structField.Name
}

func getValue(structField reflect.StructField, env environ) (val string, err error) {
	key := getEnvKey(structField)
	val, ok := env[key]

	if !ok {
		defVal, ok := structField.Tag.Lookup(structTagDefault)
		if ok {
			return defVal, nil
		}

		if tag := structField.Tag.Get(structTagOption); strings.Contains(tag, tagOptionOptional) {
			return val, nil
		}

		return val, varUndefinedError(key)
	}

	return val, nil
}

func unmarshalInt(fieldVal reflect.Value, structField reflect.StructField, env environ) error {
	val, err := getValue(structField, env)
	if err != nil {
		return err
	}

	if val == "" {
		return nil
	}

	i, err := strconv.Atoi(val)
	if err != nil {
		return err
	}

	fieldVal.SetInt(int64(i))
	return nil
}

func unmarshalFloat(fieldVal reflect.Value, structField reflect.StructField, env environ, bitsize int) error {
	val, err := getValue(structField, env)
	if err != nil {
		return err
	}

	if val == "" {
		return nil
	}

	i, err := strconv.ParseFloat(val, bitsize)
	if err != nil {
		return err
	}

	fieldVal.SetFloat(i)
	return nil
}

var boolErr = errors.New("pass string 'true' or 'false' for boolean fields")

func varUndefinedError(f string) error {
	return errors.New("variable not found in environment: " + f)
}

func unmarshalBool(fieldVal reflect.Value, structField reflect.StructField, env environ) error {
	val, err := getValue(structField, env)
	if err != nil {
		return err
	}

	var vbool bool
	switch val {
	case "true":
		vbool = true
	case "false":
		vbool = false
	case "":
		return nil
	default:
		return boolErr
	}

	fieldVal.SetBool(vbool)
	return nil
}

func unmarshalDuration(fieldVal reflect.Value, structField reflect.StructField, env environ) error {
	val, err := getValue(structField, env)
	if err != nil {
		return err
	}

	if val == "" {
		return nil
	}

	d, err := time.ParseDuration(val)
	if err != nil {
		return err
	}

	fieldVal.SetInt(int64(d))
	return nil
}

func unmarshalString(fieldVal reflect.Value, structField reflect.StructField, env environ) error {
	val, err := getValue(structField, env)
	if err != nil {
		return err
	}

	fieldVal.SetString(val)
	return nil
}

func unmarshalTextUnmarshaler(fieldVal reflect.Value, structField reflect.StructField, env environ) error {
	val, err := getValue(structField, env)
	if err != nil {
		return err
	}

	if val == "" {
		return nil
	}

	textUnmarshaler := fieldVal.Addr().Interface().(encoding.TextUnmarshaler)
	textUnmarshaler.UnmarshalText([]byte(val))
	return nil
}

func appendToStringSlice(fieldVal reflect.Value, sliceVal string) error {
	fieldVal.Set(reflect.Append(fieldVal, reflect.ValueOf(sliceVal)))
	return nil
}

func appendToTextUnmarshalerSlice(fieldVal reflect.Value, sliceVal string) error {
	sliceElem := reflect.New(fieldVal.Type().Elem())
	textUnmarshaler := sliceElem.Interface().(encoding.TextUnmarshaler)
	textUnmarshaler.UnmarshalText([]byte(sliceVal))
	fieldVal.Set(reflect.Append(fieldVal, sliceElem.Elem()))
	return nil
}

func appendToIntSlice(fieldVal reflect.Value, sliceVal string) error {
	val, err := strconv.Atoi(sliceVal)
	if err != nil {
		return err
	}
	fieldVal.Set(reflect.Append(fieldVal, reflect.ValueOf(val)))
	return nil
}

func appendToFloatSlice(fieldVal reflect.Value, sliceVal string, bitSize int) error {
	val, err := strconv.ParseFloat(sliceVal, bitSize)
	if err != nil {
		return err
	}
	fieldVal.Set(reflect.Append(fieldVal, reflect.ValueOf(val)))
	return nil
}

func appendToBoolSlice(fieldVal reflect.Value, sliceVal string) error {
	var val bool
	switch sliceVal {
	case "true":
		val = true
	case "false":
		val = false
	default:
		return boolErr
	}
	fieldVal.Set(reflect.Append(fieldVal, reflect.ValueOf(val)))
	return nil
}

func appendToDurationSlice(fieldVal reflect.Value, sliceVal string) error {
	val, err := time.ParseDuration(sliceVal)
	if err != nil {
		return err
	}
	fieldVal.Set(reflect.Append(fieldVal, reflect.ValueOf(val)))
	return nil
}

func unmarshalSlice(fieldVal reflect.Value, structField reflect.StructField, env environ) error {
	envKey := getEnvKey(structField)
	envNames := make([]string, 0)

	for envName := range env {
		if strings.HasPrefix(envName, envKey) {
			envNames = append(envNames, envName)
		}
	}
	sort.Strings(envNames)

	if tag := structField.Tag.Get(structTagOption); len(envNames) == 0 && !strings.Contains(tag, tagOptionOptional) {
		return varUndefinedError(envKey)
	}

	var err error
	for _, envName := range envNames {
		val, ok := env[envName]
		if !ok {
			continue
		}
		if isTextUnmarshaler(structField.Type.Elem()) {
			err = appendToTextUnmarshalerSlice(fieldVal, val)
			if err != nil {
				return err
			}
			continue
		}
		switch structField.Type.Elem().Kind() {
		case reflect.String:
			err = appendToStringSlice(fieldVal, val)
		case reflect.Int, reflect.Int64:
			if structField.Type.Elem().PkgPath() == "time" && structField.Type.Elem().Name() == "Duration" {
				err = appendToDurationSlice(fieldVal, val)
			} else {
				err = appendToIntSlice(fieldVal, val)
			}
		case reflect.Bool:
			err = appendToBoolSlice(fieldVal, val)
		case reflect.Float64:
			err = appendToFloatSlice(fieldVal, val, 64)
		case reflect.Float32:
			err = appendToFloatSlice(fieldVal, val, 32)
		}
		if err != nil {
			return err
		}

	}
	return nil
}

func unmarshalSingleField(fieldVal reflect.Value, structField reflect.StructField, env environ) error {
	if !fieldVal.CanSet() { // unexported field can not be set
		return nil
	}
	// special case for structs that implement TextUnmarshaler interface
	if isTextUnmarshaler(structField.Type) {
		return unmarshalTextUnmarshaler(fieldVal, structField, env)
	}
	switch structField.Type.Kind() {
	case reflect.Int, reflect.Int64:
		if structField.Type.PkgPath() == "time" && structField.Type.Name() == "Duration" {
			return unmarshalDuration(fieldVal, structField, env)
		}
		return unmarshalInt(fieldVal, structField, env)
	case reflect.Float64:
		return unmarshalFloat(fieldVal, structField, env, 64)
	case reflect.Float32:
		return unmarshalFloat(fieldVal, structField, env, 32)
	case reflect.String:
		return unmarshalString(fieldVal, structField, env)
	case reflect.Bool:
		return unmarshalBool(fieldVal, structField, env)
	case reflect.Slice:
		return unmarshalSlice(fieldVal, structField, env)
	}
	return nil
}

func unmarshalAllStructFields(structVal reflect.Value, env environ) error {
	for i := 0; i < structVal.NumField(); i++ {
		if err := unmarshalSingleField(structVal.Field(i), structVal.Type().Field(i), env); err != nil {
			return err
		}
	}
	return nil
}

func getStructValue(v interface{}) reflect.Value {
	str := reflect.ValueOf(v)
	for {
		if str.Kind() == reflect.Struct {
			break
		}
		str = str.Elem()
	}
	return str
}

func makeSureValueIsInitialized(v interface{}) {
	if reflect.TypeOf(v).Elem().Kind() != reflect.Ptr {
		return
	}
	if reflect.ValueOf(v).Elem().IsNil() {
		reflect.ValueOf(v).Elem().Set(reflect.New(reflect.TypeOf(v).Elem().Elem()))
	}
}

func makeSureTypeIsSupported(v interface{}) (reflect.Type, error) {
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		return nil, errors.New("we need a pointer")
	}
	if reflect.TypeOf(v).Elem().Kind() == reflect.Ptr && reflect.TypeOf(v).Elem().Elem().Kind() == reflect.Struct {
		return reflect.TypeOf(v).Elem().Elem(), nil
	} else if reflect.TypeOf(v).Elem().Kind() == reflect.Struct && reflect.ValueOf(v).Elem().CanAddr() {
		return reflect.TypeOf(v).Elem(), nil
	}
	return nil, errors.New("we need a pointer to struct or pointer to pointer to struct")
}

func isSupportedStructField(k reflect.StructField) bool {
	// special case for types that implement TextUnmarshaler interface
	if isTextUnmarshaler(k.Type) {
		return true
	}
	switch k.Type.Kind() {
	case reflect.String:
		return true
	case reflect.Bool:
		return true
	case reflect.Int, reflect.Int64:
		return true
	case reflect.Float64:
		return true
	case reflect.Float32:
		return true
	case reflect.Slice:
		// special case for types that implement TextUnmarshaler interface
		if isTextUnmarshaler(k.Type.Elem()) {
			return true
		}
		switch k.Type.Elem().Kind() {
		case reflect.String:
			return true
		case reflect.Bool:
			return true
		case reflect.Int, reflect.Int64:
			return true
		case reflect.Float64:
			return true
		case reflect.Float32:
			return true
		default:
			return false
		}
	default:
		return false
	}
}

func makeSureStructFieldTypesAreSupported(structType reflect.Type) error {
	for i := 0; i < structType.NumField(); i++ {
		if !isSupportedStructField(structType.Field(i)) {
			return fmt.Errorf("unsupported struct field type: %v", structType.Field(i).Type)
		}
	}
	return nil
}

type environ map[string]string

func getAllEnvironNames(envList []string) (map[string]struct{}, error) {
	envNames := make(map[string]struct{})

	for _, kv := range envList {
		split := strings.SplitN(kv, "=", 2)
		if len(split) != 2 {
			return nil, fmt.Errorf("unknown environ condition - env variable not in k=v format: %v", kv)
		}
		envNames[split[0]] = struct{}{}
	}

	return envNames, nil
}

func newEnviron() (environ, error) {

	envNames, err := getAllEnvironNames(os.Environ())
	if err != nil {
		return nil, err
	}

	env := make(environ)

	for name := range envNames {
		env[name] = os.ExpandEnv(os.Getenv(name))
	}

	return env, nil
}
