package idobata

import (
    "bufio"
    "bytes"
    "mime/multipart"
    "net/http"
    "net/url"
    "os"
)

type Hook struct {
    Url string
}

func NewHook(url string) *Hook {
    return &Hook{Url: url}
}

func (hook *Hook) PostMessage(msg string) (*http.Response, error) {
    body := url.Values{"source": []string{msg}}
    return http.PostForm(hook.Url, body)
}

func (hook *Hook) PostHtml(html string) (*http.Response, error) {
    body := url.Values{
      "source": []string{html},
      "format": []string{"html"},
    }
    return http.PostForm(hook.Url, body)
}

func createBodyFromFile(filename string, buffer *bytes.Buffer) (string, error) {
    multipartWriter := multipart.NewWriter(buffer)
    defer multipartWriter.Close()

    writer, err := multipartWriter.CreateFormFile("image", filename)
    if err != nil { return "", err }

    source, err := os.Open(filename)
    if err != nil { return "", err }
    defer source.Close()

    bufferedWriter := bufio.NewWriter(writer)
    _, err = bufferedWriter.ReadFrom(source)
    if err != nil { return "", err }

    err = bufferedWriter.Flush()
    if err != nil { return "", err }

    return multipartWriter.FormDataContentType(), nil
}

func (hook *Hook) PostImageFile(filename string) (*http.Response, error) {
    buffer := new(bytes.Buffer)

    contentType, err := createBodyFromFile(filename, buffer)
    if err != nil { return nil, err }

    return http.Post(hook.Url, contentType, buffer)
}
