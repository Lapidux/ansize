# Ansize - Lapidux Fork

Original Repository (upstream): https://github.com/jhchen/ansize

Converts an image to (optionally non-binary) ANSI art like so:

## Example:
### Original Image:
![putty.png](https://raw.githubusercontent.com/Lapidux/ansize/master/examples/putty.png)
### Binary ANSI:
![putty-binary.png](https://raw.githubusercontent.com/Lapidux/ansize/master/examples/putty-output-bin.png)
### Non-binary ANSI:
![pikachu-original-example.png](https://raw.githubusercontent.com/Lapidux/ansize/master/examples/putty-output.png)

Check out the `examples/` folder for some image samples and their corresponding output. Ex.

    cat examples/pikachu.ansi

The original maintainer optimized for images with dark backgrounds and used 0's and 1's for the character set. Lapidux made minor changes to support loading text files as a source of characters instead of using random binary. This change was inspired by an old website called commits.io, now defunct.

### Basic technical explanation:
1. Shrink image to desired size (From upstream)
2. For each pixel, find the corresponding color in ANSI's limited color palette (From upstream)
3. Set the foreground to that color (From upstream)
4. Print a random 0 or 1 if no characters are provided, otherwise sequentially use the provided characters

### Note:
This project currently uses a file extension whitelist when provided with a directory, to stop it parsing binary files. It's likely that a number of file extensions are not on the whitelist, so will not be parsed. Feel free to add them yourself (and submit a pull request to help us all out!). The whitelist is currently:
    `"java", "txt", "go", "py", "asm", "aspx", "bat", "htm", "html", "inc", "js", "jsp", "php", "src", "r", "cpp", "c"`

## Installation

    go get github.com/Lapidux/ansize

## Usage

    ansize [-f <characters file> OR -d <directory containing files>] [-w <width of output, default is currently 100>] <input image> <output ANSI file>

## Development

1. Install go
2. Set your $GOPATH
3. Install github.com/nfnt/resize
4. Clone ansize
5. Build ansize

On a Mac with Homebrew the commands are

    brew install go
    mkdir /usr/local/lib/go
    export GOPATH=/usr/local/lib/go
    go get github.com/nfnt/resize
    git clone git@github.com:lapidux/ansize.git
    go build ansize.go
