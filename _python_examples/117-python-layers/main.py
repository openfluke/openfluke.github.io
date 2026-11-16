"""117. Python C-ABI — MHA place and prepare_grid."""
from __future__ import annotations

import welvet as w


def main() -> None:
    w.assert_engine_version()
    kinds = w.list_layer_types()
    print("layer_types", len(kinds), "sample", kinds[:8])

    g = w.create_grid(backend="cpu_tiled")
    try:
        w.place(
            g,
            "mha",
            {
                "dim": 8,
                "DModel": 8,
                "NumHeads": 2,
                "SeqLen": 4,
                "dtype": w.DType.FLOAT32,
            },
        )
        # batch=1, seq=4, dim=8
        out = w.forward(g, [0.1] * 32, shape=[1, 4, 8])
        print("mha_out_len", len(out))

        # For qnn / mlx grids, call prepare_grid after place_* before forward:
        #   info = w.prepare_grid(g)
        #   print("prepare", info)
        print("prepare_grid", "call after place_* when backend is qnn or mlx")
    finally:
        w.free_grid(g)


if __name__ == "__main__":
    main()
