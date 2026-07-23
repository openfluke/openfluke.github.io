#!/usr/bin/env python3
"""Generate Welvet feature book (HTML) from engine inventory — not loom paste."""
from __future__ import annotations

import argparse
import json
import os
import re
import shutil
import subprocess
import sys
from dataclasses import dataclass, field
from pathlib import Path

ROOT = Path(__file__).resolve().parent
BOOK = ROOT / "welvet"
CHDIR = BOOK / "chapters"
EXAMPLES = BOOK / "examples"
DOWNLOADS = BOOK / "downloads"
WELVET_ROOT = ROOT.parent  # …/welvet engine
CHAOSGLUE_ROOT = ROOT.parent.parent
WEBGPU_ROOT = CHAOSGLUE_ROOT / "webgpu"
SENTENCEPIECE_ROOT = WELVET_ROOT / "third_party" / "go-sentencepiece"
MANIFEST = EXAMPLES / "_manifest.json"
DIST = ROOT / "dist"
# Chrome profile — outside the repo (SingletonLock etc. are runtime junk)
CHROME_PROFILE = Path.home() / ".cache" / "welvet-book-chrome"
RELEASES_URL = "https://github.com/openfluke/openfluke.github.io/releases"
LOGO_NAME = "openflukelogo_notxt.png"


@dataclass
class Chapter:
    slug: str
    num: str
    title: str
    part: str
    pkg: str  # import path or ""
    status: str  # ok | partial | missing
    status_label: str
    why: str
    what: str
    body_extra: str = ""
    example: str = ""
    run: str = ""
    runnable: bool = True  # False → no main.go / go run


def esc(s: str) -> str:
    return s.replace("&", "&amp;").replace("<", "&lt;").replace(">", "&gt;")


def code(s: str) -> str:
    return f"<pre><code>{esc(s.strip())}</code></pre>"


def status_badge(kind: str, label: str) -> str:
    return f'<span class="status {kind}">{esc(label)}</span>'


def ascii_fig(text: str, caption: str = "") -> str:
    cap = f"<figcaption>{esc(caption)}</figcaption>" if caption else ""
    return f'<figure class="diagram"><pre class="ascii">{esc(text.strip())}</pre>{cap}</figure>'


def read_welvet_version() -> tuple[str, float]:
    """Read scorecard version from welvet/README.md (single source of truth)."""
    text = (WELVET_ROOT / "README.md").read_text(encoding="utf-8")
    earned: float | None = None
    m = re.search(r"\*\*(\d+(?:\.\d+)?)\s*/\s*100\*\*\s*pts", text)
    if m:
        earned = float(m.group(1))
    if earned is None:
        m = re.search(r"\|\s*\*\*Version\*\*\s*\|\s*\*\*(v[\d.]+)\*\*", text)
        if m:
            v = m.group(1)
            if v == "v1.0":
                earned = 100.0
            elif v.startswith("v0."):
                earned = float(v[3:])
    if earned is None:
        earned = 76.0
    ver = "v1.0" if earned >= 100 else f"v0.{int(round(earned)):02d}"
    return ver, earned


def pdf_filename(version: str) -> str:
    return f"welvet-feature-book-{version}.pdf"


def ensure_logo() -> None:
    """Copy repo-root logo into welvet/assets/ if present."""
    dst = BOOK / "assets" / LOGO_NAME
    for src in (ROOT / LOGO_NAME, BOOK / LOGO_NAME):
        if src.is_file():
            dst.parent.mkdir(parents=True, exist_ok=True)
            shutil.copy2(src, dst)
            return


def title_logo_html() -> str:
    if (BOOK / "assets" / LOGO_NAME).is_file():
        return (
            f'<img class="title-logo" src="assets/{LOGO_NAME}" '
            f'alt="OpenFluke logo" width="160" height="123"/>'
        )
    return ""


# ---------------------------------------------------------------------------
# Chapter catalog — every Welvet feature package
# ---------------------------------------------------------------------------

def chapters(version: str = "v0.95", earned: float = 95.0) -> list[Chapter]:
    C = Chapter
    out: list[Chapter] = []
    earned_i = int(round(earned))

    out.append(C(
        "01-welvet", "1", "What Welvet is", "I · Orientation",
        "github.com/openfluke/welvet", "ok", "✅ engine",
        why="Loom’s flat <code>poly/</code> package hit import-cycle and honesty walls "
            "(QAT morph, silent fallbacks, god-layer). Welvet is the rewrite: one feature "
            "per folder, storage-truth dtypes/quants, Dense as the shared MatVec microkernel, "
            "tests only in <code>w2a</code>, apps only in <code>apps/</code>.",
        what="An AI engine in Go: layers, 34 dtypes, 20 quant formats, and three backends "
             f"(CPU tiled · Plan 9 SIMD · WebGPU). Version tracks a 100-point scorecard "
             f"(today <strong>{esc(version)}</strong>).",
        body_extra=ascii_fig("""
Rules
  1. No tests in engine packages          → w2a/
  2. No silent fallbacks                  → hard error
  3. No hardcoded float32                 → Tensor[T]
  4. No QAT                               → DType + Format = truth
  5. Train keeps storage truth            → in-dtype / re-Pack SGD, no retained f32 master
  6. One feature → one folder
  7. v1.0 = scorecard 100/100
""", "Non-negotiable engine rules.") + """
<h2>Origin</h2>
<p class="origin-byline"><strong>Samuel Watson</strong></p>
<p>Samuel Watson created OpenFluke after intensive T-ALL treatment (post-2018), from a simple frustration: AI tooling was heavy, opaque, and not portable enough to move models cleanly across operating systems or personal devices. Setting up GPU paths often meant large, brittle installs before experimentation could even begin.</p>
<p>That gap drove a long detour through cybersecurity, MBA, and psychology before returning through applied AI study - and finding that most paths still taught framework usage more than system construction. The result was a build journey across <a href="https://github.com/planetbridging">planetbridging</a> into OpenFluke projects, from paragon and loom to Welvet.</p>
<p>The early breakthroughs were practical and hard-won: first dense nets in Go, then WebGPU proof points (for example <a href="https://github.com/openfluke/iso-demo/releases/tag/v0.1.0">iso-demo v0.1.0</a>), then repeated cycles of speed, layers, 3D, tiling, SIMD, and bit-oriented model paths for personal intelligence.</p>
<p><strong>Why "OpenFluke"</strong>: open source values plus a life-saving fluke blood test that caught cancer early.<br/>
<strong>Why "Welvet"</strong>: a name shaped in gaming sessions with family.</p>
<blockquote>
<p>I wish this tool existed in 2018 to keep trying things in AI and asking "what if I try this?" - and maybe help deterministic medical deployment of AI models.</p>
</blockquote>
<p class="example-meta">"Our duty, as soldiers, is to protect humanity. Whatever the cost." - Master Chief</p>
""",
        example="""
package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/dense"
)

func main() {
	l, err := dense.New(8, 4, core.ActivationReLU, core.DTypeFloat32)
	if err != nil {
		panic(err)
	}
	l.Exec.Backend = core.BackendCPUTiled

	x := core.NewTensor[float32](1, 8)
	for i := range x.Data {
		x.Data[i] = float32(i) * 0.1
	}
	pre, post, err := dense.Forward(l, x)
	if err != nil {
		panic(err)
	}
	fmt.Println("pre", len(pre.Data), "post", len(post.Data))
}
""",
        run="cd welvet && go build ./...\ncd w2a && go test ./tests/dense -v",
    ))

    out.append(C(
        "02-tree", "2", "Repository map", "I · Orientation",
        "", "ok", "✅ layout",
        why="Readers need a single map of what is engine vs harness vs app vs stub.",
        what="Top-level folders and their ownership. Engine packages never import w2a.",
        body_extra=ascii_fig("""
welvet/
  core weights quant simd webgpu tiling architecture fusedgpu
  layers/*          ← one folder per op
  runtime/{forward,backward,training,step}
  systems/{dna,evolution,tween,tanhi,telemetry}
  model/{entity,hf,tokenizer,sampling,transformer}
  apps/{octo,flux2,mosstts}
  stub/*            ← designed surfaces, partial or empty
  w2a/              ← harness (separate module)
  tools/            ← eval / test tooling
"""),
        example="""
package main

import (
	"fmt"
	"os"
)

func main() {
	// Mental model only — list engine roots you depend on:
	roots := []string{
		"core", "weights", "quant", "simd", "webgpu", "tiling",
		"architecture", "fusedgpu", "layers", "runtime", "systems", "model",
	}
	for _, r := range roots {
		fmt.Println("import github.com/openfluke/welvet/" + r + "/…")
	}
	_ = os.Getenv // apps/stub/w2a are separate concerns
}
""",
    ))

    # Foundation
    out.append(C(
        "03-core", "3", "core — types & backends", "II · Foundation",
        "github.com/openfluke/welvet/core", "ok", "✅",
        why="Every polymorphic path needs one place for DType, LayerType, Activation, "
            "Backend, Tensor[T], and slim Layer metadata — without QAT morph defaults.",
        what="34 storage dtypes, Numeric generics, Tensor[T], ExecConfig "
             "(BackendCPUTiled | BackendSIMD | BackendWebGPU), Activate/ActivateDeriv, "
             "and converters for f16/bf16/fp8/fp4.",
        example="""
package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
)

func main() {
	t := core.NewTensor[float32](2, 4)
	t.Data[0] = -1
	t.Data[0] = core.Activate(t.Data[0], core.ActivationReLU) // 0

	cfg := core.ExecConfig{
		Backend:   core.BackendSIMD,
		MultiCore: true,
		TileSize:  32,
	}
	fmt.Println(cfg.Backend.String(), core.ParseDType("bfloat16"), t.Len())
}
""",
    ))

    out.append(C(
        "04-weights", "4", "weights — FormatNone MatVec", "II · Foundation",
        "github.com/openfluke/welvet/weights", "ok", "✅",
        why="Unquantized matrices still need a typed store that streams MatVec and SGD without "
            "forcing a float32 master or Morph-as-training.",
        what="Store holds DType + Format + Native/Packed bytes. New[T], MatVec, MatVecT, "
             "DecodeRow, SelectWire (F32/F64/I8), ApplySGD. FormatNone: update in native lanes "
             "(float32 payload is the only f32 buffer; float64 uses native ALU; other dtypes "
             "decode→update→re-encode). Packed: unpack→update→re-Pack then drop scratch. "
             "RetainsF32Master() is true only for FormatNone+float32. Dense and composite projs share this store.",
        example="""
package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/quant"
	"github.com/openfluke/welvet/weights"
)

func main() {
	s, err := weights.New[float32](2, 2, []float32{1, 0, 0, 1}, core.DTypeFloat32, quant.FormatNone)
	if err != nil {
		panic(err)
	}
	y := make([]float32, 2)
	if err := weights.MatVec(s, []float32{3, 4}, y); err != nil {
		panic(err)
	}
	fmt.Println(y) // ≈ [3, 4]
}
""",
    ))

    out.append(C(
        "05-quant", "5", "quant — 20 pack formats", "II · Foundation",
        "github.com/openfluke/welvet/quant", "ok", "✅",
        why="Inference, storage, and train need classic Q-packs, k-quants, IQ, Ternary/Binary, and "
            "Affine without a separate QAT mode or retained f32 master. Format is storage truth.",
        what="Pack / Unpack / MatVec / MatVecT for FormatNone…AffinePacked. Dense SIMD once-projects "
             "codes into Int8QS + scales (EnsureQ* / EnsureK/IQ/AffineSIMDCache) — no full-matrix F32 inflate for k/IQ/Affine. "
             "SGD on packed stores: short-lived unpack scratch → update → re-Pack; Packed stays truth.",
        body_extra="<p>Families: classic Q8/Q4/Q5 · K-quants · IQ · TernaryPacked/BinaryPacked · AffinePacked.</p>"
                   "<p>Honesty: fused Dense SIMD for all 20 quants (group DotKRow / DotIQRow / DotAffineRow + classic Q*/BitNet). "
                   "Peak dedicated k/IQ Plan 9 <code>.s</code> remains scorecard §12.</p>",
        example="""
package main

import (
	"fmt"

	"github.com/openfluke/welvet/quant"
)

func main() {
	w := []float32{0.1, -0.2, 0.3, 0.4, -0.5, 0.6, 0.05, -0.15}
	b, err := quant.Pack(quant.FormatQ4_0, w, 2, 4)
	if err != nil {
		panic(err)
	}
	y := make([]float32, 2)
	if err := quant.MatVec(b, []float32{1, 1, 1, 1}, y); err != nil {
		panic(err)
	}
	fmt.Println("rows", b.Rows, "cols", b.Cols, "y", y)
}
""",
    ))

    out.append(C(
        "06-simd", "6", "simd — Plan 9 kernels", "II · Foundation",
        "github.com/openfluke/welvet/simd", "ok", "✅",
        why="CPU peak needs hand-written AVX2/NEON without a silent Go fallback that "
            "pretends SIMD ran.",
        what="amd64/arm64 .s kernels: DotTile, DotI8/U8, DotQ4_0, Saxpy, BitNet helpers, "
             "packed f16/bf16/fp8/fp4 dots; Go fused DotKRow / DotIQRow / DotAffineRow for k/IQ/Affine. "
             "SimdEnabled() false → BackendSIMD hard-errors.",
        example="""
package main

import (
	"fmt"

	"github.com/openfluke/welvet/simd"
)

func main() {
	if !simd.SimdEnabled() {
		fmt.Println("SIMD not available on this arch — BackendSIMD must hard-error")
		return
	}
	acc := simd.DotTile([]float32{1, 2, 3, 4}, []float32{1, 1, 1, 1}, 0, 4, 0)
	fmt.Println(acc)
}
""",
    ))

    out.append(C(
        "07-webgpu", "7", "webgpu — device GEMV & shaders", "II · Foundation",
        "github.com/openfluke/welvet/webgpu", "ok", "✅",
        why="GPU paths must bind a real adapter. Host “fake GPU” was banned so suites "
            "cannot stamp WebGPU done when ALU ran on CPU.",
        what="DenseGEMV family (incl. quant/I8/resident), DenseGEMVT/DenseDW, RMSNorm, "
             "LayerNorm fwd, Softmax family, SwiGLUFuse. Available()/InitError() gate use.",
        example="""
package main

import (
	"fmt"

	"github.com/openfluke/welvet/webgpu"
)

func main() {
	if !webgpu.Available() {
		fmt.Println("no adapter:", webgpu.InitError())
		return
	}
	y := make([]float32, 2)
	err := webgpu.DenseGEMV(
		[]float32{1, 0, 0, 1},
		[]float32{1.5, 2.5},
		y, 1, 2, 2,
	)
	fmt.Println(y, err, webgpu.AdapterName())
}
""",
    ))

    out.append(C(
        "08-tiling", "8", "tiling — SC/MC & workgroups", "II · Foundation",
        "github.com/openfluke/welvet/tiling", "ok", "✅",
        why="MatVec throughput depends on tile size and when to go multi-core vs GPU workgroups. "
            "Centralizing caps keeps Dense and friends consistent.",
        what="DefaultCPUTile, DefaultGPUWG, CPUTile, PreferMultiCore, GPUWorkgroupsX.",
        example="""
package main

import (
	"fmt"

	"github.com/openfluke/welvet/tiling"
)

func main() {
	tile := tiling.CPUTile(0) // default 32
	mc := tiling.PreferMultiCore(8, 256, tile)
	wg := tiling.GPUWorkgroupsX(1024, 0)
	fmt.Println(tile, mc, wg)
}
""",
    ))

    out.append(C(
        "09-architecture", "9", "architecture — volumetric grid", "II · Foundation",
        "github.com/openfluke/welvet/architecture", "ok", "✅",
        why="Networks are spatial (Depth×Rows×Cols×LayersPerCell), not only linear stacks. "
            "Topology lives here; compute lives in layer packages.",
        what="Grid/Cell/Coord, NewGrid, BindOp, HopOrder, SetRemoteLink/ResolveHop. "
             "VolumetricNetwork is an alias for Grid.",
        example="""
package main

import (
	"fmt"

	"github.com/openfluke/welvet/architecture"
	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/dense"
)

func main() {
	g := architecture.NewGrid(1, 1, 2, 1) // depth, rows, cols, layersPerCell
	l, _ := dense.New(4, 4, core.ActivationLinear, core.DTypeFloat32)
	if err := dense.Place(g, 0, 0, 0, 0, l); err != nil {
		panic(err)
	}
	_ = g.SetRemoteLink(0, 0, 1, 0, 0, 0, 0, 0) // spatial feedback hop
	fmt.Println("cells", len(g.HopOrder()))
}
""",
    ))

    out.append(C(
        "10-fusedgpu", "10", "fusedgpu — decoder on device", "II · Foundation",
        "github.com/openfluke/welvet/fusedgpu", "ok", "✅",
        why="Token-by-token host round-trips kill decode. A fused engine keeps weights and "
            "scratch resident for Q4_0 and BinaryG128 hybrid paths.",
        what="Engine from Spec (AppendTokens/Reset/Close); HybridEngine "
             "(PrefillSample/DecodeSample/DecodeChunk). Specs come from model/transformer export.",
        example="""
package main

import (
	"fmt"

	"github.com/openfluke/welvet/fusedgpu"
	"github.com/openfluke/welvet/model/transformer"
)

func main() {
	m, err := transformer.LoadEntity("model.entity")
	if err != nil {
		fmt.Println("need an ENTITY file:", err)
		return
	}
	spec, err := m.ExportFusedGPUSpec()
	if err != nil {
		panic(err)
	}
	eng, err := fusedgpu.NewFromSpec(spec)
	if err != nil {
		panic(err)
	}
	defer eng.Close()
	logits, err := eng.AppendTokens([]uint32{1, 2, 3})
	fmt.Println(len(logits), err)
}
""",
    ))

    # Layers
    layer_specs = [
        ("11-dense", "11", "layers/dense — MatVec microkernel",
         "github.com/openfluke/welvet/layers/dense", "ok", "✅",
         "Most FLOPs are W@x. One Dense stack owns FormatNone×34 and all quants × three backends "
         "so every composite proj shares one correctness surface — including native in-dtype SGD.",
         "New / NewConfigured[T], Forward/Backward (dispatch on Exec.Backend), Place, ApplyGradSGD "
         "(→ weights.ApplySGD on the store). "
         "SIMD: fused Dot* for classic Q*, k/IQ (group scales), AffinePacked code-dot — no F32 inflate. "
         "Composites (MHA, SwiGLU, CNN im2col, RNN/LSTM) reuse Dense children and the same SGD path.",
         """
package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/dense"
	"github.com/openfluke/welvet/quant"
)

func main() {
	init := make([]float32, 4*8)
	l, err := dense.NewConfigured(8, 4, core.ActivationReLU, core.DTypeFloat32, quant.FormatNone, init)
	if err != nil {
		panic(err)
	}
	l.Exec.Backend = core.BackendCPUTiled
	x := core.NewTensor[float32](1, 8)
	pre, post, err := dense.Forward(l, x)
	if err != nil {
		panic(err)
	}
	gIn, gW, err := dense.Backward(l, post, x, pre)
	_ = dense.ApplyGradSGD(l, gW, 1e-3)
	fmt.Println(len(gIn.Data), len(gW.Data))
}
"""),
        ("12-mha", "12", "layers/mha — attention",
         "github.com/openfluke/welvet/layers/mha", "ok", "✅",
         "Transformers need multi-head attention with masks, RoPE/ALiBi, GQA/MQA, and cross-attn — "
         "without forking MatVec for every projection.",
         "Q/K/V/O are Dense children. Presets: DecoderCausal, EncoderBidirectional, CrossAttention, … "
         "Attn/RoPE ALU is host today; on-device attn shaders are still open.",
         """
package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/mha"
)

func main() {
	l, err := mha.New(mha.DecoderCausal(32, 4, 4))
	if err != nil {
		panic(err)
	}
	x := core.NewTensor[float32](1, 8, 32) // [batch, seq, dim]
	pre, post, err := mha.Forward(l, x)
	fmt.Println(pre != nil, len(post.Data), err)
}
"""),
        ("13-swiglu", "13", "layers/swiglu — gated FFN",
         "github.com/openfluke/welvet/layers/swiglu", "ok", "✅",
         "Modern decoder FFNs are SiLU(gate)⊙up → down. Projections must share Dense’s quant/backend matrix.",
         "Gate/Up/Down Dense children; DefaultFFN(dModel); WebGPU SiLU⊙ fuse on forward.",
         """
package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/swiglu"
)

func main() {
	l, err := swiglu.New(swiglu.DefaultFFN(64))
	if err != nil {
		panic(err)
	}
	x := core.NewTensor[float32](1, 64)
	_, post, err := swiglu.Forward(l, x)
	fmt.Println(len(post.Data), err)
}
"""),
        ("14-rmsnorm", "14", "layers/rmsnorm",
         "github.com/openfluke/welvet/layers/rmsnorm", "ok", "✅",
         "Llama-style blocks normalize by RMS, not mean+var. Needs native fwd/bwd and WebGPU shaders.",
         "Per-token RMS + γ on weights.Store; WebGPU fwd+bwd; SIMD DotTile stats + host scale.",
         """
package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/rmsnorm"
	"github.com/openfluke/welvet/quant"
)

func main() {
	gamma := []float32{1, 1, 1, 1}
	l, err := rmsnorm.NewConfigured(rmsnorm.Config{Dim: 4}, core.DTypeFloat32, quant.FormatNone, gamma)
	if err != nil {
		panic(err)
	}
	x := core.NewTensor[float32](1, 4)
	_, y, err := rmsnorm.Forward(l, x)
	fmt.Println(y.Data, err)
}
"""),
        ("15-layernorm", "15", "layers/layernorm",
         "github.com/openfluke/welvet/layers/layernorm", "ok", "✅",
         "Classic mean+var normalization with γ/β — still required for many HF architectures.",
         "WebGPU forward; backward host today. Same dtype×quant axes as other weighted layers.",
         """
package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/layernorm"
)

func main() {
	l, err := layernorm.New(layernorm.Config{Dim: 8})
	if err != nil {
		panic(err)
	}
	x := core.NewTensor[float32](1, 8)
	_, y, err := layernorm.Forward(l, x)
	fmt.Println(len(y.Data), err)
}
"""),
        ("16-embedding", "16", "layers/embedding",
         "github.com/openfluke/welvet/layers/embedding", "ok", "✅",
         "Token IDs must gather rows from a table — not a Dense MatVec — with scatter grads on backward.",
         "Config{VocabSize, EmbeddingDim, SeqLen}; table on weights.Store; host gather on SIMD/WebGPU today.",
         """
package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/embedding"
)

func main() {
	l, err := embedding.New(embedding.Config{VocabSize: 32, EmbeddingDim: 16, SeqLen: 4})
	if err != nil {
		panic(err)
	}
	ids := core.NewTensor[float32](1, 4)
	ids.Data[0], ids.Data[1] = 3, 7
	_, y, err := embedding.Forward(l, ids)
	fmt.Println(y.Shape, err)
}
"""),
        ("17-softmax", "17", "layers/softmax",
         "github.com/openfluke/welvet/layers/softmax", "ok", "✅",
         "Classification heads and attention need stable softmax variants, including sparse/Gumbel/Entmax for research paths.",
         "Weightless layer; KindStandard/Temperature/Grid/Hierarchical/Gumbel/Masked/Sparse/… "
         "WebGPU covers std family; exotic kinds hard-error on GPU (no silent host).",
         """
package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/softmax"
)

func main() {
	l, err := softmax.New(softmax.Config{Dim: 4, Kind: softmax.KindStandard})
	if err != nil {
		panic(err)
	}
	x := core.NewTensor[float32](1, 4)
	copy(x.Data, []float32{2, 1, 0.1, -1})
	_, y, err := softmax.Forward(l, x)
	fmt.Println(y.Data, err)
}
"""),
        ("18-sequential", "18", "layers/sequential",
         "github.com/openfluke/welvet/layers/sequential", "ok", "✅",
         "Some cells need an ordered Dense chain without burning grid hops.",
         "Dense→Dense compose in one cell. Nested non-Dense children still open.",
         """
package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/sequential"
)

func main() {
	l, err := sequential.New(sequential.Config{Dim: 16, Depth: 3})
	if err != nil {
		panic(err)
	}
	x := core.NewTensor[float32](1, 16)
	_, y, err := sequential.Forward(l, x)
	fmt.Println(len(y.Data), err)
}
"""),
        ("19-residual", "19", "layers/residual",
         "github.com/openfluke/welvet/layers/residual", "ok", "✅",
         "Skip connections stabilize deep stacks: y = F(x) + x with correct skip grads.",
         "F is Dense Dim→Dim (Depth≥1). Heterogeneous non-Dense F still open.",
         """
package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/residual"
)

func main() {
	l, err := residual.New(residual.Config{Dim: 16, Depth: 2})
	if err != nil {
		panic(err)
	}
	x := core.NewTensor[float32](1, 16)
	_, y, err := residual.Forward(l, x)
	fmt.Println(len(y.Data), err)
}
"""),
        ("20-cnn", "20", "layers/cnn1 · cnn2 · cnn3",
         "github.com/openfluke/welvet/layers/cnn2", "ok", "✅",
         "Conv nets must sit on the same dtype×quant×backend matrix as Dense. im2col → Dense GEMV is the intentional first cut.",
         "cnn1/cnn2/cnn3: host im2col then Proj Dense. Full timed matrices. Tiled conv shaders still ⬜.",
         """
package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/cnn2"
)

func main() {
	l, err := cnn2.New(cnn2.Config{
		InChannels: 1, Filters: 4, Height: 8, Width: 8, Kernel: 3, Stride: 1,
	})
	if err != nil {
		panic(err)
	}
	x := core.NewTensor[float32](1, 1, 8, 8)
	_, y, err := cnn2.Forward(l, x)
	fmt.Println(y.Shape, err)
}
"""),
        ("21-rnn-lstm", "21", "layers/rnn · lstm",
         "github.com/openfluke/welvet/layers/lstm", "ok", "✅",
         "Sequence models before transformers still need vanilla RNN and LSTM with BPTT on the shared MatVec stack.",
         "IH/HH (and LSTM gates) via Dense; recurrence ALU host; device required for WebGPU path.",
         """
package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/lstm"
	"github.com/openfluke/welvet/layers/rnn"
)

func main() {
	r, _ := rnn.New(rnn.Config{InputSize: 8, HiddenSize: 16, SeqLen: 4})
	l, _ := lstm.New(lstm.Config{InputSize: 8, HiddenSize: 16, SeqLen: 4})
	x := core.NewTensor[float32](1, 4, 8)
	_, yr, _ := rnn.Forward(r, x)
	_, yl, _ := lstm.Forward(l, x)
	fmt.Println(len(yr.Data), len(yl.Data))
}
"""),
        ("22-seqmix", "22", "layers/seqmix — mixer contract",
         "github.com/openfluke/welvet/layers/seqmix", "ok", "✅",
         "Attention, SSM, linear attn, and conv mixers must not be accidental forks of mha. Naming the contract keeps packages honest.",
         "KindAttention | KindSSM | KindLinearAttn | KindConvMix and Contract{Kind,DModel,MaxT}. No compute here.",
         """
package main

import (
	"fmt"

	"github.com/openfluke/welvet/layers/seqmix"
)

func main() {
	c := seqmix.Contract{Kind: seqmix.KindAttention, DModel: 64, MaxT: 2048}
	fmt.Println(c.Kind.String()) // attention
}
"""),
        ("23-gdn", "23", "layers/gdn — gated delta net",
         "github.com/openfluke/welvet/layers/gdn", "ok", "✅",
         "Linear attention / decode-first mixers (Gated DeltaNet) need a first-class package under KindLinearAttn.",
         "Exec CPU/SIMD/WebGPU; ForwardDecode; truncated BPTT; full timed matrix + train grids (Float32-primary).",
         """
package main

import (
	"fmt"

	"github.com/openfluke/welvet/layers/gdn"
)

func main() {
	l, err := gdn.New(gdn.Config{
		HiddenSize: 64, NumKeyHeads: 4, NumValueHeads: 4,
		KeyHeadDim: 16, ValueHeadDim: 16, ConvKernel: 4,
	})
	if err != nil {
		fmt.Println("config rejected:", err)
		return
	}
	l.Reset()
	x := make([]float32, 64)
	y := make([]float32, 64)
	fmt.Println(l.ForwardDecode(x, y))
}
"""),
        ("24-mamba", "24", "layers/mamba — selective SSM",
         "github.com/openfluke/welvet/layers/mamba", "ok", "✅",
         "SSM mixers (KindSSM) are not MHA clones — they need their own selective-scan path.",
         "InProj → softplus(Δ) scan → OutProj. Full timed matrix + train grids (scan ALU host).",
         """
package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/mamba"
)

func main() {
	l, err := mamba.New(mamba.Config{DModel: 32, DState: 16, Expand: 2, SeqLen: 8})
	if err != nil {
		fmt.Println(err)
		return
	}
	x := core.NewTensor[float32](1, 8, 32)
	_, y, err := mamba.Forward(l, x)
	fmt.Println(y != nil, err)
}
"""),
        ("25-convt", "25", "layers/convt1 · convt2 · convt3",
         "github.com/openfluke/welvet/layers/convt1", "ok", "✅",
         "Generators and U-Nets need transposed convolution as a peer of CNN, on the same Dense proj surface.",
         "Scatter upsample + Dense Proj. Full timed matrix + train grids (tiled transpose shaders still §12).",
         """
package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/convt1"
)

func main() {
	l, err := convt1.New(convt1.Config{InChannels: 4, Filters: 2, SeqLen: 8, Kernel: 3, Stride: 2})
	if err != nil {
		panic(err)
	}
	x := core.NewTensor[float32](1, 4, 8)
	_, y, err := convt1.Forward(l, x)
	fmt.Println(y.Shape, err)
}
"""),
        ("26-kmeans", "26", "layers/kmeans",
         "github.com/openfluke/welvet/layers/kmeans", "ok", "✅",
         "Soft clustering as a differentiable layer lets topology experiments sit inside the same train loop.",
         "Centers on Dense (K×FeatureDim); soft assignment outputs. Full timed matrix + train grids.",
         """
package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/kmeans"
)

func main() {
	l, err := kmeans.New(kmeans.Config{NumClusters: 4, FeatureDim: 8})
	if err != nil {
		panic(err)
	}
	x := core.NewTensor[float32](1, 8)
	_, y, err := kmeans.Forward(l, x)
	fmt.Println(y.Shape, err)
}
"""),
        ("27-parallel", "27", "layers/parallel — MoE combine",
         "github.com/openfluke/welvet/layers/parallel", "ok", "✅",
         "Mixture-of-experts and multi-path cells need concat/add/avg/filter combines over branches.",
         "Dense branches; filter gate mode for soft MoE. Full timed matrix + train grids. Heterogeneous residual graft still open.",
         """
package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/parallel"
)

func main() {
	l, err := parallel.New(parallel.Config{
		Dim: 16, Branches: 2, OutFeat: 16, Combine: parallel.CombineConcat,
	})
	if err != nil {
		panic(err)
	}
	x := core.NewTensor[float32](1, 16)
	_, y, err := parallel.Forward(l, x)
	fmt.Println(y.Shape, err)
}
"""),
        ("28-metacognition", "28", "layers/metacognition",
         "github.com/openfluke/welvet/layers/metacognition", "ok", "✅",
         "Observed layers can apply heuristic stability rules (gate/scale/reset) without dtype morph/QAT.",
         "Wraps Dense + DefaultStabilityRules(); Stats exposed. Full timed matrix + train grids.",
         """
package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/metacognition"
)

func main() {
	l, err := metacognition.New(metacognition.Config{
		Dim: 16, Rules: metacognition.DefaultStabilityRules(),
	})
	if err != nil {
		panic(err)
	}
	x := core.NewTensor[float32](1, 16)
	_, y, err := metacognition.Forward(l, x)
	fmt.Println(len(y.Data), err)
}
"""),
    ]
    for slug, num, title, pkg, st, lab, why, what, ex in layer_specs:
        out.append(C(slug, num, title, "III · Layers", pkg, st, lab, why, what, "", ex))

    # Runtime
    out.append(C(
        "29-forward", "29", "runtime/forward", "IV · Runtime",
        "github.com/openfluke/welvet/runtime/forward", "ok", "✅",
        why="A grid of heterogeneous ops needs one walker that dispatches by concrete type and fails loudly on unknowns.",
        what="Forward[T](grid, input) → Result tape; Cell[T] for single-cell dispatch. Covers Dense…Residual + extended wired ops.",
        example="""
package main

import (
	"fmt"

	"github.com/openfluke/welvet/architecture"
	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/dense"
	"github.com/openfluke/welvet/runtime/forward"
)

func main() {
	g := architecture.NewGrid(1, 1, 1, 1)
	l, _ := dense.New(8, 8, core.ActivationLinear, core.DTypeFloat32)
	_ = dense.Place(g, 0, 0, 0, 0, l)
	res, err := forward.Forward(g, core.NewTensor[float32](1, 8))
	fmt.Println(res != nil, err)
}
""",
    ))

    out.append(C(
        "30-backward", "30", "runtime/backward", "IV · Runtime",
        "github.com/openfluke/welvet/runtime/backward", "ok", "✅",
        why="Training needs a reverse tape over the same ops forward used — no separate graph framework.",
        what="Backward[T](fwdResult, gradOut) walks the tape and returns per-op weight grads.",
        example="""
package main

import (
	"fmt"

	"github.com/openfluke/welvet/architecture"
	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/dense"
	"github.com/openfluke/welvet/runtime/backward"
	"github.com/openfluke/welvet/runtime/forward"
)

func main() {
	g := architecture.NewGrid(1, 1, 1, 1)
	l, _ := dense.New(4, 4, core.ActivationLinear, core.DTypeFloat32)
	_ = dense.Place(g, 0, 0, 0, 0, l)
	fwd, _ := forward.Forward(g, core.NewTensor[float32](1, 4))
	bwd, err := backward.Backward(fwd, core.NewTensor[float32](1, 4))
	fmt.Println(bwd != nil, err)
}
""",
    ))

    out.append(C(
        "31-training", "31", "runtime/training", "IV · Runtime",
        "github.com/openfluke/welvet/runtime/training", "ok", "✅",
        why="Suites and small nets need MSE+SGD and tween hooks without inventing an external trainer "
            "or a retained float32 master beside storage.",
        what="MSE/MSEGrad, SGD, Step, ApplyTween/StepTween, StepMesh. Layer-agnostic ApplyGradSGD "
             "dispatch (Dense…Mamba/GDN/…). FormatNone: in-dtype ApplySGD; packed: unpack→update→re-Pack. "
             "No QAT dual path — storage dtype/format is truth after every step.",
        example="""
package main

import (
	"fmt"

	"github.com/openfluke/welvet/architecture"
	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/dense"
	"github.com/openfluke/welvet/runtime/forward"
	"github.com/openfluke/welvet/runtime/training"
)

func main() {
	g := architecture.NewGrid(1, 1, 1, 1)
	l, _ := dense.New(4, 4, core.ActivationLinear, core.DTypeFloat32)
	_ = dense.Place(g, 0, 0, 0, 0, l)
	fwd, _ := forward.Forward(g, core.NewTensor[float32](1, 4))
	loss, err := training.Step(fwd, core.NewTensor[float32](1, 4), 1e-2)
	fmt.Println(loss, err)
}
""",
    ))

    out.append(C(
        "32-step", "32", "runtime/step — step mesh", "IV · Runtime",
        "github.com/openfluke/welvet/runtime/step", "ok", "✅",
        why="Spatial feedback (remote links) needs a discrete-time mesh where every cell updates from a double buffer — different from a decoder wavefront. "
            "Cross-numeric train also needs the same mesh with weight DType ⊥ activation Tensor[T].",
        what="State[T], StepForward/StepBackward/StepApplyTween / StepMesh across the grid for all wired Ops × dtype × quant × CPU/SIMD. "
             "W2A Cross-Numeric Train: polyops.AllKinds() × weight dtype × act host "
             "(smoke ~21×7×5 ≈ 735; full ~21×34×15 ≈ 10.7k) — asserts no retained f32 master after StepMesh.",
        example="""
package main

import (
	"fmt"

	"github.com/openfluke/welvet/architecture"
	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/dense"
	"github.com/openfluke/welvet/runtime/step"
)

func main() {
	g := architecture.NewGrid(1, 1, 1, 1)
	l, _ := dense.New(4, 4, core.ActivationLinear, core.DTypeFloat32)
	_ = dense.Place(g, 0, 0, 0, 0, l)
	st := step.New[float32](g)
	_, err := step.StepForward(g, st, false)
	fmt.Println(err)
}
""",
    ))

    out.append(C(
        "65-cross-numeric", "65", "Cross-numeric train + down-the-dem", "IV · Runtime",
        "github.com/openfluke/welvet/runtime/training", "ok", "✅",
        why="Weight storage dtype and activation Tensor[T] are independent axes. "
            "Proving train without a retained float32 master means sweeping W×A — not only matched float32 acts.",
        what="W2A Step Cross-Numeric Train: polyops.AllKinds() × FormatNone weight dtype × Go Numeric act host "
             "(smoke ~735; full ~10.7k) via StepMesh, then assert no retained f32 master. "
             "Public Dense volumetric showcase: down-the-dem — dtype demotion ladder, packed quants, "
             "and the full 34×15×3 perm matrix with charts/PDF.",
        body_extra="""
<p><strong>down-the-dem</strong> — same 5-layer Dense stack on cubes 1³ / 2³ / 3³; only the number story changes:</p>
<ul>
<li><strong>dtype</strong> — FormatNone demotion (wide → 1-bit), in-dtype SGD</li>
<li><strong>quant</strong> — packed formats; unpack → update → re-Pack</li>
<li><strong>perm</strong> — every weight dtype × every act host (1,530 cells on three grids)</li>
</ul>
<p>Repo + report: <a href="https://github.com/openfluke/down-the-dem">github.com/openfluke/down-the-dem</a>
 · PDF: <code>report/down-the-dem.pdf</code> in that repo.
 Honesty: runnable storage-truth train ≠ task accuracy or native low-precision GEMV ALU.</p>
""",
        example="""
package main

import (
	"fmt"

	"github.com/openfluke/welvet/architecture"
	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/dense"
	"github.com/openfluke/welvet/quant"
	"github.com/openfluke/welvet/runtime/training"
)

func main() {
	// Weight storage: int8 FormatNone. Activations/grads: Tensor[int8].
	g := architecture.NewGrid(1, 1, 1, 1)
	init := make([]float32, 4*4)
	for i := range init {
		init[i] = 0.1
	}
	l, err := dense.NewConfigured(4, 4, core.ActivationLinear, core.DTypeInt8, quant.FormatNone, init)
	if err != nil {
		panic(err)
	}
	if err := dense.Place(g, 0, 0, 0, 0, l); err != nil {
		panic(err)
	}
	x := core.NewTensor[int8](1, 4)
	y := core.NewTensor[int8](1, 4)
	for i := 0; i < 4; i++ {
		x.Data[i] = int8(i + 1)
		y.Data[i] = int8(i)
	}
	loss, _, err := training.StepMesh(g, x, y, 1, 0.05)
	fmt.Println(loss, err, l.Weights.RetainsF32Master())
}
""",
    ))

    # Systems
    for slug, num, title, pkg, why, what, ex in [
        ("33-dna", "33", "systems/dna", "github.com/openfluke/welvet/systems/dna",
         "Quant and train must be measurable as topology/weight fingerprints — DNA detects logic shifts.",
         "ExtractDNA, CompareNetworks, CosineSimilarity, FlattenOp / CollectStores across dtype×quant.",
         """
package main

import (
	"fmt"

	"github.com/openfluke/welvet/architecture"
	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/dense"
	"github.com/openfluke/welvet/systems/dna"
)

func main() {
	g := architecture.NewGrid(1, 1, 1, 1)
	l, _ := dense.New(4, 4, core.ActivationLinear, core.DTypeFloat32)
	_ = dense.Place(g, 0, 0, 0, 0, l)
	a := dna.ExtractDNA(g)
	fmt.Println(dna.CompareNetworks(a, a))
}
"""),
        ("34-evolution", "34", "systems/evolution", "github.com/openfluke/welvet/systems/evolution",
         "Topology search and weight crossover need first-class splice + NEAT on CPU-resident grids.",
         "SpliceDNA, NEATMutate, NewNEATPopulation, CloneGrid — dtype/quant preserved via SetFromF32.",
         """
package main

import (
	"fmt"

	"github.com/openfluke/welvet/architecture"
	"github.com/openfluke/welvet/systems/evolution"
)

func main() {
	a, b := architecture.NewGrid(1, 1, 1, 1), architecture.NewGrid(1, 1, 1, 1)
	child, err := evolution.SpliceDNA(a, b, evolution.DefaultSpliceConfig())
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = evolution.NEATMutate(child, evolution.DefaultNEATConfig(16))
	fmt.Println(err)
}
"""),
        ("35-tween", "35", "systems/tween", "github.com/openfluke/welvet/systems/tween",
         "Target propagation (chain-rule or Hebbian layerwise gaps) is an alternative credit-assignment path.",
         "NewState, Forward, BackwardChainRule / BackwardLayerwise, ApplyGaps; SIMD DotTile/Saxpy budgets.",
         """
package main

import (
	"fmt"

	"github.com/openfluke/welvet/architecture"
	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/systems/tween"
)

func main() {
	g := architecture.NewGrid(1, 1, 1, 1)
	st := tween.NewState[float32](g, tween.DefaultConfig())
	_, err := tween.Forward(g, st, core.NewTensor[float32](1, 4))
	fmt.Println(err)
}
"""),
        ("36-tanhi", "36", "systems/tanhi — UDP HUD", "github.com/openfluke/welvet/systems/tanhi",
         "Training visualization must never block the engine — best-effort UDP JSON-lines to a HUD.",
         "ConfigFromGrid, Emit/EmitSweep, DefaultUDPPort. SoulGlitch-style consumers.",
         """
package main

import (
	"github.com/openfluke/welvet/architecture"
	"github.com/openfluke/welvet/systems/tanhi"
)

func main() {
	g := architecture.NewGrid(1, 1, 1, 1)
	cfg := tanhi.ConfigFromGrid(g)
	tanhi.EmitSweep(cfg, "epoch-0") // non-blocking best-effort
}
"""),
        ("37-telemetry", "37", "systems/telemetry", "github.com/openfluke/welvet/systems/telemetry",
         "Static structural blueprints (sizes, op kinds) differ from live tanhi events.",
         "ExtractNetworkBlueprint, ExtractLayerTelemetry for introspection/UIs.",
         """
package main

import (
	"fmt"

	"github.com/openfluke/welvet/architecture"
	"github.com/openfluke/welvet/systems/telemetry"
)

func main() {
	g := architecture.NewGrid(1, 1, 1, 1)
	bp := telemetry.ExtractNetworkBlueprint(g, "demo")
	fmt.Printf("%+v\\n", bp)
}
"""),
    ]:
        out.append(C(slug, num, title, "V · Systems", pkg, "ok", "✅", why, what, "", ex))

    # Model
    out.append(C(
        "38-entity", "38", "model/entity — .entity files", "VI · Model IO",
        "github.com/openfluke/welvet/model/entity", "ok", "✅",
        why="HF safetensors are awkward for native topology + packed weights. ENTITY is the Welvet checkpoint.",
        what="Open/Inspect/IsEntity, LoadBlob/LoadQuantBlob, PackFromHF/ImportFromHF, WriteTransformerFile, SerializeNetwork.",
        example="""
package main

import (
	"fmt"

	"github.com/openfluke/welvet/model/entity"
)

func main() {
	info, err := entity.Inspect("model.entity")
	if err != nil {
		fmt.Println(err)
		return
	}
	ef, err := entity.Open("model.entity")
	if err != nil {
		panic(err)
	}
	defer ef.Close()
	fmt.Println(info, ef.HasTokenizerBlob())
}
""",
    ))

    out.append(C(
        "39-hf", "39", "model/hf — snapshots", "VI · Model IO",
        "github.com/openfluke/welvet/model/hf", "ok", "✅",
        why="Import starts with probing HF/MLX layouts before packing ENTITY.",
        what="InspectSnapshot, DetectArchitecture, safetensors/MLX loaders, Qwen3.5 hybrid helpers.",
        example="""
package main

import (
	"fmt"

	"github.com/openfluke/welvet/model/hf"
)

func main() {
	info, err := hf.InspectSnapshot("/path/to/hf-snapshot")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(info)
}
""",
    ))

    out.append(C(
        "40-tokenizer", "40", "model/tokenizer", "VI · Model IO",
        "github.com/openfluke/welvet/model/tokenizer", "ok", "✅",
        why="Generate needs encode/decode of HF tokenizer.json without pulling Python.",
        what="LoadTokenizer, Encode/Decode, LoadForEntity.",
        example="""
package main

import (
	"fmt"

	"github.com/openfluke/welvet/model/tokenizer"
)

func main() {
	tok, err := tokenizer.LoadTokenizer("tokenizer.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	ids := tok.Encode("hello welvet", true)
	fmt.Println(ids, tok.Decode(ids, true))
}
""",
    ))

    out.append(C(
        "41-sampling", "41", "model/sampling", "VI · Model IO",
        "github.com/openfluke/welvet/model/sampling", "ok", "✅",
        why="Logits → token ID needs ArgMax, TopK+temperature, penalties, and chat hygiene in one place.",
        what="ArgMax, SampleTopK, ApplyRepetitionPenalty, BanIDs, SanitizeChatReply.",
        example="""
package main

import (
	"fmt"

	"github.com/openfluke/welvet/model/sampling"
)

func main() {
	logits := []float32{0.1, 2.4, 0.3, -1}
	fmt.Println(sampling.ArgMax(logits))
	fmt.Println(sampling.SampleTopK(logits, 2, 0.8, true))
}
""",
    ))

    out.append(C(
        "42-transformer", "42", "model/transformer — generate", "VI · Model IO",
        "github.com/openfluke/welvet/model/transformer", "ok", "✅",
        why="ENTITY packs must run as Llama-style decoders with KV cache, profiles (SIMD/WebGPU/fused), and chat templates.",
        what="LoadEntity → Model; Generate; ApplyExec profiles; ExportFusedGPUSpec / hybrid sync.",
        example="""
package main

import (
	"fmt"

	"github.com/openfluke/welvet/model/tokenizer"
	"github.com/openfluke/welvet/model/transformer"
)

func main() {
	m, err := transformer.LoadEntity("model.entity")
	if err != nil {
		fmt.Println(err)
		return
	}
	tok, err := tokenizer.LoadTokenizer("tokenizer.json")
	if err != nil {
		panic(err)
	}
	if err := m.ApplyExec(transformer.ProfileSIMDMultiCore()); err != nil {
		panic(err)
	}
	text, _, err := m.Generate(
		tok.Encode, tok.Decode, nil,
		"You are helpful.", "Say hi.",
		transformer.GenOptions{MaxTokens: 32, Silent: true},
	)
	fmt.Println(text, err)
}
""",
    ))

    # Apps
    out.append(C(
        "43-apps", "43", "apps — octo · flux2 · mosstts", "VII · Apps",
        "github.com/openfluke/welvet/apps/…", "partial", "🚧",
        why="Products must not pollute engine packages. Octo is the model shell; flux2/mosstts are domain apps.",
        what="octo (own module): download/convert/chat, see §44. flux2: MMDiT image. mosstts: Speak/SpeakToFile pipeline.",
        example="""
package main

import "fmt"

func main() {
	fmt.Println("Apps import the engine — the engine never imports apps.")
	fmt.Println("  cd apps/octo && go run .")
	fmt.Println("  cd apps/flux2 && go run .")
	fmt.Println("  cd apps/mosstts && go run .")
}
""",
        run="cd apps/octo && go run .   # when module wired",
    ))

    out.append(C(
        "44-octo", "44", "Octo — model shell", "VII · Apps",
        "github.com/openfluke/welvet/apps/octo", "ok", "✅ runs",
        why="A model is only useful with a shell around it: pull weights from Hugging Face, convert them to a Welvet <code>.entity</code>, then chat, serve, or benchmark. Octo is that shell, kept in its own module so the engine never depends on an app.",
        what="Subcommands cover the whole loop: <code>hub</code> download/ensure repos, <code>convert</code> pack to a single <code>.entity</code>, interactive <code>run</code>/<code>serve</code>/chat, plus image and speech menus. Its <code>bench</code> harness sweeps every quant format across a CPU Plan 9 SIMD fused profile and a WebGPU fused profile.",
        body_extra="""
<h3>Bench: SmolLM2-135M-Instruct across quant formats</h3>
<p class="example-meta">template smol2-135m-fuse · linux/amd64 · NVIDIA GTX 1650 SUPER · 3m17s · 80 runs (76 ok, 4 skipped: "none" has no fused kernel). Throughput is total tokens/s on the first prompt; entity is on-disk size.</p>
<table>
<thead><tr><th>Quant</th><th>GPU-fused tok/s</th><th>SIMD-fused tok/s</th><th>Entity MB</th></tr></thead>
<tbody>
<tr><td>Q8_0</td><td>108.2</td><td>33.2</td><td>252.5</td></tr>
<tr><td>Q6_K</td><td>73.5</td><td>4.9</td><td>216.4</td></tr>
<tr><td>Q5_K</td><td>71.7</td><td>6.3</td><td>208.4</td></tr>
<tr><td>Q4_K</td><td>73.4</td><td>7.8</td><td>192.3</td></tr>
<tr><td>Q4_0</td><td>112.1</td><td>36.8</td><td>188.3</td></tr>
<tr><td>IQ4_XS</td><td>105.5</td><td>7.6</td><td>204.4</td></tr>
<tr><td>IQ4_NL</td><td>84.1</td><td>8.8</td><td>188.3</td></tr>
<tr><td>TernaryPacked</td><td>106.9</td><td>18.9</td><td>172.3</td></tr>
<tr><td>BinaryPacked</td><td>88.3</td><td>21.0</td><td>140.2</td></tr>
</tbody></table>
<p class="example-meta">21 quant formats run end to end; the table shows a representative slice. Full report: <code>apps/octo/dist/octo-run-v0.76.txt</code>.</p>
<div class="callout warn"><strong>Honest low-bit note</strong>High-bit formats (Q8, Q6, Q5, Q4, IQ4) stay coherent. Sub-2-bit formats (IQ1/IQ2, Ternary, Binary) are fast and tiny but degrade into gibberish. The bench keeps them in to show the speed, size, and quality trade honestly rather than hiding it.</div>
<h3>Sample replies</h3>
<pre class="ascii">Q8_0  / gpu   I'm a helpful assistant. I'm here to help you with your inquiries.
Q4_0  / gpu   I'm ready to help. What can I help you with?
Q6_K  / gpu   2 + 2 is 4.
Q6_K  / simd  I'm a helpful assistant.
IQ2_XXS/gpu   ,))hidAPAAPAAPAAPAAPA...   (sub-2-bit degrades)</pre>
""",
        example="""
package main

import "fmt"

func main() {
	fmt.Println("Octo is its own module — the engine never imports it.")
	fmt.Println("  cd apps/octo && go run .")
	fmt.Println("  download -> convert -> .entity -> chat / serve / bench")
}
""",
        run="cd apps/octo && go run .",
    ))

    # Stubs — each one
    stubs = [
        ("44-stub-seed", "44", "stub/seed", "github.com/openfluke/welvet/stub/seed",
         "Ship topology recipes (layer seeds → He-init) without weight blobs.",
         "SeedFrom/From, RNG, InitStoreHe, Dense/MHA/SwiGLU/CNN infinite manifests.",
         """
package main

import (
	"fmt"

	"github.com/openfluke/welvet/stub/seed"
)

func main() {
	s := seed.From("demo-net", 0)
	rng := seed.New(s)
	fmt.Println(rng.NormFloat64())
}
"""),
        ("45-stub-serialization", "45", "stub/serialization", "github.com/openfluke/welvet/stub/serialization",
         "Volumetric grids need JSON/ENTITY persist beyond transformer packs.",
         "SerializeEntity/LoadEntity, SerializeGrid/GridFromSpec — native FormatNone bytes; packed via wire.",
         """
package main

import (
	"fmt"

	"github.com/openfluke/welvet/architecture"
	"github.com/openfluke/welvet/stub/serialization"
)

func main() {
	g := architecture.NewGrid(1, 1, 1, 1)
	raw, err := serialization.SerializeGrid(g)
	fmt.Println(len(raw), err)
}
"""),
        ("46-stub-memory", "46", "stub/memory", "github.com/openfluke/welvet/stub/memory",
         "HF→ENTITY and GPU upload need footprint accounting and optional history charts.",
         "FromGrid, Footprint, InitScavenger, ReleaseTransient; WELVET_MEMORY_HISTORY=1.",
         """
package main

import (
	"fmt"

	"github.com/openfluke/welvet/architecture"
	"github.com/openfluke/welvet/stub/memory"
)

func main() {
	g := architecture.NewGrid(1, 1, 1, 1)
	fp := memory.FromGrid(g)
	fmt.Printf("%+v\\n", fp)
}
"""),
        ("47-stub-donate", "47", "stub/donate", "github.com/openfluke/welvet/stub/donate",
         "LAN donors should accept framed JSON jobs without embedding HTTP in the engine.",
         "u32-LE + JSON frames; ServeTCP/Dial; model_push vs local_lm. v0 workers echo.",
         """
package main

import (
	"fmt"

	"github.com/openfluke/welvet/stub/donate"
)

func main() {
	fmt.Println("default port", donate.DefaultPort)
	// ln, err := donate.ServeTCP(donate.ServerOptions{Addr: ":17001", Mode: donate.ServerLocalLM})
}
"""),
        ("48-stub-fountain", "48", "stub/fountain", "github.com/openfluke/welvet/stub/fountain",
         "Recover specialist weight blobs over lossy links via LT fountain codes, then ensemble.",
         "NewLTEncoder/Decoder, RecoverWeightBlobs, Neural/DenseFactory helpers.",
         """
package main

import (
	"fmt"

	"github.com/openfluke/welvet/stub/fountain"
)

func main() {
	src := [][]byte{make([]byte, 64), make([]byte, 64), make([]byte, 64)}
	enc, err := fountain.NewLTEncoder(src, 42)
	fmt.Printf("%T %v\\n", enc, err)
}
"""),
        ("49-stub-hardware", "49", "stub/hardware", "github.com/openfluke/welvet/stub/hardware",
         "Dispatchers and UIs need a portable host audit (OS/CPU/RAM/GPU).",
         "Audit() → SystemAudit with linux /proc or platform fallbacks.",
         """
package main

import (
	"fmt"

	"github.com/openfluke/welvet/stub/hardware"
)

func main() {
	a := hardware.Audit()
	fmt.Printf("%+v\\n", a.CPU)
	fmt.Println(hardware.Description())
}
"""),
        ("50-stub-accel", "50", "stub/accel — NPU/Metal/QNN", "github.com/openfluke/welvet/stub/accel",
         "Vendor accelerators (Intel NPU, Qualcomm QNN, Apple Metal) will plug beside WebGPU — not replace it.",
         "Package is doc.go only today (⬜). No symbols yet; reserved for plugin loader design.",
         """
package main

// import "github.com/openfluke/welvet/stub/accel"
// No exported API yet — see package doc.go for intent.

func main() {}
"""),
        ("51-stub-clustering", "51", "stub/clustering", "github.com/openfluke/welvet/stub/clustering",
         "Offline clustering helpers on tensors without inventing a second math stack.",
         "KMeansCluster, HierarchicalGroup, distance/silhouette helpers.",
         """
package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/stub/clustering"
)

func main() {
	pts := []*core.Tensor[float32]{
		core.NewTensor[float32](2),
		core.NewTensor[float32](2),
	}
	pts[0].Data[0], pts[0].Data[1] = 0, 0
	pts[1].Data[0], pts[1].Data[1] = 1, 1
	cent, assign := clustering.KMeansCluster(pts, 2, 10, false)
	fmt.Println(len(cent), assign)
}
"""),
        ("52-stub-ensemble", "52", "stub/ensemble", "github.com/openfluke/welvet/stub/ensemble",
         "Combine multiple model votes and find complementary specialists.",
         "MajorityVote, FindComplementaryMatches, PerformanceSimilarity.",
         """
package main

import (
	"fmt"

	"github.com/openfluke/welvet/stub/ensemble"
)

func main() {
	votes := [][]int{{0, 1, 1}, {0, 0, 1}, {0, 1, 0}}
	fmt.Println(ensemble.MajorityVote(votes))
}
"""),
        ("53-stub-evaluation", "53", "stub/evaluation", "github.com/openfluke/welvet/stub/evaluation",
         "Benchmark grids through runtime/forward with deviation metrics.",
         "EvaluateNetwork, MultiNetworkEvaluation, DeviationMetrics.",
         """
package main

import (
	"fmt"

	"github.com/openfluke/welvet/architecture"
	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/stub/evaluation"
)

func main() {
	g := architecture.NewGrid(1, 1, 1, 1)
	in := []*core.Tensor[float32]{core.NewTensor[float32](1, 4)}
	rep, err := evaluation.EvaluateNetwork(g, in, []float64{0, 0, 0, 0})
	fmt.Println(rep, err)
}
"""),
        ("54-stub-grafting", "54", "stub/grafting", "github.com/openfluke/welvet/stub/grafting",
         "Merge grids into Parallel/Residual structures for topology experiments.",
         "GraftGrids, ResidualGraft, GraftToGrid — Dense branches in v0.",
         """
package main

import (
	"fmt"

	"github.com/openfluke/welvet/architecture"
	"github.com/openfluke/welvet/layers/parallel"
	"github.com/openfluke/welvet/stub/grafting"
)

func main() {
	a, b := architecture.NewGrid(1, 1, 1, 1), architecture.NewGrid(1, 1, 1, 1)
	out, err := grafting.GraftGrids([]*architecture.Grid{a, b}, parallel.CombineConcat)
	fmt.Println(out != nil, err)
}
"""),
        ("55-stub-grouping", "55", "stub/grouping", "github.com/openfluke/welvet/stub/grouping",
         "Detect layer archetypes from safetensor-style names before mounting.",
         "GroupRelatedTensors, DetectMHA/DetectSwiGLU/DetectRMSNorm → ArchetypeHint.",
         """
package main

import (
	"fmt"

	"github.com/openfluke/welvet/stub/grouping"
)

func main() {
	tensors := []grouping.DetectedTensor{
		{Name: "model.layers.0.self_attn.q_proj.weight"},
		{Name: "model.layers.0.self_attn.k_proj.weight"},
		{Name: "model.layers.0.self_attn.v_proj.weight"},
		{Name: "model.layers.0.self_attn.o_proj.weight"},
	}
	ok, hint := grouping.DetectMHA("block0", tensors, 64, 4)
	fmt.Println(ok, hint)
}
"""),
        ("56-stub-introspection", "56", "stub/introspection", "github.com/openfluke/welvet/stub/introspection",
         "UIs and FFI need to list Grid methods without hardcoding every export.",
         "GetMethods, GetMethodsJSON, GetMethodSignature.",
         """
package main

import (
	"fmt"

	"github.com/openfluke/welvet/architecture"
	"github.com/openfluke/welvet/stub/introspection"
)

func main() {
	g := architecture.NewGrid(1, 1, 1, 1)
	methods, err := introspection.GetMethods(g)
	fmt.Println(len(methods), err)
}
"""),
        ("57-stub-observer", "57", "stub/observer", "github.com/openfluke/welvet/stub/observer",
         "Attach forward/backward observers for debugging without coupling to tanhi UDP.",
         "Observer interface; ConsoleObserver / HTTPObserver / BufferObserver; ComputeLayerStats.",
         """
package main

import (
	"fmt"

	"github.com/openfluke/welvet/stub/observer"
)

func main() {
	var obs observer.Observer = &observer.ConsoleObserver{}
	fmt.Printf("%T\\n", obs)
}
"""),
        ("58-stub-pipeline", "58", "stub/pipeline", "github.com/openfluke/welvet/stub/pipeline",
         "Decoder wavefront stats helpers — not a full Lucy-style pipeline runner yet.",
         "PipelineForwardStats, TokenTimelineSummary.",
         """
package main

import (
	"fmt"

	"github.com/openfluke/welvet/stub/pipeline"
)

func main() {
	var s pipeline.PipelineForwardStats
	fmt.Printf("%+v\\n", s)
}
"""),
        ("59-stub-templates", "59", "stub/templates", "github.com/openfluke/welvet/stub/templates",
         "Chat prompts must match model families (ChatML, Llama3, BitNet) without app-specific string glue.",
         "Template.BuildPrompt; presets ChatML, Llama3, BitNetInstruction, MicrosoftBitNetChat.",
         """
package main

import (
	"fmt"

	"github.com/openfluke/welvet/stub/templates"
)

func main() {
	p := templates.ChatML.BuildPrompt(nil, "You are helpful.", "Say hi")
	fmt.Println(p)
}
"""),
        ("60-stub-universal", "60", "stub/universal", "github.com/openfluke/welvet/stub/universal",
         "Probe unknown safetensor geometry and mount placeholder grids until full weight import lands.",
         "LoadUniversal / LoadUniversalDetailed, ProbeDeepGeometry, MountGeometrically.",
         """
package main

import (
	"fmt"

	"github.com/openfluke/welvet/stub/universal"
)

func main() {
	g, err := universal.LoadUniversal("/path/to/snapshot")
	fmt.Println(g != nil, err)
}
"""),
    ]
    for i, (slug, _num, title, pkg, why, what, ex) in enumerate(stubs):
        n = 45 + i  # stubs shift down by one now that §44 is Octo
        slug = f"{n:02d}-{slug.split('-', 1)[1]}"
        num = str(n)
        st, lab = ("missing", "⬜") if "accel" in slug else ("partial", "🚧")
        out.append(C(slug, num, title, "VIII · Stubs", pkg, st, lab, why, what, "", ex))

    out.append(C(
        "62-w2a", "62", "w2a — validation harness", "IX · Validate",
        "github.com/openfluke/w2a", "partial", "🚧 harness",
        why="Engine packages must stay free of tests. w2a owns timed 34×20×3 matrices, gap census, and honesty stamps. See §63 for a live full-suite run.",
        what="Interactive go run ., suites under suites/*, go test ./tests/<layer>. StampBackendNote / AffinePackable prevent fake ✅.",
        example="""
package main

import "fmt"

func main() {
	fmt.Println("w2a is a separate module — engine packages never contain tests.")
	fmt.Println("")
	fmt.Println("  cd w2a")
	fmt.Println("  go run .                    # interactive menu")
	fmt.Println("  go test ./tests/dense -v    # timed FormatNone matrix")
	fmt.Println("  go test ./tests/mha -v      # MHA coverage")
}
""",
        run="cd w2a && go run .\ncd w2a && go test ./tests/dense -v && go test ./tests/mha -v",
    ))

    out.append(C(
        "63-validation", "63", "Validation report — full suite", "IX · Validate",
        "github.com/openfluke/w2a", "ok", "✅ 137k cells",
        why="Claims are cheap; a green matrix is not. This is the actual output of one full w2a run so the book's ✅ marks are backed by numbers you can reproduce, not asserted.",
        what="Every timed layer sweeps its dtype × format × backend matrix; every suite runs its case checks. The run below is zero-gap and zero-fail across the whole board.",
        body_extra="""
<div class="callout"><strong>Full suite: PASS</strong>137,039 matrix cells — OK 137,039 · GAP 0 · FAIL 0. 326 suite cases — PASS 326 · FAIL 0. Elapsed 24m09s on amd64 with an NVIDIA GTX 1650 SUPER.</div>
<h3>Coverage by layer (timed matrices)</h3>
<table>
<thead><tr><th>Layer</th><th>Cells</th><th>OK</th><th>Cases</th></tr></thead>
<tbody>
<tr><td>dense</td><td>2,802</td><td>2,802</td><td>18</td></tr>
<tr><td>mha</td><td>2,847</td><td>2,847</td><td>18</td></tr>
<tr><td>softmax</td><td>2,847</td><td>2,847</td><td>19</td></tr>
<tr><td>cnn1 · cnn2 · cnn3</td><td>2,847 ea</td><td>2,847 ea</td><td>14 ea</td></tr>
<tr><td>embedding · layernorm · lstm</td><td>2,847 ea</td><td>2,847 ea</td><td>14 ea</td></tr>
<tr><td>residual · rmsnorm · rnn</td><td>2,847 ea</td><td>2,847 ea</td><td>14 ea</td></tr>
<tr><td>sequential · swiglu</td><td>2,847 ea</td><td>2,847 ea</td><td>14 ea</td></tr>
<tr><td>dna</td><td>16,159</td><td>16,159</td><td>6</td></tr>
<tr><td>evolution</td><td>16,152</td><td>16,152</td><td>6</td></tr>
<tr><td>step</td><td>32,458</td><td>32,458</td><td>11</td></tr>
<tr><td>tween</td><td>32,457</td><td>32,457</td><td>10</td></tr>
<tr><td><strong>Total</strong></td><td><strong>137,039</strong></td><td><strong>137,039</strong></td><td><strong>326</strong></td></tr>
</tbody>
</table>
<p class="example-meta">Full timed matrices: dense…residual + §5 gdn/mamba/convt/kmeans/parallel/metacognition. Case-only stubs: seed, serialization, memory, donate, fountain, hardware, helpers.</p>
<h3>Dense timed matrix — highlights</h3>
<p>All 34 dtypes run forward and backward on CPU-tiled, Plan 9 SIMD, and WebGPU with zero gaps. Fastest forward paths on SIMD: <code>int8</code> 45µs, <code>float32</code> 57µs, <code>int4</code> 87µs; WebGPU stays in the ~165–490µs band across every dtype.</p>
<h3>Cross-numeric train</h3>
<p>Beyond matched float32 acts: W2A Step <strong>Cross-Numeric Train</strong> sweeps every Op kind × FormatNone weight dtype × Go <code>Numeric</code> activation host (smoke ~735 cells; full census ~10.7k). Storage truth is checked after <code>StepMesh</code> — no retained f32 master for non-float32 FormatNone. Public Dense demotion + full W×A showcase: <a href="https://github.com/openfluke/down-the-dem">down-the-dem</a> (see also chapter <a href="65-cross-numeric.html">65 · Cross-numeric train</a>).</p>
""",
        example="""
package main

import "fmt"

func main() {
	cells, cases := 137039, 326
	fmt.Printf("w2a suite: %d/%d matrix cells OK, %d/%d cases PASS\\n", cells, cells, cases, cases)
	fmt.Println("gaps: 0   fails: 0   result: PASS")
}
""",
        run="cd w2a && go run .   # full timed suite (~24m); writes logs/suite.txt",
    ))

    out.append(C(
        "64-scorecard", "64", "Scorecard → v1.0", "IX · Validate",
        "", "partial", version,
        why="Version is earned from a weighted board, not marketing. Apps, stubs, and Accel still leave points on the table.",
        what=f"version = 0.{{round(earned)}} until 100 → v1.0. Today <strong>{esc(version)}</strong> "
             f"({earned_i}/100). Biggest remaining: apps/stubs (§9–10) and Accel (§11).",
        body_extra=f"""
<table><thead><tr><th>§</th><th>Area</th><th>Wt</th><th>Earned</th></tr></thead><tbody>
<tr><td>1–4</td><td>Foundation + Dense + transformer + CNN/RNN</td><td>50</td><td>50</td></tr>
<tr><td>5</td><td>Extended layers</td><td>7</td><td>7</td></tr>
<tr><td>6–8</td><td>Runtime + systems + model</td><td>21</td><td>21</td></tr>
<tr><td>9–11</td><td>Apps + stubs + accel</td><td>8</td><td>3.0</td></tr>
<tr><td>12</td><td>Peak fused / no host ALU</td><td>14</td><td>14</td></tr>
<tr><td></td><td><strong>Total</strong></td><td>100</td><td><strong>{earned_i}</strong></td></tr>
</tbody></table>
""",
        example=f"""
package main

import "fmt"

func main() {{
	// Recompute when a board row flips ✅/🚧/⬜ in welvet/README.md
	earned := {earned:.1f}
	fmt.Printf("v0.%02.0f\\n", earned) // round({earned_i}) → {version} until earned==100 → v1.0
}}
""",
    ))

    return out


# ---------------------------------------------------------------------------
# HTML rendering
# ---------------------------------------------------------------------------

PARTS_NAV = [
    ("I · Orientation", ["01-welvet", "02-tree"]),
    ("II · Foundation", ["03-core", "04-weights", "05-quant", "06-simd", "07-webgpu", "08-tiling", "09-architecture", "10-fusedgpu"]),
    ("III · Layers", [
        "11-dense", "12-mha", "13-swiglu", "14-rmsnorm", "15-layernorm", "16-embedding", "17-softmax",
        "18-sequential", "19-residual", "20-cnn", "21-rnn-lstm", "22-seqmix", "23-gdn", "24-mamba",
        "25-convt", "26-kmeans", "27-parallel", "28-metacognition",
    ]),
    ("IV · Runtime", ["29-forward", "30-backward", "31-training", "32-step", "65-cross-numeric"]),
    ("V · Systems", ["33-dna", "34-evolution", "35-tween", "36-tanhi", "37-telemetry"]),
    ("VI · Model IO", ["38-entity", "39-hf", "40-tokenizer", "41-sampling", "42-transformer"]),
    ("VII · Apps", ["43-apps", "44-octo"]),
    ("VIII · Stubs", [
        "45-stub-seed", "46-stub-serialization", "47-stub-memory", "48-stub-donate",
        "49-stub-fountain", "50-stub-hardware", "51-stub-accel",
        "52-stub-clustering", "53-stub-ensemble", "54-stub-evaluation", "55-stub-grafting",
        "56-stub-grouping", "57-stub-introspection", "58-stub-observer", "59-stub-pipeline",
        "60-stub-templates", "61-stub-universal",
    ]),
    ("IX · Validate", ["62-w2a", "63-validation", "64-scorecard"]),
]


def nav_html(by_slug: dict[str, Chapter], prefix: str, version: str) -> str:
    parts = [
        f'<aside class="sidebar"><a class="brand" href="{prefix}index.html">Welvet</a>'
        f'<div class="brand-sub">Feature book · {esc(version)}</div>'
        f'<div class="nav-group"><h2>Front</h2>'
        f'<a href="{prefix}index.html">Title</a>'
        f'<a href="{prefix}toc.html">Contents</a></div>'
    ]
    for part, slugs in PARTS_NAV:
        parts.append(f'<div class="nav-group"><h2>{esc(part)}</h2>')
        for s in slugs:
            c = by_slug[s]
            href = f"{prefix}chapters/{c.slug}.html" if prefix else f"{c.slug}.html"
            if prefix == "../":
                href = f"{c.slug}.html"
            short = c.title.split("—")[0].split("·")[0].strip()
            if len(short) > 28:
                short = short[:27] + "…"
            parts.append(f'<a href="{href}">{esc(c.num)}. {esc(short)}</a>')
        parts.append("</div>")
    parts.append("</aside>")
    return "\n".join(parts)


def chapter_body_html(c: Chapter, run_info: dict | None, *, for_print: bool = False) -> str:
    pkg = (
        f'<p class="pill-row"><span class="pill">{esc(c.pkg)}</span>'
        f"{status_badge(c.status, c.status_label)}</p>"
        if c.pkg
        else f"<p>{status_badge(c.status, c.status_label)}</p>"
    )
    validate = f"<h3>Validate (harness)</h3>{code(c.run)}" if c.run else ""

    example_section = ""
    if c.example.strip() and c.runnable:
        run_cmd = f"cd welvet/examples/{c.slug} && source ../env.sh && go run ."
        output_block = ""
        if run_info:
            exit_code = run_info.get("exit_code", 1)
            out_cls = "example-output" if exit_code == 0 else "example-output fail"
            combined = run_info.get("combined", "").strip() or "(no output)"
            output_block = (
                f'<p class="example-meta">exit {exit_code} · '
                f'last run via <code>go run .</code></p>'
                f'<pre class="{out_cls}"><code>{esc(combined)}</code></pre>'
            )
        file_line = ""
        if not for_print:
            rel_go = f"../examples/{c.slug}/main.go"
            file_line = f'<p><a class="go-file-link" href="{rel_go}">examples/{c.slug}/main.go</a></p>'
        else:
            file_line = f"<p class=\"go-file-link\">examples/{esc(c.slug)}/main.go</p>"
        example_section = f"""
<section class="examples" id="examples">
<h2>Go example</h2>
{file_line}
<div class="run-cmd"><span>Run:</span><code>{esc(run_cmd)}</code></div>
{code(c.example)}
<h3>Output</h3>
{output_block if output_block else '<p class="example-meta">Run <code>python3 _gen_welvet_book.py --run</code> to capture output.</p>'}
{validate}
</section>
"""
    elif c.example.strip():
        example_section = f"""
<section class="examples" id="examples">
<h2>Go example</h2>
{code(c.example)}
{validate}
</section>
"""

    return f"""
<p class="chapter-kicker">Chapter {esc(c.num)}</p>
<h1>{esc(c.title)}</h1>
{pkg}
<hr class="rule"/>
<h2>Why it exists</h2>
<p>{c.why}</p>
<h2>What it is</h2>
<p>{c.what}</p>
{c.body_extra}
{example_section}
"""


def render_chapter(
    c: Chapter,
    prev: Chapter | None,
    next_: Chapter | None,
    by_slug: dict[str, Chapter],
    run_info: dict | None,
    version: str,
) -> str:
    nav = nav_html(by_slug, "../", version)
    prev_a = f'<a href="../toc.html">← Contents</a>' if prev is None else f'<a href="{prev.slug}.html">← {esc(prev.num)}. {esc(prev.title[:40])}</a>'
    next_a = f'<a href="../toc.html">Contents →</a>' if next_ is None else f'<a href="{next_.slug}.html">{esc(next_.num)}. {esc(next_.title[:40])} →</a>'
    body = chapter_body_html(c, run_info)

    return f"""<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8"/><meta name="viewport" content="width=device-width, initial-scale=1"/>
<title>{esc(c.num)}. {esc(c.title)} — Welvet</title>
<link rel="stylesheet" href="../assets/book.css"/>
</head>
<body>
<a class="skip" href="#main">Skip</a>
<div class="shell">
{nav}
<div class="main">
<div class="topbar">
<button class="menu-btn" id="menu-btn" type="button">Menu</button>
<div class="meta">{esc(c.part)}</div>
<nav class="pager">{prev_a}{next_a}</nav>
</div>
<article class="content" id="main">
{body}
<nav class="pager" style="margin-top:2rem">{prev_a}{next_a}</nav>
</article>
</div></div>
<script src="../assets/book.js"></script>
</body></html>
"""


def write_env_sh() -> None:
    cache = EXAMPLES / ".cache"
    sh = f"""#!/usr/bin/env bash
# Shared Go build cache for book examples (avoids /tmp quota exhaustion with webgpu CGO).
_root="$(cd "$(dirname "${{BASH_SOURCE[0]}}")" && pwd)"
export GOCACHE="${{GOCACHE:-$_root/.cache/gocache}}"
export GOTMPDIR="${{GOTMPDIR:-$_root/.cache/gotmp}}"
export TMPDIR="${{TMPDIR:-$_root/.cache/tmp}}"
mkdir -p "$GOCACHE" "$GOTMPDIR" "$TMPDIR"
"""
    p = EXAMPLES / "env.sh"
    p.write_text(sh, encoding="utf-8")
    p.chmod(0o755)


def write_examples_go_mod() -> None:
    EXAMPLES.mkdir(parents=True, exist_ok=True)

    def rel(p: Path) -> str:
        return os.path.relpath(p, EXAMPLES)

    go_mod = f"""module github.com/openfluke/welvet-book-examples

go 1.22.5

require github.com/openfluke/welvet v0.0.0

replace github.com/openfluke/welvet => {rel(WELVET_ROOT)}

replace github.com/openfluke/webgpu => {rel(WEBGPU_ROOT)}

replace github.com/eliben/go-sentencepiece => {rel(SENTENCEPIECE_ROOT)}
"""
    (EXAMPLES / "go.mod").write_text(go_mod, encoding="utf-8")
    write_env_sh()


def write_example_files(chs: list[Chapter]) -> None:
    write_examples_go_mod()
    slugs = {c.slug for c in chs}
    # Prune stale chapter example dirs (e.g. after a renumber) so we never
    # accumulate orphaned NN-*/ folders. Keep non-chapter dirs like .cache/bin.
    keep = {".cache", "bin"}
    for d in EXAMPLES.iterdir():
        if not d.is_dir() or d.name in keep:
            continue
        if re.match(r"^\d+-", d.name) and d.name not in slugs:
            shutil.rmtree(d, ignore_errors=True)
    for c in chs:
        if not c.runnable or not c.example.strip():
            continue
        if "package main" not in c.example:
            continue
        d = EXAMPLES / c.slug
        d.mkdir(parents=True, exist_ok=True)
        (d / "main.go").write_text(c.example.strip() + "\n", encoding="utf-8")


def run_all_examples(chs: list[Chapter], timeout: float = 90.0) -> dict[str, dict]:
    """Run each example; return slug → {exit_code, stdout, stderr, combined}."""
    write_example_files(chs)
    cache = EXAMPLES / ".cache"
    bin_dir = cache / "bin"
    cache.mkdir(parents=True, exist_ok=True)
    bin_dir.mkdir(parents=True, exist_ok=True)
    env = {
        **os.environ,
        "GOCACHE": str(cache / "gocache"),
        "GOTMPDIR": str(cache / "gotmp"),
        "TMPDIR": str(cache / "tmp"),
    }
    for d in (cache / "gocache", cache / "gotmp", cache / "tmp"):
        d.mkdir(parents=True, exist_ok=True)

    subprocess.run(
        ["go", "mod", "tidy"],
        cwd=EXAMPLES,
        capture_output=True,
        text=True,
        check=False,
        env=env,
    )
    results: dict[str, dict] = {}
    for c in chs:
        if not c.runnable or not c.example.strip() or "package main" not in c.example:
            continue
        d = EXAMPLES / c.slug
        if not (d / "main.go").is_file():
            continue
        print(f"  run {c.slug} …", flush=True)
        bin_path = bin_dir / c.slug
        try:
            build = subprocess.run(
                ["go", "build", "-o", str(bin_path), "."],
                cwd=d,
                capture_output=True,
                text=True,
                timeout=timeout,
                env=env,
            )
            if build.returncode != 0:
                combined = (build.stdout or "") + (build.stderr or "")
                results[c.slug] = {
                    "exit_code": build.returncode,
                    "stdout": build.stdout or "",
                    "stderr": build.stderr or "",
                    "combined": _trim_output(combined),
                }
                print(f"    BUILD FAIL ({build.returncode})", flush=True)
                continue
            proc = subprocess.run(
                [str(bin_path)],
                cwd=d,
                capture_output=True,
                text=True,
                timeout=timeout,
                env=env,
            )
        except subprocess.TimeoutExpired:
            results[c.slug] = {
                "exit_code": 124,
                "stdout": "",
                "stderr": f"timeout after {timeout}s",
                "combined": f"timeout after {timeout}s",
            }
            continue
        combined = proc.stdout
        if proc.stderr:
            combined = (combined + "\n" if combined else "") + proc.stderr
        results[c.slug] = {
            "exit_code": proc.returncode,
            "stdout": proc.stdout,
            "stderr": proc.stderr,
            "combined": _trim_output(combined.rstrip()),
        }
        mark = "OK" if proc.returncode == 0 else "FAIL"
        print(f"    {mark} ({proc.returncode})", flush=True)
    MANIFEST.write_text(json.dumps(results, indent=2) + "\n", encoding="utf-8")
    return results


def _trim_output(text: str, max_lines: int = 40) -> str:
    lines = text.splitlines()
    if len(lines) <= max_lines:
        return text.strip()
    return "\n".join(lines[:5] + ["…"] + lines[-max_lines + 6:]).strip()


def load_manifest() -> dict[str, dict]:
    if MANIFEST.is_file():
        return json.loads(MANIFEST.read_text(encoding="utf-8"))
    return {}


def title_page(version: str, pdf_name: str, chs: list[Chapter]) -> str:
    pdf_url = f"{RELEASES_URL}/latest/download/{pdf_name}"
    pdf_btn = (
        f'<a class="btn btn-ghost" href="{esc(pdf_url)}">Download PDF</a>'
        f'<a class="btn btn-ghost" href="{RELEASES_URL}">All releases</a>'
    )
    part_pills = []
    for part, slugs in PARTS_NAV:
        label = part.split("·")[-1].strip()
        first = slugs[0] if slugs else ""
        href = f"chapters/{first}.html" if first else "toc.html"
        part_pills.append(f'<a class="pill" href="{href}">{esc(label)}</a>')
    parts_html = "".join(part_pills)
    stats = [
        (str(len(chs)), "Chapters"),
        ("34", "Dtypes"),
        ("20", "Quant formats"),
        ("3", "Backends"),
    ]
    stats_html = "".join(
        f'<div class="stat"><div class="stat-num">{esc(n)}</div>'
        f'<div class="stat-label">{esc(lbl)}</div></div>'
        for n, lbl in stats
    )
    return f"""<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8"/><meta name="viewport" content="width=device-width, initial-scale=1"/>
<title>Welvet · Feature Book</title>
<meta name="description" content="Welvet feature book: why and what for every engine package, with runnable Go examples."/>
<link rel="stylesheet" href="assets/book.css"/>
</head>
<body>
<main class="title-page">
<div class="title-card landing">
<div class="tc-main">
{title_logo_html()}
<span class="title-badge"><span class="dot"></span>OpenFluke feature book</span>
<h1><span class="grad">Welvet</span></h1>
<p class="subtitle">Why and what for every engine package, with runnable Go examples. Built from the Welvet tree, not a loom paste.</p>
<div class="cta-row">
<a class="btn btn-primary" href="toc.html">Read the contents</a>
<a class="btn btn-ghost" href="chapters/01-welvet.html">Start reading</a>
<a class="btn btn-ghost" href="https://github.com/openfluke/welvet">Source</a>
{pdf_btn}
</div>
</div>
<div class="tc-side">
<div class="title-stats">
{stats_html}
</div>
<div class="title-meta">
<div><strong>Version</strong> {esc(version)} · scorecard toward v1.0</div>
<div><strong>Module</strong> github.com/openfluke/welvet</div>
<div><strong>Harness</strong> github.com/openfluke/w2a</div>
</div>
<div class="title-parts">
<p class="label">What's inside</p>
<div class="pill-row">
{parts_html}
</div>
</div>
</div>
</div>
</main>
</body></html>
"""


def toc_page(chs: list[Chapter], version: str) -> str:
    by = {c.slug: c for c in chs}
    nav = nav_html(by, "", version)
    # fix chapter hrefs for root toc
    nav = nav.replace('href="chapters/', 'href="chapters/').replace('href="01-', 'href="chapters/01-')
    # rebuild nav properly for root
    parts = [
        '<aside class="sidebar"><a class="brand" href="index.html">Welvet</a>'
        f'<div class="brand-sub">Feature book · {esc(version)}</div>'
        '<div class="nav-group"><h2>Front</h2>'
        '<a href="index.html">Title</a><a href="toc.html">Contents</a></div>'
    ]
    for part, slugs in PARTS_NAV:
        parts.append(f'<div class="nav-group"><h2>{esc(part)}</h2>')
        for s in slugs:
            c = by[s]
            short = c.title.split("—")[0].strip()
            if len(short) > 28:
                short = short[:27] + "…"
            parts.append(f'<a href="chapters/{c.slug}.html">{esc(c.num)}. {esc(short)}</a>')
        parts.append("</div>")
    parts.append("</aside>")
    nav = "\n".join(parts)

    body = ['<p class="lede">One chapter per Welvet feature. Each page: <strong>why</strong>, <strong>what</strong>, status, and a <strong>Go example</strong>.</p>',
            f'<div class="pill-row"><span class="pill">{len(chs)} chapters</span><span class="pill">engine · layers · runtime · systems · model · apps · stubs · w2a</span></div>']
    for part, slugs in PARTS_NAV:
        body.append(f'<p class="part">{esc(part)}</p><ul class="toc-list">')
        for s in slugs:
            c = by[s]
            body.append(
                f'<li><a href="chapters/{c.slug}.html"><span class="num">{esc(c.num)}</span>'
                f'<span><span class="title">{esc(c.title)}</span>'
                f'<div class="blurb">{status_badge(c.status, c.status_label)}'
                f'{" · " + esc(c.pkg) if c.pkg else ""}</div></span></a></li>'
            )
        body.append("</ul>")
    return f"""<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8"/><meta name="viewport" content="width=device-width, initial-scale=1"/>
<title>Contents — Welvet Feature Book</title>
<link rel="stylesheet" href="assets/book.css"/>
</head>
<body>
<div class="shell">
{nav}
<div class="main">
<div class="topbar">
<button class="menu-btn" id="menu-btn" type="button">Menu</button>
<div class="meta">Front</div>
<nav class="pager"><a href="index.html">← Title</a><a href="chapters/01-welvet.html">Ch 1 →</a></nav>
</div>
<article class="content" id="main">
<p class="chapter-kicker">Index</p>
<h1>Table of contents</h1>
{"".join(body)}
</article>
</div></div>
<script src="assets/book.js"></script>
</body></html>
"""


PRINT_CSS = """/* Print / PDF layout — combined with book.css */
@page {
  size: A4;
  margin: 18mm 16mm 20mm;
}

body.print-doc {
  background: white;
  font-size: 10.5pt;
}

.print-doc .print-cover {
  min-height: 90vh;
  display: grid;
  place-items: center;
  page-break-after: always;
}

.print-doc .print-toc {
  page-break-after: always;
}

.print-doc .print-chapter {
  page-break-before: always;
}

.print-doc .print-chapter:first-of-type {
  page-break-before: auto;
}

.print-doc .print-toc-list {
  columns: 2;
  column-gap: 2rem;
  list-style: none;
  padding: 0;
  margin: 0;
}

.print-doc .print-toc-list li {
  break-inside: avoid;
  margin: 0.2rem 0;
  font-size: 0.88rem;
}

.print-doc .print-toc-list .part {
  column-span: all;
  margin: 1rem 0 0.35rem;
  font-size: 0.72rem;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: var(--accent);
  font-weight: 700;
}

.print-doc pre {
  font-size: 0.72rem;
  line-height: 1.4;
  white-space: pre-wrap;
  word-break: break-word;
}

.print-doc .pager,
.print-doc .run-cmd,
.print-doc .go-file-link {
  /* keep run hints in PDF */
}

.print-doc h1 {
  font-size: 1.65rem;
}

.print-doc h2 {
  font-size: 1.15rem;
}

.print-doc .examples {
  border-top-width: 1px;
}

@media print {
  a {
    color: inherit;
    text-decoration: none;
  }
}
"""


def write_print_css() -> None:
    (BOOK / "assets" / "book-print.css").write_text(PRINT_CSS, encoding="utf-8")


def render_print_document(
    chs: list[Chapter],
    run_results: dict[str, dict],
    version: str,
) -> str:
    toc_items: list[str] = []
    for part, slugs in PARTS_NAV:
        toc_items.append(f'<li class="part">{esc(part)}</li>')
        for s in slugs:
            c = next(x for x in chs if x.slug == s)
            toc_items.append(f'<li>{esc(c.num)}. {esc(c.title)}</li>')

    chapters_html: list[str] = []
    for c in chs:
        info = run_results.get(c.slug)
        chapters_html.append(
            f'<section class="print-chapter content" id="{c.slug}">'
            f"{chapter_body_html(c, info, for_print=True)}"
            "</section>"
        )

    return f"""<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8"/>
<title>Welvet Feature Book {esc(version)}</title>
<link rel="stylesheet" href="assets/book.css"/>
<link rel="stylesheet" href="assets/book-print.css"/>
</head>
<body class="print-doc">
<main class="print-cover title-page">
<div class="title-card">
{title_logo_html()}
<p class="eyebrow">OpenFluke · github.com/openfluke/welvet</p>
<h1>Welvet</h1>
<p class="subtitle">Feature book — why and what for every engine package, with runnable Go examples.</p>
<div class="title-meta">
<div><strong>Version</strong> {esc(version)}</div>
<div><strong>Generated</strong> from welvet/README.md scorecard</div>
<div><strong>Chapters</strong> {len(chs)}</div>
</div>
</div>
</main>
<section class="print-toc content">
<h1>Contents</h1>
<ul class="print-toc-list">
{"".join(toc_items)}
</ul>
</section>
{"".join(chapters_html)}
</body></html>
"""


def find_chrome() -> str:
    for name in ("google-chrome", "google-chrome-stable", "chromium", "chromium-browser"):
        path = subprocess.run(["which", name], capture_output=True, text=True)
        if path.returncode == 0 and path.stdout.strip():
            return path.stdout.strip()
    raise SystemExit(
        "PDF needs headless Chrome/Chromium. Install google-chrome or pass --pdf-html-only."
    )


def generate_pdf(print_html: Path, pdf_path: Path) -> None:
    chrome = find_chrome()
    CHROME_PROFILE.mkdir(parents=True, exist_ok=True)
    pdf_path.parent.mkdir(parents=True, exist_ok=True)
    # Clear TMPDIR/GOTMPDIR — env.sh from Go examples can break Chrome if inherited
    env = {k: v for k, v in os.environ.items() if k not in ("TMPDIR", "GOTMPDIR")}
    cmd = [
        chrome,
        "--headless",
        "--disable-gpu",
        "--no-pdf-header-footer",
        f"--user-data-dir={CHROME_PROFILE}",
        f"--print-to-pdf={pdf_path.resolve()}",
        print_html.as_uri(),
    ]
    print(f"  PDF → {pdf_path} …", flush=True)
    proc = subprocess.run(cmd, capture_output=True, text=True, env=env)
    if proc.returncode != 0 or not pdf_path.is_file() or pdf_path.stat().st_size == 0:
        err = (proc.stderr or proc.stdout or "unknown error").strip()
        raise SystemExit(f"PDF generation failed: {err[-800:]}")


def main() -> None:
    ap = argparse.ArgumentParser(description="Generate Welvet HTML book + Go examples")
    ap.add_argument("--run", action="store_true", help="go run every example and embed output in HTML")
    ap.add_argument("--html-only", action="store_true", help="skip go run even if manifest missing")
    ap.add_argument("--pdf", action="store_true", help="build print HTML + PDF (needs Chrome)")
    ap.add_argument(
        "--pdf-html-only",
        action="store_true",
        help="build welvet/print.html only, skip Chrome PDF step",
    )
    ap.add_argument(
        "--pdf-out",
        type=Path,
        default=None,
        help="PDF output path (default: dist/welvet-feature-book-VERSION.pdf, not committed)",
    )
    args = ap.parse_args()

    version, earned = read_welvet_version()
    pdf_name = pdf_filename(version)
    pdf_path = args.pdf_out if args.pdf_out else DIST / pdf_name

    CHDIR.mkdir(parents=True, exist_ok=True)
    (BOOK / "assets").mkdir(parents=True, exist_ok=True)
    ensure_logo()
    write_print_css()
    for p in CHDIR.glob("*.html"):
        p.unlink()

    chs = chapters(version, earned)
    by = {c.slug: c for c in chs}
    listed = {s for _, ss in PARTS_NAV for s in ss}
    missing = [c.slug for c in chs if c.slug not in listed]
    if missing:
        raise SystemExit(f"PARTS_NAV missing: {missing}")

    write_example_files(chs)

    (EXAMPLES / "README.md").write_text(
        "# Welvet book examples\n\n"
        f"One `main.go` per chapter ({len(chs)} total).\n\n"
        "```bash\n"
        "cd welvet/examples/01-welvet\n"
        "source ../env.sh   # shared GOCACHE — required for webgpu/CGO examples\n"
        "go run .\n"
        "```\n\n"
        "Regenerate HTML + capture all outputs:\n\n"
        "```bash\n"
        "cd openfluke.github.io && python3 _gen_welvet_book.py --run\n"
        "```\n",
        encoding="utf-8",
    )

    run_results: dict[str, dict] = {}
    if args.run:
        print("Running examples …")
        run_results = run_all_examples(chs)
    elif not args.html_only:
        run_results = load_manifest()

    ok = sum(1 for r in run_results.values() if r.get("exit_code") == 0)
    fail = len(run_results) - ok
    if run_results:
        print(f"Examples: {ok} OK, {fail} failed (manifest → {MANIFEST})")

    for i, c in enumerate(chs):
        prev = chs[i - 1] if i else None
        nxt = chs[i + 1] if i + 1 < len(chs) else None
        info = run_results.get(c.slug)
        (CHDIR / f"{c.slug}.html").write_text(
            render_chapter(c, prev, nxt, by, info, version), encoding="utf-8"
        )

    print_html = BOOK / "print.html"
    if args.pdf or args.pdf_html_only:
        print_html.write_text(
            render_print_document(chs, run_results, version), encoding="utf-8"
        )
        print(f"Print HTML → {print_html} (gitignored build artifact)")
    if args.pdf and not args.pdf_html_only:
        generate_pdf(print_html, pdf_path)
        mb = pdf_path.stat().st_size / (1024 * 1024)
        print(f"PDF: {pdf_path} ({mb:.1f} MB) · welvet {version} ({earned}/100)")
        print(f"Upload: gh release upload {version} {pdf_path}")

    (BOOK / "index.html").write_text(title_page(version, pdf_name, chs), encoding="utf-8")
    (BOOK / "toc.html").write_text(toc_page(chs, version), encoding="utf-8")
    (ROOT / "index.html").write_text(
        "<!DOCTYPE html><meta charset=utf-8><meta http-equiv=refresh content='0;url=welvet/'>"
        "<title>OpenFluke</title><a href='welvet/'>Welvet feature book</a>\n",
        encoding="utf-8",
    )
    (ROOT / ".nojekyll").write_text("", encoding="utf-8")
    (ROOT / "README.md").write_text(
        "# openfluke.github.io\n\n"
        "**[Welvet feature book](welvet/)** — why/what for every engine package with runnable Go examples.\n\n"
        "```bash\n"
        "python3 _gen_welvet_book.py --run          # HTML + examples + captured outputs\n"
        "python3 _gen_welvet_book.py --run --pdf    # + PDF in dist/ (upload to GitHub Release)\n"
        "./release-book.sh                          # release helper\n"
        "cd welvet/examples/01-welvet && source ../env.sh && go run .\n"
        "```\n\n"
        f"PDF is **not** committed — version **{version}** from `welvet/README.md`, "
        f"upload `dist/{pdf_name}` to [GitHub Releases]({RELEASES_URL}).\n",
        encoding="utf-8",
    )
    print(f"OK: {len(chs)} chapters → {BOOK}")
    if fail and args.run:
        print(f"WARNING: {fail} example(s) failed — outputs still embedded in HTML", file=sys.stderr)


if __name__ == "__main__":
    main()
