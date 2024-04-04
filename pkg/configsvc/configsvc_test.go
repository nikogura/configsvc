package configsvc

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMarshalData(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	expectedBytes := []byte(`{
  "foo": "bar",
  "wip": "zoz"
}
`)

	expected := map[string][]byte{
		"test.json": expectedBytes,
	}

	path := "../../testdata"

	actual, err := MarshalData(path)
	if err != nil {
		t.Errorf("failed reading data from %s: %s", path, err)
	}

	assert.Equal(t, expected, actual, "returned data does not match expectations.")

}

func TestInfoHandler(t *testing.T) {
	expectedBytes := []byte(`{
  "foo": "bar",
  "wip": "zoz"
}
`)

	path := "../../testdata"

	data, err := MarshalData(path)
	if err != nil {
		t.Errorf("failed reading data from %s: %s", path, err)
	}

	staticData = data

	req := httptest.NewRequest(http.MethodGet, "http://configsvc:8888/test.json", nil)

	w := httptest.NewRecorder()

	InfoHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	actualBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("failed reading response body: %s", err)
	}

	assert.Equal(t, expectedBytes, actualBytes, "returned bytes do not meet expecations")

}
