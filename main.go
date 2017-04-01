package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"

	"github.com/pwed/disttube/ffmpeg"
)

//Compile templates on start
var templates = template.Must(template.ParseFiles("tmpl/upload.html"))

//Display the named template
func display(w http.ResponseWriter, tmpl string, data interface{}) {
	templates.ExecuteTemplate(w, tmpl+".html", data)
}

//A struct containing the configuration parameters of the running program
type config struct {
	Port     string `json:"port"`
	VideoDir string `json:"video_dir"`
	TempDir  string `json:"temp_dir"`
}
//Default config values
var conf = config{
	":8080",
	"videos/",
	"temp/"}

//Set up the environment before running the application
func init() {
	configFile, err := os.Open("config.json")
	defer configFile.Close()
	if err != nil {

	}
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&conf); err != nil {
		fmt.Printf("parsing config file: %v", err.Error())
	}
	os.Mkdir(conf.VideoDir, 0777)
	os.Mkdir(conf.TempDir, 0777)
}

//main :)
func main() {
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

	mx.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	http.ListenAndServe(conf.Port, mx)
}

//Handles multi file video uploads and passes them on to the video transcoder
func uploadHandler(w http.ResponseWriter, r *http.Request) {
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

		//get the *fileHeaders
		files := m.File["myfiles"]
		for i, _ := range files {
			//for each fileHeader, get a handle to the actual file
			file, err := files[i].Open()
			defer file.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			//create destination file making sure the path is writable.
			tempFile := conf.TempDir + files[i].Filename
			filename := files[i].Filename
			dst, err := os.Create(tempFile)
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
			perminentFile := conf.VideoDir + strings.TrimSuffix(filename, filepath.Ext(filename))

			go ffmpeg.SequentialBatchIngest("ffmpeg", tempFile, perminentFile)

			fmt.Println(perminentFile)
		}
		//display success message.
		display(w, "upload", "Upload successful.")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
