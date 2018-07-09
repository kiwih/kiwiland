package kiwiserver

import (
	"bytes"
	"os/exec"
	"strings"
)

//RunCECCommand will run a CEC command through the cec-client program
func RunCECCommand(command string) (string, error) {
	cmd := exec.Command("cec-client", "-s", "-d", "1")
	cmd.Stdin = strings.NewReader(command)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out.String(), err
}

//TVGetStatus returns the status of the TV using the built-in cec-client command
//response looks like this: opening a connection to the CEC adapter...\npower status: on\n
func TVGetStatus() (string, error) {
	return RunCECCommand("pow 0")
}

//TVTurnOn turns the TV on using the built-in cec-client command
func TVTurnOn() (string, error) {
	return RunCECCommand("on 0")
}

//TVTurnOff turns the TV off using the built-in cec-client command
func TVTurnOff() (string, error) {
	return RunCECCommand("standby 0")
}

//TVSelectHDMI4 will set the TV to use HDMI4 using the raw tx command
//See http://www.cec-o-matic.com/ to decode
func TVSelectHDMI4() (string, error) {
	return RunCECCommand("tx 4F:82:40:00")
}

//TVSelectHDMI2 will set the TV to use HDMI4 using the raw tx command
//See http://www.cec-o-matic.com/ to decode
func TVSelectHDMI2() (string, error) {
	return RunCECCommand("tx 4F:82:20:00")
}

//TVSelectHDMI1 will set the TV to use HDMI4 using the raw tx command
//See http://www.cec-o-matic.com/ to decode
func TVSelectHDMI1() (string, error) {
	return RunCECCommand("tx 4F:82:10:00")
}

//TVVolumeUp will set the TV to use HDMI4 using the raw tx command
//See http://www.cec-o-matic.com/ to decode
func TVVolumeUp() (string, error) {
	return RunCECCommand("tx 4F:44:41")
}

//TVVolumeDown will set the TV to use HDMI4 using the raw tx command
//See http://www.cec-o-matic.com/ to decode
func TVVolumeDown() (string, error) {
	return RunCECCommand("tx 4F:44:42")
}
