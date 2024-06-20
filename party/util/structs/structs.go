package structs

import (
	"errors"
	"gitlab.top.slotssprite.com/my/api-layout/party/util"
	"reflect"
)

// TagType tag type
type TagType uint8

const (
	JSON TagType = iota
	GORM
)

// ToString tag type to string
func (t TagType) ToString() string {
	switch t {
	case JSON:
		return "json"
	case GORM:
		return "gorm"
	}
	return ""
}

type Option interface {
	apply(s *Struct)
}

type funcOption struct {
	f func(s *Struct)
}

func (f *funcOption) apply(s *Struct) {
	f.f(s)
}

// WithIgnoreOmitempty ignore structure omitempty tag
func WithIgnoreOmitempty() Option {
	return &funcOption{
		f: func(s *Struct) {
			s.isIgnoreOmitempty = true
		},
	}
}

// WithKeyToSnakeCase map key to snake case
func WithKeyToSnakeCase() Option {
	return &funcOption{
		f: func(s *Struct) {
			s.isKeyToSnakeCase = true
		},
	}
}

// WithOnlyPrimaryDataType only primary data type
func WithOnlyPrimaryDataType() Option {
	return &funcOption{
		f: func(s *Struct) {
			s.onlyPrimaryDataType = true
		},
	}
}

// WithIgnoreKeys ignore keys
func WithIgnoreKeys(keys ...string) Option {
	return &funcOption{
		f: func(s *Struct) {
			s.ignoreKeys = keys
		},
	}
}

// Struct struct
type Struct struct {
	raw                 interface{}
	value               reflect.Value
	tagName             string
	isIgnoreOmitempty   bool     // 忽略 struct omitempty tag
	isKeyToSnakeCase    bool     // key 转下划线
	onlyPrimaryDataType bool     // 仅基本数据类型
	ignoreKeys          []string // 忽略的keys
}

var DefaultIgnoreKeys = []string{"created_timestamp", "created_at", "updated_at"}

// New new struct
func New(tagType TagType, s interface{}) (*Struct, error) {
	v := reflect.ValueOf(s)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, errors.New("not a structure")
	}
	tagName := tagType.ToString()
	if tagName == "" {
		return nil, errors.New("not support tag type")
	}
	return &Struct{
		raw:     s,
		value:   v,
		tagName: tagName,
	}, nil
}

// structFields get struct fields
func (s *Struct) structFields() []reflect.StructField {
	t := s.value.Type()
	sf := make([]reflect.StructField, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" {
			continue
		}
		if tag := field.Tag.Get(s.tagName); tag == "-" {
			continue
		}
		sf = append(sf, field)
	}
	return sf
}

// ProtoStructToGormMap proto struct to gorm map
func (s *Struct) ProtoStructToGormMap(options ...Option) map[string]interface{} {
	options = append(options, WithIgnoreOmitempty(), WithKeyToSnakeCase(), WithOnlyPrimaryDataType())
	return s.StructToMap(options...)
}

// GormStructToMap ...
func (s *Struct) GormStructToMap(options ...Option) map[string]interface{} {
	options = append(options, WithKeyToSnakeCase(), WithOnlyPrimaryDataType())
	return s.StructToMap(options...)
}

// StructToMap struct to map[string]interface{}
func (s *Struct) StructToMap(options ...Option) map[string]interface{} {
	out := make(map[string]interface{})
	for _, option := range options {
		option.apply(s)
	}
	for _, field := range s.structFields() {
		name := field.Name
		val := s.value.FieldByName(name)

		tagName, tagOpts := parseTag(field.Tag.Get(s.tagName))
		if s.tagName == GORM.ToString() {
			tagName = parseGormColumn(tagName)
		}
		if tagName != "" {
			if s.isKeyToSnakeCase {
				name = util.CamelToUnderline(tagName)
			} else {
				name = tagName
			}
		} else {
			if s.isKeyToSnakeCase {
				name = util.CamelToUnderline(name)
			}
		}

		if tagOpts.Has("omitempty") && !s.isIgnoreOmitempty {
			zero := reflect.Zero(val.Type()).Interface()
			current := val.Interface()
			if reflect.DeepEqual(current, zero) {
				continue
			}
		}

		if s.onlyPrimaryDataType {
			if val.Kind() == reflect.Ptr {
				val = val.Elem()
			}
			switch val.Kind() {
			case reflect.Bool,
				reflect.String,
				reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint,
				reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int,
				reflect.Float32, reflect.Float64,
				reflect.Complex64, reflect.Complex128:
				out[name] = val.Interface()
			}
		} else {
			out[name] = val.Interface()
		}
	}
	for _, key := range s.ignoreKeys {
		if _, ok := out[key]; ok {
			delete(out, key)
		}
	}
	return out
}
