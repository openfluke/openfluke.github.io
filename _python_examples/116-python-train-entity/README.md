# 116. Python C-ABI — train SGD and .entity round-trip

**Part:** XVII · Python C-ABI (welvet on PyPI)  
**Package:** `welvet` (PyPI)  
**Status:** ok — cpu_tiled

## When

Train one Dense step from Python and checkpoint with the same .entity bytes Go uses.

## Where

From `example/welvet/116-python-train-entity`:

```sh
pip install welvet==1.1.2
python3 main.py
```

## Why

Python scripts and Go apps should share checkpoints. `serialize_entity` / `deserialize_entity` are the CABI twins of Go ENTITY IO.

## What

Runs `train_sgd`, serializes the grid, deserializes a second handle, forwards once, and frees both.

## Requirements

Same as chapter 114.

## Captured sample

```text
train {'loss': 0.25, 'status': 'ok'}
entity_bytes 550
reloaded_out [0.0, 0.0, 0.0, 0.0127]
```
