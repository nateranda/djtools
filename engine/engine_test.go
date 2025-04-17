package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQUncompress(t *testing.T) {
	// valid blob
	compBlob := []byte{0, 0, 0, 138, 120, 156, 115, 120, 218, 209, 192, 0, 4, 142, 81, 231, 189, 29, 128, 52, 35, 3, 4, 48, 73, 37, 238, 92, 46, 120, 233, 227, 129, 63, 255, 33, 160, 13, 36, 147, 192, 192, 144, 217, 242, 116, 222, 190, 123, 81, 142, 77, 48, 149, 72, 154, 40, 208, 199, 192, 0, 0, 75, 220, 41, 20}
	uncompBlob := []byte{64, 229, 136, 128, 0, 0, 0, 0, 65, 90, 207, 75, 64, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 2, 26, 97, 185, 167, 17, 210, 241, 192, 252, 255, 255, 255, 255, 255, 255, 255, 134, 1, 0, 0, 0, 96, 0, 0, 105, 132, 229, 158, 190, 222, 90, 65, 130, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 26, 97, 185, 167, 17, 210, 241, 192, 252, 255, 255, 255, 255, 255, 255, 255, 134, 1, 0, 0, 0, 96, 0, 0, 105, 132, 229, 158, 190, 222, 90, 65, 130, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	blob, err := qUncompress(compBlob)
	assert.Nil(t, err, "Error on valid blob should be nil")
	assert.Equal(t, uncompBlob, blob, "These blobs should be equal.")

	// blob under 5 bytes
	compBlob = []byte{0, 0, 0, 0}
	_, err = qUncompress(compBlob)
	assert.NotNil(t, err, "Blob of under 5 bytes should throw an error.")

	// corrupted blob
	compBlob = []byte{0, 0, 0, 138, 121, 156, 115, 120, 218, 209, 192, 0, 4, 142, 81, 231, 189, 29, 128, 52, 35, 3, 4, 48, 73, 37, 238, 92, 46, 120, 233, 227, 129, 63, 255, 33, 160, 13, 36, 147, 192, 192, 144, 217, 242, 116, 222, 190, 123, 81, 142, 77, 48, 149, 72, 154, 40, 208, 199, 192, 0, 0, 75, 220, 41, 20}
	_, err = qUncompress(compBlob)
	assert.NotNil(t, err, "Corrupted blob should throw an error.")

	// invalid length
	compBlob = []byte{0, 0, 0, 139, 120, 156, 115, 120, 218, 209, 192, 0, 4, 142, 81, 231, 189, 29, 128, 52, 35, 3, 4, 48, 73, 37, 238, 92, 46, 120, 233, 227, 129, 63, 255, 33, 160, 13, 36, 147, 192, 192, 144, 217, 242, 116, 222, 190, 123, 81, 142, 77, 48, 149, 72, 154, 40, 208, 199, 192, 0, 0, 75, 220, 41, 20}
	_, err = qUncompress(compBlob)
	assert.NotNil(t, err, "Blob with invalid length should throw an error.")
}
