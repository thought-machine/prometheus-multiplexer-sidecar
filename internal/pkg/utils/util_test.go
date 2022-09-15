package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateContainerToPortMap(t *testing.T) {

	testCases := []struct {
		name          string
		flags         []string
		errString     string
		resultMapping map[string]int
	}{
		{"returns an error when input is nil",
			nil,
			"received an empty container",
			nil,
		},
		{
			"returns an error when input is empty",
			[]string{},
			"received an empty container",
			nil,
		},
		{
			"returns an error when input doesn't have colon",
			[]string{"port124"},
			"incorrect number of elements",
			nil,
		},

		{"returns an error when input doesn't have container name",
			[]string{":1234"},
			"missing container name",
			nil,
		},

		{"returns an error when input doesn't have port value",
			[]string{"port:"},
			"missing port value",
			nil,
		},
		{"returns an error when input has non-integer port value",
			[]string{"port:123w"},
			"failed to parse string to int",
			nil,
		},
		{"generates correct mapping with correct inputs",
			[]string{"port1:123", "port2:456", "portYeah:114534"},
			"",
			map[string]int{
				"port1":    123,
				"port2":    456,
				"portYeah": 114534,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := GenerateContainerToPortMap(tc.flags)
			assert.Equal(t, tc.resultMapping, m)
			if len(tc.errString) != 0 {
				assert.ErrorContains(t, err, tc.errString)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGzipToCompressData(t *testing.T) {
	var testCases = []struct {
		name           string
		content        []byte
		expectedResult []byte
	}{
		{"decompresses gzip formatted data successfully",
			[]byte{0x1f, 0x8b, 0x8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0x4a, 0x4c, 0x4a, 0x4e, 0x49, 0x4d, 0x33, 0x34, 0x34, 0x4c, 0x4c, 0x4c, 0x4c, 0x0, 0x0, 0x0, 0x0, 0xff, 0xff, 0x1, 0x0, 0x0, 0xff, 0xff, 0x3d, 0x6f, 0x9a, 0xaf, 0xd, 0x0, 0x0, 0x0},
			[]byte("abcdef111aaa`"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := GzipToCompressData(tc.content)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestCompressDataToGzip(t *testing.T) {
	var testCases = []struct {
		name           string
		content        []byte
		expectedResult []byte
	}{
		{"compresses gzip formatted data successfully",
			[]byte("abcdef111aaa`"),
			[]byte{0x1f, 0x8b, 0x8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0x4a, 0x4c, 0x4a, 0x4e, 0x49, 0x4d, 0x33, 0x34, 0x34, 0x4c, 0x4c, 0x4c, 0x4c, 0x0, 0x0, 0x0, 0x0, 0xff, 0xff, 0x1, 0x0, 0x0, 0xff, 0xff, 0x3d, 0x6f, 0x9a, 0xaf, 0xd, 0x0, 0x0, 0x0},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := CompressDataToGzip(tc.content)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestValidateLabelName(t *testing.T) {
	var testCases = []struct {
		name        string
		labelName   string
		expectedErr error
	}{
		{"test return nil when input label name is correct",
			"container_1",
			nil},

		{
			"test return an error when input label name is incorrect",
			"container-1",
			errInvalidLabelName,
		},
	}
	for _, tc := range testCases {
		err := ValidateLabelName(tc.labelName)
		t.Run(tc.name, func(t *testing.T) {
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}
}
