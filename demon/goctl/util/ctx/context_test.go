package ctx

import (
	"fmt"
	"testing"
)
var (
	protoDir = "../../genproject/parser/protocol"
)
func TestPrepare(t *testing.T) {
	resp, err := Prepare(protoDir)
	if err != nil {
		t.Errorf("Prepare error: %v", err)
		return
	}
	fmt.Printf("Prepare return projectContext: %v", resp)
}
