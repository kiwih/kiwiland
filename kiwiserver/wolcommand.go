package kiwiserver

import (
	"bytes"
	"os/exec"
	"strings"
)

//ToshibaWOL runs the rpi wakeonlan command
//wakeonlan 04:7D:7B:5B:FE:4D
func ToshibaWOL() (string, error) {
	cmd := exec.Command("wakeonlan", "04:7D:7B:5B:FE:4D")
	cmd.Stdin = strings.NewReader("")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out.String(), err
}
