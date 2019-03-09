package searcharray_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/coc1961/mapindex"
	"github.com/coc1961/searcharray"
)

func TestSearchArray_Set(t *testing.T) {
	data := ReadDataAsObject()

	sa := searcharray.NewSearchArray()
	idx := []string{"Bin", "Segment", "Brand", "Issuer", "Country"}

	sa.Set(data, idx)

	start := time.Now()
	res, err := sa.Find(
		sa.Q("Country", "AR"),
		sa.Q("Issuer", int(1)),
		sa.Q("Brand", int(5)),
	)
	fmt.Printf("Set Elapsed Time %s, Found Records %d, Total Records %d\n", time.Since(start), len(res), len(data))
	if err != nil || len(res) != 6 {
		t.Errorf("Error Find\n")
	}

	start = time.Now()
	res, err = sa.Find(
		sa.Q("Bin", int(290107)),
	)
	fmt.Printf("Set Elapsed Time %s, Found Records %d, Total Records %d\n", time.Since(start), len(res), len(data))
	if err != nil || len(res) != 1 {
		t.Errorf("Error Find\n")
	}

	// Get Bin Index
	ib := sa.Index("Bin")

	//Convert to int
	ints := make([]int, 0, len(ib.Values()))
	for _, v := range ib.Values() {
		ints = append(ints, v.(int))
	}

	//Sort Bins
	sort.Ints(ints)

	//Search 290107 Bin
	found := sort.Search(len(ints), func(i int) bool {
		if ints[i] < 200000 {
			return false
		}
		return true
	})

	cont := 0
	for i := found; i < len(ints); i++ {
		v := ints[i]
		if v > 295000 {
			break
		}
		r, _ := sa.Find(sa.Q("Bin", v))
		cont += len(r)
	}

	fmt.Println(cont)

}

//Bin  estructura de bin
type Bin struct {
	Bin     int
	Segment int
	Brand   int
	Issuer  int
	Country string
}

func (a *Bin) GetValue(item interface{}, indexField string) mapindex.IndexValue {
	bin := item.(*Bin)
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
