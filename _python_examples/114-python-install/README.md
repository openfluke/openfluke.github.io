# 114. Python C-ABI — install and Dense forward

**Part:** XVII · Python C-ABI (welvet on PyPI)  
**Package:** `welvet` (PyPI)  
**Status:** ok — needs Linux x86_64 or macOS arm64 wheel

## When

Install Welvet for Python and run a first Dense forward without a Go toolchain.

## Where

From `example/welvet/114-python-install`:

```sh
pip install welvet==1.1.2
python3 main.py
```

Editable from the monorepo:

```sh
PYTHONPATH=../../w2a/python/src python3 main.py
```

## Why

See the Python package / CABI contract version and the Go engine pin, then place a Dense layer on `cpu_tiled`.

## What

Asserts the CABI version, creates a grid, places Dense 8→4, and prints one forward.

## Requirements

Python 3.10+, `welvet==1.1.2` from PyPI (ships `welvet.so` / `welvet.dylib`). Go engine behind the ABI: **v1.7.3**.

## Captured sample

```text
python_package 1.1.2
cabi_contract 1.1.2
go_engine v1.7.3
input [0.05, 0.1, 0.15, 0.2, 0.25, 0.3, 0.35, 0.4]
output [0.0, 0.0, 0.0, 0.0]
```
