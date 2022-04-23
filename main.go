package main

import (
	"bufio"
	"bytes"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"os"
)

// мой вариант хендлера на основе echo, а не пакета http
// предложенный в учебнике вариант через пакет template не работает должным образом
// содержимое файла в формате html переносим в байтовый срез, который используем как аргумент для метода HTML
func newHandler(c echo.Context) error {
	f, err := os.Open("view.html")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	wr := bytes.Buffer{} //
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		wr.WriteString(sc.Text())
	}
	return c.HTML(http.StatusOK, wr.String())
}

// на удивление рабочий вариант на простых элементах, без использования пакета bytes
func testHandler(c echo.Context) error {
	f, err := os.Open("new.html")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	var wr string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		wr += sc.Text()
	}
	return c.HTML(http.StatusOK, wr)
}

// делалось ради проверки

func main() {
	e := echo.New()
	e.GET("/guestbook/new", testHandler)
	e.GET("/guestbook", newHandler)
	e.Logger.Fatal(e.Start(":8080"))

}

// вои и все
