package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func stdinIsTTY() bool {
	fi, err := os.Stdin.Stat()
	return err == nil && fi.Mode()&os.ModeCharDevice != 0
}

func promptLine(in io.Reader, out io.Writer, question, def string) (string, error) {
	if def != "" {
		fmt.Fprintf(out, "%s [%s] ", question, def)
	} else {
		fmt.Fprintf(out, "%s ", question)
	}
	line, err := bufio.NewReader(in).ReadString('\n')
	if err != nil && len(strings.TrimSpace(line)) == 0 {
		return "", err
	}
	ans := strings.TrimSpace(line)
	if ans == "" {
		return def, nil
	}
	return ans, nil
}

func templateOverrideEnv(recipeID string) string {
	key := strings.ToUpper(strings.ReplaceAll(recipeID, "/", "_"))
	return "SBX_" + key + "_TEMPLATE"
}
