"""116. Python C-ABI — train SGD and .entity round-trip."""
from __future__ import annotations

import welvet as w


def main() -> None:
    w.assert_engine_version()
    g = w.create_grid(backend="cpu_tiled")
    try:
        w.place_dense(
            g,
            {"in": 8, "out": 4, "act": "relu", "dtype": w.DType.FLOAT32},
        )
        x = [0.05 * (i + 1) for i in range(8)]
        y = [0.0, 0.0, 0.0, 1.0]
        result = w.train_sgd(g, x, y, lr=0.05)
        print("train", result)

        blob = w.serialize_entity(g)
        print("entity_bytes", len(blob))
    finally:
        w.free_grid(g)

    g2 = w.deserialize_entity(blob)
    try:
        out = w.forward(g2, x)
        print("reloaded_out", [round(v, 4) for v in out])
    finally:
        w.free_grid(g2)


if __name__ == "__main__":
    main()
