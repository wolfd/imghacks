package main

import (
	"image"
	"image/png"
	"log"
	"os"
	"sort"
)

// A data structure to hold a key/value pair.
type Pair struct {
	Key   int
	Value uint32
}

// A slice of Pairs that implements sort.Interface to sort by Value.
type PairList []Pair

func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

// A function to turn a map into a PairList, then sort and return it.
func sortMapByValue(r map[int]uint32) PairList {
	p := make(PairList, len(r))
	i := 0
	for k, v := range r {
		p[i] = Pair{k, v}
		i++
	}
	sort.Stable(p)

	return p
}

func main() {
	reader, err := os.Open("data/input.png")
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	m, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	bounds := m.Bounds()

	rows := make(map[int]uint32)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := m.At(x, y).RGBA()
			// Sum all channels plus weight the current y value
			rows[y] += r + g + b + a + (uint32(y) * 100)
		}
	}

	// Sort the lines based on the values assigned above
	sorted := sortMapByValue(rows)

	n := image.NewRGBA(image.Rect(0, 0, bounds.Max.X, bounds.Max.Y))
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			y_swap := sorted[y].Key
			p := m.At(x, y_swap)
			// Actually move the pixels
			n.Set(x, y, p)
		}
	}

	writer, err := os.Create("data/output.png")
	if err != nil {
		log.Fatal(err)
	}
	defer writer.Close()

	png.Encode(writer, n)
}
