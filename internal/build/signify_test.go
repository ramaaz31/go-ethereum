// Copyright 2020 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// signFile reads the contents of an input file and signs it (in armored format)
// with the key provided, placing the signature into the output file.

package build

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"testing"
	"time"
)

var (
	testSecKey = "RWRCSwAAAABVN5lr2JViGBN8DhX3/Qb/0g0wBdsNAR/APRW2qy9Fjsfr12sK2cd3URUFis1jgzQzaoayK8x4syT4G3Gvlt9RwGIwUYIQW/0mTeI+ECHu1lv5U4Wa2YHEPIesVPyRm5M="
	testPubKey = "RWTAPRW2qy9FjsBiMFGCEFv9Jk3iPhAh7tZb+VOFmtmBxDyHrFT8kZuT"
)

func TestSignify(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	rand.Seed(time.Now().UnixNano())

	data := make([]byte, 1024)
	rand.Read(data)
	tmpFile.Write(data)

	if err = tmpFile.Close(); err != nil {
		t.Fatal(err)
	}

	err = SignifySignFile(tmpFile.Name(), fmt.Sprintf("%s.sig", tmpFile.Name()), testSecKey)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(fmt.Sprintf("%s.sig", tmpFile.Name()))

	// if signify-openbsd is present, check the signature.
	// signify-openbsd will be present in CI.
	if runtime.GOOS == "linux" {
		cmd := exec.Command("which", "signify-openbsd")
		if err = cmd.Run(); err == nil {
			// Write the public key into the file to pass it as
			// an argument to signify-openbsd
			pubKeyFile, err := ioutil.TempFile("", "")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(pubKeyFile.Name())
			defer pubKeyFile.Close()
			pubKeyFile.WriteString("untrusted comment: signify public key\n")
			pubKeyFile.WriteString(testPubKey)
			pubKeyFile.WriteString("\n")

			cmd := exec.Command("signify-openbsd", "-V", "-p", pubKeyFile.Name(), "-x", fmt.Sprintf("%s.sig", tmpFile.Name()), "-m", tmpFile.Name())
			if output, err := cmd.CombinedOutput(); err != nil {
				fmt.Println(string(output))
				t.Fatalf("could not verify the file: %v", err)
			}
		}
	}
}