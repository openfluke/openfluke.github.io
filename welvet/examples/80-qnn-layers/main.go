package main

import (
	"fmt"

	"github.com/openfluke/welvet/nativeqnn"
)

func main() {
	if !nativeqnn.Available() {
		fmt.Println("SKIP: need windows/arm64 + libwelvet_qnn (+ QAIRT)")
		return
	}
	sess, err := nativeqnn.DefaultSession()
	must(err)

	demoDense(sess)
	demoSoftmax(sess)
	demoRms(sess)
	demoLN(sess)
	demoCnn2(sess)
	demoCnn3(sess)
	demoConvT2(sess)
	demoConvT3(sess)
	demoEmbed(sess)
	demoResidual(sess)
	demoSwiGlu(sess)
	demoMha(sess)
	fmt.Println("qnn-layers done")
}

func demoDense(s *nativeqnn.Session) {
	const batch, in, out = 2, 4, 4
	W := fill(in * out)
	x, y := fill(batch*in), make([]float32, batch*out)
	must(s.DensePrepareF32(batch, in, out, W, nil))
	must(s.DenseExecuteF32(x, y))
	fmt.Println("dense", len(y))
}

func demoSoftmax(s *nativeqnn.Session) {
	const rows, cols = 2, 8
	x, y := fill(rows*cols), make([]float32, rows*cols)
	must(s.SoftmaxForwardF32(rows, cols, x, y))
	fmt.Println("softmax", len(y))
}

func demoRms(s *nativeqnn.Session) {
	const b, seq, dim = 1, 2, 4
	g, x, y := ones(dim), fill(b*seq*dim), make([]float32, b*seq*dim)
	must(s.RmsNormForwardF32(b, seq, dim, 1e-5, g, x, y))
	fmt.Println("rmsnorm", len(y))
}

func demoLN(s *nativeqnn.Session) {
	const b, seq, dim = 1, 2, 4
	g, beta := ones(dim), make([]float32, dim)
	x, y := fill(b*seq*dim), make([]float32, b*seq*dim)
	must(s.LayerNormForwardF32(b, seq, dim, 1e-5, g, beta, x, y))
	fmt.Println("layernorm", len(y))
}

func demoCnn2(s *nativeqnn.Session) {
	// HTP path today is 1×1 (W = cout*cin), same as w2a go_v_qnn.
	const n, h, w, cin, cout = 1, 4, 4, 2, 2
	W, bias := fill(cout*cin), make([]float32, cout)
	x, y := fill(n*h*w*cin), make([]float32, n*h*w*cout)
	must(s.Cnn2ForwardF32(n, h, w, cin, cout, W, bias, x, y))
	fmt.Println("cnn2", len(y))
}

func demoCnn3(s *nativeqnn.Session) {
	const n, d, h, w, cin, cout = 1, 2, 2, 2, 2, 2
	W, bias := fill(cout*cin), make([]float32, cout)
	x, y := fill(n*d*h*w*cin), make([]float32, n*d*h*w*cout)
	must(s.Cnn3ForwardF32(n, d, h, w, cin, cout, W, bias, x, y))
	fmt.Println("cnn3", len(y))
}

func demoConvT2(s *nativeqnn.Session) {
	const n, h, w, cin, cout = 1, 4, 4, 2, 2
	W, bias := fill(cout*cin), make([]float32, cout)
	x, y := fill(n*h*w*cin), make([]float32, n*h*w*cout)
	must(s.ConvT2ForwardF32(n, h, w, cin, cout, W, bias, x, y))
	fmt.Println("convt2", len(y))
}

func demoConvT3(s *nativeqnn.Session) {
	const n, d, h, w, cin, cout = 1, 2, 2, 2, 2, 2
	W, bias := fill(cout*cin), make([]float32, cout)
	x, y := fill(n*d*h*w*cin), make([]float32, n*d*h*w*cout)
	must(s.ConvT3ForwardF32(n, d, h, w, cin, cout, W, bias, x, y))
	fmt.Println("convt3", len(y))
}

func demoEmbed(s *nativeqnn.Session) {
	const vocab, dim, n = 16, 4, 4
	table := fill(vocab * dim)
	ids := []int32{0, 1, 2, 3}
	y := make([]float32, n*dim)
	must(s.EmbeddingForwardF32(vocab, dim, n, table, ids, y))
	fmt.Println("embedding", len(y))
}

func demoResidual(s *nativeqnn.Session) {
	const n = 16
	a, b, y := fill(n), fill(n), make([]float32, n)
	must(s.ResidualForwardF32(n, a, b, y))
	fmt.Println("residual", len(y))
}

func demoSwiGlu(s *nativeqnn.Session) {
	const batch, in, hidden = 1, 4, 8
	wG, wU, wD := fill(hidden*in), fill(hidden*in), fill(in*hidden)
	x, y := fill(batch*in), make([]float32, batch*in)
	must(s.SwiGluForwardF32(batch, in, hidden, wG, wU, wD, x, y))
	fmt.Println("swiglu", len(y))
}

func demoMha(s *nativeqnn.Session) {
	const batch, seq, dim, heads = 1, 2, 8, 2
	w := fill(dim * dim)
	x, y := fill(batch*seq*dim), make([]float32, batch*seq*dim)
	must(s.MhaForwardF32(batch, seq, dim, heads, w, w, w, w, x, y))
	fmt.Println("mha", len(y))
}

func fill(n int) []float32 {
	out := make([]float32, n)
	for i := range out {
		out[i] = 0.01 * float32(i+1)
	}
	return out
}

func ones(n int) []float32 {
	out := make([]float32, n)
	for i := range out {
		out[i] = 1
	}
	return out
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
