package tls_test

import (
	"errors"
	"io"
	"testing"

	tls "github.com/excitedplus1s/utlscm"
	"github.com/excitedplus1s/utlscm/dicttls"
)

func TestGREASEECHWrite(t *testing.T) {
	for _, testsuite := range []rawECHTestSuite{rawECH_HKDFSHA256_AES128GCM} {

		gech := &tls.GREASEEncryptedClientHelloExtension{}

		n, err := gech.Write(testsuite.raw[4:]) // skip the first 4 bytes which are the extension type and length
		if err != nil {
			t.Fatalf("Failed to write GREASE ECH extension: %s", err)
		}

		if n != len(testsuite.raw[4:]) {
			t.Fatalf("Failed to write all GREASE ECH extension bytes: %d != %d", n, len(testsuite.raw[4:]))
		}

		var gechBytes []byte = make([]byte, 1024)
		n, err = gech.Read(gechBytes)
		if err != nil && !errors.Is(err, io.EOF) {
			t.Fatalf("Failed to read GREASE ECH extension: %s", err)
		}

		if n != len(testsuite.raw) {
			t.Fatalf("GREASE ECH Read length mismatch: %d != %d", n, len(testsuite.raw))
		}

		// manually check fields in the GREASE ECH extension
		if len(gech.CandidateCipherSuites) != 1 ||
			gech.CandidateCipherSuites[0].KdfId != testsuite.kdfID ||
			gech.CandidateCipherSuites[0].AeadId != testsuite.aeadID {
			t.Fatalf("GREASE ECH Read cipher suite mismatch")
		}

		if len(gech.EncapsulatedKey) != int(testsuite.encapsulatedKeyLength) {
			t.Fatalf("GREASE ECH Read encapsulated key length mismatch")
		}

		if len(gech.CandidatePayloadLens) != 1 || gech.CandidatePayloadLens[0] != testsuite.payloadLength {
			t.Fatalf("GREASE ECH Read payload length mismatch")
		}
	}
}

type rawECHTestSuite struct {
	kdfID                 uint16
	aeadID                uint16
	encapsulatedKeyLength uint16
	payloadLength         uint16

	raw []byte
}

var (
	rawECH_HKDFSHA256_AES128GCM rawECHTestSuite = rawECHTestSuite{
		kdfID:                 dicttls.HKDF_SHA256,
		aeadID:                dicttls.AEAD_AES_128_GCM,
		encapsulatedKeyLength: 32,
		payloadLength:         208 - 16,
		raw: []byte{
			0xfe, 0x0d, 0x00, 0xfa, 0x00, 0x00, 0x01, 0x00,
			0x01, 0x77, 0x00, 0x20, 0x3d, 0x3e, 0xe0, 0xa6,
			0x1f, 0x46, 0x4f, 0x89, 0x5f, 0x39, 0x4a, 0xfd,
			0x6e, 0xbc, 0x7f, 0x4e, 0xe2, 0x5a, 0xdc, 0x4e,
			0xda, 0x9a, 0x9f, 0x5f, 0x2b, 0xf5, 0x21, 0x0e,
			0xc6, 0x33, 0x64, 0x32, 0x00, 0xd0, 0xae, 0xff,
			0x25, 0xd6, 0x4a, 0x23, 0x3a, 0x13, 0x5b, 0xdc,
			0xe4, 0xaf, 0x6c, 0xb8, 0xaf, 0x66, 0x57, 0xbd,
			0x44, 0x2d, 0xca, 0xb6, 0xbb, 0xaf, 0xda, 0x8a,
			0x6b, 0x12, 0xb2, 0x42, 0xf1, 0x3d, 0xf6, 0x26,
			0xd4, 0x82, 0x30, 0x40, 0xd4, 0x53, 0x06, 0x7c,
			0xf1, 0x10, 0xf3, 0x80, 0x16, 0x95, 0xa7, 0xfb,
			0x08, 0x76, 0x82, 0x85, 0x86, 0xb4, 0x3a, 0x7b,
			0xea, 0xfb, 0xaa, 0xc3, 0xe0, 0x51, 0xcf, 0x42,
			0xf6, 0xa0, 0x15, 0x0e, 0x26, 0x4d, 0x37, 0x35,
			0x95, 0x4d, 0xce, 0xf6, 0xd6, 0x58, 0x78, 0x67,
			0x42, 0xd3, 0xc6, 0xac, 0xb5, 0xe9, 0x3e, 0xb6,
			0x02, 0x87, 0x66, 0xb3, 0xb2, 0x56, 0x99, 0xb2,
			0xdb, 0x8c, 0x3b, 0x04, 0xf1, 0x7c, 0x85, 0x5b,
			0xc3, 0x93, 0x8e, 0xdb, 0x5d, 0x87, 0x66, 0xfb,
			0x66, 0x54, 0xf3, 0xec, 0x25, 0xe5, 0x70, 0x3c,
			0xd5, 0x0e, 0x8e, 0xd5, 0xd2, 0xbb, 0x24, 0x2b,
			0xb5, 0x01, 0xa0, 0x5e, 0xba, 0x45, 0xaf, 0x68,
			0x96, 0x8a, 0x83, 0x90, 0x20, 0x5b, 0x8c, 0x7d,
			0x24, 0x00, 0x2f, 0x08, 0x7f, 0x29, 0x8c, 0x32,
			0x5e, 0x57, 0xb5, 0x64, 0xaa, 0x0b, 0xf4, 0x42,
			0x54, 0xdc, 0xe5, 0xd4, 0x08, 0xf4, 0x4d, 0x27,
			0x5d, 0x90, 0x52, 0x32, 0x22, 0xc8, 0xb6, 0xd8,
			0x80, 0xa6, 0x30, 0xa0, 0x20, 0x98, 0x2c, 0x0b,
			0x3e, 0x55, 0x4a, 0x09, 0xa9, 0x09, 0xa4, 0x99,
			0x89, 0x02, 0x6e, 0xab, 0xe3, 0xa1, 0xe9, 0xb8,
			0x58, 0x20, 0xcc, 0xc8, 0xb0, 0x73,
		},
	}
)
