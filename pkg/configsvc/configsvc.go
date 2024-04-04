package configsvc

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// DEFAULT_DATA_PATH default path in which to find data
const DEFAULT_DATA_PATH = "/opt/data"

// CONFIG_DATA_PATH_ENV_VAR Env var from which to load data path
const CONFIG_DATA_PATH_ENV_VAR = "CONFIG_DATA_PATH"

// staticData singleton to store data in memory
var staticData map[string][]byte

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

// MarshalData reads files in CONFIG_DATA_PATH and loads them into memory in a map keyed by the filename with the value being the bytes of the file contents.
func MarshalData(path string) (data map[string][]byte, err error) {
	data = make(map[string][]byte)
	if path == "" {
		path = DEFAULT_DATA_PATH
	}

	// open path, or DEFAULT_DATA_PATH
	logrus.Debugf("Opening %s", path)
	files, err := os.ReadDir(path)
	if err != nil {
		err = errors.Wrapf(err, "failed reading path %s", path)
		return data, err
	}

	dotfile := regexp.MustCompile(`^\..+`)

	// iterate over the files found
	for _, f := range files {
		// skip dotfiles
		if dotfile.MatchString(f.Name()) {
			continue
		}

		// read the file
		filePath := fmt.Sprintf("%s/%s", path, f.Name())
		logrus.Debugf("Opening %s", filePath)

		// read the file.  We really don't care about the contents.
		b, err := os.ReadFile(filePath)
		if err != nil {
			err = errors.Wrapf(err, "failed reading file %s", filePath)
			return data, err
		}

		// populate static info for each path.  Yes, storing huge files in memory will get expensive, but that's really not what is intended by this tool.

		logrus.Debugf("Setting key %q", f.Name())
		data[f.Name()] = b
	}

	return data, err
}

// Server runs the http server to return data
func Server(address string, path string) (err error) {
	if path == "" {
		path = os.Getenv(CONFIG_DATA_PATH_ENV_VAR)
	}

	data, err := MarshalData(path)
	if err != nil {
		return err
	}

	staticData = data

	logrus.Infof("Starting ConfigSvc on %s reading path %q for data\n", address, path)

	http.HandleFunc("/", InfoHandler)
	err = http.ListenAndServe(address, nil)

	return err
}

// InfoHandler Receives the HTTP Request, and writes the bytes from 'staticData[path]' back to the client.  Returns 400 if the path doesn't exist.
func InfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	path := r.URL.Path

	path = strings.TrimLeft(path, "/")

	logrus.Infof("Serving data from path %s", path)

	data, ok := staticData[path]
	if ok {
		logrus.Debugf("Found data at path %s", path)
		_, err := w.Write(data)
		if err != nil {
			logrus.Errorf("failed serving data from path %s: %s", path, err)
		}
		return
	}

	logrus.Debugf("no data at path %s", path)
	w.WriteHeader(http.StatusBadRequest)
}
