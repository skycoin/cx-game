package UI_Injector

// import (
// 	"html/template"
// 	"log"
// 	"net/http"
// 	"time"
// )

// type Post struct {
// 	User, Time, Text string
// }

// var t = template.Must(template.New("").Parse(page))

// var lastTime = time.Now()

// func produceTime() string {
// 	lastTime = lastTime.Add(time.Second)
// 	return lastTime.Format("15:04:05")
// }

// func postsHandler(w http.ResponseWriter, r *http.Request) {
// 	// Put up some random data for demonstration:
// 	data := map[string]interface{}{"posts": []Post{
// 		{User: "Bob", Time: produceTime(), Text: "The weather is nice."},
// 		{User: "Alice", Time: produceTime(), Text: "It's raining."},
// 	}}
// 	var err error
// 	switch r.URL.Path {
// 	case "/posts/":
// 		err = t.Execute(w, data)
// 	case "/posts/more":
// 		err = t.ExecuteTemplate(w, "batch", data)
// 	}
// 	if err != nil {
// 		log.Printf("Template execution error: %v", err)
// 	}
// }

// func main() {
// 	http.HandleFunc("/posts/", postsHandler)
// 	panic(http.ListenAndServe(":8080", nil))
// }

// const page = `<html><body><h2>Posts</h2>
// {{block "batch" .}}
// 	{{range .posts}}
// 		<div><b>{{.Time}} {{.User}}:</b> {{.Text}}</div>
// 	{{end}}
// 	<div id="nextBatch"></div>
// {{end}}
// <button onclick="loadMore()">Load more</button>
// <script>
// 	function loadMore() {
// 		var e = document.getElementById("nextBatch");
// 		var xhr = new XMLHttpRequest();
// 		xhr.onreadystatechange = function() {
// 			if (xhr.readyState == 4 && xhr.status == 200) {
// 				e.outerHTML = xhr.responseText;
// 			}
// 		}
// 		xhr.open("GET", "/posts/more", true);
// 		try { xhr.send(); } catch (err) { /* handle error */ }
// 	}
// </script>
// </body></html>`
