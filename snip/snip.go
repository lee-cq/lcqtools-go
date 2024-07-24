package snip

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

func Snip(ef string) error {
	// dir := os.TempDir()
	cmd := exec.Command(ef, "snip", "-o", path.Join(os.TempDir(), "tmp-lcqtools-ocr-snip.png"))
	if err := cmd.Run(); err != nil {
		return err
	}
	if cmd.ProcessState.ExitCode() == 0 {
		return nil
	}
	return fmt.Errorf("")
}
