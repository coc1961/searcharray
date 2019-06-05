package searcharray_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/coc1961/searcharray"
	"github.com/coc1961/searcharray/mapindex"
)

func TestSearchArray_Set(t *testing.T) {
	data := ReadDataAsObject()

	sa := searcharray.NewSearchArray()
	idx := []string{"Bin", "Segment", "Brand", "Issuer", "Country"}

	start := time.Now()
	sa.Set(data, idx)
	fmt.Printf("Set Data Elapsed Time %s\n", time.Since(start))

	start = time.Now()
	res := make([]*Bin, 0)
	_, err := sa.Find(func(i int) error { res = append(res, data[i].(*Bin)); return nil },
		sa.Q("Country", "AR"),
		sa.Q("Issuer", int(1)),
		sa.Q("Brand", int(5)),
	)
	fmt.Printf("Set Elapsed Time %s, Found Records %d, Total Records %d\n", time.Since(start), len(res), len(data))
	if err != nil || len(res) != 6 {
		t.Errorf("Error Find\n")
	}

	start = time.Now()
	res = make([]*Bin, 0)
	_, err = sa.Find(func(i int) error { res = append(res, data[i].(*Bin)); return nil },
		sa.Q("Bin", int(290107)),
	)
	fmt.Printf("Set Elapsed Time %s, Found Records %d, Total Records %d\n", time.Since(start), len(res), len(data))
	if err != nil || len(res) != 1 {
		t.Errorf("Error Find\n")
	}

	start = time.Now()

	//Find 25000 Bins
	cont := 0
	res = make([]*Bin, 0)
	for i := 200000; i < 225000; i++ {
		r, _ := sa.Find(func(i int) error { res = append(res, data[i].(*Bin)); return nil }, sa.Q("Bin", i))
		cont += len(r)
	}

	fmt.Printf("Set Elapsed Time %s, Found Records %d, Total Records %d\n", time.Since(start), cont, len(data))

	if cont != 25000 {
		t.Errorf("Error Find Bins\n")
	}

	res = make([]*Bin, 0)
	_, err = sa.Find(func(i int) error { res = append(res, data[i].(*Bin)); return nil },
		sa.Q("Country", "AR"),
		sa.Q("Issuer", int(1)),
		sa.Q("Brand", int(5)),
	)
	if err != nil || len(res) != 6 {
		t.Errorf("Error Find\n")
	}

	wg := sync.WaitGroup{}

	ch1 := int32(0)
	ch2 := int32(0)
	ch3 := int32(0)

	for a := 0; a < 1000; a++ {
		wg.Add(1)
		go func(aa, id int) {
			time.Sleep(time.Duration(aa) * time.Millisecond)
			a := make([]*Bin, 0)
			b, c := sa.Find(func(i int) error { res = append(res, data[i].(*Bin)); return nil }, sa.Q("Bin", id))
			_ = a
			_ = b
			_ = c
			if len(b) > 0 {
				if atomic.LoadInt32(&ch2) > 0 {
					atomic.AddInt32(&ch3, 1)
				} else {
					atomic.AddInt32(&ch1, 1)
				}
			} else {
				atomic.AddInt32(&ch2, 1)
			}
			wg.Done()
		}(a, res[2].Bin)
	}

	time.Sleep(300 * time.Millisecond)

	start = time.Now()

	for a := 0; a < 100; a++ {
		wg.Add(1)
		go func(aa, id int) {
			time.Sleep(time.Duration(aa) * time.Millisecond)
			b, c := sa.Find(nil, sa.Q("Bin", id))
			_ = a
			_ = b
			_ = c
			if len(b) > 0 {
				if atomic.LoadInt32(&ch2) > 0 {
					atomic.AddInt32(&ch3, 1)
				} else {
					atomic.AddInt32(&ch1, 1)
				}
			} else {
				atomic.AddInt32(&ch2, 1)
			}
			wg.Done()
		}(a, res[2].Bin)
	}

	end := time.Since(start)
	wg.Wait()
	fmt.Printf("Set Elapsed Time %s, Found Records %d, Total Records %d , No Cero (prev) %d , No Cero (pos)  %d , Cero %d\n", end, len(res), len(data), ch1, ch3, ch2)

}

//Bin  estructura de bin
type Bin struct {
	Bin     int
	Segment int
	Brand   int
	Issuer  int
	Country string
}

func (a *Bin) GetValue(indexField string) mapindex.IndexValue {
	bin := a
	switch indexField {
	case "Bin":
		return mapindex.IndexValue(bin.Bin)
	case "Segment":
		return mapindex.IndexValue(bin.Segment)
	case "Brand":
		return mapindex.IndexValue(bin.Brand)
	case "Issuer":
		return mapindex.IndexValue(bin.Issuer)
	case "Country":
		return mapindex.IndexValue(bin.Country)
	}
	return mapindex.IndexValue(nil)
}

//ReadDataAsObject datos como estructura
func ReadDataAsObject() []searcharray.ArrayItem {
	//fmt.Println("Iniciando Lectura de Json")
	jsonFile, err := os.Open("test/data.json")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var data []Bin
	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	rt := make([]searcharray.ArrayItem, 0, len(data))
	for i := 0; i < len(data); i++ {
		data[i].Bin = i + 200000
		rt = append(rt, &data[i])
	}
	return rt

}
