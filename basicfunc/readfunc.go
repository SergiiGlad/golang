package basicfunc

import (
    "bufio"
    "fmt"
    "strings"
    "os"
    "strconv"
)


func ReadString( quest string) string  {

    reader := bufio.NewReader(os.Stdin)

    fmt.Println(  quest )
    text, _ := reader.ReadString('\n')

    // convert CRLF to LF

	text = strings.Replace( text, "\n", "", -1 )

    return text
}


func ReadInt ( q string) (v int, err error ) {

       return strconv.Atoi( ReadString( q ) )
}