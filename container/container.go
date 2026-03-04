package container

import (
	"errors"
	"os"
	"io/ioutil"
	"strings"

	"github.com/11notes/go-eleven/util"
)

type Container struct{}

// tries to get a secret either from environment variable or from a secrets file set by environment variable
func (c *Container) GetSecret(env string, envPath string) (string, error){
	if value, ok := os.LookupEnv(env); ok {
		return value, nil
	}else{
		if value, ok := os.LookupEnv(envPath); ok {
			bytes, err := ioutil.ReadFile(value)
			if err != nil {
				return "", err
			}
			return strings.TrimSpace(string(bytes)), nil
		}else{
			return "", errors.New(env + " and " + envPath + " do not exist!")
		}
	}
}

// merges default entrypoint command and user provided command
func (c *Container) MergeCommand(d []string) []string{
	if(len(os.Args) > 0){
		args := os.Args[1:]
		for _, value := range args{
			d = append(d, value)
		}
	}
	return(d)
}

// replaces variables inside a file
func (c *Container) FileContentReplace(file string, v map[string]interface{}) error{
	// open file
	text, err := (&util.Util{}).ReadFile(file)
	if err != nil {
		return err
	}

	text = (&util.Util{}).StringReplaceVar(text, v)

	// write file
	err = (&util.Util{}).WriteFile(file, text)
	if err != nil {
		return err
	}

	return nil
}

// replaces all environment variables inside a file
func (c *Container) FileContentReplaceEnv(file string) error{
	env := map[string]any{}
	for _, e := range os.Environ() {
		key := strings.Split(e, "=")[0]
		value := os.Getenv(key)
		env[key] = value
	}

	return c.FileContentReplace(file, env)
}

// converts an environment variable to a file with the file content being the value of the variable
func (c *Container) EnvToFile(env string, path string) error{
	if value, ok := os.LookupEnv(env); ok {
		return (&util.Util{}).WriteFile(path, value)
	}else{
		return errors.New(env + " does not exist!")
	}
}