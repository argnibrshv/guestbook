package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

// Guestbook - структура, используемая при отображении шаблона view.html.
type Guestbook struct {
	SignaturesCount int
	Signatures      []string
}

// viewHandler читает записи гостевой книги и выводит их вместе со счетчиком записей.
func viewHandler(writer http.ResponseWriter, request *http.Request) {
	signatures := getString("signatures.txt")
	guestbook := Guestbook{
		SignaturesCount: len(signatures),
		Signatures:      signatures,
	}
	html, err := template.ParseFiles("view.html")
	checkError(err)
	err = html.Execute(writer, guestbook)
	checkError(err)
}

// newHandler отображает форму для ввода записи.
func newHandler(writer http.ResponseWriter, request *http.Request) {
	html, err := template.ParseFiles("new.html")
	checkError(err)
	err = html.Execute(writer, nil)
	checkError(err)
}

// createHandler получает запрос POST с новой записью и добавляет ее в файл signatures.
func createHandler(writer http.ResponseWriter, request *http.Request) {
	signature := request.FormValue("signature")
	options := os.O_WRONLY | os.O_APPEND | os.O_CREATE
	file, err := os.OpenFile("signatures.txt", options, os.FileMode(0600))
	checkError(err)
	_, err = fmt.Fprintln(file, signature)
	checkError(err)
	err = file.Close()
	checkError(err)
	http.Redirect(writer, request, "/guestbook", http.StatusFound)
}

func main() {
	http.HandleFunc("/guestbook", viewHandler)
	http.HandleFunc("/guestbook/new", newHandler)
	http.HandleFunc("/guestbook/create", createHandler)
	err := http.ListenAndServe(":8080", nil)
	checkError(err)
}

// getStrings возвращает сегмент строк, прочитанный из fileName, по одной строке на каждую строку файла.
func getString(fileName string) []string {
	file, err := os.Open(fileName)
	if os.IsNotExist(err) {
		return nil
	}
	checkError(err)
	defer file.Close()
	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	checkError(scanner.Err())
	return lines
}

// check вызывает log.Fatal для любых ошибок, отличных от nil.
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
