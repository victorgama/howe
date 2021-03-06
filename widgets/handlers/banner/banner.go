package banner

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/arsham/rainbow/rainbow"
	"github.com/fatih/color"
	"github.com/lukesampson/figlet/figletlib"
	"github.com/victorgama/howe/helpers"
	"github.com/victorgama/howe/widgets"
)

const font = `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`

const defaultFigletsDir = "/usr/share/howe"

var availableColors = map[string]color.Attribute{
	"red":     color.FgRed,
	"green":   color.FgGreen,
	"yellow":  color.FgYellow,
	"blue":    color.FgBlue,
	"magenta": color.FgMagenta,
	"cyan":    color.FgCyan,
	"white":   color.FgWhite,
	// Value of this next key isn't really used, but the key must be present.
	"rainbow": color.FgWhite,
}

func loadDefaultFiglet() ([]byte, error) {
	fontBytes, err := base64.StdEncoding.DecodeString(font)
	if err != nil {
		return nil, err
	}

	zr, err := gzip.NewReader(bytes.NewBuffer(fontBytes))
	if err != nil {
		return nil, err
	}

	outBuf := bytes.NewBuffer([]byte{})
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(outBuf, zr)
	if err != nil {
		return nil, err
	}

	return outBuf.Bytes(), nil
}

func loadFiglet(from string) (*figletlib.Font, error) {
	var buf []byte
	var err error

	if from == "" {
		buf, err = loadDefaultFiglet()
		if err != nil {
			return nil, err
		}
	} else {
		if filepath.Ext(from) == "" {
			from = from + ".flf"
		}

		if !filepath.IsAbs(from) {
			from = filepath.Join(defaultFigletsDir, from)
		}

		dat, err := ioutil.ReadFile(from)
		if err != nil {
			if os.IsNotExist(err) {
				helpers.ReportError(fmt.Sprintf("banner: Falling back to default font due to error: %s", err))
				return loadFiglet("")
			}
			return nil, err
		}
		buf = dat
	}
	return figletlib.ReadFontFromBytes(buf)
}

func colorizeOutput(value, colorName string) string {
	if colorName == "rainbow" {
		buffer := []byte{}
		output := bytes.NewBuffer(buffer)
		r := rainbow.Light{
			Reader: strings.NewReader(value),
			Writer: output,
		}
		r.Paint()
		return string(output.Bytes())
	}
	foreground := color.New(availableColors[colorName]).SprintFunc()
	return foreground(value)
}

func handle(payload map[string]interface{}, output chan interface{}, wait *sync.WaitGroup) {
	toWrite, err := helpers.TextOrCommand("banner", payload)
	if err != nil {
		output <- err
		wait.Done()
		return
	}

	fontNameOrPath := ""
	if font, ok := payload["font"]; ok {
		if strFont, ok := font.(string); ok {
			fontNameOrPath = strFont
		} else {
			output <- fmt.Errorf("font property must be a string")
			wait.Done()
			return
		}
	}

	fontColor := "magenta"
	if color, ok := payload["color"]; ok {
		if strColor, ok := color.(string); ok {
			fontColor = strings.ToLower(strColor)
		} else {
			output <- fmt.Errorf("color property must be a string")
			wait.Done()
			return
		}
	}

	_, valid := availableColors[fontColor]

	if !valid {
		colors := make([]string, len(availableColors))
		for k := range availableColors {
			colors = append(colors, k)
		}
		output <- fmt.Errorf("invalid color; valid values are " + strings.Join(colors, ", "))
		wait.Done()
		return
	}

	font, err := loadFiglet(fontNameOrPath)
	if err != nil {
		output <- err
		wait.Done()
		return
	}

	output <- colorizeOutput(strings.TrimSuffix(figletlib.SprintMsg(toWrite, font, 80, font.Settings(), "left"), "\n"), fontColor)
	wait.Done()
	return
}

func init() {
	widgets.Register("banner", handle)
}
