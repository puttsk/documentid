package documentid_test

import (
	"testing"
	"time"

	"github.com/puttsk/documentid"
)

type GenerateDocumentIDTestcase struct {
	Prefix   string
	Suffix   string
	Number   int64
	Zeroes   int64
	Time     time.Time
	Expected string
	Args     map[string]interface{}
}

func TestGetGenerateDocumentID(t *testing.T) {
	testTime := time.Date(2022, time.April, 10, 0, 0, 0, 0, time.Local)
	cases := []GenerateDocumentIDTestcase{
		{Prefix: "%(y)-%(m)-%(d)-", Suffix: "", Number: 10, Zeroes: 5, Time: testTime, Expected: "2022-04-10-00010"},
		{Prefix: "R%(date)-", Suffix: "-%(y)", Number: 10, Zeroes: 7, Time: testTime, Expected: "R20220410-0000010-2022"},
		{Prefix: "TEST%(date)-", Suffix: "-%(noargs)", Number: 10, Zeroes: 7, Time: testTime, Expected: "TEST20220410-0000010-"},
		{Prefix: "%(pre)%(d)-", Suffix: "-%(hello)-%(date)", Number: 10, Zeroes: 7, Time: testTime, Expected: "12010-0000010-HELLO-20220410", Args: map[string]interface{}{"pre": 120, "hello": "HELLO"}}, // Test variable override
		{Prefix: "TEST%(d)-", Suffix: "-%(y)", Number: 10, Zeroes: 7, Time: testTime, Expected: "TEST10-0000010-2022", Args: map[string]interface{}{"d": "date", "y": 0}},                                 // Test variable override
	}

	for _, c := range cases {
		actual := documentid.GenerateDocumentID(c.Prefix, c.Suffix, c.Time, c.Zeroes, c.Number, c.Args)
		if actual != c.Expected {
			t.Errorf("Invalid result: Actual: %s, Expected: %s, Input: %+v", actual, c.Expected, c)
		}
	}
}

func BenchmarkGenerateDocumentID(b *testing.B) {
	t := time.Now()
	for n := 0; n < b.N; n++ {
		documentid.GenerateDocumentID("DOC-%(date)-", "-TEST%(y)", t, 5, int64(n), nil)
	}
}

func BenchmarkReplaceVars(b *testing.B) {
	for n := 0; n < b.N; n++ {
		documentid.ReplaceVars("%(a)-%(b)", map[string]interface{}{"a": 123, "b": "test"})
	}
}
