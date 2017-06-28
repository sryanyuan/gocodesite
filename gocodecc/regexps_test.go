package gocodecc

import (
	"strings"
	"testing"
)

func TestMatchAtPeople(t *testing.T) {
	substrs := mentionPeopleReg.FindAllString("hgrwauorge @12342a  @45bdf6 @hello @your1z @hhahagw hgwra", -1)
	if nil == substrs {
		t.FailNow()
	}
	t.Error(len(substrs), substrs)
}

func TestFind(t *testing.T) {
	strs := []string{"@12342a", "@45bdf6", "@hello", "@your1z", "@hhahagw"}
	for _, v := range strs {
		if !strings.Contains("hgrwauorge @12342a  @45bdf6 @hello @your1z @hhahagw hgwra", v) {
			t.Error(v)
		}
	}
	t.Error("OK")
}
