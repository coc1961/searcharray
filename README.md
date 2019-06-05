# SearchArray

This framework allows to perform searches on arrays of data structures.


Example:

```Go
type TestStruct type {
    Field1 string 
    Field2 int
}


func main(){

    var myArray []*TestStruct

    //Fill the array
    myArray=....
    
    //initialize the searcharray structure
    sa := searcharray.NewSearchArray()

    //create indexes of the fields Field1 and Field2   
    idx := []string{"Field1", "Field1"}

    //This function allows access to the attributes of the instances of TestStruct
    fnGetFieldValue := func(ind int, indexField string) mapindex.IndexValue {
		obj := myArray[ind]
        switch indexField {
        case "Field1":
            return mapindex.IndexValue(obj.Field1)
        case "Field2":
            return mapindex.IndexValue(obj.Field2)
        }
        return mapindex.IndexValue(nil)
    }
    
    //Initialize with the array and the indices to create
    sa.Set(fnGetFieldValue, len(myArray), idx)

    //I look for the records according to the filter 
    res, _, err := sa.Find( func (i int) error { fmt.Print("Record found",data[i]); return nil }
        sa.Q("Field1", "Hello"),
        sa.Q("Field2", 20),
    )

}
```

