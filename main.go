package main

import (
	"fmt"
	"github.com/d-tsuji/clipboard"
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
	"log"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

var regExp = regexp.MustCompile(`\s+`)

const (
	wbURL = `https://www.wildberries.ru/catalog/%d/detail.aspx`
	logo  = `
██╗    ██╗██████╗     ███████╗ ██████╗██╗   ██╗    
██║    ██║██╔══██╗    ██╔════╝██╔════╝██║   ██║    
██║ █╗ ██║██████╔╝    ███████╗██║     ██║   ██║    
██║███╗██║██╔══██╗    ╚════██║██║     ██║   ██║    
╚███╔███╔╝██████╔╝    ███████║╚██████╗╚██████╔╝    
 ╚══╝╚══╝ ╚═════╝     ╚══════╝ ╚═════╝ ╚═════╝     
 ██████╗ ██████╗ ███████╗███╗   ██╗███████╗██████╗ 
██╔═══██╗██╔══██╗██╔════╝████╗  ██║██╔════╝██╔══██╗
██║   ██║██████╔╝█████╗  ██╔██╗ ██║█████╗  ██████╔╝
██║   ██║██╔═══╝ ██╔══╝  ██║╚██╗██║██╔══╝  ██╔══██╗
╚██████╔╝██║     ███████╗██║ ╚████║███████╗██║  ██║
 ╚═════╝ ╚═╝     ╚══════╝╚═╝  ╚═══╝╚══════╝╚═╝  ╚═╝
`
)

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
			continue
		}
		scu[num] = struct{}{}
	}

	return Data{SCU: scu}, nil
}

func main() {
	fmt.Println(logo)
	fmt.Println(`
HOW USE:
1. SELECT TEXT
2. PRESS CTRL+SHIFT+F
`)

	fmt.Println("HISTORY:")

	for {
		if pressed := hook.AddEvents(robotgo.KeyF, robotgo.Ctrl, robotgo.Shift); pressed {
			text := takeSelectedText()
			data, err := ParseText(text)
			if err != nil {
				log.Println(err)
			}

			if len(data.SCU) > 0 {
				fmt.Print("OPEN:\t")
			}

			var wg sync.WaitGroup
			wg.Add(len(data.SCU))
			for scu, _ := range data.SCU {
				go openSCUInNewTab(&wg, scu)
			}

			wg.Wait()

			fmt.Print("\n")
		}
	}
}

func openSCUInNewTab(wg *sync.WaitGroup, scu uint64) {
	defer wg.Done()
	fmt.Printf("%d\t", scu)
	err := openBrowser(fmt.Sprintf(wbURL, scu))
	if err != nil {
		fmt.Println("can't open browser")
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
		return fmt.Errorf("unknow platform")
	}

	return cmd.Start()
}

func takeSelectedText() string {
	err := robotgo.KeyTap(robotgo.KeyC, robotgo.Ctrl)
	if err != nil {
		log.Println("can't copy selected text")
		return ""
	}

	robotgo.MilliSleep(50)

	text, err := clipboard.Get()
	if err != nil {
		log.Println("can't take data from buffer")
		return ""
	}

	return text
}
