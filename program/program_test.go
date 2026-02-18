package program

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zenizh/go-capturer"
)

func TestOptions_Version(t *testing.T) {
	var program Options

	exitValue := -1
	exitFunc = func(x int) {
		exitValue = x
	}
	defer func() { exitFunc = os.Exit }()

	out := capturer.CaptureStdout(func() {
		_, err := program.Parse([]string{"--version", "list"})
		assert.NoError(t, err)
	})

	assert.Equal(t, 0, exitValue)
	assert.Equal(t, "unknown\n", out)
}

func TestOptions_ParseLogin(t *testing.T) {
	var program Options

	ctx, err := program.Parse([]string{"login", "--client-id=test", "--client-secret=secret"})
	assert.NoError(t, err)
	assert.Equal(t, "login", ctx.Command())
	assert.Equal(t, "test", program.ClientID)
	assert.Equal(t, "secret", program.ClientSecret)
	assert.Equal(t, 8085, program.CallbackPort)
}

func TestOptions_ParseList(t *testing.T) {
	var program Options

	ctx, err := program.Parse([]string{"list", "--from=2026-01-01", "--to=2026-02-01"})
	assert.NoError(t, err)
	assert.Equal(t, "list", ctx.Command())
	assert.Equal(t, "2026-01-01", program.List.From)
	assert.Equal(t, "2026-02-01", program.List.To)
}

func TestOptions_ParseDownload(t *testing.T) {
	var program Options

	ctx, err := program.Parse([]string{"download", "--output-dir=./out"})
	assert.NoError(t, err)
	assert.Equal(t, "download", ctx.Command())
	assert.Equal(t, "./out", program.Download.OutputDir)
}
