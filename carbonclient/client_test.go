package carbonclient

import (
	"io/ioutil"
	"testing"
)

func TestParseRes(t *testing.T) {
	buffer, err := ioutil.ReadFile("sample_notifications_data.json")
	if err != nil {
		t.Error(err.Error())
	}

	_, err = parseNotificationsResBody(buffer)
	if err != nil {
		t.Error(err.Error())
	}
}
