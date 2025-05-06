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
	"log"
	"os"
	"strings"
)

// Credentials holds the product regcodes
type Regcodes []string

func (rc *Regcodes) separator() byte {
	return ':'
}

func (rc *Regcodes) locations() []string {
	return []string{
		"/run/secrets/SUSEregcodes",
		"/run/secrets/credentials.d/SUSEregcodes",
	}
}

func (rc Regcodes) onLocationsNotFound() bool {
	env_regcode := os.Getenv("SCC_REGCODES")

	if env_regcode != "" {
		rc = strings.SplitN(env_regcode, string(rc.separator()), -1)
		return true
	}

	return false
}

func (rc Regcodes) setValues(key, value string) {
	log.Printf("setValues with %v: %v", key, value)
	rc = append(rc, value)
}

func (rc *Regcodes) afterParseCheck() error {
	return nil
}
