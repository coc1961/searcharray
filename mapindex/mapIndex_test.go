package mapindex_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/coc1961/searcharray/mapindex"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestNewIndex(t *testing.T) {
	idx := mapindex.NewIndex()
	arr := make([]mapindex.IndexValue, 300000)

	var value mapindex.IndexValue
	var value1 mapindex.IndexValue
	for i := 0; i < len(arr); i++ {
		v := RandStringRunes(2)
		arr[i] = v
		idx.Add(v, i)
		if value == nil {
			value = v
		} else if value1 == nil {
			value1 = v
		}

	}

	if len(idx.Get(mapindex.IndexValue(value))) < 1 {
		t.Errorf("Error Find\n")
	}

	if len(idx.Get(mapindex.IndexValue([]interface{}{value, value1}))) <= len(idx.Get(mapindex.IndexValue(value))) {
		t.Errorf("Error Find\n")
	}

	cont := 0
	ar := idx.Get(mapindex.IndexValue(value))

	for _, i := range ar {
		if cont != 0 {
			idx.Delete(arr[i], i)
		}
		cont++
	}
	if len(idx.Get(mapindex.IndexValue(value))) != 1 {
		t.Errorf("Error Delete\n")
	}

	for _, i := range idx.Get(mapindex.IndexValue(value)) {
		idx.Delete(arr[i], i)
	}
	if len(idx.Get(mapindex.IndexValue(value))) != 0 {
		t.Errorf("Error Delete Ultimo\n")
	}

}

//var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var letterRunes = []rune("ABCDEFGHIJKLMN")

//RandStringRunes RandStringRunes
func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
