package input

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type FileInput struct {
	Filename string
	Meta     FileMeta
}

func (file_input FileInput) GetData() ([]byte, error) {
	data, err := os.ReadFile(file_input.Filename)
	if err != nil {
		fmt.Println("[ERROR]", err)
		return nil, errors.New("[ERROR] Failed to read file")
	}
	return data, nil
}

func ExpandInput(filename string) ([]Input, error) {
	inputs := make([]Input, 0)
	stat_obj, err := os.Lstat(filename)
	if err != nil {
		fmt.Println("[ERROR] ", err)
		return nil, errors.New("[ERROR] Failed to open location")
	}
	if (stat_obj.Mode() & os.ModeSymlink) == os.ModeSymlink {
		fmt.Println("[WARNING] Excluding symlink: " + filename)
	} else if stat_obj.IsDir() {
		subdirs, err := ioutil.ReadDir(filename)
		if err != nil {
			fmt.Println("[ERROR] Couldn't list directory " + filename)
			return nil, errors.New("[ERROR] Failed to read directory")
		}
		for _, subdir := range subdirs {
			sub_inputs, err := ExpandInput(filename + "/" + subdir.Name())
			if err != nil {
				fmt.Println("[ERROR] Failed to expand subdirectories of " + subdir.Name())
				return nil, errors.New("[ERROR] Subdirectory error")
			}
			inputs = append(inputs, sub_inputs...)
		}
	} else {
		// TODO This may be a good place to verify that files are readable or error out
		inputs = append(inputs, FileInput{
			Filename: filename,
			Meta: FileMeta{
				Mode: stat_obj.Mode(),
			},
		})
	}
	return inputs, nil
}
