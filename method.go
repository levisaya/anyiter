package anyiter

import "reflect"

// SafeMethod is an interface wrapper for reflect.Method, with Type and Func modified to return SafeType and SafeFunc
type SafeMethod interface {
	ReflectMethod() reflect.Method
	Name() string
	PkgPath() string

	Type() SafeType
	Func() SafeValue
	Index() int
}

type safeMethod struct {
	method reflect.Method
}

// NewSafeMethod wraps a reflect.Method in the SafeMethod interface
func NewSafeMethod(reflectMethod reflect.Method) SafeMethod {
	return &safeMethod{method: reflectMethod}
}

func (s safeMethod) ReflectMethod() reflect.Method {
	return s.method
}

func (s safeMethod) Name() string {
	panic("implement me")
}

func (s safeMethod) PkgPath() string {
	panic("implement me")
}

func (s safeMethod) Type() SafeType {
	panic("implement me")
}

func (s safeMethod) Func() SafeValue {
	panic("implement me")
}

func (s safeMethod) Index() int {
	panic("implement me")
}
