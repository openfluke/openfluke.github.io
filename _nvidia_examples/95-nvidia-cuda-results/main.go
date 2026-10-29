// Pointers to CUDA validation matrices; no GPU run is implied.
package main

import "fmt"

func main() {
	fmt.Println("NVIDIA CUDA results pointers (not a live benchmark)")
	fmt.Println()
	fmt.Println("Engine pin: welvet.cuda 0.79 · model ABI 1 · Go binding: nativecuda (Welvet v1.4.0)")
	fmt.Println("Usage tour: 89 mount · 90 Mamba dtypes · 91 train modes · 92 SyncToGo/ENTITY · 93 JSON · 94 mixed")
	fmt.Println()
	fmt.Println("Live matrices (sibling checkouts):")
	fmt.Println("  cd ../../w2a && go test ./tests/go_v_cuda -count=1 -timeout 180m")
	fmt.Println("  cd ../../welvet.cuda/examples/cuda_smoke && go run .")
	fmt.Println()
	fmt.Println("Install shared library into Go:")
	fmt.Println("  cd ../../welvet.cuda && ./scripts/linux_install_go.sh")
	fmt.Println()
	fmt.Println("Small usage tour without the large matrix:")
	fmt.Println("  ../run_nvidia.sh")
}
