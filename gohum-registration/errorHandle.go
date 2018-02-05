package main

import (
	"fmt"
	"os"
  //"log"
  //"log"
)

// type error interface {
//     Error() string
// }

type errorStruct struct {
    message			error
    package_name	string
    file_name 		string
    function_name 	string
}

func (s *errorStruct) Error() string{
	return fmt.Sprintf("Package:%s\nFile = %s\nFunction = %s\nERROR message = %s\n", s.package_name, s.file_name, s.function_name, s.message)
}


func goHumError(mess error, pack string, fil string, function string) *errorStruct {
	return &errorStruct {
	message : mess,
	package_name : pack,
	file_name : fil,
	function_name : function,
	}
}

func main() {
	_, err := os.Open("file.go") // For read access.
//if err != nil {
//  fmt.Println(err)
//  fmt.Println("--------")
//
//  //log.Fatal(err)
//}

	var instanceClass = goHumError(err, "main", "errorHandle.go", "main()" )

	fmt.Println(instanceClass)
}
