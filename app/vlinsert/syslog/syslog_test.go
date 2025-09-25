package syslog

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"github.com/VictoriaMetrics/VictoriaLogs/app/vlinsert/insertutil"
)

func TestSyslogLineReader_Success(t *testing.T) {
	f := func(data string, linesExpected []string) {
		t.Helper()

		r := bytes.NewBufferString(data)
		slr := getSyslogLineReader(r)
		defer putSyslogLineReader(slr)

		var lines []string
		for slr.nextLine() {
			lines = append(lines, string(slr.line))
		}
		if err := slr.Error(); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if !reflect.DeepEqual(lines, linesExpected) {
			t.Fatalf("unexpected lines read;\ngot\n%q\nwant\n%q", lines, linesExpected)
		}
	}

	f("", nil)
	f("\n", nil)
	f("\n\n\n", nil)

	f("foobar", []string{"foobar"})
	f("foobar\n", []string{"foobar\n"})
	f("\n\nfoo\n\nbar\n\n", []string{"foo\n\nbar\n\n"})

	f(`Jun  3 12:08:33 abcd systemd: Starting Update the local ESM caches...`, []string{"Jun  3 12:08:33 abcd systemd: Starting Update the local ESM caches..."})

	f(`Jun  3 12:08:33 abcd systemd: Starting Update the local ESM caches...

48 <165>Jun  4 12:08:33 abcd systemd[345]: abc defg<123>1 2023-06-03T17:42:12.345Z mymachine.example.com appname 12345 ID47 [exampleSDID@32473 iut="3" eventSource="Application 123 = ] 56" eventID="11211"] This is a test message with structured data.

`, []string{
		"Jun  3 12:08:33 abcd systemd: Starting Update the local ESM caches...",
		"<165>Jun  4 12:08:33 abcd systemd[345]: abc defg",
		`<123>1 2023-06-03T17:42:12.345Z mymachine.example.com appname 12345 ID47 [exampleSDID@32473 iut="3" eventSource="Application 123 = ] 56" eventID="11211"] This is a test message with structured data.`,
	})
}

func TestSyslogLineReader_Failure(t *testing.T) {
	f := func(data string) {
		t.Helper()

		r := bytes.NewBufferString(data)
		slr := getSyslogLineReader(r)
		defer putSyslogLineReader(slr)

		if slr.nextLine() {
			t.Fatalf("expecting failure to read the first line")
		}
		if err := slr.Error(); err == nil {
			t.Fatalf("expecting non-nil error")
		}
	}

	// invalid format for message size
	f("12foo bar")

	// too big message size
	f("123 aa")
	f("1233423432 abc")
}

func TestProcessStreamInternal_Success(t *testing.T) {
	f := func(data string, currentYear int, timestampsExpected []int64, resultExpected string) {
		t.Helper()

		MustInit()
		defer MustStop()

		globalTimezone = time.UTC
		globalCurrentYear.Store(int64(currentYear))

		tlp := &insertutil.TestLogMessageProcessor{}
		r := bytes.NewBufferString(data)
		if err := processStreamInternal(r, "", false, "1.2.3.4", tlp); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if err := tlp.Verify(timestampsExpected, resultExpected); err != nil {
			t.Fatal(err)
		}
	}

	data := `Jun  3 12:08:33 abcd systemd: Starting Update the local ESM caches...

Sep 19 08:26:10 host CEF:0|Security|threatmanager|1.0|100|worm successfully stopped|10|src=10.0.0.1 dst=2.1.2.2 spt=1232
48 <165>Jun  4 12:08:33 abcd systemd[345]: abc defg<123>1 2023-06-03T17:42:12.345Z mymachine.example.com appname 12345 ID47 [exampleSDID@32473 iut="3" eventSource="Application 123 = ] 56" eventID="11211"] This is a test message with structured data.
`
	currentYear := 2023
	timestampsExpected := []int64{1685794113000000000, 1695111970000000000, 1685880513000000000, 1685814132345000000}
	resultExpected := `{"format":"rfc3164","hostname":"abcd","app_name":"systemd","_msg":"Starting Update the local ESM caches...","remote_ip":"1.2.3.4"}
{"format":"rfc3164","hostname":"host","app_name":"CEF","cef.version":"0","cef.device_vendor":"Security","cef.device_product":"threatmanager","cef.device_version":"1.0","cef.device_event_class_id":"100","cef.name":"worm successfully stopped","cef.severity":"10","cef.extension.src":"10.0.0.1","cef.extension.dst":"2.1.2.2","cef.extension.spt":"1232","remote_ip":"1.2.3.4"}
{"priority":"165","facility_keyword":"local4","level":"notice","facility":"20","severity":"5","format":"rfc3164","hostname":"abcd","app_name":"systemd","proc_id":"345","_msg":"abc defg","remote_ip":"1.2.3.4"}
{"priority":"123","facility_keyword":"solaris-cron","level":"error","facility":"15","severity":"3","format":"rfc5424","hostname":"mymachine.example.com","app_name":"appname","proc_id":"12345","msg_id":"ID47","exampleSDID@32473.iut":"3","exampleSDID@32473.eventSource":"Application 123 = ] 56","exampleSDID@32473.eventID":"11211","_msg":"This is a test message with structured data.","remote_ip":"1.2.3.4"}`
	f(data, currentYear, timestampsExpected, resultExpected)
}

func TestProcessStreamInternal_Failure(t *testing.T) {
	f := func(data string) {
		t.Helper()

		MustInit()
		defer MustStop()

		tlp := &insertutil.TestLogMessageProcessor{}
		r := bytes.NewBufferString(data)
		if err := processStreamInternal(r, "", false, "1.2.3.4", tlp); err == nil {
			t.Fatalf("expecting non-nil error")
		}
	}

	// invalid format for message size
	f("12foo bar")

	// too big message size
	f("123 foo")
	f("123456789 bar")
}
