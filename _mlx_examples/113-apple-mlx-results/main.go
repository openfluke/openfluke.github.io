// Pointers to MLX validation matrices; no GPU run is implied.
package main

import "fmt"

func main() {
	fmt.Println("Apple MLX results pointers (not a live benchmark)")
	fmt.Println()
	fmt.Println("Engine pin: welvet.mlx 1.0.0 · Go binding: nativemlx (Welvet v1.7.0)")
	fmt.Println("Usage tour: 110 SyncToMLX · 111 caps · 112 Dense/Sequential · 113 pointers")
	fmt.Println()
	fmt.Println("Live matrices (sibling checkouts):")
	fmt.Println("  cd ../../w2a && go test ./tests/go_v_mlx -count=1 -timeout 60m")
	fmt.Println("  cd ../../welvet.mlx && ./scripts/macos_build_w2a.sh")
	fmt.Println()
	fmt.Println("Install shared library into Go:")
	fmt.Println("  cd ../../welvet.mlx && ./scripts/macos_install_go.sh")
	fmt.Println()
	fmt.Println("Small usage tour without the large matrix:")
	fmt.Println("  ../run_mlx.sh")
}
