package anyiter

import (
	"errors"
	"reflect"
)

// SafeType is a method-for-method recreation of the reflect.Type interface, but with the methods that can panic
// modified to return errors instead in the cases where they would panic.
type SafeType interface {
	// ReflectType returns the underlying reflect.Type
	ReflectType() reflect.Type

	// Align returns the alignment in bytes of a value of
	// this type when allocated in memory.
	Align() int

	// FieldAlign returns the alignment in bytes of a value of
	// this type when used as a field in a struct.
	FieldAlign() int

	// Method returns the i'th method in the type's method set.
	// It errors if i is not in the range [0, NumMethod()).
	//
	// For a non-interface type T or *T, the returned Method's Type and Func
	// fields describe a function whose first argument is the receiver,
	// and only exported methods are accessible.
	//
	// For an interface type, the returned Method's Type field gives the
	// method signature, without a receiver, and the Func field is nil.
	//
	// Methods are sorted in lexicographic order.
	Method(int) (SafeMethod, error)

	// MethodByName returns the method with that name in the type's
	// method set and a boolean indicating if the method was found.
	//
	// For a non-interface type T or *T, the returned Method's Type and Func
	// fields describe a function whose first argument is the receiver.
	//
	// For an interface type, the returned Method's Type field gives the
	// method signature, without a receiver, and the Func field is nil.
	MethodByName(string) (reflect.Method, bool)

	// NumMethod returns the number of methods accessible using Method.
	//
	// Note that NumMethod counts unexported methods only for interface types.
	NumMethod() int

	// Name returns the type's name within its package for a defined type.
	// For other (non-defined) types it returns the empty string.
	Name() string

	// PkgPath returns a defined type's package path, that is, the import path
	// that uniquely identifies the package, such as "encoding/base64".
	// If the type was predeclared (string, error) or not defined (*T, struct{},
	// []int, or A where A is an alias for a non-defined type), the package path
	// will be the empty string.
	PkgPath() string

	// Size returns the number of bytes needed to store
	// a value of the given type; it is analogous to unsafe.Sizeof.
	Size() uintptr

	// String returns a string representation of the type.
	// The string representation may use shortened package names
	// (e.g., base64 instead of "encoding/base64") and is not
	// guaranteed to be unique among types. To test for type identity,
	// compare the Types directly.
	String() string

	// Kind returns the specific kind of this type.
	Kind() reflect.Kind

	// Implements reports whether the type implements the interface type u.
	Implements(u SafeType) bool

	// AssignableTo reports whether a value of the type is assignable to type u.
	AssignableTo(u SafeType) bool

	// ConvertibleTo reports whether a value of the type is convertible to type u.
	ConvertibleTo(u SafeType) bool

	// Comparable reports whether values of this type are comparable.
	Comparable() bool

	// Bits returns the size of the type in bits.
	// It panics if the type's Kind is not one of the
	// sized or unsized Int, Uint, Float, or Complex kinds.
	Bits() int

	// ChanDir returns a channel type's direction.
	// It errors if the type's Kind is not Chan.
	ChanDir() (reflect.ChanDir, error)

	// IsVariadic reports whether a function type's final input parameter
	// is a "..." parameter. If so, t.In(t.NumIn() - 1) returns the parameter's
	// implicit actual type []T.
	//
	// For concreteness, if t represents func(x int, y ... float64), then
	//
	//	t.NumIn() == 2
	//	t.In(0) is the reflect.Type for "int"
	//	t.In(1) is the reflect.Type for "[]float64"
	//	t.IsVariadic() == true
	//
	// IsVariadic errors if the type's Kind is not Func.
	IsVariadic() (bool, error)

	// Elem returns a type's element type.
	// It errors if the type's Kind is not Array, Chan, Map, Ptr, or Slice.
	Elem() (SafeType, error)

	// Field returns a struct type's i'th field.
	// It errors if the type's Kind is not Struct.
	// It errors if i is not in the range [0, NumField()).
	Field(i int) (reflect.StructField, error)

	// FieldByIndex returns the nested field corresponding
	// to the index sequence. It is equivalent to calling Field
	// successively for each index i.
	// It errors if the type's Kind is not Struct.
	FieldByIndex(index []int) (reflect.StructField, error)

	// FieldByName returns the struct field with the given name
	// and a boolean indicating if the field was found.
	FieldByName(name string) (reflect.StructField, bool)

	// FieldByNameFunc returns the struct field with a name
	// that satisfies the match function and a boolean indicating if
	// the field was found.
	//
	// FieldByNameFunc considers the fields in the struct itself
	// and then the fields in any embedded structs, in breadth first order,
	// stopping at the shallowest nesting depth containing one or more
	// fields satisfying the match function. If multiple fields at that depth
	// satisfy the match function, they cancel each other
	// and FieldByNameFunc returns no match.
	// This behavior mirrors Go's handling of name lookup in
	// structs containing embedded fields.
	FieldByNameFunc(match func(string) bool) (reflect.StructField, bool)

	// In returns the type of a function type's i'th input parameter.
	// It errors if the type's Kind is not Func.
	// It errors if i is not in the range [0, NumIn()).
	In(i int) (SafeType, error)

	// Key returns a map type's key type.
	// It errors if the type's Kind is not Map.
	Key() (SafeType, error)

	// Len returns an array type's length.
	// It errors if the type's Kind is not Array.
	Len() (int, error)

	// NumField returns a struct type's field count.
	// It errors if the type's Kind is not Struct.
	NumField() (int, error)

	// NumIn returns a function type's input parameter count.
	// It errors if the type's Kind is not Func.
	NumIn() (int, error)

	// NumOut returns a function type's output parameter count.
	// It errors if the type's Kind is not Func.
	NumOut() (int, error)

	// Out returns the type of a function type's i'th output parameter.
	// It errors if the type's Kind is not Func.
	// It errors if i is not in the range [0, NumOut()).
	Out(i int) (SafeType, error)
}

type safeType struct {
	reflectType reflect.Type
}

// NewSafeType wraps an existing reflect.Type in the SafeType interface
func NewSafeType(reflectType reflect.Type) SafeType {
	return &safeType{reflectType: reflectType}
}

func (s safeType) ReflectType() reflect.Type {
	return s.reflectType
}

func (s safeType) Align() int {
	return s.reflectType.Align()
}

func (s safeType) FieldAlign() int {
	return s.reflectType.FieldAlign()
}

func (s safeType) Method(i int) (SafeMethod, error) {
	if i > s.reflectType.NumMethod() {
		return nil, errors.New("index out of range")
	}

	return NewSafeMethod(s.reflectType.Method(i)), nil
}

func (s safeType) MethodByName(s2 string) (reflect.Method, bool) {
	return s.reflectType.MethodByName(s2)
}

func (s safeType) NumMethod() int {
	return s.reflectType.NumMethod()
}

func (s safeType) Name() string {
	return s.reflectType.Name()
}

func (s safeType) PkgPath() string {
	return s.reflectType.PkgPath()
}

func (s safeType) Size() uintptr {
	return s.reflectType.Size()
}

func (s safeType) String() string {
	return s.reflectType.String()
}

func (s safeType) Kind() reflect.Kind {
	return s.reflectType.Kind()
}

func (s safeType) Implements(u SafeType) bool {
	return s.reflectType.Implements(u.ReflectType())
}

func (s safeType) AssignableTo(u SafeType) bool {
	return s.reflectType.AssignableTo(u.ReflectType())
}

func (s safeType) ConvertibleTo(u SafeType) bool {
	return s.reflectType.ConvertibleTo(u.ReflectType())
}

func (s safeType) Comparable() bool {
	return s.reflectType.Comparable()
}

func (s safeType) Bits() int {
	return s.reflectType.Bits()
}

func (s safeType) ChanDir() (reflect.ChanDir, error) {
	if s.reflectType.Kind() != reflect.Chan {
		return reflect.ChanDir(0), errors.New("type is not Chan")
	}
	return s.reflectType.ChanDir(), nil
}

func (s safeType) IsVariadic() (bool, error) {
	if s.reflectType.Kind() != reflect.Func {
		return false, errors.New("type is not Func")
	}
	return s.reflectType.IsVariadic(), nil
}

func (s safeType) Elem() (SafeType, error) {
	switch s.reflectType.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice:
		return NewSafeType(s.ReflectType().Elem()), nil
	}
	return nil, errors.New("type is not Array, Chan, Map, Ptr or Slice")
}

func (s safeType) Field(i int) (reflect.StructField, error) {
	if s.reflectType.Kind() != reflect.Struct {
		return reflect.StructField{}, errors.New("type is not Struct")
	}

	if i > s.reflectType.NumField() {
		return reflect.StructField{}, errors.New("index out of range")
	}

	return s.reflectType.Field(i), nil
}

func (s safeType) FieldByIndex(index []int) (reflect.StructField, error) {
	if s.reflectType.Kind() != reflect.Struct {
		return reflect.StructField{}, errors.New("type is not Struct")
	}
	return s.reflectType.FieldByIndex(index), nil
}

func (s safeType) FieldByName(name string) (reflect.StructField, bool) {
	if s.reflectType.Kind() != reflect.Struct {
		return reflect.StructField{}, false
	}
	return s.reflectType.FieldByName(name)
}

func (s safeType) FieldByNameFunc(match func(string) bool) (reflect.StructField, bool) {
	if s.reflectType.Kind() != reflect.Struct {
		return reflect.StructField{}, false
	}
	return s.reflectType.FieldByNameFunc(match)
}

func (s safeType) In(i int) (SafeType, error) {
	if s.reflectType.Kind() != reflect.Func {
		return nil, errors.New("type is not Func")
	}

	if i > s.reflectType.NumIn() {
		return nil, errors.New("index out of range")
	}

	return NewSafeType(s.reflectType.In(i)), nil
}

func (s safeType) Key() (SafeType, error) {
	if s.reflectType.Kind() != reflect.Map {
		return nil, errors.New("type is not Map")
	}
	return NewSafeType(s.reflectType.Key()), nil
}

func (s safeType) Len() (int, error) {
	if s.reflectType.Kind() != reflect.Array {
		return 0, errors.New("type is not Array")
	}
	return s.reflectType.Len(), nil
}

func (s safeType) NumField() (int, error) {
	if s.reflectType.Kind() != reflect.Struct {
		return 0, errors.New("type is not Struct")
	}

	return s.reflectType.NumField(), nil
}

func (s safeType) NumIn() (int, error) {
	if s.reflectType.Kind() != reflect.Func {
		return 0, errors.New("type is not Func")
	}
	return s.reflectType.NumIn(), nil
}

func (s safeType) NumOut() (int, error) {
	if s.reflectType.Kind() != reflect.Func {
		return 0, errors.New("type is not Func")
	}
	return s.reflectType.NumOut(), nil
}

func (s safeType) Out(i int) (SafeType, error) {
	if s.reflectType.Kind() != reflect.Func {
		return nil, errors.New("type is not Func")
	}

	if i > s.reflectType.NumOut() {
		return nil, errors.New("index out of range")
	}

	return NewSafeType(s.reflectType.Out(i)), nil
}
