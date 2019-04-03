package eidas

import (
	"encoding/hex"
	"testing"
)

const defaultCAName = "Financial Conduct Authority"
const defaultCAID = "GB-FCA"

func TestSimple(t *testing.T) {
	pspAS := "305b3013060604008e4601063009060704008e4601060330440606040081982702303a301330110607040081982701010c065053505f41530c1b46696e616e6369616c20436f6e6475637420417574686f726974790c0647422d464341"
	d, err := Serialize([]string{"PSP_AS"}, defaultCAName, defaultCAID)
	if err != nil {
		t.Error(err)
	}
	if hex.EncodeToString(d) != pspAS {
		t.Error("Mismatch with PSP_AS")
	}
}

func TestAll(t *testing.T) {
	type testData struct{
		Expected string
		Roles []string
	}
	expected := []testData{testData{
		Expected: "305b3013060604008e4601063009060704008e4601060330440606040081982702303a301330110607040081982701010c065053505f41530c1b46696e616e6369616c20436f6e6475637420417574686f726974790c0647422d464341",
		Roles: []string{"PSP_AS"},
	}, testData{
		Expected: "305b3013060604008e4601063009060704008e4601060330440606040081982702303a301330110607040081982701020c065053505f50490c1b46696e616e6369616c20436f6e6475637420417574686f726974790c0647422d464341",
		Roles: []string{"PSP_PI"},
	}, testData{
		Expected: "305b3013060604008e4601063009060704008e4601060330440606040081982702303a301330110607040081982701030c065053505f41490c1b46696e616e6369616c20436f6e6475637420417574686f726974790c0647422d464341",
		Roles: []string{"PSP_AI"},
	}, testData{
		Expected: "305b3013060604008e4601063009060704008e4601060330440606040081982702303a301330110607040081982701040c065053505f49430c1b46696e616e6369616c20436f6e6475637420417574686f726974790c0647422d464341",
		Roles: []string{"PSP_IC"},
	}, testData{
		Expected: "306c3013060604008e4601063009060704008e4601060330550606040081982702304b302430220607040081982701010c065053505f41530607040081982701020c065053505f50490c1b46696e616e6369616c20436f6e6475637420417574686f726974790c0647422d464341",
		Roles: []string{"PSP_AS", "PSP_PI"},
	}}
	for _, e := range expected {
		// Check our serialization matches theirs.
		s, err := Serialize(e.Roles, defaultCAName, defaultCAID)
		if err != nil {
			t.Error(err)
		}
		if hex.EncodeToString(s) != e.Expected {
			t.Errorf("Mismatch with roles: %v", e.Roles)
		}

		// Check we can extract the roles, name and ID correctly.
		d, err := hex.DecodeString(e.Expected)
		if err != nil {
			t.Error(err)
		}
		roles, name, id, err := Extract(d)
		if err != nil {
			t.Error(err)
		}
		for i, r := range roles {
			if e.Roles[i] != r {
				t.Errorf("Expected role: %s but got %s", e.Roles[i], r)
			}
		}
		if name != defaultCAName {
			t.Errorf("Expected CA name: %s but got %s", defaultCAName, name)
		}
		if id != defaultCAID {
			t.Errorf("Expected CA id: %s but got %s", defaultCAID, id)
		}
	}
}
