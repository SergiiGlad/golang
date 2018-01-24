package messages

import (
	"fmt"
	"html/template"
	"net/http"
)

//HandlerOfMessages This func should process any messages end point
func HandlerOfMessages(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("client/index.html")
	if err != nil {
		fmt.Fprintf(w, "%s", "Error")
	} else {
		tmpl.Execute(w, r)
	}
	fmt.Fprintf(w, "Hello any one! Message %s", r.URL.Path[1:])

}
