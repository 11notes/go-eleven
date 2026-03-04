package util

import (
	"os"
	"os/exec"
	"syscall"
	"strings"
	"bufio"
	"io"
	"errors"
	"regexp"
	"fmt"
	"crypto/sha256"
	"time"
)

type Util struct{}

// reads a file if it exists and returns the content of the file
func (c *Util) ReadFile(path string) (string, error){
	bytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// writes contents to a file
func (c *Util) WriteFile(path string, txt string) error{
	err := os.WriteFile(path, []byte(txt), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// checks if the command line argument exists (case-sensitive)
func (c *Util) CommandLineArgumentExists(f string) bool{
	if(len(os.Args) > 1){
		for _, a := range os.Args[1:] {
			if(f == a){
				return true
			}
		}
	}

	return false
}

// checks if an environment variable exists and if not assigns a default value
func (c *Util) GetEnv(key string, fallback string) string{
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// checks if a file containing an environment variable exists and if not assigns a default value
func (c *Util) GetEnvFile(path string, fallback string) string{
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		value, err := c.ReadFile(path)
		if err != nil {
			return fallback
		}
		return value
	}
	return fallback
}

// run an external program and return output
func (c *Util) Run(bin string, params []string) (string, error){
	cmd := exec.Command(bin, params...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid:true}

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	out := []string{}
	go func() {
		stdoutScanner := bufio.NewScanner(io.MultiReader(stdout,stderr))
		for stdoutScanner.Scan() {
			out = append(out, stdoutScanner.Text())
		}
	}()

	err := cmd.Start()
	if err != nil {
		return "", errors.New(err.Error() + strings.Join(out, " "))
	}
	err = cmd.Wait()
	if err != nil {
		return "", errors.New(err.Error() + strings.Join(out, " "))
	}

	return strings.Join(out, " "), nil
}

// replace all variables in a string
func (c *Util) StringReplaceVar(str string, r map[string]interface{}) string{
	// replace all variables
	for key, value := range r{
		str = string(regexp.MustCompile(fmt.Sprintf(`\${%s}`, key)).ReplaceAllString(str, fmt.Sprintf("%s", value)))
	}

	// replace all not set variables with an empty string
	empty := regexp.MustCompile(`\$\{[A-Z_a-z]+\}`).FindAllString(str, -1)
	for _, e := range empty {
		str = string(regexp.MustCompile(fmt.Sprintf(`%s`, e)).ReplaceAllString(str, ""))
	}

	return str
}

// check if string is present in file
func (c *Util) FileContains(file string, str string) (bool, error){
	// open file
	text, err := c.ReadFile(file)
	if err != nil {
		return false, err
	}

	// check for string
	return strings.Contains(text, str), nil
}

// create random default password made up of four blocks containing 5 random characters
func (c *Util) Password() string{
	str := fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%x", time.Now().Unix()))))
	m := regexp.MustCompile(`.{1,5}`).FindAllString(str, -1)
	return strings.Join(m[0:3], ".")
}

// replace variables in a file
func (c *Util) FileReplaceStrings(file string, str map[string]interface{}) (bool, error){
	// set initial state
	replaced := false

	// open file
	text, err := c.ReadFile(file)
	if err != nil {
		return false, err
	}

	// replace all variables
	for key, value := range str{
		if strings.Contains(text, key) {
			replaced = true
			text = strings.ReplaceAll(text, key, fmt.Sprintf("%s", value))
		}
	}

	// write file
	err = c.WriteFile(file, text)
	if err != nil {
		return false, err
	}

	// return
	return replaced, nil
}