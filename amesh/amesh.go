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

type Image struct {
    Timestamp time.Time
    Topography image.Image
    Boundary image.Image
    Radar image.Image
}

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

func LatestImage()  (*Image, error) {
    lastIndex, _ := index.LatestIndex()
    timestamp, _ := strToTime(lastIndex)
    radarUrl := fmt.Sprintf(constants.MESH_URL_FORMAT, lastIndex)

    topographyImage, err := downloadImage(jpeg.Decode, constants.TOPOGRAPHY_URL)
    if err != nil { return nil, err }

    boundaryImage, err := downloadImage(png.Decode, constants.BOUNDARY_URL)
    if err != nil { return nil, err }

    radarImage, err := downloadImage(gif.Decode, radarUrl)
    if err != nil { return nil, err }

    image := Image{
      Timestamp: timestamp,
      Topography: topographyImage,
      Boundary: boundaryImage,
      Radar: radarImage,
    }

    return &image, err
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

func (image *Image) Composite() (image.Image, error) {
    return composeImages(image.Topography, image.Boundary, image.Radar), nil
}

func norm(x, y int) int {
    return x * x + y * y
}

func (image *Image) RainingRatio(point image.Point, radius int) int {
    area := 0
    count := 0
    radiusSquared := radius * radius
    rect := image.Radar.Bounds()
    for y := rect.Min.Y; y < rect.Max.Y; y++ {
          for x := rect.Min.X; x < rect.Max.X; x++ {
              _, _, _, a := image.Radar.At(x, y).RGBA()
              if norm(point.X - x, point.Y - y) <= radiusSquared {
                  area += 1
                  if a > 0 {
                      count += 1
                  }
              }
          }
    }
    return count * 100 / area
}
