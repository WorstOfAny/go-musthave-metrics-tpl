package storage

import(
	"testing"
	"github.com/stretchr/testify/assert"

)

func TestNewStorage (t *testing.T) {
	t.Run("Create int64 Storage", func(t *testing.T) {
		s := NewStorage[*int64]()

		if s.ds == nil {
			t.Fatal("Map didn't initialized")
		}

		s.ds["myInt"] = new(int64)
	})
}

func TestSet(t *testing.T) {
	storage := NewStorage[*int64]()

	type want struct {
		storedValue int64
	}

	tests := []struct {
		name string
		itemName string
		itemValue int64
		want want
	}{
		{
			name: "Test Set myInt",
			itemName: "myInt",
			itemValue: int64(100),
			want: want{storedValue: int64(100)},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage.Set(test.itemName, &test.itemValue)
			value, _ := storage.Get(test.itemName)
			assert.Equal(t, &test.want.storedValue, value)
		})
	}
}

func TestGet(t *testing.T) {
	storage := NewStorage[*int64]()
	myInt := int64(1)
	storage.Set("myInt", &myInt)

	type want struct {
		returnValue *int64
		valueExist bool
	}

	tests := []struct {
		name string
		itemName string
		want want
	}{
		{
			name: "Test Item that exist",
			itemName: "myInt",
			want: want{ returnValue: &myInt, valueExist: true },
		},
		{
			name: "Test Item that not exist",
			itemName: "testInt",
			want: want{ returnValue: nil, valueExist: false },
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
	storage := NewStorage[*int64]()
	testInt := int64(2)

	storage.Set("testInt", &testInt)
	tests := []struct {
		name string
		itemName string
	}{
		{
			name: "Test testInt Remove",
			itemName: "testInt",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage.Remove(test.itemName)
			_, ok := storage.Get(test.itemName)
			assert.Equal(t, false, ok)
		})
	}
}

func TestAll(t *testing.T) {
	st := NewStorage[int64]()

	st.Set("1", 1)
	st.Set("2", 2)

	t.Run("Should iterate through inner map values", func(t *testing.T) {
		for v:= range st.All() {
			assert.IsType(t, int64(1), v)
		}
		for v := range st.All() {
			assert.NotNil(t, v)
			break
		}
	})
}
