package searcharray

import (
	"errors"
	"sort"

	"github.com/coc1961/searcharray/mapindex"
)

//NewSearchArray NewSearchArray
func NewSearchArray() *SearchArray {
	internal := internalSearchArray{
		index: make(map[string]*mapindex.Index),
	}
	arr := SearchArray{}
	arr.data = &internal
	return &arr
}

//FnGetFieldValue obtiene el valor de i¡un campo
type FnGetFieldValue func(ind int, indexField string) mapindex.IndexValue

//FnResult recibe un registro encontrado
type FnResult func(int) error

//ArrayItem ArrayItem
type ArrayItem interface {
	GetValue(indexField string) mapindex.IndexValue
}

//Q query
type Q struct {
	K string
	V mapindex.IndexValue
}

//SearchArray SearchArray
type SearchArray struct {
	data *internalSearchArray
}

type internalSearchArray struct {
	index map[string]*mapindex.Index
}

func (a *SearchArray) read() *internalSearchArray {
	return a.data
}

//Set set
func (a *SearchArray) Set(fn FnGetFieldValue, len int, indexField []string) {
	aData := a.read()
	aData.index = make(map[string]*mapindex.Index)

	ind := newIndexer(fn, len, indexField)
	for _, i := range ind {
		idx := mapindex.NewIndex()
		idx.Idx = i.Index
		aData.index[i.Name] = idx
	}
}

type arraySort [][]int

func (a arraySort) Len() int           { return len(a) }
func (a arraySort) Less(i, j int) bool { return len(a[i]) < len(a[j]) }
func (a arraySort) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

//Find Find
func (a *SearchArray) Find(fn FnResult, querys ...Q) ([]int, error) {
	aData := a.read()

	var ret []int

	records := make([][]int, 0)

	for _, q := range querys {
		index := aData.index[q.K]
		if index == nil {
			return nil, errors.New("Invalid Index " + q.K)
		}
		ind := index.Get(q.V)

		records = append(records, ind)
	}

	sort.Sort(arraySort(records))

	for _, ind := range records {
		if ret == nil {
			ret = ind
			continue
		}
		ret = intersection(ret, ind)
		if ret == nil {
			break
		}
	}

	if ret == nil {
		ret = make([]int, 0)
	}
	if fn != nil {
		for _, r := range ret {
			if err := fn(r); err != nil {
				return nil, err
			}
		}
	}
	return ret, nil
}

//Index index
func (a *SearchArray) Index(indexName string) *mapindex.Index {
	aData := a.read()

	return aData.index[indexName]
}

//Q q
func (a *SearchArray) Q(indexName string, searchValue mapindex.IndexValue) Q {
	return Q{K: indexName, V: searchValue}
}

// Retorna la intersección de dos array de ints
func intersection(a, b []int) (c []int) {
	m := make(map[int]bool, 0)

	for _, item := range a {
		m[item] = true
	}

	for _, item := range b {
		if _, ok := m[item]; ok {
			c = append(c, item)
		}
	}
	return
}
