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
    "io"
    "strconv"
    "time"

    "github.com/mattsan/emattsan-go/amesh/constants"
    "github.com/mattsan/emattsan-go/amesh/index"
)

type decoder func (io.Reader) (image.Image, error)

func downloadImage(decode decoder, path string) (image.Image, error) {
    resp, err := http.Get(path)
    if err != nil { return nil, err }
    return decode(resp.Body)
}

func loadImageFromFile(decode decoder, filename string) (image.Image, error) {
    file, err := os.Open(filename)
    if err != nil { return nil, err }
    defer file.Close()

    return decode(file)
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
    topographyImage, err := downloadImage(jpeg.Decode, topography)
    if err != nil { return nil, err }
    boundaryImage, err := downloadImage(png.Decode, boundary)
    if err != nil { return nil, err }
    radarImage, err := downloadImage(gif.Decode, radar)
    if err != nil { return nil, err }

    return composeImages(topographyImage, boundaryImage, radarImage), nil
}

func strsToInts(ss ...string) ([]int, error) {
    is := make([]int, len(ss))

    for index, s := range ss {
        i, err := strconv.ParseInt(s, 10, 64)
        if err != nil { return nil, err }
        is[index] = int(i)
    }

    return is, nil
}

func strToTime(s string) (time.Time, error) {
    is, err := strsToInts(s[0:4], s[4:6], s[6:8], s[8:10], s[10:12])
    if err != nil { return time.Time{}, err }

    datetime := time.Date(
        is[0],
        time.Month(is[1]),
        is[2],
        is[3],
        is[4],
        0,
        0,
        time.Local,
    )

    return datetime, nil
}

func LatestImage()  (image.Image, time.Time, error) {
    lastIndex, _ := index.LatestIndex()
    datetime, _ := strToTime(lastIndex)
    meshUrl := fmt.Sprintf(constants.MESH_URL_FORMAT, lastIndex)
    image, err := composite(constants.TOPOGRAPHY_URL, constants.BOUNDARY_URL, meshUrl)
    return image, datetime, err
}
