package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	evdev "github.com/gvalkov/golang-evdev"
)

const (
	mouseDevPath     = `/dev/input/event5`
	mouseXInputID    = 9
	mouseXInputParam = 155
	mouseButtonX     = 275
	mouseButtonY     = 276
)

func newSensitivityCommand(x, y uint64) *exec.Cmd {
	return exec.Command("xinput",
		"set-prop", strconv.FormatUint(mouseXInputID, 10), strconv.FormatUint(mouseXInputParam, 10),
		strconv.FormatUint(x, 10)+`.000000`, `0.000000`, `0.000000`, `0.000000`,
		strconv.FormatUint(y, 10)+`.000000`, `0.000000`, `0.000000`, `0.000000`,
		`1.000000`,
	)
}

func setXY(x, y uint64) {
	cmd := newSensitivityCommand(x, y)
	fmt.Println(cmd.Args)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("error setting x%d,y%d: %v\n", x, y, err)
		fmt.Println(string(output))
	}
}

func main() {

	mouseDev, err := evdev.Open(mouseDevPath)
	if err != nil {
		fmt.Printf("failed to open mouse input device: %v", err)
		os.Exit(1)
	}
	for {
		ev, err := mouseDev.ReadOne()
		if err != nil {
			fmt.Printf("error reading %s: %v", mouseDev.Name, err)
			continue
		}
		fmt.Println(ev.String())
		switch ev.Code {
		case mouseButtonX:
			if ev.Value == 0 {
				setXY(1, 1)
			} else if ev.Value == 1 {
				setXY(1, 0)
			}
		case mouseButtonY:
			if ev.Value == 0 {
				setXY(1, 1)
			} else if ev.Value == 1 {
				setXY(0, 1)
			}
		}
	}

	select {}
}
