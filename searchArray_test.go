package searcharray_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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
	res, _, err := sa.Find(
		sa.Q("Country", "AR"),
		sa.Q("Issuer", int(1)),
		sa.Q("Brand", int(5)),
	)
	fmt.Printf("Set Elapsed Time %s, Found Records %d, Total Records %d\n", time.Since(start), len(res), len(data))
	if err != nil || len(res) != 6 {
		t.Errorf("Error Find\n")
	}

	start = time.Now()
	res, _, err = sa.Find(
		sa.Q("Bin", int(290107)),
	)
	fmt.Printf("Set Elapsed Time %s, Found Records %d, Total Records %d\n", time.Since(start), len(res), len(data))
	if err != nil || len(res) != 1 {
		t.Errorf("Error Find\n")
	}

	start = time.Now()

	//Find 25000 Bins
	cont := 0
	for i := 200000; i < 225000; i++ {
		r, _, _ := sa.Find(sa.Q("Bin", i))
		cont += len(r)
	}

	fmt.Printf("Set Elapsed Time %s, Found Records %d, Total Records %d\n", time.Since(start), cont, len(data))

	if cont != 25000 {
		t.Errorf("Error Find Bins\n")
	}

	start = time.Now()
	var recs []int
	res, recs, err = sa.Find(
		sa.Q("Country", "AR"),
		sa.Q("Issuer", int(1)),
		sa.Q("Brand", int(5)),
	)
	if err != nil || len(res) != 6 {
		t.Errorf("Error Find\n")
	}

	sa.Delete(recs[2])

	res, recs, err = sa.Find(
		sa.Q("Country", "AR"),
		sa.Q("Issuer", int(1)),
		sa.Q("Brand", int(5)),
	)
	if err != nil || len(res) != 5 {
		t.Errorf("Error Delete\n")
	}

	b := &Bin{
		Bin:     111,
		Country: "AR",
		Issuer:  int(1),
		Brand:   int(5),
	}

	sa.Add(b)

	res, recs, err = sa.Find(
		sa.Q("Country", "AR"),
		sa.Q("Issuer", int(1)),
		sa.Q("Brand", int(5)),
	)
	if err != nil || len(res) != 6 {
		t.Errorf("Error Delete\n")
	}
	fmt.Printf("Set Elapsed Time %s, Found Records %d, Total Records %d\n", time.Since(start), len(res), len(data))

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
