package main

import (
	"fmt"
	"github.com/d-tsuji/clipboard"
	hook "github.com/robotn/gohook"
	"log"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var regExp = regexp.MustCompile(`\s+`)

const wbURL = `https://www.wildberries.ru/catalog/%d/detail.aspx`

type Data struct {
	SCU map[uint64]struct{}
}

func ParseText(text string) (Data, error) {
	cleanText := regExp.ReplaceAllString(text, " ")

	nums := strings.Fields(cleanText)

	scu := make(map[uint64]struct{})

	for _, numStr := range nums {
		num, err := strconv.ParseUint(numStr, 10, 64)
		if err != nil {
			return Data{}, fmt.Errorf("invalid number: %s", numStr)
		}
		scu[num] = struct{}{}
	}

	return Data{SCU: scu}, nil
}

func main() {
	eventChan := hook.Start()
	defer hook.End()

	keyPressedChan := make(chan bool)
	go checkHotKeyPress(eventChan, keyPressedChan)

	for pressed := range keyPressedChan {
		if pressed {
			text := takeTextFromBuffer()
			data, err := ParseText(text)
			if err != nil {
				log.Println(err)
			}

			for scu, _ := range data.SCU {
				go openSCUInNewTab(scu)
			}
		}
	}

}

func openSCUInNewTab(scu uint64) {
	// Открытие URL в браузере
	err := openBrowser(fmt.Sprintf(wbURL, scu))
	if err != nil {
		fmt.Println("Ошибка при открытии браузера:", err)
	}
}

func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/C", "start", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "darwin": // для macOS
		cmd = exec.Command("open", url)
	default:
		return fmt.Errorf("неизвестная операционная система")
	}

	return cmd.Start()
}

func checkHotKeyPress(evenetChan <-chan hook.Event, result chan<- bool) {
	for event := range evenetChan {
		if event.Keycode == 58 && event.Kind == hook.KeyUp {
			result <- true
		}
	}
}

func takeTextFromBuffer() string {
	text, err := clipboard.Get()
	if err != nil {
		log.Println("can't take data from buffer", err)
		return ""
	}

	return text
}
