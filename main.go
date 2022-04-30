package main

import (
	"bufio"
	"bytes"
	"github.com/labstack/echo"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

type TemplateRegistry struct {
	templates *template.Template
}

// Implement e.Renderer interface
func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// мой вариант хендлера на основе echo, а не пакета http
// предложенный в учебнике вариант через пакет template не работает должным образом
// содержимое файла в формате html переносим в байтовый срез, который используем как аргумент для метода HTML

func newHandler(c echo.Context) error {
	f, err := os.Open("new.html")
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

/* на удивление рабочий вариант на простых элементах, без использования пакета bytes
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
 делалось ради проверки */
//

func getString(fileName string) []string {
	var lines []string
	file, err := os.Open(fileName)
	if os.IsNotExist(err) {
		return nil
	}
	check(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	check(scanner.Err())
	return lines
}

func templateHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "view.html", map[string]interface{}{
		"name":  "DARKFANTASY",
		"price": "300 bucks",
	})
}

func main() {
	e := echo.New()
	e.Renderer = &TemplateRegistry{
		templates: template.Must(template.ParseGlob("*.html")),
	}
	//e.GET("/guestbook/new", testHandler)
	e.GET("/guestbook", newHandler)
	e.GET("/", templateHandler)
	e.Logger.Fatal(e.Start(":8080"))

}
