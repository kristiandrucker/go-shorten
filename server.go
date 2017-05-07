package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo"
	"html/template"
	"io"
	"math/rand"
	"net/http"
	"time"
	"flag"
	"fmt"
)

type Template struct {
	templates *template.Template
}

type Url struct {
	LongUrl string
	Code string
}

var hostname = flag.String("hostname", "localhost", "Hostname for URLs")

var r *rand.Rand

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	fmt.Println("Going to sleep for 30s. Wait for db to start.")
	time.Sleep(time.Duration(time.Second * 30))
	fmt.Println("Going to start the app.")
	r = rand.New(rand.NewSource(time.Now().UnixNano()))

	db, _ := gorm.Open("mysql", "root@/go-shorten")
	db.AutoMigrate(&Url{})

	e := echo.New()
	t := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	e.Renderer = t
	e.GET("/", index)
	e.POST("/submit", submit)
	e.GET("/u/:urlCode", getUrl)
	e.Start(":1323")
	defer db.Close()
}

func index(c echo.Context) error {
	return c.Render(http.StatusOK, "index", nil)
}

func submit(c echo.Context) error {
	db, _ := gorm.Open("mysql", "root@/go-shorten")
	var url Url
	url.LongUrl = c.FormValue("url")
	url.Code = RandomString(10)
	if db.NewRecord(&url) == true {
		db.Create(&url)
	} else {
		db.First(&url)
	}
	db.Close()
	return c.Render(http.StatusOK, "index", *hostname+":1323/u/"+url.Code)
}

func getUrl(c echo.Context) error {
	db, _ := gorm.Open("mysql", "root@/go-shorten")
	var url Url
	url.Code = c.Param("urlCode")
	db.Find(&url)
	return c.Redirect(http.StatusTemporaryRedirect, url.LongUrl)
}

func RandomString(strlen int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}
	return string(result)
}
