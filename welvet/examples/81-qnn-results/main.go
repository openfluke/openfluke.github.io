package main

import "fmt"

// Snapshot of w2a/suites/go_v_qnn/RESULTS.md (2026-10-23, Snapdragon X).
// This example is host-only — no DLL required.
func main() {
	fmt.Println("welvet.qnn 1.0.0 · go_v_qnn large matrix")
	fmt.Println("host: Windows ARM64 · shapes: dense 32x512x512 …")
	fmt.Println()
	fmt.Println("NPU DLL wins (subset vs Go f32):")
	rows := []struct{ layer, dt, goT, npu, pct string }{
		{"dense", "u8", "3.05ms", "164µs", "95%"},
		{"dense", "f32", "3.05ms", "201µs", "93%"},
		{"softmax", "f16", "4.31ms", "446µs", "90%"},
		{"rmsnorm", "u8", "9.40ms", "897µs", "90%"},
		{"layernorm", "f16", "6.82ms", "553µs", "92%"},
		{"cnn2", "i8", "12.56ms", "763µs", "94%"},
		{"swiglu", "u16", "1.92ms", "352µs", "82%"},
		{"mha", "f16", "4.51ms", "372µs", "92%"},
	}
	for _, r := range rows {
		fmt.Printf("  %-10s %-4s  Go %s  NPU %s  → %s\n", r.layer, r.dt, r.goT, r.npu, r.pct)
	}
	fmt.Println()
	fmt.Println("note: call timer ≫ DLL exec on online layers (compose each call)")
	fmt.Println("fail: lstm f32/f16 execute code=2")
	fmt.Println("full: w2a/suites/go_v_qnn/RESULTS.md")
}
