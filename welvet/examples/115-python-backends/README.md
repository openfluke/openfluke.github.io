# 115. Python C-ABI — list backends and pick one

**Part:** XVII · Python C-ABI (welvet on PyPI)  
**Package:** `welvet` (PyPI)  
**Status:** ok — accelerators soft-SKIP without device

## When

Probe which backends this host can run, then create a grid on the first usable one.

## Where

From `example/welvet/115-python-backends`:

```sh
pip install welvet==1.1.2
python3 main.py
```

## Why

`list_backends()` reports library presence; a device may still be missing. Prefer probe → try → soft-SKIP → fall through to CPU.

## What

Prints every backend id, then tries cuda / metal / vulkan / rocm / native_cpu / simd / cpu_tiled until one forward succeeds.

## Requirements

Same as chapter 114. Missing GPU libs print `--`; create+forward failures print `SKIP …` and continue.

## Captured sample

```text
backends
  ok cpu_tiled
  ok simd
  ok webgpu
  ok native_cpu
  ok cuda
  -- metal
  ok rocm
  ok vulkan
  -- qnn
  -- mlx
SKIP cuda: forward … no CUDA-capable device is detected
SKIP vulkan: …
SKIP rocm: …
chosen native_cpu
forward_len 16 head […]
```

(Exact SKIP lines depend on the host; a machine with a live NVIDIA GPU will choose `cuda`.)
