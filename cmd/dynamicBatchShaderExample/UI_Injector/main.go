package UI_Injector

// import (
// 	"html/template"
// 	"log"
// 	"net/http"
// )

// type ObjectOffset struct {
// 	Name string
// 	x    string
// 	y    string
// 	z    string
// }

// var details ObjectOffset

// func main() {
// 	tmpl := template.Must(template.ParseFiles("forms.html"))

// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method != http.MethodPost {
// 			tmpl.Execute(w, nil)
// 			return
// 		}

// 		details = ObjectOffset{
// 			Name: r.FormValue("Name"),
// 			x:    r.FormValue("X"),
// 			y:    r.FormValue("Y"),
// 			z:    r.FormValue("Z"),
// 		}
// 		// if r.FormValue("Name")== "Camera"{
// 		// }

// 		// do something with details
// 		_ = details
// 		log.Println(details)

// 		tmpl.Execute(w, struct{ Success bool }{true})
// 	})

// 	http.ListenAndServe(":8080", nil)
// }
