package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/getlantern/systray"
)

var _gitRepositories = make([]string, 0)
var _gitMenus = make([](*systray.MenuItem), 0)
var _executableDir string

type gitStatus struct {
	unpushed int
	uncommit int
}

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	getExeDir()
	systray.SetIcon(getIcon(path.Join(_executableDir, "assets/icon.png")))
	getGitRepositories()
	setupMenus()

	go func() {
		for {
			status := getGitStatus()
			if status.uncommit+status.unpushed < 1 {
				systray.SetTitle("")
			} else {
				systray.SetTitle(strconv.Itoa(status.uncommit) + "|" + strconv.Itoa(status.unpushed))
			}
			systray.SetTooltip(strconv.Itoa(status.uncommit) + "change(s) not committed\n" + strconv.Itoa(status.unpushed) + "change(s) not pushed")

			// update every 20 sec(s)
			time.Sleep(20 * time.Second)
		}
	}()
}

func onExit() {
	// Cleaning stuff here.
}

// getExeDir get directory of the executable file
func getExeDir() {
	exePath, err := os.Executable()
	if err != nil {
		onError(err)
	}
	_executableDir = path.Dir(exePath)
}

// getGitRepositories
// read git repositories in config.txt, put then into _gitRepositories
// the content should be directories of the git repositories,
// each line should be one git repository, like:
// ▶️ lines begin with double slash with be treated as comments
// ▶️ /Users/username/git/myproject1
// ▶️ /Users/username/git/myproject2
// ▶️ (empty line at end of of the file)
func getGitRepositories() {
	// open file
	file, err := os.Open(path.Join(_executableDir, "config.txt"))
	if err != nil {
		onError(err)
	}
	defer file.Close()

	// read contents
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// ignore empty lines and lines start with "//"
		if strings.Trim(line, "\t ") == "" || strings.HasPrefix(line, "//") {
			continue
		}
		_gitRepositories = append(_gitRepositories, line)
	}

	if err := scanner.Err(); err != nil {
		onError(err)
	}
}

func setupMenus() {
	for _, line := range _gitRepositories {
		_gitMenus = append(_gitMenus, systray.AddMenuItem(path.Base(line), line))
	}

	systray.AddSeparator()
	menuQuit := systray.AddMenuItem("Quit", "Quits")

	go func() {
		for {
			select {
			case <-menuQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func getGitStatus() gitStatus {
	var total gitStatus

	for index, fullPath := range _gitRepositories {
		var curStatus gitStatus

		menu := _gitMenus[index]
		if menu == nil {
			continue
		}

		// git cherry -v
		gitCmd := exec.Command("git", "cherry", "-v")
		gitCmd.Dir = fullPath
		cmdOutput, err := gitCmd.Output()
		if err != nil {
			menu.SetTitle("[e] " + fullPath)
			continue
		}
		cmdOutputStr := string(cmdOutput[:])
		curStatus.unpushed = len(strings.Split(cmdOutputStr, "\n")) - 1
		total.unpushed += curStatus.unpushed

		// git status -s
		gitCmd = exec.Command("git", "status", "-s")
		gitCmd.Dir = fullPath
		cmdOutput, err = gitCmd.Output()
		if err != nil {
			menu.SetTitle("[e] " + fullPath)
			continue
		}
		cmdOutputStr = string(cmdOutput[:])
		curStatus.uncommit = len(strings.Split(cmdOutputStr, "\n")) - 1
		total.uncommit += curStatus.uncommit

		menu.SetTitle("[" + strconv.Itoa(curStatus.uncommit) + "|" + strconv.Itoa(curStatus.unpushed) + "] " + path.Base(fullPath))
		menu.SetTooltip(fullPath + "\n\n" + strconv.Itoa(curStatus.uncommit) + " change(s) not committed\n" + strconv.Itoa(curStatus.unpushed) + " change(s) not pushed")
	}

	return total
}

func getIcon(s string) []byte {
	b, err := ioutil.ReadFile(s)
	if err != nil {
		onError(err)
	}
	return b
}

func onError(err error) {
	fmt.Println(err)
	os.Exit(-1)
}
