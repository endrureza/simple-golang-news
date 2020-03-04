package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux" // For Router
	"github.com/jinzhu/gorm" // For ORM and Database Connection
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/unrolled/render"
)

// News is a struct model 
type News struct {
	gorm.Model
	Author string
	Body string `gorm:"type:text"`
}

var r = render.New()
var db *gorm.DB
var router = mux.NewRouter()

func init() {
	var err error
	db, err = gorm.Open("mysql", "root:root@(localhost)/kumparan?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		panic(err)
	}

	// Migrate the schema
	db.AutoMigrate(&News{})
}

// Index is main route handler
func Index(res http.ResponseWriter, req *http.Request) {
	r.JSON(res, http.StatusOK, map[string]interface{}{
		"status": "200",
		"message": "v1.0",
	})
}
 
// DisplayNews is to display latest news
func DisplayNews(res http.ResponseWriter, req *http.Request) {
	news := &News{}
	r.JSON(res, http.StatusOK, map[string]interface{}{
		"status": "200",
		"message": "success",
		"data": db.Find(news),
	})
}

// StoreNews is to add new record of news
func StoreNews(res http.ResponseWriter, req *http.Request) {
	reqBody, _ := ioutil.ReadAll(req.Body)
	var news News
	json.Unmarshal(reqBody, &news)

	db.Create(&news)

	r.JSON(res, http.StatusOK, map[string]interface{}{
		"status": "200",
		"message": "success",
		"data": &news,
	})
}

// HandleRequests is a function to handle all route request
func HandleRequests() {
	router.HandleFunc("/", Index).Methods("GET")
	router.HandleFunc("/news", DisplayNews).Methods("GET")
	router.HandleFunc("/news", StoreNews).Methods("POST")
	log.Fatal(http.ListenAndServe(":5000",router))
}

func main() {
	fmt.Println("Application is running...")

	HandleRequests()
}