// This file is part of Switchconf.
//
// Switchconf is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Switchconf is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Switchconf.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

var (
	configFile   = flag.String("f", "", "path list file to use")
	verbose   = flag.Bool("v", false, "Details of what's going on")
	switchconfRe = regexp.MustCompile("^[ \t]*(?P<comment>[^ \t]+)[ \t]*switchconf:[ \t]*(?P<context>.*)$")

	disabledTag = ":off:"
	switchConfTag = ":switchconf:"
)

// Read a whole file into the memory and store it as array of lines
func readLines(path string) (lines []string, err error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	lines = strings.Split(string(content), "\n")
	return
}

func writeLines(lines []string, path string) (err error) {
	file, err := os.Create(path)
	if err != nil {
		return
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	content := strings.Join(lines, "\n")
	_, err = writer.WriteString(content)
	writer.Flush()
	return
}

func commentOut(line, comment string) string {
	disCom := fmt.Sprintf("%s%s", comment, disabledTag)
	Logf("DisCom: %s\n", disCom)
	if strings.Index(line, disCom) == 0 {
		// Already commented out, just return it
		return line
	}
	return fmt.Sprintf("%s%s", disCom, line)
}

func uncomment(line, comment string) string {
	re, err := regexp.Compile(fmt.Sprintf("^%s:off:", comment))
	if err != nil {
		return line
	}

	return re.ReplaceAllString(line, "")
}

func processFile(path string, file ConfigFile, context string) (err error) {
	Logf("Processing %s\n", path)
	lines, err := readLines(path)
	if err != nil {
		return
	}
	curcontext := ""
	reStr := fmt.Sprintf("^%s%s[ \t]*(?P<context>.*)$", file.Comment, switchConfTag)
	switchconfRe, err := regexp.Compile(reStr)
	for i, line := range lines {
		matches := switchconfRe.FindStringSubmatch(line)
		if len(matches) > 0 {
			curcontext = matches[1]
			if err != nil {
				return
			}
			continue
		}
		if context == curcontext {
			lines[i] = uncomment(line, file.Comment)
		} else if curcontext != "" {
			lines[i] = commentOut(line, file.Comment)
		}
	}
	writeLines(lines, path)
	return
}

func main() {
	//processFile("foo.txt", "intel")
	flag.Parse()
	files, err := readConfigFile(*configFile)
	if err != nil {
		panic(err)
	}
	if len(flag.Args()) != 1 {
		fmt.Println("One single context name required.")
		return
	}
	for fileName, fileInfo := range files {
		err = processFile(fileName, fileInfo, flag.Args()[0])
		if err != nil {
			panic(err)
		}
	}
}
