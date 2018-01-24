package messages

import (
	"fmt"
	"html/template"
	"net/http"
)

// import (
// 	"fmt"
// 	"html/template"
// 	"net/http"
// )

func handlerOfMessages(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("client/index.html")
	if err != nil {
		fmt.Fprintf(w, "%s", "Error")
	} else {
		tmpl.Execute(w, r)
	}
	//fmt.Fprintf(w, "Hello Home! %s", r.URL.Path[1:]) )

}
