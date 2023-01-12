package logger

import "testing"

func TestPrint(t *testing.T) {
	log, _ := InitJSONLogger("./", true)
	log.Infow("a", "b", "c")
}
