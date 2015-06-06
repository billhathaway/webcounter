package webcounter

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

const (
	srcImage      = "https://raw.githubusercontent.com/SDGophers/2015-04-challenge/master/images/numbers.png"
	glyphSize     = 100
	srcImageWidth = 3
)

var (
	srcImageChars = []byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', ',', '0', '.'}
)

// Controller provides a web-counter service
type Controller struct {
	counts map[string]int
	rects  map[byte]image.Rectangle
	img    image.Image
	sync.Mutex
}

// numToImage converts the int val to an image
func (c *Controller) numToImage(val int) image.Image {
	sval := strconv.Itoa(val)
	counterImage := image.NewNRGBA(image.Rect(0, 0, (glyphSize*len(sval))-1, glyphSize-1))
	for i := 0; i < len(sval); i++ {
		rect, ok := c.rects[sval[i]]
		if !ok {
			panic("did not find glyph rectange for " + sval[i:i+1])
		}
		// TODO: (billh) seems like there should be a better way than copying each pixel
		for xOff := 0; xOff+rect.Min.X <= rect.Max.X; xOff++ {
			for yOff := 0; yOff+rect.Min.Y <= rect.Max.Y; yOff++ {
				pixel := c.img.At(rect.Min.X+xOff, rect.Min.Y+yOff)
				counterImage.Set(xOff+i*glyphSize, yOff, pixel)
			}
		}
	}
	return counterImage
}

// loadSourceImage retrieves the image that is the source for our glyphs
func loadSourceImage() (image.Image, error) {
	resp, err := http.Get(srcImage)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code fetching src image at %s - %d ", srcImage, resp.StatusCode)
	}
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	return img, nil
}

// New creates a new Controller
func New() (*Controller, error) {
	img, err := loadSourceImage()
	if err != nil {
		return nil, err
	}
	t := &Controller{counts: make(map[string]int), img: img, rects: make(map[byte]image.Rectangle)}
	for i, ch := range srcImageChars {
		x0 := i % srcImageWidth * glyphSize
		x1 := x0 + glyphSize - 1
		row := i / srcImageWidth
		y0 := row * glyphSize
		y1 := y0 + glyphSize - 1
		t.rects[ch] = image.Rect(x0, y0, x1, y1)
	}
	return t, nil
}

// Get retrieves the count and increments it by one
func (c *Controller) Get(id string) int {
	c.Lock()
	count := c.counts[id]
	c.counts[id] = count + 1
	c.Unlock()
	return count
}

// Delete resets the count
func (c *Controller) Delete(id string) {
	c.Lock()
	delete(c.counts, id)
	c.Unlock()
}

// render sends the output back to the client in a format appropriate to the request suffix
// currently supported suffixes are:
// txt - print as plain text
// png - render as PNG image (default)
// jpg/jpeg - render as a JPEG image
func (c *Controller) render(w http.ResponseWriter, count int, suffix string) {
	switch suffix {
	case "txt":
		w.Header().Add("Content-type", "text/plain")
		fmt.Fprintf(w, "%d", count)
	case "jpg", "jpeg":
		w.Header().Add("Content-type", "image/jpeg")
		jpeg.Encode(w, c.numToImage(count), &jpeg.Options{})
	case "gif":
		w.Header().Add("Content-type", "image/gif")
		gif.Encode(w, c.numToImage(count), &gif.Options{})
	default:
		w.Header().Add("Content-type", "image/png")
		png.Encode(w, c.numToImage(count))
	}
}

func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := strings.Trim(r.URL.Path, "/")
	if id == "" || id == "favicon.ico" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var suffix string
	idx := strings.Index(id, ".")

	if idx >= 0 {
		suffix = id[idx+1:]
		id = strings.TrimSuffix(id, "."+suffix)
	}
	switch r.Method {
	case "GET":
		count := c.Get(id)
		c.render(w, count, suffix)
	case "DELETE":
		c.Delete(id)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}
