package index

import (
    "regexp"
    "net/http"
    "bytes"

    "emattsan-go/amesh/constants"
)

func scanDigits(src []byte) [][]byte {
    matcher := regexp.MustCompile(`[\d]+`)
    return matcher.FindAll(src, 100)
}

func fetchIndices(url string) ([][]byte, error) {
    resp, err := http.Get(url)
    if err != nil { return nil, err }
    defer resp.Body.Close()

    buffer := new(bytes.Buffer)
    _, err = buffer.ReadFrom(resp.Body)
    if err != nil { return nil, err }

    return scanDigits(buffer.Bytes()), nil
}

func LatestIndex() (string, error) {
    indices, err := fetchIndices(constants.INDICES_URL)
    if err != nil { return "", err }

    maxIndex := []byte{}

    for _, index := range indices {
        if bytes.Compare(maxIndex, index) < 0 {
            maxIndex = index
        }
    }

    return string(maxIndex), nil
}
