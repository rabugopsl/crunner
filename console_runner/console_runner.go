package console_runner

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Instruction represents a single user input passed to underlying executable
type Instruction struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
	Wait  int         `json:"wait,omitempty"`
}

// Returns an array of Instructions parsed from the JSON file passed as parameter
func LoadInputParams(src string) ([]Instruction, error) {

	file, err := os.Open(src)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var params []Instruction
	jsonDec := json.NewDecoder(file)
	jsonDec.Decode(&params)

	return params, nil
}

// RunCommand runs the specified command using the args and proceeds to deliver params as user input
func RunCommand(debug bool, params []Instruction, command string, args ...string) error {

	if debug {
		fmt.Printf("<--Run Data\n")
		fmt.Printf("\t<--Instructions: \n\t\t%v\n", params)
		fmt.Printf("\t<--Command: \t%v\n", command)
		fmt.Printf("\t<--Command Arguments: \t%v\n", args)
	}
	cmd := exec.Command(command, args...)
	cmdstdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	cmdstdout, err := cmd.StdoutPipe()
	defer cmdstdout.Close()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	index := 0
	inchan := make(chan string)
	outchan := make(chan string)
	readychan := make(chan bool)
	errchan := make(chan error)

	go writeToStdin(cmdstdin, inchan, errchan)
	go readFromStdout(cmdstdout, readychan, outchan, errchan)

	for {
		select {
		case err := <-errchan:
			return err
		case out := <-outchan:
			fmt.Printf("%s", out)
			if index < len(params) && strings.Contains(out, params[index].Key) {
				if debug {
					fmt.Printf("<--Delivering data: %s\n", params[index].Value.(string))
				}
				inchan <- params[index].Value.(string)
				time.Sleep(time.Duration(params[index].Wait) * time.Second)
				index++
			}
		}
	}
	cmd.Wait()

	return nil
}

func readFromStdout(cmdstdout io.ReadCloser, readychan chan<- bool, outchan chan string, errchan chan<- error) {
	outbuff := make([]byte, 1000)
	for {
		for outbytes, err := cmdstdout.Read(outbuff); outbytes > 0 && err != io.EOF; outbytes, err = cmdstdout.Read(outbuff) {
			if err != nil {
				errchan <- err
			}
			outchan <- string(outbuff[:outbytes])
		}
	}
}

func writeToStdin(cmdstdin io.WriteCloser, inchan <-chan string, errchan chan<- error) {
	defer cmdstdin.Close()
	for {
		if _, err := io.WriteString(cmdstdin, <-inchan+"\n"); err != nil {
			errchan <- err
		}
	}
}
