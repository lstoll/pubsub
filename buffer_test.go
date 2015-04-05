package pubsub

import (
	"reflect"
	"testing"
)

func TestBufferWrite(t *testing.T) {
	buffer := NewBuffer(3, 1)

	want := []interface{}{}
	if got := buffer.Read(); !reflect.DeepEqual(want, got) {
		t.Errorf("want empty buffer read %v, got %v", want, got)
	}

	buffer.Write("one")
	want = append(want, "one")
	if got := buffer.Read(); !reflect.DeepEqual(want, got) {
		t.Errorf("want sparse buffer read %v, got %v", want, got)
	}

	buffer.Write("two")
	buffer.Write("three")
	want = append(want, "two", "three")
	if got := buffer.Read(); !reflect.DeepEqual(want, got) {
		t.Errorf("want full buffer read %v, got %v", want, got)
	}

	buffer.Write("four")
	want = append(want[1:], "four")
	if got := buffer.Read(); !reflect.DeepEqual(want, got) {
		t.Errorf("want wrapped buffer read %v, got %v", want, got)
	}
}

func TestBufferReadTo(t *testing.T) {
	buffer := NewBuffer(3, 1)
	donec := make(chan struct{})
	want := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I"}

	got := []string{}
	var rfn ReaderFunc = func(v interface{}) bool {
		got = append(got, v.(string))
		if len(got) == len(want) {
			close(donec)
			return false
		}
		return true
	}
	buffer.ReadTo(rfn)

	for _, v := range want {
		buffer.Write(v)
	}

	<-donec
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want buffer read to %v, got %v", want, got)
	}
}

func TestBufferWriteSlice(t *testing.T) {
	buffer := NewBuffer(3, 1)
	donec := make(chan struct{})
	want := []interface{}{"A", "B", "C", "D", "E"}

	got := []interface{}{}
	var rfn ReaderFunc = func(v interface{}) bool {
		got = append(got, v)
		if len(got) == len(want) {
			close(donec)
			return false
		}
		return true
	}
	buffer.ReadTo(rfn)

	buffer.WriteSlice(want)

	<-donec
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want buffer read to %v, got %v", want, got)
	}

	got = buffer.Read()
	if !reflect.DeepEqual(want[2:], got) {
		t.Errorf("want buffer read %v, got %v", want[2:], got)
	}
}

func TestBufferFullReadTo(t *testing.T) {
	buffer := NewBuffer(3, 1)
	donec := make(chan struct{})
	data := []interface{}{"A", "B", "C", "D", "E", "F", "G", "H", "I"}

	buffer.WriteSlice(data[:5])

	got := []interface{}{}
	var rfn ReaderFunc = func(v interface{}) bool {
		got = append(got, v)
		if len(got) == len(data[5:]) {
			close(donec)
			return false
		}
		return true
	}

	want := data[2:5]
	if s := buffer.FullReadTo(rfn); !reflect.DeepEqual(want, s) {
		t.Errorf("want full read %v, got %v", want, s)
	}

	buffer.WriteSlice(data[5:])

	<-donec
	want = data[5:]
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want full read func %v, got %v", want, got)
	}

}
