"""114. Python C-ABI — install and Dense forward."""
from __future__ import annotations

import welvet as w


def main() -> None:
    w.assert_engine_version()
    print("python_package", w.__version__)
    print("cabi_contract", w.engine_version())
    print("go_engine", "v1.7.3")  # engine behind this CABI wheel

    g = w.create_grid(backend="cpu_tiled")
    try:
        w.place_dense(
            g,
            {"in": 8, "out": 4, "act": "relu", "dtype": w.DType.FLOAT32},
        )
        x = [0.05 * (i + 1) for i in range(8)]
        y = w.forward(g, x)
        print("input", [round(v, 4) for v in x])
        print("output", [round(v, 4) for v in y])
    finally:
        w.free_grid(g)


if __name__ == "__main__":
    main()
