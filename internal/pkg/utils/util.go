package utils

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/prometheus/common/model"
)

var errInvalidLabelName = errors.New("received invalid label name")

// GenerateContainerToPortMap generates a map structure for the flag `container_to_port_map`.
func GenerateContainerToPortMap(containerPortList []string) (map[string]int, error) {
	if len(containerPortList) == 0 {
		return nil, fmt.Errorf("received an empty container:port list")
	}

	containerPortMap := make(map[string]int)
	for _, entry := range containerPortList {
		s := strings.Split(entry, ":")
		if len(s) != 2 {
			return nil, fmt.Errorf("failed to parse \"container_name\":\"port\" entry: %s: incorrect number of elements", entry)
		}
		if len(s[0]) == 0 {
			return nil, fmt.Errorf("missing container name for entry '%s'", entry)
		}

		if len(s[1]) == 0 {
			return nil, fmt.Errorf("missing port value for entry '%s'", entry)
		}
		var err error
		containerPortMap[s[0]], err = strconv.Atoi(s[1])
		if err != nil {
			return nil, fmt.Errorf("failed to parse string to int for entry %s: %w", entry, err)
		}
	}
	return containerPortMap, nil
}

// CompressDataToGzip zips the given byte slice using gzip encoding
func CompressDataToGzip(data []byte) (compressedData []byte) {
	if len(data) == 0 {
		return nil
	}
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write(data); err != nil {
		return nil
	}

	if err := gz.Flush(); err != nil {
		return nil
	}

	if err := gz.Close(); err != nil {
		return nil
	}
	compressedData = b.Bytes()
	return compressedData
}

// GzipToCompressData unzips the given gzip encoding byte slice.
func GzipToCompressData(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, nil
	}
	zipBuf := bytes.NewBuffer(data)
	r, err := gzip.NewReader(zipBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to create new gzip reader: %w", err)
	}

	var unzipBuf bytes.Buffer

	if _, err := unzipBuf.ReadFrom(r); err != nil {
		return nil, fmt.Errorf("failed to read from gzip reader: %w", err)
	}

	return unzipBuf.Bytes(), nil
}

// ValidateLabelName checks if the provided label name is valid.
func ValidateLabelName(labelName string) error {
	if !model.LabelName(labelName).IsValid() {
		return fmt.Errorf("invalid label name %s: %w", labelName, errInvalidLabelName)
	}
	return nil
}
