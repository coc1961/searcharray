# SearchArray

This framework allows to perform searches on arrays of data structures.


Example:

```Go
type TestStruct type {
    Field1 string 
    Field2 int
}

//This function allows access to the attributes of the instances of TestStruct
func (a *TestStruct) GetValue(item interface{}, indexField string) mapindex.IndexValue {
    obj := item.(*TestStruct)
    switch indexField {
    case "Field1":
        return mapindex.IndexValue(obj.Field1)
    case "Field2":
        return mapindex.IndexValue(obj.Field2)
    }
    return mapindex.IndexValue(nil)
}


func main(){

    var myArray []*TestStruct
    myArray=....
    
    //initialize the searcharray structure
    sa := searcharray.NewSearchArray()

    //create indexes of the fields Field1 and Field2   
    idx := []string{"Field1", "Field1"}

    //Initialize with the array and the indices to create
    sa.Set(data, idx)

    //I look for the records according to the filter 
    res, _, err := sa.Find( func (i int) error { fmt.Print("Record found",data[i]); return nil }
        sa.Q("Field1", "Hello"),
        sa.Q("Field2", 20),
    )

    //In res I get the results of the search

}
```

