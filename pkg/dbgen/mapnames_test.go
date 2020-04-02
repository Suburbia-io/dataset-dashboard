package dbgen

import "testing"

func TestSQLNameToGo(t *testing.T) {
	type TestCase struct {
		In  string
		Out string
	}

	cases := []TestCase{
		{"api_key", "APIKey"},
		{"user_id", "UserID"},
		{"email", "Email"},
		{"created_at", "CreatedAt"},
		{"sftp_username", "SFTPUsername"},
	}

	for _, tc := range cases {
		out := SQLNameToGo(tc.In)
		if out != tc.Out {
			t.Fatalf("%s != %s", out, tc.Out)
		}
	}
}
