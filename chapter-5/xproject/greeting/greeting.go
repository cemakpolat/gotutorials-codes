package greeting

import (
	"fmt"
	"io"
	"os"
)

// SayHello prints a greeting message to the console
func SayHello(name string) {
	fmt.Printf("Hello, %s!\n", name)
}

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		os.Stdout = old
	}()
	f()
	w.Close()
	out, _ := io.ReadAll(r)
	return string(out)
}
