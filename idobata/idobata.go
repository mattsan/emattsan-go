package idobata

import (
    "bytes"
    "io/ioutil"
    "mime/multipart"
    "net/http"
    "net/url"
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

func (hook *Hook) PostImageFile(filename string) (*http.Response, error) {
    image, err := ioutil.ReadFile(filename)

    if err != nil { return nil, err }

    buffer := new(bytes.Buffer)

    multipartWriter := multipart.NewWriter(buffer)
    writer, err := multipartWriter.CreateFormFile("image", filename)

    if err != nil { return nil, err }

    writer.Write(image)

    multipartWriter.Close()

    return http.Post(hook.Url, multipartWriter.FormDataContentType(), buffer)
}
