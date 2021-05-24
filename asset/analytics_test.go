package asset

import (
	"testing"
)

func TestAsset_Analytics_EncodeVersion(t *testing.T) {
	s, err := encodeVersion("1.24.0")
	if err != nil || s != "Alh" {
		t.Fatal("Invalid signature")
	}

	s, err = encodeVersion("12.0")
	if err != nil || s != "AM" {
		t.Fatal("Invalid signature")
	}

	s, err = encodeVersion("43.21.26")
	if err != nil || s != "///" {
		t.Fatal("Invalid signature")
	}

	s, err = encodeVersion("0.0.0")
	if err != nil || s != "AAA" {
		t.Fatal("Invalid signature")
	}
}

func TestAsset_Analytics_InvalidVersion(t *testing.T) {
	_, err := encodeVersion("44.45.46")
	if err == nil || err.Error() != "version must be smaller than 43.21.26" {
		t.Fatal("Error expected")
	}
}

func TestAsset_Analytics_Signature(t *testing.T) {
	// FIXME: use setUp/tearDown
	sdkVersionBackup := sdkVersion
	techVersionBackup := techVersion

	sdkVersion = "2.3.4"
	techVersion = "8.0"

	analyticsSignature = sdkAnalyticsSignature()

	// FIXME: use setUp/tearDown
	sdkVersion = sdkVersionBackup
	techVersion = techVersionBackup

	if analyticsSignature != "AQJ1uAI" {
		t.Fatal("Invalid signature", analyticsSignature)
	}
}
