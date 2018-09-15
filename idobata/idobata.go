package idobata

import (
    "bytes"
    "io"
    "mime/multipart"
    "net/http"
    "os"
)

type Form interface {
    TakePartIn(multipart *multipart.Writer) error
}

type Source struct {
    Value string
}

type Format struct {
    Value string
}

type ImageFile struct {
    Filename string
}

type Image struct {
    Reader io.Reader
    Filename string
}

func (source *Source) TakePartIn(multipartWriter *multipart.Writer) error {
    return multipartWriter.WriteField("source", source.Value)
}

func (format *Format) TakePartIn(multipartWriter *multipart.Writer) error {
    return multipartWriter.WriteField("format", format.Value)
}

func (imageFile *ImageFile) TakePartIn(multipartWriter *multipart.Writer) error {
    form, err := multipartWriter.CreateFormFile("image", imageFile.Filename)
    if err != nil { return err }

    file, err := os.Open(imageFile.Filename)
    if err != nil { return err }
    defer file.Close()

    _, err = io.Copy(form, file)
    return err
}

func (image *Image) TakePartIn(multipartWriter *multipart.Writer) error {
    form, err := multipartWriter.CreateFormFile("image", image.Filename)
    if err != nil { return err }

    _, err = io.Copy(form, image.Reader)
    return err
}

type Hook struct {
    Url string
}

func NewHook(url string) *Hook {
    return &Hook{Url: url}
}

func (hook *Hook) Post(forms ...Form) (*http.Response, error) {
    body := new(bytes.Buffer)

    multipartWriter := multipart.NewWriter(body)
    defer multipartWriter.Close()

    for _, form := range forms {
        if err := form.TakePartIn(multipartWriter); err != nil {
            return nil, err
        }
    }

    multipartWriter.Close()

    contentType := multipartWriter.FormDataContentType()
    return http.Post(hook.Url, contentType, body)
}

func (hook *Hook) PostText(text string) (*http.Response, error) {
    return hook.Post(&Source{Value: text})
}

func (hook *Hook) PostHtml(html string) (*http.Response, error) {
    return hook.Post(&Source{Value: html}, &Format{Value: "html"})
}

func (hook *Hook) PostMarkdown(md string) (*http.Response, error) {
    return hook.Post(&Source{Value: md}, &Format{Value: "markdown"})
}

func (hook *Hook) PostImageFile(filename string) (*http.Response, error) {
    return hook.Post(&ImageFile{Filename: filename})
}

func (hook* Hook) PostImage(reader io.Reader, filename string) (*http.Response, error) {
    return hook.Post(&Image{Reader: reader, Filename: filename})
}
