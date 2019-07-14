// +build go1.12

package razorcore

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var (
	casesOpt   = "casesOpt"
	testOpt    = "testOpt"
	testOptGen = "testOptGen"
)

func TestGenerateOpt(t *testing.T) {
	casedir, _ := filepath.Abs(filepath.Dir("./" + casesOpt + "/"))
	testGenDir, _ := filepath.Abs(filepath.Dir("./" + testOptGen + "/"))
	sap := string(filepath.Separator)

	visit := func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() { // regular file
			if strings.HasPrefix(filepath.Base(path), ".") {
				return nil
			}
			name := strings.Replace(path, ".gohtml", ".go", 1)
			cmp := strings.Replace(name, sap+casesOpt+sap, sap+testOpt+sap, -1)
			log := strings.Replace(name, sap+casesOpt+sap, sap+testOptGen+sap, -1)

			if !exists(cmp) {
				t.Error("No cmp:", cmp)
			} else if !exists(log) {
				t.Error("No log:", log)
			} else {
				//compare the log file and cmp file
				_cmp, _e1 := ioutil.ReadFile(cmp)
				_log, _e2 := ioutil.ReadFile(log)
				if _e1 != nil || _e2 != nil {
					t.Error("Reading")
				} else if string(_cmp) != string(_log) {
					t.Error("MISMATCH:", log, cmp)
				} else {
					t.Log("PASS")
				}
			}
		}
		return nil
	}
	QuickMode = false
	option := Option{}
	GenFolder(casedir, testGenDir, option)
	err := filepath.Walk(casedir, visit)
	if err != nil {
		t.Error("walk")
	}
}
