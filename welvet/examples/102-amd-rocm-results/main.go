// Pointers to ROCm validation matrices; no GPU run is implied.
package main

import "fmt"

func main() {
	fmt.Println("AMD ROCm results pointers (not a live benchmark)")
	fmt.Println()
	fmt.Println("Engine pin: welvet.rocm 0.79 · model ABI 1 · Go binding: nativerocm (Welvet v1.5.0)")
	fmt.Println("Usage tour: 96 mount · 97 Mamba dtypes · 98 train modes · 99 SyncToGo/ENTITY · 100 JSON · 101 mixed")
	fmt.Println()
	fmt.Println("Live matrices (sibling checkouts):")
	fmt.Println("  cd ../../w2a && go test ./tests/go_v_rocm -count=1 -timeout 180m")
	fmt.Println("  cd ../../welvet.rocm/examples/rocm_smoke && go run .")
	fmt.Println()
	fmt.Println("Install shared library into Go:")
	fmt.Println("  cd ../../welvet.rocm && ./scripts/linux_install_go.sh")
	fmt.Println()
	fmt.Println("Small usage tour without the large matrix:")
	fmt.Println("  ../run_amd.sh")
}
