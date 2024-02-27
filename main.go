package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"github.com/gorilla/mux"
	"path/filepath"
	"strconv"
	"encoding/json"
	"io/ioutil"
	"github.com/thedevsaddam/renderer"
)

var tpl *template.Template
var render = renderer.New(renderer.Options{
	ParseGlobPattern: "./templates/*.html",
})
var GALLERY_PATH = "D:/Photo/p"
var DB_PATH = GALLERY_PATH + "/db"
var IMAGE_PER_PAGE = 50
var image_type_list = [2]string{".jpg", ".png"}

type Alblum struct {
    Thumbnail string `json:"thumbnail"`
	Name      string `json:"Name"`
	Path      string `json:"path"`
}

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Page struct {
	page_num string
	is_current string
}
type Gallery map[string]Alblum

var gallery Gallery



func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))

	db_file, err := os.Open(DB_PATH)
    // if we os.Open returns an error then handle it
    if err != nil {
        fmt.Println(err)
    }
	defer db_file.Close()
	byteValue, _ := ioutil.ReadAll(db_file)
	
    json.Unmarshal([]byte(byteValue), &gallery)

	
}

func main() {

	router := mux.NewRouter()
	
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./bootstrap-5.3.2/"))))
	router.HandleFunc("/s/{alblum}/{image}", func(w http.ResponseWriter, r *http.Request) {
		// Figure out what the resolved prefix was.
		p := mux.Vars(r)
		// image := mux.Vars(r)["image"]
		prefix := fmt.Sprintf("/s/%s/", p["alblum"])
		file_dir := fmt.Sprintf("%s", gallery[p["alblum"]].Path)
		// Strip it the normal way.
		fileServer := http.FileServer(http.Dir(file_dir))
		http.StripPrefix(prefix, fileServer).ServeHTTP(w, r)
	})

	router.Handle("/favicon.ico", http.NotFoundHandler())
	router.HandleFunc("/", GalleryPage)
	router.HandleFunc("/a/{alblum_id}", AlblumPage)
	router.HandleFunc("/a/{alblum_id}/{image_id}", ImagePage)
	
	http.ListenAndServe(":8080", router)
}

func generatePageList(data_length int, current_page int) *[]map[string]string {
	total_page := (data_length / IMAGE_PER_PAGE) + 4
	page_list := make([]map[string]string, total_page)
    for i := range page_list {
		is_current := ""
		
		if i == 0 {
			page_list[i] = map[string]string{
				"page_num": "←",
				"is_current": "",
			}
		} else if i == total_page - 2{
			page_list[i] = map[string]string{
				"page_num": "→",
				"is_current": "",
				
			}
		} else if i == total_page - 1{
			page_list[i] = map[string]string{
				"page_num": "⌂",
				"is_current": "",
				
			}
		} else {
			if i + 1 == current_page{
				is_current = "true"
			}
			
			page_list[i] = map[string]string{
				"page_num": strconv.Itoa(i),
				"is_current": is_current,
			}
		} 

    }

	return &page_list

}

func ImagePage(w http.ResponseWriter, req *http.Request){
	vars := mux.Vars(req)
	image_json := map[string]interface{}{
		"alblum_name": gallery[vars["alblum_id"]],
		"alblum_id": vars["alblum_id"],
		"image_id": vars["image_id"],
	}
	tpl.ExecuteTemplate(w, "image_1_page.html", image_json)
}
func GalleryPage(w http.ResponseWriter, req *http.Request ){

	alblum_list := make([]interface{}, len(gallery))

	{
		index := 0
		for id, alblum_item := range gallery {

			alblum_list[index] = map[string]string{
				"alblum_id": id,
				"thumbnail": alblum_item.Thumbnail,
			}
			index++
		} 		
	}
	

	file, err := os.Open(GALLERY_PATH)
 	if err != nil {
 		log.Fatalf("failed opening directory: %s", err)
 	}
 	defer file.Close()

	current_page, err := strconv.Atoi(req.URL.Query().Get("p"))
	if err != nil {
		current_page = 1
	}
	if req.URL.Query().Get("p") == "⌂"{
		http.Redirect(w, req, "/", http.StatusSeeOther)

	} 
	last_image := current_page * IMAGE_PER_PAGE
	
	
	if len(alblum_list) < last_image{
		last_image = len(alblum_list)
	}
	alblum_list = alblum_list[((current_page - 1) * IMAGE_PER_PAGE):last_image]
	
 	gallery_json := map[string]interface{}{
 		"title": "my gallery",
		"alblums": alblum_list,
		"page_list": generatePageList(len(gallery)+1, current_page),
 	}
	tpl.ExecuteTemplate(w, "gallery.html", gallery_json)
	// tpls := []string{"templates/gallery.html", "templates/template.tmpl"}
	// render.Template(w, http.StatusOK, tpls, gallery_json)
 }

 func AlblumPage(w http.ResponseWriter, req *http.Request ){
	vars := mux.Vars(req)
	alblum_data := gallery[vars["alblum_id"]]
	// fmt.Println("opening alblum: " + GALLERY_PATH + "/" + alblum_data.Name)
	file, err := os.Open(alblum_data.Path)
 	if err != nil {
 		log.Fatalf("failed opening alblum: %s", err)
 	}
 	defer file.Close()

	current_page, err := strconv.Atoi(req.URL.Query().Get("p"))
	if err != nil {
		current_page = 1
	}
	if req.URL.Query().Get("p") == "⌂"{
		http.Redirect(w, req, "/", http.StatusSeeOther)

	} 
	// else if req.URL.Query().Get("p") == "←" {
	// 	if current_page != 1 {
	// 		current_page--
	// 	}
		
	// } else if req.URL.Query().Get("p") == "→" {
	// 	current_page++
	// }
	
 	file_list, _ := file.Readdirnames(0)
	image_list := make([]string, len(file_list))
	alblum_list := make([]string, len(file_list))
	image_counter := 0
	alblum_counter := 0

	for _, file_name := range file_list {
		file_info, err := os.Stat(filepath.Join(alblum_data.Path, file_name))
		if err != nil{	
			log.Fatalf("failed opening image: %s", err)
		} else if file_info.IsDir(){
			log.Print("adding path", file_name)
			alblum_list[alblum_counter] = file_name
			alblum_counter++
			continue
		}

		for _, file_type := range image_type_list{
			if strings.Contains(file_name, file_type){
				image_list[image_counter] = file_name
				image_counter++
				break
			} 			
		//	if err != nil {
		//		log.Fatalf("failed opening directory: %s", err)
		//	}

		//	if file_info.IsDir(){
		//		alblum_list = append(alblum_list, file_name)
		//	}
		}
		

	}

	image_list = image_list[:image_counter]
	
	if image_counter != 0 {
		first_image := (current_page - 1) * IMAGE_PER_PAGE
		last_image := current_page * IMAGE_PER_PAGE
		if last_image > image_counter{
			last_image = image_counter
		}
		image_list = image_list[first_image:last_image]
	}


 	gallery_json := map[string]interface{}{
		"title": alblum_data.Name,
 		"alblum_name": alblum_data.Name,
		"alblum_id": vars["alblum_id"],
		"images": image_list,
		"page_list": generatePageList(len(file_list)+1, current_page),
 	}
	tpl.ExecuteTemplate(w, "alblum.html", gallery_json)
 }

// func AlblumPage(w http.ResponseWriter, req *http.Request) {
// 	if req.Method == http.MethodPost {
// 		fmt.Println("post")
// 	}
	
// 	req_url_token := strings.Split(req.URL.Path, "/")
// 	req_url_len := len(req_url_token)
// 	// fmt.Println(req_url_token)
// 	file, err := os.Open(GALLERY_PATH+"/"+req_url_token[2])
// 	if err != nil {
// 		log.Fatalf("failed opening directory: %s", err)
// 	}
// 	defer file.Close()
	
// 	image_list, _ := file.Readdirnames(0) // 0 to read all files and folders
// 	alblum_list := make([]interface{}, len(image_list))
// 	for index :=  0; index <= len(image_list); index++{
// 		alblum := map[string]string{
// 			"alblum_name": req_url_token[2],
// 			"image_name": image_list[index],
// 		}
// 		alblum_list[index] = alblum
		
// 	} 		

// 	// switch req_url_len{
// 	// case 3:
// 	// 	gallery := map[string]interface{}{
// 	// 		"name": "twitter",
// 	// 		"alblum": image_list,
// 	// 	}
// 	// }
// 	if req_url_len == 5 {
// 		tpl.ExecuteTemplate(w, "image_1_page.html", req_url_token[4])
// 		return
// 	} 
// 	tpl.ExecuteTemplate(w, "alblum.html", alblum_list)
	
// } 
