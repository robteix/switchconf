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
	"errors"
	"fmt"
	"io/ioutil"
	"launchpad.net/goyaml"
	"os"
	"os/user"
)

type ConfigFile struct {
	Comment string
	Enabled bool
}

func pickConfigFile() string {
	defaultConfFileLocations := []string{"/etc/switchconf"}
	if curUser, err := user.Current(); err == nil {
		homeConfLoc := fmt.Sprintf("%s/.switchconf", curUser.HomeDir)
		defaultConfFileLocations = []string{homeConfLoc, "/etc/switchconf"}
	}
	for _, i := range defaultConfFileLocations {
		if _, err := os.Stat(i); err == nil {
			return i
		}
	}

	return ""
}

func readConfigFile(path string) (files map[string]ConfigFile, err error) {

	// if path is empty, then we look for default locations
	if path == "" {
		if path = pickConfigFile(); path == "" {
			err = errors.New("No configuration file found (~/.switchconf or /etc/switchconf)")
			return
		}
	}
	if (*verbose) {
		Logf("Using configuration file %s\n", path)
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	err = goyaml.Unmarshal([]byte(data), &files)

	if err != nil {
		panic(err)
	}
	for k, v := range files {
		if v.Comment == "" {
			cf := v
			cf.Comment = "#"
			files[k] = cf
		}
	}
	return
}
