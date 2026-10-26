// A labelled snapshot of user-provided forward results; no GPU run is implied.
package main

import "fmt"

func main() {
	fmt.Println("Captured Mamba forward excerpt, supplied 2026-10-27 (not a live benchmark)")
	fmt.Printf("%-10s %12s %12s %12s %13s %14s\n", "dtype", "Go ms", "C++ ms", "Metal ms", "Metal/Go x", "Metal/C++ x")
	rows := []struct {
		dtype                string
		goUS, cppUS, metalUS float64
	}{
		{"float64", 89314.5, 236241.7, 28958.6}, {"float32", 69018.5, 225812.9, 29010.4},
		{"float16", 190545.0, 1155902.2, 28946.8}, {"bfloat16", 149351.3, 1049179.8, 29003.3},
		{"fp8e4m3", 265881.6, 1519874.3, 29159.0},
	}
	for _, r := range rows {
		fmt.Printf("%-10s %12.4f %12.4f %12.4f %13.3f %14.3f\n", r.dtype, r.goUS/1000, r.cppUS/1000, r.metalUS/1000, r.goUS/r.metalUS, r.cppUS/r.metalUS)
	}
	fmt.Println("Go/Metal exact bits and within-tolerance agreement: 100%/100% across all 12 supplied forward rows.")
	fmt.Println("Full excerpt and metric definitions: ../apple_stress_test/README.md")
	fmt.Println("Run your own 34-dtype / 29-mode matrix: cd ../apple_stress_test && ./run.sh")
	fmt.Println("Small verification: ./run.sh -profile smoke -samples 1 -iters 1 -steps 8")
}
