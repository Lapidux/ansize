package main

import (
	"bufio"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const (
	ANSI_BASIC_BASE  int     = 16
	ANSI_COLOR_SPACE uint32  = 6
	ANSI_FOREGROUND  string  = "38"
	ANSI_RESET       string  = "\x1b[0m"
	DEFAULT_WIDTH    int     = 100
	PROPORTION       float32 = 0.46
	RGBA_COLOR_SPACE uint32  = 1 << 16
)

var BANNED_CHARACTERS = [...]string{" ", "\n", "\r"}

var characterBuffer []byte

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func toAnsiCode(c color.Color) string {
	r, g, b, _ := c.RGBA()
	code := int(ANSI_BASIC_BASE + toAnsiSpace(r)*36 + toAnsiSpace(g)*6 + toAnsiSpace(b))
	if code == ANSI_BASIC_BASE {
		return ANSI_RESET
	}
	return "\033[" + ANSI_FOREGROUND + ";5;" + strconv.Itoa(code) + "m"
}

func toAnsiSpace(val uint32) int {
	return int(float32(ANSI_COLOR_SPACE) * (float32(val) / float32(RGBA_COLOR_SPACE)))
}

func writeAnsiImage(img image.Image, file *os.File, width int) {
	m := resize.Resize(uint(width), uint(float32(width)*PROPORTION), img, resize.Lanczos3)
	var current, previous string
	bounds := m.Bounds()
	var currentCharCounter = 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			for isBadCharacter(string(characterBuffer[currentCharCounter])) {
				currentCharCounter++
			}
			if currentCharCounter == len(characterBuffer) {
				currentCharCounter = 0
			}
			current = toAnsiCode(m.At(x, y))
			if current != previous {
				fmt.Print(current)
				file.WriteString(current)
			}
			if ANSI_RESET != current {
				char := string(characterBuffer[currentCharCounter])
				fmt.Print(char)
				file.WriteString(char)
			} else {
				fmt.Print(string(characterBuffer[currentCharCounter]))
				file.WriteString(string(characterBuffer[currentCharCounter]))
			}
			currentCharCounter++
		}
		fmt.Print("\n")
		file.WriteString("\n")
	}
	fmt.Print(ANSI_RESET)
	file.WriteString(ANSI_RESET)
}

func isBadCharacter(char string) bool {
	for _, banned := range BANNED_CHARACTERS {
		if banned == char {
			return true
		}
	}

	return false
}

func readCharactersFromFile(filename string) []byte {
	dat, err := ioutil.ReadFile(filename)
	check(err)
	return dat
}

func main() {

	characterBuffer = readCharactersFromFile("examples/test.txt")

	rand.Seed(time.Now().UTC().UnixNano())
	if len(os.Args) < 3 {
		fmt.Println("Usage: ansize <image> <output> [width]")
		return
	}
	imageName, outputName := os.Args[1], os.Args[2]
	var width int = DEFAULT_WIDTH
	if len(os.Args) >= 4 {
		var err error
		width, err = strconv.Atoi(os.Args[3])
		if err != nil {
			fmt.Println("Invalid width " + os.Args[3] + ". Please enter an integer.")
			return
		}
	}
	imageFile, err := os.Open(imageName)
	if err != nil {
		fmt.Println("Could not open image " + imageName)
		return
	}
	outFile, err := os.Create(outputName)
	if err != nil {
		fmt.Println("Could not open " + outputName + " for writing")
		return
	}
	defer imageFile.Close()
	defer outFile.Close()
	imageReader := bufio.NewReader(imageFile)
	img, _, err := image.Decode(imageReader)
	if err != nil {
		fmt.Println("Could not decode image")
		return
	}

	writeAnsiImage(img, outFile, width)
}
