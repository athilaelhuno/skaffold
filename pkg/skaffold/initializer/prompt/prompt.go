/*
Copyright 2020 The Skaffold Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package prompt

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/AlecAivazis/survey.v1"
)

// For testing
var (
	BuildConfigFunc = buildConfig
)

func buildConfig(image string, choices []string) (string, error) {
	var selectedBuildConfig string
	prompt := &survey.Select{
		Message:  fmt.Sprintf("Choose the builder to build image %s", image),
		Options:  choices,
		PageSize: 15,
	}
	err := survey.AskOne(prompt, &selectedBuildConfig, nil)
	if err != nil {
		return "", err
	}

	return selectedBuildConfig, nil
}

func WriteSkaffoldConfig(out io.Writer, pipeline []byte, filePath string) (bool, error) {
	fmt.Fprintln(out, string(pipeline))

	reader := bufio.NewReader(os.Stdin)
confirmLoop:
	for {
		fmt.Fprintf(out, "Do you want to write this configuration to %s? [y/n]: ", filePath)

		response, err := reader.ReadString('\n')
		if err != nil {
			return true, errors.Wrap(err, "reading user confirmation")
		}

		response = strings.ToLower(strings.TrimSpace(response))
		switch response {
		case "y", "yes":
			break confirmLoop
		case "n", "no":
			return true, nil
		}
	}
	return false, nil
}
