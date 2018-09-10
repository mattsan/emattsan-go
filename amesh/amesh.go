package amesh

import (
    "fmt"
    "image"
    "image/draw"
    "image/jpeg"
    "image/png"
    "image/gif"
    "net/http"
    "os"

    "emattsan-go/amesh/constants"
    "emattsan-go/amesh/index"
)

func downloadJpeg(path string) (image.Image, error) {
    resp, err := http.Get(path)
    if err != nil { return nil, err }
    return jpeg.Decode(resp.Body)
}

func downloadPng(path string) (image.Image, error) {
    resp, err := http.Get(path)
    if err != nil { return nil, err }
    return png.Decode(resp.Body)
}

func downloadGif(path string) (image.Image, error) {
    resp, err := http.Get(path)
    if err != nil { return nil, err }
    return gif.Decode(resp.Body)
}

func loadJpegFromFile(filename string) (image.Image, error) {
    file, err := os.Open(filename)
    if err != nil { return nil, err }
    defer file.Close()

    return jpeg.Decode(file)
}

func loadPngFromFile(filename string) (image.Image, error) {
    file, err := os.Open(filename)
    if err != nil { return nil, err }
    defer file.Close()

    return png.Decode(file)
}

func loadGifFromFile(filename string) (image.Image, error) {
    file, err := os.Open(filename)
    if err != nil { return nil, err }
    defer file.Close()

    return gif.Decode(file)
}

func saveToJpegFile(filename string, image image.Image) error {
    file, err := os.Create(filename)
    if err != nil { return err }
    defer file.Close()

    return jpeg.Encode(file, image, &jpeg.Options{Quality: 100})
}

func composeImages(topographyImage, boundaryImage, radarImage image.Image) image.Image {
    topographyRect := image.Rectangle{image.Point{0, 0}, topographyImage.Bounds().Size()}
    boundaryRect := image.Rectangle{image.Point{0, 0}, boundaryImage.Bounds().Size()}
    radarRect := image.Rectangle{image.Point{0, 0}, radarImage.Bounds().Size()}

    resultImage := image.NewRGBA(topographyRect)

    draw.Draw(resultImage, topographyRect, topographyImage, image.Point{0, 0}, draw.Src)
    draw.Draw(resultImage, radarRect, radarImage, image.Point{0, 0}, draw.Over)
    draw.Draw(resultImage, boundaryRect, boundaryImage, image.Point{0, 0}, draw.Over)

    return resultImage
}

func composite(topography, boundary, radar string) (image.Image, error) {
    topographyImage, err := downloadJpeg(topography)
    if err != nil { return nil, err }
    boundaryImage, err := downloadPng(boundary)
    if err != nil { return nil, err }
    radarImage, err := downloadGif(radar)
    if err != nil { return nil, err }

    return composeImages(topographyImage, boundaryImage, radarImage), nil
}

func LatestImage()  (image.Image, error) {
    lastIndex, _ := index.LatestIndex()
    return composite(constants.Topography, constants.Boundary, fmt.Sprintf(constants.IMAGE_URL_FORMAT, lastIndex))
}
