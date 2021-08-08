package shamir

import "testing"

func TestTables(t *testing.T) {
	for i := 1; i < 256; i++ {
		logV := logTable[i]
		expV := expTable[logV]
		assert(t, expV == uint8(i), "bad: %d log: %d exp: %d", i, logV, expV)
	}
}
