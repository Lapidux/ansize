package main

import (
	"bufio"
	"flag"
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
	"path/filepath"
	"strconv"
	"strings"
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
var FILETYPE_WHITELIST = [...]string{"java", "txt", "go", "py", "asm", "aspx", "bat", "htm", "html", "inc", "js", "jsp", "php", "src", "r", "cpp", "c"}

var characterBuffer []byte

func check(e error) {
	if e != nil {
		panic(e)
	}
}

/*
Thanks to Siong-Ui Te: https://siongui.github.io/2018/03/10/go-set-of-all-elements-in-two-arrays/
*/
func Union(a, b []byte) []byte {
	m := make(map[byte]bool)

	for _, item := range a {
		m[item] = true
	}

	for _, item := range b {
		if _, ok := m[item]; !ok {
			a = append(a, item)
		}
	}
	return a
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

				if currentCharCounter == len(characterBuffer) {
					currentCharCounter = 0
				}

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
				fmt.Print(" ")
				file.WriteString(" ")
			}
			currentCharCounter++

			if currentCharCounter == len(characterBuffer) {
				currentCharCounter = 0
			}
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

	rand.Seed(time.Now().UTC().UnixNano())

	// variables declaration
	var file, directory string
	var width int

	// flags declaration using flag package
	flag.StringVar(&file, "f", "", "Specify file to read characters from.")
	flag.StringVar(&directory, "d", "", "Specify directory containing files to read characters from.")
	flag.IntVar(&width, "w", DEFAULT_WIDTH, "Specify width of output, defaults to "+string(DEFAULT_WIDTH))

	flag.Parse() // after declaring flags we need to call it

	if len(flag.Args()) != 2 {
		fmt.Println("Usage ([] denotes optional):\n ansize [-f <characters file> OR -d <directory containing files>] [-w <width of output, default is " + string(DEFAULT_WIDTH) + ">] <input image> <output ANSI file>")
		return
	}

	if file != "" {

		if directory != "" {
			fmt.Println("You may only specify a file, a directory, or neither. You may not specify both.")
			return
		}

		characterBuffer = readCharactersFromFile(file)
	} else if directory != "" {
		filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("Error walking the path %q: %v\n", path, err)
				return nil
			}

			if !info.IsDir() && strings.Contains(info.Name(), ".") {
				var splitName = strings.Split(info.Name(), ".")
				var fileType = splitName[len(splitName)-1]
				for _, fType := range FILETYPE_WHITELIST {
					if fType == fileType {
						characterBuffer = Union(characterBuffer, readCharactersFromFile(path))
					}
				}

			}

			return nil
		})

	} else {
		characterBuffer = []byte("01")
	}

	imageName, outputName := flag.Args()[0], flag.Args()[1]

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
