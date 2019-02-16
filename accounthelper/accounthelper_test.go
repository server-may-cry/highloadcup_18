package accounthelper_test

import (
	"testing"

	"github.com/server-may-cry/highloadcup_18/accounthelper"
)

func TestExtractDomain(t *testing.T) {
	var emailDomainTestData = []struct {
		in  string
		out string
	}{
		{"bla@bla.bla", "bla.bla"},
		//{"\"bla@bla\"@b.b", "b.b"},
	}
	for _, testData := range emailDomainTestData {
		out := accounthelper.ExtractDomain(testData.in)
		if out != testData.out {
			t.Errorf("got %q, want %q", out, testData.out)
		}
	}
}

func TestExtractPhoneCode(t *testing.T) {
	out := accounthelper.ExtractPhoneCode("0(123)456")
	expected := "123"
	if out != expected {
		t.Errorf("got %q, want %q", out, expected)
	}
	out = accounthelper.ExtractPhoneCode("")
	expected = ""
	if out != expected {
		t.Errorf("got %q, want %q", out, expected)
	}
}

func TestExtractYearOfBirth(t *testing.T) {
	out := accounthelper.ExtractYearOfBirth(1545503558)
	expected := int16(2018)
	if out != expected {
		t.Errorf("got %d, want %d", out, expected)
	}
}

func TestPhoneToInt(t *testing.T) {
	out := accounthelper.PhoneToInt("1(234)5678901")
	expected := int64(12345678901)
	if out != expected {
		t.Errorf("got %d, want %d", out, expected)
	}
}

func TestPhoneToString(t *testing.T) {
	out := accounthelper.PhoneToString(12345678901)
	expected := "1(234)5678901"
	if out != expected {
		t.Errorf("got %s, want %s", out, expected)
	}
}
