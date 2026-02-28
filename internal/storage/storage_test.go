package storage

import(
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

type myInt int64

func (i myInt) String() string {
	return fmt.Sprintf("%d", i)
}

func (i myInt) Key() string {
	return i.String()
}

func TestNewStorage (t *testing.T) {
	t.Run("Create int64 Storage", func(t *testing.T) {
		s := NewStorage[*myInt]()

		if s.ds == nil {
			t.Fatal("Map didn't initialized")
		}

		s.ds["0"] = new(myInt)
	})
}

func TestSet(t *testing.T) {
	storage := NewStorage[*myInt]()

	type want struct {
		storedValue myInt
	}

	tests := []struct {
		name string
		itemValue myInt
		want want
	}{
		{
			name: "Test Set myInt",
			itemValue: myInt(100),
			want: want{storedValue: myInt(100)},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage.Set(test.itemValue.Key(), &test.itemValue)
			value, _ := storage.Get(test.itemValue.Key())
			assert.Equal(t, &test.want.storedValue, value)
		})
	}
}

func TestGet(t *testing.T) {
	storage := NewStorage[myInt]()
	m := myInt(1)
	storage.Set(m.Key(), m)

	type want struct {
		returnValue myInt
		valueExist bool
	}

	tests := []struct {
		name string
		itemName string
		want want
	}{
		{
			name: "Test Item that exist",
			itemName: "1",
			want: want{ returnValue: m, valueExist: true },
		},
		{
			name: "Test Item that not exist",
			itemName: "0",
			want: want{ returnValue: myInt(0), valueExist: false },
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value, ok := storage.Get(test.itemName)
			assert.Equal(t, test.want.returnValue, value)
			assert.Equal(t, test.want.valueExist, ok)
		})
	}
}

func TestRemove(t *testing.T) {
	storage := NewStorage[*myInt]()
	testInt := myInt(2)

	storage.Set(testInt.Key(), &testInt)
	tests := []struct {
		name string
		itemName string
	}{
		{
			name: "Test testInt Remove",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage.Remove(testInt.Key())
			_, ok := storage.Get(testInt.Key())
			assert.Equal(t, false, ok)
		})
	}
}

func TestAll(t *testing.T) {
	st := NewStorage[myInt]()

	st.Set("1", myInt(1))
	st.Set("2", myInt(2))

	t.Run("Should iterate through inner map values", func(t *testing.T) {
		for v:= range st.All() {
			assert.IsType(t, myInt(1), v)
		}
		for v := range st.All() {
			assert.NotNil(t, v)
			break
		}
	})
}
