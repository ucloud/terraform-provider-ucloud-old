package ucloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/hashcode"
)

// ifaceToStringSlice used for converting terraform attribute of TypeString embedded in TypeList to a string slice.
// it expected interface{} type as []interface{}, usually get the value from `d.Get` of terraform resource data.
func ifaceToStringSlice(iface interface{}) []string {
	s := []string{}

	for _, i := range iface.([]interface{}) {
		s = append(s, i.(string))
	}

	return s
}

func hashStringArray(arr []string) string {
	var buf bytes.Buffer

	for _, s := range arr {
		buf.WriteString(fmt.Sprintf("%s-", s))
	}

	return fmt.Sprintf("%d", hashcode.String(buf.String()))
}

func getAbsPath(filePath string) (string, error) {
	if strings.HasPrefix(filePath, "~") {
		usr, err := user.Current()
		if err != nil {
			return "", fmt.Errorf("get current user got an error: %#v", err)
		}

		if usr.HomeDir != "" {
			filePath = strings.Replace(filePath, "~", usr.HomeDir, 1)
		}
	}

	return filepath.Abs(filePath)
}

func writeToFile(filePath string, data interface{}) error {
	absPath, err := getAbsPath(filePath)
	if err != nil {
		return err
	}

	os.Remove(absPath)

	var bs []byte
	switch data.(type) {
	case string:
		bs = []byte(data.(string))
	default:
		bs, err = json.MarshalIndent(data, "", "\t")
		if err != nil {
			return fmt.Errorf("MarshalIndent data %#v and got an error: %#v", data, err)
		}
	}

	ioutil.WriteFile(absPath, bs, 422)
	return nil
}

func checkStringIn(val string, availables []string) error {
	for _, choice := range availables {
		if val == choice {
			return nil
		}
	}

	return fmt.Errorf("should be one of %s, got %s", strings.Join(availables, ","), val)
}

func checkIntIn(val int, availables []int) error {
	for _, choice := range availables {
		if val == choice {
			return nil
		}
	}

	return fmt.Errorf("should be one of %v, got %d", availables, val)
}
func timestampToString(ts int) string {
	return time.Unix(int64(ts), 0).Format(time.RFC3339)
}

func stringToTimestamp(ts string) (int, error) {
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return 0, err
	}
	return int(t.Unix()), nil
}
