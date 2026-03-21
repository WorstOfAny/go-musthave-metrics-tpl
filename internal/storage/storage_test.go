package storage

import(
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
	"os"
	"bufio"
	"encoding/json"
	"time"
	"context"
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
		tmpFile, err := os.CreateTemp("", "*.json")
		defer os.Remove(tmpFile.Name())

		if err != nil { t.Fatal(err) }

		s, err := NewStorage[*myInt](tmpFile, false)

		if err != nil { t.Fatal(err) }

		if s.ds == nil {
			t.Fatal("Map didn't initialized")
		}

		s.ds["0"] = new(myInt)
	})
}

func TestSet(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "*.json")
	defer os.Remove(tmpFile.Name())

	if err != nil { t.Fatal(err) }

	storage, err := NewStorage[*myInt](tmpFile, false)

	if err != nil { t.Fatal(err) }

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
	tmpFile, err := os.CreateTemp("", "*.json")
	defer os.Remove(tmpFile.Name())

	if err != nil { t.Fatal(err) }

	storage, err := NewStorage[myInt](tmpFile, false)
	if err != nil { t.Fatal(err) }

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
	tmpFile, err := os.CreateTemp("", "*.json")
	defer os.Remove(tmpFile.Name())

	if err != nil { t.Fatal(err) }

	storage, err := NewStorage[*myInt](tmpFile, false)
	if err != nil { t.Fatal(err) }

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
	tmpFile, err := os.CreateTemp("", "*.json")
	defer os.Remove(tmpFile.Name())

	if err != nil { t.Fatal(err) }

	st, err := NewStorage[myInt](tmpFile, false)
	if err != nil { t.Fatal(err) }

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

func TestRestore(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "*.json")
	defer os.Remove(tmpFile.Name())

	if err != nil { t.Fatal(err) }

	writer := bufio.NewWriter(tmpFile)

	firstVal := myInt(1)
	marshalFirstVal, err := json.Marshal(&firstVal)
	if err != nil { t.Fatal(err) }

	secondVal := myInt(2)
	marshalSecondVal, err := json.Marshal(&secondVal)
	if err != nil { t.Fatal(err) }

	_, err = writer.Write(marshalFirstVal)
	if err != nil { t.Fatal(err) }

	err = writer.WriteByte('\n')
	if err != nil { t.Fatal(err) }

	_, err = writer.Write(marshalSecondVal)
	if err != nil { t.Fatal(err) }

	err = writer.WriteByte('\n')
	if err != nil { t.Fatal(err) }

	writer.Flush()

	tmpFile.Seek(0,0)

	st, err := NewStorage[*myInt](tmpFile, true)
	if err != nil { t.Fatal(err) }

	t.Run("Should restore storage from file", func(t *testing.T) {
		val, exist := st.Get("1")
		assert.Equal(t, &firstVal, val)
		assert.NotNil(t, val)
		assert.True(t, exist)

		val, exist = st.Get("2")
		assert.Equal(t, &secondVal, val)
		assert.NotNil(t, val)
		assert.True(t, exist)
	})
}

func TestWriteToFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "*.json")
	defer os.Remove(tmpFile.Name())

	if err != nil { t.Fatal(err) }

	st, err := NewStorage[myInt](tmpFile, false)
	if err != nil { t.Fatal(err) }

	st.Set("1", myInt(1))
	st.Set("2", myInt(2))

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	errCh := make(chan error, 2)

	go st.WriteToFile(ctx, errCh, 5 * time.Second)

	time.Sleep(6 * time.Second)
	t.Run("Should write to file", func(t *testing.T) {
		tmpFile.Seek(0, 0)

		scanner := bufio.NewScanner(tmpFile)

		assert.True(t, scanner.Scan())
		assert.NoError(t, scanner.Err())
		
	})
}
