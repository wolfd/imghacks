package main

import (
	"fmt"
	"image"
	"image/color"
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
	// Input file
	input := fmt.Sprintf("data/%v.png", os.Args[1])
	// Hack mode
	mode := os.Args[2]

	reader, err := os.Open(input)
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

	fmt.Printf("mode: %v\n", mode)

	n := image.NewRGBA(image.Rect(0, 0, bounds.Max.X, bounds.Max.Y))

	switch mode {
	case "sortboxwaves":
		fmt.Println("sortbox and screw up")
		boxWidth := bounds.Max.Y / 16
		for box := 0; box < boxWidth; box++ {
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				for x := bounds.Min.X + (box * boxWidth); x < bounds.Min.X+((box+1)*boxWidth) || x < bounds.Max.X; x++ {
					r, g, b, a := m.At(x, y).RGBA()
					// Sum all channels plus weight the current y value
					rows[y] += r + g + b + a + (uint32(y)*10 + uint32(x)*25)
				}
			}

			// Sort the lines based on the values assigned above
			sorted := sortMapByValue(rows)

			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				for x := bounds.Min.X + (box * boxWidth); x < bounds.Min.X+((box+1)*boxWidth) || x < bounds.Max.X; x++ {
					y_swap := sorted[y].Key
					p := m.At((x+box)%bounds.Max.X, y_swap)
					// Actually move the pixels
					r, g, b, a := p.RGBA()

					colors := [3]uint8{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)}
					a_ := uint8(a >> 8)
					c := box
					q := color.NRGBA{
						R: colors[c%3],
						G: colors[(c+1)%3],
						B: colors[(c+2)%3],
						A: a_,
					}

					n.Set(x, y, q)
				}
			}
		}
	case "sortbox":
		fmt.Println("sortbox and screw up")
		boxWidth := bounds.Max.Y / 16
		for box := 0; box < boxWidth; box++ {
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				for x := bounds.Min.X + (box * boxWidth); x < bounds.Min.X+((box+1)*boxWidth) || x < bounds.Max.X; x++ {
					r, g, b, a := m.At(x, y).RGBA()
					// Sum all channels plus weight the current y value
					rows[y] += r + g + b + a + (uint32(y) * 100)
				}
			}

			// Sort the lines based on the values assigned above
			sorted := sortMapByValue(rows)

			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				for x := bounds.Min.X + (box * boxWidth); x < bounds.Min.X+((box+1)*boxWidth) || x < bounds.Max.X; x++ {
					y_swap := sorted[y].Key
					p := m.At(x, y_swap)
					// Actually move the pixels
					n.Set(x, y, p)
				}
			}
		}
	case "waves":
		fmt.Println("waves and screw up")
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, a := m.At(x, y).RGBA()
				// Sum all channels plus weight the current y value
				rows[y] += r + g + b + a + (uint32(y) * 100)
			}
		}

		// Sort the lines based on the values assigned above
		sorted := sortMapByValue(rows)

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				y_swap := sorted[y].Key
				p := m.At((x+(y*y)/(x+1))%bounds.Max.X, y_swap)
				// Actually move the pixels
				n.Set(x, y, p)
			}
		}
	case "sort":
		fmt.Println("sort and screw up")
		fallthrough
	default:
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, a := m.At(x, y).RGBA()
				// Sum all channels plus weight the current y value
				rows[y] += r + g + b + a + (uint32(y) * 100)
			}
		}

		// Sort the lines based on the values assigned above
		sorted := sortMapByValue(rows)

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				y_swap := sorted[y].Key
				p := m.At(x, y_swap)
				// Actually move the pixels
				n.Set(x, y, p)
			}
		}
	}

	writer, err := os.Create("data/output.png")
	if err != nil {
		log.Fatal(err)
	}
	defer writer.Close()

	png.Encode(writer, n)
}
