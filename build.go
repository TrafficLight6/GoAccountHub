//go:build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/magefile/mage/sh"
)

const (
	appName     = "GoAccountHub"
	versionFile = "GAHversion"
	outDir      = "bin"
	frontendDir = "GAHFrontend"
)

// platforms is the list of platforms to build for.
var platforms = []struct {
	GOOS   string
	GOARCH string
}{
	{"windows", "amd64"},
	{"windows", "arm64"},
	{"linux", "amd64"},
	{"linux", "arm64"},
	{"darwin", "amd64"},
	{"darwin", "arm64"},
}

func Build() error {
	out := "gah"
	if runtime.GOOS == "windows" {
		out += ".exe"
	}
	fmt.Printf("Building %s for %s/%s...\n", out, runtime.GOOS, runtime.GOARCH)
	return sh.Run("go", "build", "-o", out, ".")
}

// All build the frontend and embed the dist into the native binary.
func All() error {
	if err := buildFrontend(); err != nil {
		return err
	}
	return Build()
}

// AllCross build the frontend and embed the dist into cross-built binaries.
func AllCross() error {
	if err := buildFrontend(); err != nil {
		return err
	}
	return Cross()
}

// buildFrontend builds the frontend and makes sure the output is really there.
func buildFrontend() error {
	if err := Frontend(); err != nil {
		return err
	}
	index := filepath.Join(frontendDir, "dist", "index.html")
	if _, err := os.Stat(index); err != nil {
		return fmt.Errorf("frontend build output missing, %s not found: %w", index, err)
	}
	return nil
}

// Frontend build the frontend and embed the dist into the binary.
func Frontend() error {
	fmt.Println("Installing frontend dependencies...")
	if err := sh.Run("npm", "ci", "--prefix", frontendDir); err != nil {
		return err
	}
	fmt.Println("Building frontend...")
	return sh.Run("npm", "run", "build", "--prefix", frontendDir)
}

// Cross build all supported platforms.
func Cross() error {
	version, err := readVersion()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	for _, p := range platforms {
		ext := ""
		if p.GOOS == "windows" {
			ext = ".exe"
		}
		out := filepath.Join(outDir, fmt.Sprintf("%s-%s-%s-%s%s", appName, version, p.GOOS, p.GOARCH, ext))
		fmt.Printf("Building %s...\n", out)
		if err := crossBuild(out, p.GOOS, p.GOARCH); err != nil {
			return err
		}
	}
	fmt.Printf("Done, %d binaries written to %s/\n", len(platforms), outDir)
	return nil
}

// crossBuild cross build binary for a specific platform.
// Avoids depending on the caller's environment.
func crossBuild(out, goos, goarch string) error {
	env := map[string]string{
		"CGO_ENABLED": "0",
		"GOOS":        goos,
		"GOARCH":      goarch,
	}
	return sh.RunWith(env, "go", "build", "-trimpath", "-ldflags", "-s -w", "-o", out, ".")
}

func readVersion() (string, error) {
	b, err := os.ReadFile(versionFile)
	if err != nil {
		return "", fmt.Errorf("read %s: %w", versionFile, err)
	}
	return strings.TrimSpace(string(b)), nil
}
