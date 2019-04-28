package mapindex

import (
	"errors"
	"reflect"
)

/*************************************************************
 ** Index
 **************************************************************/

//NewIndex NewIndex
func NewIndex() *Index {
	idx := Index{Idx: make(map[IndexValue][]int, 0)}
	return &idx
}

//IndexValue IndexValue
type IndexValue interface{}

//Index Index
type Index struct {
	Idx map[IndexValue][]int
}

//Add add
func (i *Index) Add(val IndexValue, index int) {
	if i.Idx[val] == nil {
		i.Idx[val] = make([]int, 0)
	}
	i.Idx[val] = append(i.Idx[val], index)
}

//Get get
func (i *Index) Get(val IndexValue) []int {
	//Si es array...
	typ := reflect.ValueOf(val)
	if typ.Kind() == reflect.Slice {
		retArr := make([]int, 0)
		for ii := 0; ii < typ.Len(); ii++ {
			it := i.Idx[IndexValue(typ.Index(ii).Interface())]
			if it == nil {
				continue
			}
			ret := make([]int, len(it))
			copy(ret, it)
			retArr = append(retArr, ret...)
		}
		return retArr
	}

	// Si no array
	it := i.Idx[val]
	if it == nil {
		return make([]int, 0)
	}
	ret := make([]int, len(it))
	copy(ret, it)
	return ret
}

//Delete Delete
func (i *Index) Delete(val IndexValue, index int) error {
	arr := i.Idx[val]
	if arr == nil {
		return errors.New("Record Not Found")
	}
	for ind := 0; ind < len(arr); ind++ {
		if arr[ind] == index {
			arr = append(arr[:ind], arr[ind+1:]...)
			i.Idx[val] = arr
			if len(arr) == 0 {
				delete(i.Idx, val)
			}
			return nil
		}
	}
	return errors.New("Record Not Found")
}
