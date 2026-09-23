package smscodec

import (
	"testing"

	"github.com/warthog618/sms/encoding/tpdu"
	"github.com/warthog618/sms/encoding/ucs2"
)

func TestBuildSubmitTPDUsWithOptionsForcesUCS2(t *testing.T) {
	tpdus, _, err := BuildSubmitTPDUsWithOptions("10086", "hello", SubmitOptions{Encoding: SMSEncodingUCS2})
	if err != nil {
		t.Fatalf("BuildSubmitTPDUsWithOptions() error = %v", err)
	}
	if len(tpdus) != 1 {
		t.Fatalf("parts=%d want 1", len(tpdus))
	}

	pdu := &tpdu.TPDU{Direction: tpdu.MO}
	if err := pdu.UnmarshalBinary(tpdus[0]); err != nil {
		t.Fatalf("UnmarshalBinary() error = %v", err)
	}
	if pdu.DCS != tpdu.DcsUCS2Data {
		t.Fatalf("DCS=0x%02x want 0x%02x", byte(pdu.DCS), byte(tpdu.DcsUCS2Data))
	}
	if got, want := []byte(pdu.UD), ucs2.Encode([]rune("hello")); string(got) != string(want) {
		t.Fatalf("UD=%x want UCS2 %x", got, want)
	}
}

func TestBuildSubmitTPDUsKeepsAutoEncodingByDefault(t *testing.T) {
	tpdus, _, err := BuildSubmitTPDUsWithOptions("10086", "hello", SubmitOptions{})
	if err != nil {
		t.Fatalf("BuildSubmitTPDUs() error = %v", err)
	}
	if len(tpdus) != 1 {
		t.Fatalf("parts=%d want 1", len(tpdus))
	}

	pdu := &tpdu.TPDU{Direction: tpdu.MO}
	if err := pdu.UnmarshalBinary(tpdus[0]); err != nil {
		t.Fatalf("UnmarshalBinary() error = %v", err)
	}
	if pdu.DCS != 0x00 {
		t.Fatalf("DCS=0x%02x want auto GSM7 0x00", byte(pdu.DCS))
	}
}

func TestNormalizeDestinationNumber(t *testing.T) {
	tests := map[string]string{
		"13812345678":       "+8613812345678",
		"138 1234 5678":     "+8613812345678",
		"138-1234-5678":     "+8613812345678",
		"+86 138 1234 5678": "+8613812345678",
		"8613812345678":     "+8613812345678",
		"008613812345678":   "+8613812345678",
		"10086":             "10086",
		"+12025550123":      "+12025550123",
	}
	for input, want := range tests {
		if got := NormalizeDestinationNumber(input); got != want {
			t.Errorf("NormalizeDestinationNumber(%q)=%q want %q", input, got, want)
		}
	}
}

func TestBuildSubmitTPDUsNormalizesMainlandMobileNumber(t *testing.T) {
	tpdus, _, err := BuildSubmitTPDUsWithOptions("138 1234 5678", "test", SubmitOptions{})
	if err != nil {
		t.Fatalf("BuildSubmitTPDUsWithOptions() error = %v", err)
	}
	if len(tpdus) != 1 {
		t.Fatalf("parts=%d want 1", len(tpdus))
	}
	pdu := &tpdu.TPDU{Direction: tpdu.MO}
	if err := pdu.UnmarshalBinary(tpdus[0]); err != nil {
		t.Fatalf("UnmarshalBinary() error = %v", err)
	}
	if got, want := pdu.DA.Number(), "+8613812345678"; got != want {
		t.Fatalf("destination=%q want %q", got, want)
	}
}
