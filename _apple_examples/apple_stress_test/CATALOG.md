# Mamba stress matrix

One Mamba layer for forward passes; two independently seeded Mamba layers for training.

## Parameter dtypes (34)

float64 | float32 | float16 | bfloat16 | fp8e4m3 | fp8e5m2 | int64 | int32 | int16 | int8 | uint64 | uint32 | uint16 | uint8 | int4 | uint4 | fp4 | int2 | uint2 | ternary | binary | int | uint | uintptr | complex64 | complex128 | nf4 | fp6 | int6 | uint6 | int5 | uint5 | int3 | uint3

All inputs, outputs and interop gradients are Float32. Dtypes describe the authoritative parameter storage.

## Training modes (29)

| ID | Mode | Schedule |
|---:|---|---|
| 1 | NormalBP | Whole Stack |
| 2 | StepBP | Stack Line pipeline |
| 3 | Tween | Whole Stack |
| 4 | TweenChain | Whole Stack |
| 5 | StepTween | Stack Line pipeline |
| 6 | StepTweenChain | Stack Line pipeline |
| 7 | MeshBP | Grid mesh clock |
| 8 | MeshTween | Grid mesh clock |
| 9 | MeshTweenChain | Grid mesh clock |
| 10 | TweenSplit | Whole Stack |
| 11 | StepTweenSplit | Stack Line pipeline |
| 12 | TweenAlt | Whole Stack |
| 13 | StepTweenAlt | Stack Line pipeline |
| 14 | TweenSplitHeadProxy | Whole Stack |
| 15 | TweenSplitLinear | Whole Stack |
| 16 | TweenSplitFastProxy | Whole Stack |
| 17 | TweenSplitLinearCache | Whole Stack |
| 18 | TweenSplitHeadProxyAsync | Whole Stack |
| 19 | TweenSplitSparse | Whole Stack |
| 20 | MeshTweenSplit | Grid mesh clock |
| 21 | MeshTweenAlt | Grid mesh clock |
| 22 | MeshTweenSplitFastProxy | Grid mesh clock |
| 23 | MeshTweenSplitSparse | Grid mesh clock |
| 24 | StepTweenSplitHeadProxy | Stack Line pipeline |
| 25 | StepTweenSplitLinear | Stack Line pipeline |
| 26 | StepTweenSplitFastProxy | Stack Line pipeline |
| 27 | StepTweenSplitLinearCache | Stack Line pipeline |
| 28 | StepTweenSplitHeadProxyAsync | Stack Line pipeline |
| 29 | StepTweenSplitSparse | Stack Line pipeline |

Default all/all selection: **34 forward + 34 × 29 training = 1,020 comparison rows**, each with Go, C++ and Metal results and all three pairwise comparisons.
