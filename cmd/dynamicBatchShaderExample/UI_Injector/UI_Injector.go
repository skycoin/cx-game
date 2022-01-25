package UI_Injector

import (
	"html/template"
	"net/http"
	"strconv"
)

type UI_Injector struct {
	Object [4]ObjectOffset
}

type ObjectOffset struct {
	Name  string
	X     float32
	Y     float32
	Z     float32
	Scane int32
}

func SetUpUI() *UI_Injector {
	uiObjects := &UI_Injector{}
	return uiObjects
}

func (UiObjs *UI_Injector) ListenForChanges() UI_Injector {
	tmpl := template.Must(template.ParseFiles("./forms.html"))
	var objType int
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}
		if r.FormValue("Name") == "Camera" {
			objType = 0
		} else if r.FormValue("Name") == "Object" {
			objType = 1
		} else if r.FormValue("Name") == "Object2" {
			objType = 2
		} else if r.FormValue("Name") == "Scane" {
			objType = 3
		}
		tempX, _ := strconv.ParseFloat(r.FormValue("X"), 32)
		tempY, _ := strconv.ParseFloat(r.FormValue("Y"), 32)
		tempZ, _ := strconv.ParseFloat(r.FormValue("Z"), 32)
		tempScane, _ := strconv.ParseInt(r.FormValue("Number"), 0, 32)

		UiObjs.Object[objType] = ObjectOffset{
			Name:  r.FormValue("Name"),
			X:     float32(tempX),
			Y:     float32(tempY),
			Z:     float32(tempZ),
			Scane: int32(tempScane),
		}
		// if r.FormValue("Name")== "Camera"{
		// }

		// do something with details
		//	_ = details

		tmpl.Execute(w, struct{ Success bool }{true})
	})
	//log.Println("http data ", UiObjs)
	http.ListenAndServe(":8080", nil)
	return *UiObjs
}
