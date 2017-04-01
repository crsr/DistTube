package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pwed/disttube/ffmpeg"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

//Compile templates on start
var templates = template.Must(template.ParseFiles("tmpl/upload.html"))

//Display the named template
func display(w http.ResponseWriter, tmpl string, data interface{}) {
	templates.ExecuteTemplate(w, tmpl+".html", data)
}

const (
	temp   = "temp/"
	videos = "videos/"
)

func init() {

	os.Mkdir(temp, 0777)
	os.Mkdir(videos, 0777)
}

func main() {

	fmt.Println("test")
	mx := mux.NewRouter()

	mx.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("hello"))
		})

	mx.HandleFunc("/channel/{channel}",
		func(w http.ResponseWriter, r *http.Request) {
			v := mux.Vars(r)
			channel := v["channel"]

			w.Write([]byte("the channel is " + channel))
		})
	mx.HandleFunc("/channel/{channel}/video/{video}",
		func(w http.ResponseWriter, r *http.Request) {
			v := mux.Vars(r)
			channel := v["channel"]
			video := v["video"]

			w.Write([]byte("the channel is " + channel))
			w.Write([]byte("\nthe video is " + video))
		})

	mx.HandleFunc("/upload", uploadHandler)

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	fmt.Println("test")
	http.ListenAndServe(":8080", mx)
}

//This is where the action happens.
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("test")
	switch r.Method {
	//GET displays the upload form.
	case "GET":
		display(w, "upload", nil)

	//POST takes the uploaded file(s) and saves it to disk.
	case "POST":
		//parse the multipart form in the request
		err := r.ParseMultipartForm(100000)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//get a ref to the parsed multipart form
		m := r.MultipartForm

		//get the *fileheaders
		files := m.File["myfiles"]
		for i := range files {
			//for each fileheader, get a handle to the actual file
			file, err := files[i].Open()
			defer file.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			//create destination file making sure the path is writeable.
			tempfile := temp + files[i].Filename
			filename := files[i].Filename
			dst, err := os.Create(tempfile)
			defer dst.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			//copy the uploaded file to the destination file
			if _, err := io.Copy(dst, file); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			perminentFile := videos + strings.TrimSuffix(filename, filepath.Ext(filename))

			go ffmpeg.SequentialBatchIngest("ffmpeg", tempfile, perminentFile)

			fmt.Println(perminentFile)
		}
		//display success message.
		display(w, "upload", "Upload successful.")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
