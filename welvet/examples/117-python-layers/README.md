# 117. Python C-ABI — MHA place and prepare_grid

**Part:** XVII · Python C-ABI (welvet on PyPI)  
**Package:** `welvet` (PyPI)  
**Status:** ok — cpu · accel prepare

## When

Place a non-Dense layer from Python and remember when `prepare_grid` is required.

## Where

From `example/welvet/117-python-layers`:

```sh
pip install welvet==1.1.2
python3 main.py
```

## Why

`place(handle, kind, spec)` mirrors Go layer kinds. Some accelerators (`qnn`, `mlx`) need SyncTo* via `prepare_grid` after place and before forward.

## What

Lists layer types, places MHA with a shaped forward, and documents the prepare call for qnn/mlx.

## Requirements

Same as chapter 114. Release assets: GitHub `cabi-v1.1.2` · PyPI `welvet==1.1.2`.

## Captured sample

```text
layer_types 23 sample ['Dense', 'MultiHeadAttention', 'SwiGLU', 'RMSNorm', 'LayerNorm', 'Embedding', 'Softmax', 'Sequential']
mha_out_len 32
prepare_grid call after place_* when backend is qnn or mlx
```
