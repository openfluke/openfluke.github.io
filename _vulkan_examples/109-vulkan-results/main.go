// Pointers to Vulkan validation matrices; no GPU run is implied.
package main

import "fmt"

func main() {
	fmt.Println("Vulkan results pointers (not a live benchmark)")
	fmt.Println()
	fmt.Println("Engine pin: welvet.vulkan 0.79 · model ABI 1 · Go binding: nativevulkan (Welvet v1.6.0)")
	fmt.Println("Usage tour: 103 mount · 104 Mamba dtypes · 105 train modes · 106 SyncToGo/ENTITY · 107 JSON · 108 mixed")
	fmt.Println()
	fmt.Println("Live matrices (sibling checkouts):")
	fmt.Println("  cd ../../w2a && go test ./tests/go_v_vulkan -count=1 -timeout 180m")
	fmt.Println("  cd ../../welvet.vulkan && ./scripts/build_linux.sh && ./build/example_vulkan_smoke")
	fmt.Println()
	fmt.Println("Install shared library into Go:")
	fmt.Println("  cd ../../welvet.vulkan && ./scripts/linux_install_go.sh")
	fmt.Println()
	fmt.Println("Small usage tour without the large matrix:")
	fmt.Println("  ../run_vulkan.sh")
}
