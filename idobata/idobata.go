package idobata

import (
    "bytes"
    "image"
    "image/jpeg"
    "io"
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

func createBody(filename string, buffer *bytes.Buffer, writeImage func(io.Writer) error) (string, error) {
    multipartWriter := multipart.NewWriter(buffer)
    defer multipartWriter.Close()

    writer, err := multipartWriter.CreateFormFile("image", filename)
    if err != nil { return "", err }

    err = writeImage(writer)
    if err != nil { return "", err }

    return multipartWriter.FormDataContentType(), nil
}

func createBodyFromFile(filename string, buffer *bytes.Buffer) (string, error) {
    return createBody(filename, buffer, func(writer io.Writer) error {
        source, err := os.Open(filename)
        if err != nil { return err }
        defer source.Close()

        _, err = io.Copy(writer, source)
        return err
    })
}

func createBodyFromImage(image image.Image, buffer *bytes.Buffer) (string, error) {
    return createBody("image.jpg", buffer, func(writer io.Writer) error {
        return jpeg.Encode(writer, image, &jpeg.Options{Quality: 100})
    })
}

func (hook *Hook) PostImageFile(filename string) (*http.Response, error) {
    buffer := new(bytes.Buffer)

    contentType, err := createBodyFromFile(filename, buffer)
    if err != nil { return nil, err }

    return http.Post(hook.Url, contentType, buffer)
}

func (hook* Hook) PostImage(image image.Image) (*http.Response, error) {
    buffer := new(bytes.Buffer)

    contentType, err := createBodyFromImage(image, buffer)
    if err != nil { return nil, err }

    return http.Post(hook.Url, contentType, buffer)
}