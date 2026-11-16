"""115. Python C-ABI — list backends and pick one."""
from __future__ import annotations

import welvet as w

# Prefer a real accelerator when the sibling lib + device are present.
PREFER = ("cuda", "metal", "vulkan", "rocm", "native_cpu", "simd", "cpu_tiled")


def try_forward(backend: str) -> bool:
    g = w.create_grid(backend=backend)
    try:
        w.place_dense(
            g,
            {"in": 16, "out": 16, "act": "relu", "dtype": w.DType.FLOAT32},
        )
        out = w.forward(g, [0.05] * 16)
        print("chosen", backend)
        print("forward_len", len(out), "head", [round(v, 4) for v in out[:4]])
        return True
    except RuntimeError as exc:
        print(f"SKIP {backend}: {exc}")
        return False
    finally:
        w.free_grid(g)


def main() -> None:
    w.assert_engine_version()
    backends = w.list_backends()
    print("backends")
    for b in backends:
        mark = "ok" if b.get("available") else "--"
        print(f"  {mark} {b['id']}")

    by_id = {b["id"]: b for b in backends}
    for name in PREFER:
        if not by_id.get(name, {}).get("available"):
            continue
        if try_forward(name):
            return

    print("FAIL: no usable backend")
    raise SystemExit(1)


if __name__ == "__main__":
    main()
