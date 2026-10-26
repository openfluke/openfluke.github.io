# 88. Apple Metal — Mamba results and stress runner

**Part:** XII · Apple acceleration (welvet.metal)  
**Package:** `apple_stress_test`    
**Status:** ok — labelled captured snapshot; no GPU required

## When

Read a captured result excerpt or find the runnable full stress comparison.

## Where

From `example/welvet/88-apple-mamba-results`:

```sh
go run .
cd ../apple_stress_test && ./run.sh
```

## Why

Distinguish a supplied results snapshot from a live measurement.

## What

Prints five rows from the user-provided 2026-10-27 Mamba forward excerpt. Prints commands for the 34-dtype × 29-mode Go/C++/Metal stress runner. This chapter is host-only and does not execute Metal.

These are usage examples with small shapes. For speed measurements, use
[`../apple_stress_test`](../apple_stress_test/README.md). The results are numerical
CPU/GPU comparisons, not model prediction accuracy.

## Captured sample

```text
Captured Mamba forward excerpt, supplied 2026-10-27 (not a live benchmark)
dtype             Go ms       C++ ms     Metal ms    Metal/Go x    Metal/C++ x
float64         89.3145     236.2417      28.9586         3.084          8.158
float32         69.0185     225.8129      29.0104         2.379          7.784
float16        190.5450    1155.9022      28.9468         6.583         39.932
bfloat16       149.3513    1049.1798      29.0033         5.149         36.174
fp8e4m3        265.8816    1519.8743      29.1590         9.118         52.124
Go/Metal exact bits and within-tolerance agreement: 100%/100% across all 12 supplied forward rows.
Full excerpt and metric definitions: ../apple_stress_test/README.md
Run your own 34-dtype / 29-mode matrix: cd ../apple_stress_test && ./run.sh
Small verification: ./run.sh -profile smoke -samples 1 -iters 1 -steps 8
```
