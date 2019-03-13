package searcharray

import (
	"errors"
	"sync"
	"sync/atomic"

	"github.com/jinzhu/copier"

	"github.com/coc1961/mapindex"
)

//NewSearchArray NewSearchArray
func NewSearchArray() *SearchArray {
	internal := internalSearchArray{
		data:  make([]ArrayItem, 0),
		index: make(map[string]*mapindex.Index),
	}
	arr := SearchArray{}
	arr.data.Store(&internal)
	return &arr
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
	mu   sync.Mutex
	data atomic.Value
}

type internalSearchArray struct {
	data  []ArrayItem
	index map[string]*mapindex.Index
}

func (a *SearchArray) read() *internalSearchArray {
	return a.data.Load().(*internalSearchArray)
}

func (a *SearchArray) readClone() *internalSearchArray {
	orig := a.data.Load().(*internalSearchArray)
	/*
		data := make([]ArrayItem, len(orig.data))
		for i := 0; i < len(orig.data); i++ {
			data[i] = orig.data[i]
		}
		kk := make([]string, 0)
		for k := range orig.index {
			kk = append(kk, k)
		}

		ne := NewSearchArray()
		ne.Set(data, kk)
		return ne.data.Load().(*internalSearchArray)
	*/

	dest := internalSearchArray{
		data:  make([]ArrayItem, len(orig.data)),
		index: make(map[string]*mapindex.Index),
	}
	copier.Copy(&dest.data, &orig.data)
	/*
		for i := 0; i < len(orig.data); i++ {
			dest.data[i] = orig.data[i]
		}
	*/

	for k := range orig.index {
		dest.index[k] = mapindex.NewIndex()
		//copier.Copy(&dest.index[k].Idx, orig.index[k].Idx)
		for k1, v1 := range orig.index[k].Idx {
			dest.index[k].Idx[k1] = v1
		}

	}
	return &dest

}

func (a *SearchArray) insert(arr *internalSearchArray) {
	a.data.Store(arr)
}

//Set set
func (a *SearchArray) Set(data []ArrayItem, indexField []string) {
	aData := a.read()
	aData.data = data
	aData.index = make(map[string]*mapindex.Index)

	for _, s := range indexField {
		idx := mapindex.NewIndex()
		aData.index[s] = idx
	}

	i := 0
	for _, d := range data {
		for _, s := range indexField {
			aData.index[s].Add(d.GetValue(d, s), i)
		}
		i++
	}
}

//Find Find
func (a *SearchArray) Find(querys ...Q) ([]interface{}, []int, error) {
	aData := a.read()

	var ret []int

	for _, q := range querys {
		index := aData.index[q.K]
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
		res = append(res, aData.data[i])
	}
	return res, ret, nil
}

//Add Add
func (a *SearchArray) Add(item ArrayItem) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	aData := a.readClone()

	aData.data = append(aData.data, item)
	ind := len(aData.data) - 1
	for k, i := range aData.index {
		it := item.GetValue(item, k)
		i.Add(it, ind)
	}
	a.insert(aData)
	return nil
}

//Delete delete
func (a *SearchArray) Delete(ind int) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	aData := a.readClone()

	if ind < 0 || ind >= len(aData.data) {
		return errors.New("Invalid Record")
	}
	dat := aData.data[ind]
	aData.data = append(aData.data[:ind], aData.data[ind+1:]...)
	var err, err1 error
	for k, i := range aData.index {
		err1 = i.Delete(dat.GetValue(dat, k), ind)
		if err1 != nil {
			err = err1
		}
	}
	a.insert(aData)
	return err
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
