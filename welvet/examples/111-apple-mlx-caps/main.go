// Print MLX caps and demonstrate RequireLayer rejects.
package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/openfluke/welvet/nativemlx"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if runtime.GOOS != "darwin" || runtime.GOARCH != "arm64" {
		fmt.Println("SKIP: Apple MLX needs Apple silicon macOS and cgo")
		return nil
	}
	if !nativemlx.Available() {
		fmt.Println("SKIP: need libwelvet_mlx.dylib + Metal GPU")
		return nil
	}

	layers := nativemlx.SupportedLayers()
	fmt.Println("version", nativemlx.Version())
	fmt.Println("supported_layers", strings.Join(layers, ","))
	fmt.Println("dense/float32", nativemlx.LayerSupported("dense", "float32"))
	fmt.Println("dense/float64", nativemlx.LayerSupported("dense", "float64"))

	if err := nativemlx.RequireLayer("dense", "float32"); err != nil {
		return fmt.Errorf("expected dense/float32 ok: %w", err)
	}
	fmt.Println("RequireLayer dense/float32: ok")

	err := nativemlx.RequireLayer("dense", "float64")
	if err == nil {
		return fmt.Errorf("expected dense/float64 reject")
	}
	fmt.Println("RequireLayer dense/float64:", err)
	return nil
}
