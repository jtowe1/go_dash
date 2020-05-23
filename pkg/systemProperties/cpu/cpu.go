package cpu

import (
	"os/exec"
)

type Info struct {
	Brand string
}

func GetInfo() (*Info, error) {
	var info Info

	out, err := exec.Command("sysctl","-n",  "machdep.cpu.brand_string").Output()
	if err != nil {
		return nil, err
	}

	info.Brand = string(out)
	return &info, nil
}