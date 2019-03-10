package searcharray

import (
	"errors"

	"github.com/coc1961/mapindex"
)

//NewSearchArray NewSearchArray
func NewSearchArray() *SearchArray {
	return &SearchArray{
		data:  make([]ArrayItem, 0),
		index: make(map[string]*mapindex.Index),
	}
}

//ArrayItem ArrayItem
type ArrayItem interface {
	GetValue(item interface{}, indexField string) mapindex.IndexValue
}

//Q query
type Q struct {
	K string
	V mapindex.IndexValue
}

//SearchArray SearchArray
type SearchArray struct {
	data  []ArrayItem
	index map[string]*mapindex.Index
}

//Set set
func (a *SearchArray) Set(data []ArrayItem, indexField []string) {
	a.data = data
	a.index = make(map[string]*mapindex.Index)

	for _, s := range indexField {
		idx := mapindex.NewIndex()
		a.index[s] = idx
	}

	i := 0
	for _, d := range data {
		for _, s := range indexField {
			a.index[s].Add(d.GetValue(d, s), i)
		}
		i++
	}
}

//Find Find
func (a *SearchArray) Find(querys ...Q) ([]interface{}, []int, error) {

	var ret []int

	for _, q := range querys {
		index := a.index[q.K]
		if index == nil {
			return nil, nil, errors.New("Invalid Index " + q.K)
		}
		ind := index.Get(q.V)
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

	res := make([]interface{}, 0, len(ret))
	for _, i := range ret {
		res = append(res, a.data[i])
	}
	return res, ret, nil
}

//Add Add
func (a *SearchArray) Add(item ArrayItem) error {
	a.data = append(a.data, item)
	ind := len(a.data) - 1
	for k, i := range a.index {
		it := item.GetValue(item, k)
		i.Add(it, ind)
	}
	return nil
}

//Delete delete
func (a *SearchArray) Delete(ind int) error {
	if ind < 0 || ind >= len(a.data) {
		return errors.New("Invalid Record")
	}
	dat := a.data[ind]
	a.data = append(a.data[:ind], a.data[ind+1:]...)
	var err, err1 error
	for k, i := range a.index {
		err1 = i.Delete(dat.GetValue(dat, k), ind)
		if err1 != nil {
			err = err1
		}
	}
	return err
}

//Index index
func (a *SearchArray) Index(indexName string) *mapindex.Index {
	return a.index[indexName]
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
