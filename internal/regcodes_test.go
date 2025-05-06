// Copyright (c) 2025 SUSE LLC. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package containersuseconnect

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestRegcodes(t *testing.T) {
	rc := &Regcodes{}

	if rc.separator() != ':' {
		t.Fatal("Wrong separator")
	}

	prepareLogger()
	err := rc.afterParseCheck()
	msg := "Can't find username"
	if err == nil || err.Error() != msg {
		t.Fatal("Wrong error")
	}
	shouldHaveLogged(t, msg)

	rc.setValues("username", "suse")
	prepareLogger()
	msg = "Can't find password"
	err = rc.afterParseCheck()
	if err == nil || err.Error() != msg {
		t.Fatal("Wrong error")
	}
	shouldHaveLogged(t, msg)

	rc.setValues("password", "1234")
	err = rc.afterParseCheck()
	if err != nil {
		t.Fatal("There should not be an error")
	}

	locs := rc.locations()
	if locs[0] != "/etc/zypp/credentials.d/SCCcredentials" {
		t.Fatal("Wrong location")
	}
	if locs[1] != "/run/secrets/SCCcredentials" {
		t.Fatal("Wrong location")
	}
	if locs[2] != "/run/secrets/credentials.d/SCCcredentials" {
		t.Fatal("Wrong location")
	}

	// It should log a proper warning.
	buffer := bytes.NewBuffer([]byte{})
	log.SetOutput(buffer)
	rc.setValues("unknown", "value")
	if !strings.Contains(buffer.String(), "Warning: Unknown key 'unknown'") {
		t.Fatal("Wrong warning!")
	}
}

// In the following test we will create a mock that just wraps up the
// `RegCodes` struct and replaces its `location` function for something that
// can be tested. We test for a successful run, since all the possible errors
// have already been tested in the `configuration_test.go` file.

type RegcodesMock struct {
	rc Regcodes
}

func (mock *RegcodesMock) locations() []string {
	return []string{"testdata/regcodes.txt"}
}

func (mock *RegcodesMock) onLocationsNotFound() bool {
	return mock.rc.onLocationsNotFound()
}

func (mock *RegcodesMock) separator() byte {
	return mock.rc.separator()
}

func (mock *RegcodesMock) setValues(key, value string) {
	mock.rc.setValues(key, value)
}

func (mock *RegcodesMock) afterParseCheck() error {
	return mock.rc.afterParseCheck()
}

func TestIntegrationRegcodes(t *testing.T) {
	var regcodes Regcodes
	mock := RegcodesMock{rc: regcodes}

	err := ReadConfiguration(&mock)
	if err != nil {
		t.Fatal("This should've been a successful run")
	}

	if len(mock.rc) != 2 {
		t.Fatalf("Unexpected length of regcodes %v", len(mock.rc))
	}

	if mock.rc[0] != "10yb1x6bd159g741ad420fd5aa5083e4" {
		t.Fatal("Unexpected regcode[0]")
	}

	if mock.rc[1] != "36531d07-a283-441b-a02a-1cd9a88b0d5d" {
		t.Fatal("Unexpected regcode[1]")
	}

	if mock.rc.onLocationsNotFound() {
		t.Fatalf("It should've been false")
	}
}
