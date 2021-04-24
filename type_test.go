package anyiter_test

import (
	"errors"
	"github.com/levisaya/anyiter"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type testInterface interface {
	GetSomeField() int
}

type testStruct struct {
	someField int
}

type secondTestStruct struct {
	t testStruct
}

func (ts testStruct) GetSomeField() int { return ts.someField }
func (ts testStruct) AddToField(vals ...int) {
	for _, val := range vals {
		ts.someField += val
	}
}

func TestSafeType_ReflectType(t *testing.T) {
	typ := reflect.TypeOf(1)
	assert.Equal(t, typ, anyiter.NewSafeType(typ).ReflectType())
}

func TestSafeType_Align(t *testing.T) {
	typ := reflect.TypeOf(1)
	assert.Equal(t, typ.Align(), anyiter.NewSafeType(typ).Align())
}

func TestSafeType_FieldAlign(t *testing.T) {
	typ := reflect.TypeOf(1)
	assert.Equal(t, typ.FieldAlign(), anyiter.NewSafeType(typ).FieldAlign())
}

func TestSafeType_Method(t *testing.T) {
	t.Run("out of range", func(t *testing.T) {
		typ := reflect.TypeOf(testStruct{someField: 5})
		_, err := anyiter.NewSafeType(typ).Method(100)
		assert.Equal(t, errors.New("index out of range"), err)
	})

	t.Run("success", func(t *testing.T) {
		typ := reflect.TypeOf(testStruct{someField: 5})
		method, err := anyiter.NewSafeType(typ).Method(0)
		assert.NotNil(t, method)
		assert.Nil(t, err)
	})
}

func TestSafeType_MethodByName(t *testing.T) {
	typ := reflect.TypeOf(testStruct{someField: 5})

	_, expectedFound := typ.MethodByName("GetSomeField")
	method, ok := anyiter.NewSafeType(typ).MethodByName("GetSomeField")
	assert.NotNil(t, method)
	assert.Equal(t, expectedFound, ok)
}

func TestSafeType_NumMethod(t *testing.T) {
	typ := reflect.TypeOf(testStruct{someField: 5})

	assert.Equal(t, 2, typ.NumMethod())
	assert.Equal(t, 2, anyiter.NewSafeType(typ).NumMethod())
}

func TestSafeType_Name(t *testing.T) {
	typ := reflect.TypeOf(testStruct{someField: 5})
	assert.Equal(t, typ.Name(), anyiter.NewSafeType(typ).Name())
}

func TestSafeType_PkgPath(t *testing.T) {
	typ := reflect.TypeOf(testStruct{someField: 5})
	assert.Equal(t, typ.PkgPath(), anyiter.NewSafeType(typ).PkgPath())
}

func TestSafeType_Size(t *testing.T) {
	typ := reflect.TypeOf(testStruct{someField: 5})
	assert.Equal(t, typ.Size(), anyiter.NewSafeType(typ).Size())
}

func TestSafeType_String(t *testing.T) {
	typ := reflect.TypeOf(testStruct{someField: 5})
	assert.Equal(t, typ.String(), anyiter.NewSafeType(typ).String())
}

func TestSafeType_Kind(t *testing.T) {
	typ := reflect.TypeOf(testStruct{someField: 5})
	assert.Equal(t, typ.Kind(), anyiter.NewSafeType(typ).Kind())
}

func TestSafeType_Implements(t *testing.T) {
	typ := reflect.TypeOf(testStruct{someField: 5})
	interfaceType := anyiter.NewSafeType(reflect.TypeOf((*testInterface)(nil)).Elem())
	assert.True(t, anyiter.NewSafeType(typ).Implements(interfaceType))
}

func TestSafeType_AssignableTo(t *testing.T) {
	typ := reflect.TypeOf(testStruct{someField: 5})
	safeType := anyiter.NewSafeType(typ)
	assert.True(t, safeType.AssignableTo(safeType))
}

func TestSafeType_ConvertibleTo(t *testing.T) {
	typ := reflect.TypeOf(testStruct{someField: 5})
	safeType := anyiter.NewSafeType(typ)
	assert.True(t, safeType.ConvertibleTo(safeType))
}

func TestSafeType_Comparable(t *testing.T) {
	typ := reflect.TypeOf(testStruct{someField: 5})
	safeType := anyiter.NewSafeType(typ)
	assert.True(t, safeType.Comparable())
}

func TestSafeType_Bits(t *testing.T) {
	t.Run("invalid type", func(t *testing.T) {
		typ := reflect.TypeOf(testStruct{someField: 5})
		safeType := anyiter.NewSafeType(typ)
		_, err := safeType.Bits()
		assert.Equal(t, errors.New("not a Int*, Uint*, Float* or Complex*"), err)
	})

	t.Run("valid type", func(t *testing.T) {
		typ := reflect.TypeOf(uint8(1))
		safeType := anyiter.NewSafeType(typ)
		bits, err := safeType.Bits()
		assert.Equal(t, 8, bits)
		assert.Nil(t, err)
	})
}

func TestSafeType_ChanDir(t *testing.T) {
	t.Run("invalid type", func(t *testing.T) {
		typ := reflect.TypeOf(testStruct{someField: 5})
		safeType := anyiter.NewSafeType(typ)
		_, err := safeType.ChanDir()
		assert.Equal(t, errors.New("type is not Chan"), err)
	})

	t.Run("valid type", func(t *testing.T) {
		typ := reflect.TypeOf(make(chan int))
		safeType := anyiter.NewSafeType(typ)
		cd, err := safeType.ChanDir()
		assert.NotEqual(t, 0, cd)
		assert.Nil(t, err)
	})
}

func TestSafeType_IsVariadic(t *testing.T) {
	t.Run("invalid type", func(t *testing.T) {
		typ := reflect.TypeOf(testStruct{someField: 5})
		safeType := anyiter.NewSafeType(typ)
		_, err := safeType.IsVariadic()
		assert.Equal(t, errors.New("type is not Func"), err)
	})

	t.Run("valid type", func(t *testing.T) {
		typ := reflect.TypeOf(testStruct{someField: 5}.AddToField)
		safeType := anyiter.NewSafeType(typ)
		variadric, err := safeType.IsVariadic()
		assert.True(t, variadric)
		assert.Nil(t, err)
	})
}

func TestSafeType_Elem(t *testing.T) {
	t.Run("invalid type", func(t *testing.T) {
		typ := reflect.TypeOf(testStruct{someField: 5})
		safeType := anyiter.NewSafeType(typ)
		_, err := safeType.Elem()
		assert.Equal(t, errors.New("type is not Array, Chan, Map, Ptr or Slice"), err)
	})

	t.Run("valid type", func(t *testing.T) {
		typ := reflect.TypeOf(&testStruct{someField: 5})
		safeType := anyiter.NewSafeType(typ)
		elem, err := safeType.Elem()
		assert.NotNil(t, elem)
		assert.Nil(t, err)
	})
}

func TestSafeType_Field(t *testing.T) {
	t.Run("invalid type", func(t *testing.T) {
		typ := reflect.TypeOf(1)
		safeType := anyiter.NewSafeType(typ)
		_, err := safeType.Field(0)
		assert.Equal(t, errors.New("type is not Struct"), err)
	})

	t.Run("out of range", func(t *testing.T) {
		typ := reflect.TypeOf(testStruct{someField: 5})
		safeType := anyiter.NewSafeType(typ)
		_, err := safeType.Field(500)
		assert.Equal(t, errors.New("index out of range"), err)
	})

	t.Run("valid", func(t *testing.T) {
		typ := reflect.TypeOf(testStruct{someField: 5})
		safeType := anyiter.NewSafeType(typ)
		field, err := safeType.Field(0)
		assert.NotNil(t, field)
		assert.Nil(t, err)
	})
}

func TestSafeType_FieldByIndex(t *testing.T) {
	t.Run("invalid type", func(t *testing.T) {
		typ := reflect.TypeOf(1)
		safeType := anyiter.NewSafeType(typ)
		_, err := safeType.FieldByIndex([]int{0})
		assert.Equal(t, errors.New("type is not Struct"), err)
	})

	t.Run("out of range", func(t *testing.T) {
		typ := reflect.TypeOf(testStruct{someField: 5})
		safeType := anyiter.NewSafeType(typ)
		_, err := safeType.FieldByIndex([]int{500})
		assert.Equal(t, errors.New("index out of range"), err)
	})

	t.Run("valid", func(t *testing.T) {
		typ := reflect.TypeOf(secondTestStruct{testStruct{someField: 5}})
		safeType := anyiter.NewSafeType(typ)
		field, err := safeType.FieldByIndex([]int{0, 0})
		assert.NotNil(t, field)
		assert.Nil(t, err)
	})
}

func TestSafeType_FieldByName(t *testing.T) {
	t.Run("invalid type", func(t *testing.T) {
		typ := reflect.TypeOf(1)
		safeType := anyiter.NewSafeType(typ)
		_, ok := safeType.FieldByName("someField")
		assert.False(t, ok)
	})

	t.Run("valid", func(t *testing.T) {
		typ := reflect.TypeOf(testStruct{someField: 5})
		safeType := anyiter.NewSafeType(typ)
		field, ok := safeType.FieldByName("someField")
		assert.NotNil(t, field)
		assert.True(t, ok)
	})
}

func TestSafeType_FieldByNameFunc(t *testing.T) {
	nameFunc := func(name string) bool {
		return name == "someField"
	}

	t.Run("invalid type", func(t *testing.T) {
		typ := reflect.TypeOf(1)
		safeType := anyiter.NewSafeType(typ)
		_, ok := safeType.FieldByNameFunc(nameFunc)
		assert.False(t, ok)
	})

	t.Run("valid", func(t *testing.T) {
		typ := reflect.TypeOf(testStruct{someField: 5})
		safeType := anyiter.NewSafeType(typ)
		field, ok := safeType.FieldByNameFunc(nameFunc)
		assert.NotNil(t, field)
		assert.True(t, ok)
	})
}

func TestSafeType_In(t *testing.T) {
	t.Run("invalid type", func(t *testing.T) {
		typ := reflect.TypeOf(testStruct{someField: 5})
		safeType := anyiter.NewSafeType(typ)
		_, err := safeType.In(1)
		assert.Equal(t, errors.New("type is not Func"), err)
	})

	t.Run("out of range", func(t *testing.T) {
		typ := reflect.TypeOf(testStruct{someField: 5}.AddToField)
		safeType := anyiter.NewSafeType(typ)
		_, err := safeType.In(500)
		assert.Equal(t, errors.New("index out of range"), err)
	})

	t.Run("valid type", func(t *testing.T) {
		typ := reflect.TypeOf(testStruct{someField: 5}.AddToField)
		safeType := anyiter.NewSafeType(typ)
		retType, err := safeType.In(0)
		assert.NotNil(t, retType)
		assert.Nil(t, err)
	})
}

func TestSafeType_Key(t *testing.T) {
	t.Run("invalid type", func(t *testing.T) {
		typ := reflect.TypeOf(1)
		safeType := anyiter.NewSafeType(typ)
		_, err := safeType.Key()
		assert.Equal(t, errors.New("type is not Map"), err)
	})

	t.Run("valid", func(t *testing.T) {
		typ := reflect.TypeOf(map[string]string{})
		safeType := anyiter.NewSafeType(typ)
		retType, err := safeType.Key()
		assert.Equal(t, anyiter.NewSafeType(reflect.TypeOf("")), retType)
		assert.Nil(t, err)
	})
}

func TestSafeType_Len(t *testing.T) {
	t.Run("invalid type", func(t *testing.T) {
		typ := reflect.TypeOf(1)
		safeType := anyiter.NewSafeType(typ)
		_, err := safeType.Len()
		assert.Equal(t, errors.New("type is not Array"), err)
	})

	t.Run("valid", func(t *testing.T) {
		var a [4]int
		typ := reflect.TypeOf(a)
		safeType := anyiter.NewSafeType(typ)
		l, err := safeType.Len()
		assert.Equal(t, 4, l)
		assert.Nil(t, err)
	})
}

func TestSafeType_NumField(t *testing.T) {
	t.Run("invalid type", func(t *testing.T) {
		typ := reflect.TypeOf(1)
		safeType := anyiter.NewSafeType(typ)
		_, err := safeType.NumField()
		assert.Equal(t, errors.New("type is not Struct"), err)
	})

	t.Run("valid", func(t *testing.T) {
		typ := reflect.TypeOf(testStruct{someField: 5})
		safeType := anyiter.NewSafeType(typ)
		n, err := safeType.NumField()
		assert.Equal(t, 1, n)
		assert.Nil(t, err)
	})
}

func TestSafeType_NumIn(t *testing.T) {
	t.Run("invalid type", func(t *testing.T) {
		typ := reflect.TypeOf(testStruct{someField: 5})
		safeType := anyiter.NewSafeType(typ)
		_, err := safeType.NumIn()
		assert.Equal(t, errors.New("type is not Func"), err)
	})

	t.Run("valid type", func(t *testing.T) {
		typ := reflect.TypeOf(testStruct{someField: 5}.AddToField)
		safeType := anyiter.NewSafeType(typ)
		n, err := safeType.NumIn()
		assert.Equal(t, 1, n)
		assert.Nil(t, err)
	})
}

func TestSafeType_NumOut(t *testing.T) {
	t.Run("invalid type", func(t *testing.T) {
		typ := reflect.TypeOf(testStruct{someField: 5})
		safeType := anyiter.NewSafeType(typ)
		_, err := safeType.NumOut()
		assert.Equal(t, errors.New("type is not Func"), err)
	})

	t.Run("valid type", func(t *testing.T) {
		typ := reflect.TypeOf(testStruct{someField: 5}.AddToField)
		safeType := anyiter.NewSafeType(typ)
		n, err := safeType.NumOut()
		assert.Equal(t, 0, n)
		assert.Nil(t, err)
	})
}

func TestSafeType_Out(t *testing.T) {
	t.Run("invalid type", func(t *testing.T) {
		typ := reflect.TypeOf(testStruct{someField: 5})
		safeType := anyiter.NewSafeType(typ)
		_, err := safeType.Out(0)
		assert.Equal(t, errors.New("type is not Func"), err)
	})

	t.Run("index out of range", func(t *testing.T) {
		typ := reflect.TypeOf(testStruct{someField: 5}.GetSomeField)
		safeType := anyiter.NewSafeType(typ)
		_, err := safeType.Out(500)
		assert.Equal(t, errors.New("index out of range"), err)
	})

	t.Run("valid", func(t *testing.T) {
		typ := reflect.TypeOf(testStruct{someField: 5}.GetSomeField)
		safeType := anyiter.NewSafeType(typ)
		retType, err := safeType.Out(0)
		assert.Equal(t, anyiter.NewSafeType(reflect.TypeOf(1)), retType)
		assert.Nil(t, err)
	})
}