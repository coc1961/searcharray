package searcharray

import (
	"sync"

	"github.com/coc1961/searcharray/mapindex"
)

//Newindexer indexer
func newIndexer(fn func(ind int, indexField string) mapindex.IndexValue, len int, fields []string) []outIndex {

	// Creo los chan y conecto los flujos a través de ellos
	Index1 := make([]chan outIndex, 0, len)
	AllIndex := make(chan []outIndex, len)

	Bins := make(chan ArrayItem)
	LIndex1 := make([]chan register, 0, len)

	for _, s := range fields {
		tmp := make(chan register, 10000)
		out := make(chan outIndex, 1)
		LIndex1 = append(LIndex1, tmp)
		Index1 = append(Index1, out)

		indexer1 := indexer{
			Bins:  tmp,
			field: s,
			Index: out,
		}
		indexer1.Process()
	}
	//Creo los procesos de la red
	indexJoiner := indexJoiner{
		AllIndex: AllIndex,
		Index:    Index1,
	}
	loader := loader{
		Bins:  Bins,
		Index: LIndex1,
	}

	// Inicio los procesos
	indexJoiner.Process()
	loader.Process()

	for b := 0; b < len; b++ {
		Bins <- toArrayItem(fn, b)
	}
	close(Bins)
	ret := <-AllIndex
	return ret
}

type wrap struct {
	fn    func(ind int, indexField string) mapindex.IndexValue
	index int
}

func (w wrap) GetValue(indexField string) mapindex.IndexValue {
	return w.fn(w.index, indexField)
}

func toArrayItem(fn func(ind int, indexField string) mapindex.IndexValue, index int) ArrayItem {
	return wrap{fn: fn, index: index}
}

//indexJoiner loader
type indexJoiner struct {
	Index    []chan outIndex
	AllIndex chan<- []outIndex
}

func (i *indexJoiner) merge() <-chan outIndex {
	var wg sync.WaitGroup
	out := make(chan outIndex)

	output := func(c <-chan outIndex) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(i.Index))
	for _, ix := range i.Index {
		go output(ix)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

//Process process
func (i *indexJoiner) Process() {
	indexLog("indexJoiner starts.")
	in := i.merge()
	go func() {
		lst := make([]outIndex, 0)
		for {
			idx, ok := <-in
			if !ok {
				indexLog("indexJoiner has finished.")
				i.AllIndex <- lst
				close(i.AllIndex)
				return
			}
			lst = append(lst, idx)
		}
	}()
}

//loader loader
type loader struct {
	Bins  <-chan ArrayItem
	Index []chan register
}

//Process process
func (i *loader) Process() {
	indexLog("loader starts.")
	go func() {
		ix := 0
		for {
			bins, ok := <-i.Bins
			if !ok {
				indexLog("loader has finished.")
				for _, i := range i.Index {
					close(i)
				}
				return
			}
			for _, i := range i.Index {
				i <- register{Row: ix, Bin: bins}
			}
			ix++
		}
	}()
}

//register register
type register struct {
	Row int
	Bin ArrayItem
}

//outIndex register
type outIndex struct {
	Name  string
	Index map[mapindex.IndexValue][]int
}

//indexer indexer
type indexer struct {
	Bins  <-chan register
	field string
	Index chan outIndex
}

//Process process
func (i *indexer) Process() {
	indexLog("indexer starts.")
	go func() {
		idx := map[mapindex.IndexValue][]int{}
		for {
			bin, ok := <-i.Bins
			if !ok {
				indexLog("indexer has finished.")
				i.Index <- outIndex{Name: i.field, Index: idx}
				close(i.Index)
				return
			}
			val := bin.Bin.GetValue(i.field)
			if idx[val] == nil {
				idx[val] = make([]int, 0)
			}
			idx[val] = append(idx[val], bin.Row)
		}
	}()
}

func indexLog(txt string) {
	//fmt.Println(txt)
}
