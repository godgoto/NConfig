package NConfig

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type BaseType struct {
	Name    string
	Section string
	Key     string
	Typ     string
	Value   interface{}
}

func analysisStruct(cfg interface{}, tag string) ([]*BaseType, error) {
	var listBaseType []*BaseType
	typ := reflect.TypeOf(cfg)
	val := reflect.ValueOf(cfg) //获取reflect.Type类型

	kd := val.Kind() //获取到a对应的类别
	if kd != reflect.Struct {
		fmt.Println("expect struct")
		return nil, errors.New("expect struct")
	}
	for i := 0; i < val.NumField(); i++ {
		var btype BaseType
		btype.Name = typ.Field(i).Name
		btype.Key = typ.Field(i).Tag.Get(tag)
		btype.Typ = typ.Field(i).Type.String()
		btype.Value = val.Field(i).String()
		listBaseType = append(listBaseType, &btype)
	}
	return listBaseType, nil
}

func IsNumTF(s string) bool {
	bi, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return false
	}
	if bi > 0 {
		return true
	}
	return false
}

func typeof(v interface{}) string {
	return reflect.TypeOf(v).String()
}

func saveLocalValue(cfg interface{}, baseValue *BaseType) error {
	v := reflect.ValueOf(cfg).Elem() // the struct variable
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i) // a reflect.StructField
		if fieldInfo.Name == baseValue.Name {
			if typeof(baseValue.Value) == "string" && len(baseValue.Value.(string)) == 0 {
				continue
			}

			if baseValue.Typ == "string" {
				v.FieldByName(fieldInfo.Name).SetString(baseValue.Value.(string))
			} else if baseValue.Typ == "bool" {
				if strings.ToLower(baseValue.Value.(string)) == "true" {
					v.FieldByName(fieldInfo.Name).SetBool(true)
				} else if strings.ToLower(baseValue.Value.(string)) == "false" {
					v.FieldByName(fieldInfo.Name).SetBool(false)
				} else {
					v.FieldByName(fieldInfo.Name).SetBool(IsNumTF(baseValue.Value.(string)))
				}
			} else if baseValue.Typ == "int64" || baseValue.Typ == "int32" || baseValue.Typ == "int16" || baseValue.Typ == "int8" || baseValue.Typ == "int" {
				v.FieldByName(fieldInfo.Name).SetInt(int64(baseValue.Value.(float64)))
			} else if baseValue.Typ == "float64" || baseValue.Typ == "float32" {
				v.FieldByName(fieldInfo.Name).SetFloat(baseValue.Value.(float64))
			} else if baseValue.Typ == "[]byte" {
				v.FieldByName(fieldInfo.Name).SetBytes(baseValue.Value.([]byte))
			}
			break
		}
	}
	return nil
}

func saveApolloValue(cfg interface{}, baseValue *BaseType) error {
	v := reflect.ValueOf(cfg).Elem() // the struct variable
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i) // a reflect.StructField
		if fieldInfo.Name == baseValue.Name {
			if len(baseValue.Value.(string)) == 0 {
				continue
			}
			if baseValue.Typ == "string" {
				v.FieldByName(fieldInfo.Name).SetString(baseValue.Value.(string))
			} else if baseValue.Typ == "bool" {
				if strings.ToLower(baseValue.Value.(string)) == "true" {
					v.FieldByName(fieldInfo.Name).SetBool(true)
				} else if strings.ToLower(baseValue.Value.(string)) == "false" {
					v.FieldByName(fieldInfo.Name).SetBool(false)
				} else {
					v.FieldByName(fieldInfo.Name).SetBool(IsNumTF(baseValue.Value.(string)))
				}
			} else if baseValue.Typ == "int64" || baseValue.Typ == "int32" || baseValue.Typ == "int16" || baseValue.Typ == "int8" || baseValue.Typ == "int" {
				i64, _ := strconv.ParseInt(baseValue.Value.(string), 10, 64)
				v.FieldByName(fieldInfo.Name).SetInt(i64)
			} else if baseValue.Typ == "float64" || baseValue.Typ == "float32" {
				float, _ := strconv.ParseFloat(baseValue.Value.(string), 64)
				v.FieldByName(fieldInfo.Name).SetFloat(float)
			} else if baseValue.Typ == "[]byte" {
				v.FieldByName(fieldInfo.Name).SetBytes(baseValue.Value.([]byte))
			}
			break
		}
	}
	return nil
}

func section(list []*BaseType) {
	for k, v := range list {
		if len(v.Key) == 0 {
			continue
		}
		index := strings.Index(v.Key, ".")
		if index > -1 {
			list[k].Section = v.Key[0:index]
			list[k].Key = v.Key[index+1:]
		}
	}
}
