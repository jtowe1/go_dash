package cpu

import (
	"os/exec"
	"strings"
)

func GetInfo() (string, error) {
	out, err := exec.Command("sysctl", "machdep.cpu.brand_string").Output()
	if err != nil {
		return "", err
	}
	trimmedSlice := strings.Split(string(out), ":")
	trimmedString := strings.TrimSpace(trimmedSlice[1])
	return trimmedString, nil
}