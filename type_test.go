package anyiter_test

import (
	"errors"
	"github.com/levisaya/anyiter"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type testStruct struct {
	someField int
}

func (ts testStruct) GetSomeField() int { return ts.someField }

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
