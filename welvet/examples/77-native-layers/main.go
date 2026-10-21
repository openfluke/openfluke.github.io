package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/cnn1"
	"github.com/openfluke/welvet/layers/cnn2"
	"github.com/openfluke/welvet/layers/cnn3"
	"github.com/openfluke/welvet/layers/convt1"
	"github.com/openfluke/welvet/layers/convt2"
	"github.com/openfluke/welvet/layers/convt3"
	"github.com/openfluke/welvet/layers/dense"
	"github.com/openfluke/welvet/layers/embedding"
	"github.com/openfluke/welvet/layers/gdn"
	"github.com/openfluke/welvet/layers/kmeans"
	"github.com/openfluke/welvet/layers/layernorm"
	"github.com/openfluke/welvet/layers/lstm"
	"github.com/openfluke/welvet/layers/mamba"
	"github.com/openfluke/welvet/layers/metacognition"
	"github.com/openfluke/welvet/layers/mha"
	"github.com/openfluke/welvet/layers/parallel"
	"github.com/openfluke/welvet/layers/residual"
	"github.com/openfluke/welvet/layers/rmsnorm"
	"github.com/openfluke/welvet/layers/rnn"
	"github.com/openfluke/welvet/layers/sequential"
	"github.com/openfluke/welvet/layers/softmax"
	"github.com/openfluke/welvet/layers/swiglu"
	"github.com/openfluke/welvet/nativecpp"
	"github.com/openfluke/welvet/quant"
)

func main() {
	if !nativecpp.NativeCPUAvailable() {
		fmt.Println("SKIP: need linux/amd64 libwelvet")
		return
	}
	demos := []struct {
		name string
		fn   func() error
	}{
		{"dense", demoDense},
		{"rmsnorm", demoRMS},
		{"layernorm", demoLN},
		{"swiglu", demoSwiGLU},
		{"mha", demoMHA},
		{"softmax", demoSoftmax},
		{"embedding", demoEmbed},
		{"cnn1", demoCNN1},
		{"cnn2", demoCNN2},
		{"cnn3", demoCNN3},
		{"convt1", demoConvT1},
		{"convt2", demoConvT2},
		{"convt3", demoConvT3},
		{"rnn", demoRNN},
		{"lstm", demoLSTM},
		{"sequential", demoSeq},
		{"residual", demoRes},
		{"parallel", demoPara},
		{"kmeans", demoKMeans},
		{"mamba", demoMamba},
		{"gdn", demoGDN},
		{"metacognition", demoMeta},
	}
	ok := 0
	for _, d := range demos {
		if err := d.fn(); err != nil {
			fmt.Println("FAIL", d.name, err)
			continue
		}
		fmt.Println("OK  ", d.name)
		ok++
	}
	fmt.Printf("native-layers ok=%d/%d\n", ok, len(demos))
	if ok != len(demos) {
		panic("native-layers: some demos failed")
	}
}

func fill(shape ...int) *core.Tensor[float32] {
	t := core.NewTensor[float32](shape...)
	for i := range t.Data {
		t.Data[i] = 0.1 + float32(i%7)*0.03
	}
	return t
}

func tgtLike(x *core.Tensor[float32]) *core.Tensor[float32] {
	t := core.NewTensor[float32](x.Shape...)
	for i := range t.Data {
		t.Data[i] = 0.15
	}
	return t
}

func demoDense() error {
	l, err := dense.New(4, 4, core.ActivationLinear, core.DTypeFloat32)
	if err != nil {
		return err
	}
	l.Exec.Backend = core.BackendNativeCPU
	x := fill(2, 4)
	_, err = dense.TrainNativeMSE(l, x, tgtLike(x), 0.05)
	if err != nil {
		return err
	}
	return l.SyncWeightsFromNative()
}

func demoRMS() error {
	l, err := rmsnorm.New(rmsnorm.Config{Dim: 4})
	if err != nil {
		return err
	}
	l.Exec.Backend = core.BackendNativeCPU
	x := fill(2, 4)
	_, err = rmsnorm.TrainNativeMSE(l, x, tgtLike(x), 0.05)
	if err != nil {
		return err
	}
	return l.SyncWeightsFromNative()
}

func demoLN() error {
	l, err := layernorm.New(layernorm.Config{Dim: 4})
	if err != nil {
		return err
	}
	l.Exec.Backend = core.BackendNativeCPU
	x := fill(2, 4)
	_, err = layernorm.TrainNativeMSE(l, x, tgtLike(x), 0.05)
	if err != nil {
		return err
	}
	return l.SyncWeightsFromNative()
}

func demoSwiGLU() error {
	l, err := swiglu.New(swiglu.Config{InputDim: 4, IntermediateDim: 8})
	if err != nil {
		return err
	}
	l.Exec.Backend = core.BackendNativeCPU
	x := fill(1, 4)
	_, err = swiglu.TrainNativeMSE(l, x, tgtLike(x), 0.05)
	if err != nil {
		return err
	}
	return l.SyncWeightsFromNative()
}

func demoMHA() error {
	l, err := mha.New(mha.Config{DModel: 8, NumHeads: 2, NumKVHeads: 2, HeadDim: 4, MaxSeqLen: 8})
	if err != nil {
		return err
	}
	l.Exec.Backend = core.BackendNativeCPU
	x := fill(1, 2, 8)
	_, err = mha.TrainNativeMSE(l, x, tgtLike(x), 0.02)
	if err != nil {
		return err
	}
	return l.SyncWeightsFromNative()
}

func demoSoftmax() error {
	l, err := softmax.New(softmax.Config{Dim: 4})
	if err != nil {
		return err
	}
	l.Exec.Backend = core.BackendNativeCPU
	x := fill(2, 4)
	_, err = softmax.TrainNativeMSE(l, x, tgtLike(x), 0.05)
	return err
}

func demoEmbed() error {
	l, err := embedding.New(embedding.Config{VocabSize: 8, EmbeddingDim: 4, SeqLen: 2})
	if err != nil {
		return err
	}
	l.Exec.Backend = core.BackendNativeCPU
	x := fill(1, 2)
	x.Data[0], x.Data[1] = 1, 2
	t := fill(1, 2, 4)
	_, err = embedding.TrainNativeMSE(l, x, t, 0.05)
	if err != nil {
		return err
	}
	return l.SyncWeightsFromNative()
}

func demoCNN1() error {
	l, err := cnn1.New(cnn1.Config{InChannels: 1, Filters: 2, SeqLen: 8, Kernel: 3})
	if err != nil {
		return err
	}
	l.Exec.Backend = core.BackendNativeCPU
	x := fill(1, 1, 8)
	_, post, err := cnn1.Forward(l, x)
	if err != nil {
		return err
	}
	_, err = cnn1.TrainNativeMSE(l, x, tgtLike(post), 0.02)
	if err != nil {
		return err
	}
	return l.SyncWeightsFromNative()
}

func demoCNN2() error {
	l, err := cnn2.New(cnn2.Config{InChannels: 1, Filters: 2, Height: 4, Width: 4, Kernel: 3})
	if err != nil {
		return err
	}
	l.Exec.Backend = core.BackendNativeCPU
	x := fill(1, 1, 4, 4)
	_, post, err := cnn2.Forward(l, x)
	if err != nil {
		return err
	}
	_, err = cnn2.TrainNativeMSE(l, x, tgtLike(post), 0.02)
	if err != nil {
		return err
	}
	return l.SyncWeightsFromNative()
}

func demoCNN3() error {
	l, err := cnn3.New(cnn3.Config{InChannels: 1, Filters: 2, Depth: 2, Height: 2, Width: 2, Kernel: 2})
	if err != nil {
		return err
	}
	l.Exec.Backend = core.BackendNativeCPU
	x := fill(1, 1, 2, 2, 2)
	_, post, err := cnn3.Forward(l, x)
	if err != nil {
		return err
	}
	_, err = cnn3.TrainNativeMSE(l, x, tgtLike(post), 0.02)
	if err != nil {
		return err
	}
	return l.SyncWeightsFromNative()
}

func demoConvT1() error {
	l, err := convt1.New(convt1.Config{InChannels: 1, Filters: 2, SeqLen: 4, Kernel: 3})
	if err != nil {
		return err
	}
	l.Exec.Backend = core.BackendNativeCPU
	x := fill(1, 1, 4)
	_, post, err := convt1.Forward(l, x)
	if err != nil {
		return err
	}
	_, err = convt1.TrainNativeMSE(l, x, tgtLike(post), 0.02)
	if err != nil {
		return err
	}
	return l.SyncWeightsFromNative()
}

func demoConvT2() error {
	l, err := convt2.New(convt2.Config{InChannels: 1, Filters: 2, Height: 2, Width: 2, Kernel: 2})
	if err != nil {
		return err
	}
	l.Exec.Backend = core.BackendNativeCPU
	x := fill(1, 1, 2, 2)
	_, post, err := convt2.Forward(l, x)
	if err != nil {
		return err
	}
	_, err = convt2.TrainNativeMSE(l, x, tgtLike(post), 0.02)
	if err != nil {
		return err
	}
	return l.SyncWeightsFromNative()
}

func demoConvT3() error {
	l, err := convt3.New(convt3.Config{InChannels: 1, Filters: 2, Depth: 2, Height: 2, Width: 2, Kernel: 2})
	if err != nil {
		return err
	}
	l.Exec.Backend = core.BackendNativeCPU
	x := fill(1, 1, 2, 2, 2)
	_, post, err := convt3.Forward(l, x)
	if err != nil {
		return err
	}
	_, err = convt3.TrainNativeMSE(l, x, tgtLike(post), 0.02)
	if err != nil {
		return err
	}
	return l.SyncWeightsFromNative()
}

func demoRNN() error {
	l, err := rnn.New(rnn.Config{InputSize: 4, HiddenSize: 4, SeqLen: 2})
	if err != nil {
		return err
	}
	l.Exec.Backend = core.BackendNativeCPU
	x := fill(1, 2, 4)
	_, err = rnn.TrainNativeMSE(l, x, tgtLike(x), 0.05)
	if err != nil {
		return err
	}
	return l.SyncWeightsFromNative()
}

func demoLSTM() error {
	l, err := lstm.New(lstm.Config{InputSize: 4, HiddenSize: 4, SeqLen: 2})
	if err != nil {
		return err
	}
	l.Exec.Backend = core.BackendNativeCPU
	x := fill(1, 2, 4)
	_, err = lstm.TrainNativeMSE(l, x, tgtLike(x), 0.05)
	if err != nil {
		return err
	}
	return l.SyncWeightsFromNative()
}

func demoSeq() error {
	l, err := sequential.New(sequential.Config{Dim: 4, Depth: 2})
	if err != nil {
		return err
	}
	l.Exec.Backend = core.BackendNativeCPU
	x := fill(2, 4)
	_, err = sequential.TrainNativeMSE(l, x, tgtLike(x), 0.05)
	if err != nil {
		return err
	}
	return l.SyncWeightsFromNative()
}

func demoRes() error {
	l, err := residual.New(residual.Config{Dim: 4, Depth: 1})
	if err != nil {
		return err
	}
	l.Exec.Backend = core.BackendNativeCPU
	x := fill(2, 4)
	_, err = residual.TrainNativeMSE(l, x, tgtLike(x), 0.05)
	if err != nil {
		return err
	}
	return l.SyncWeightsFromNative()
}

func demoPara() error {
	init := make([]float32, 2*4*4)
	for i := range init {
		init[i] = 0.03 * float32(i+1)
	}
	l, err := parallel.NewConfigured(parallel.Config{Dim: 4, OutFeat: 4, Branches: 2, Combine: parallel.CombineAdd},
		core.DTypeFloat32, quant.FormatNone, init, nil)
	if err != nil {
		return err
	}
	l.Exec.Backend = core.BackendNativeCPU
	l.SyncBranchExec()
	x := fill(2, 4)
	_, err = parallel.TrainNativeMSE(l, x, tgtLike(x), 0.05)
	if err != nil {
		return err
	}
	return l.SyncWeightsFromNative()
}

func demoKMeans() error {
	l, err := kmeans.New(kmeans.Config{NumClusters: 3, FeatureDim: 4})
	if err != nil {
		return err
	}
	l.Exec.Backend = core.BackendNativeCPU
	x := fill(2, 4)
	_, post, err := kmeans.Forward(l, x)
	if err != nil {
		return err
	}
	_, err = kmeans.TrainNativeMSE(l, x, tgtLike(post), 0.05)
	if err != nil {
		return err
	}
	return l.SyncWeightsFromNative()
}

func demoMamba() error {
	l, err := mamba.New(mamba.Config{DModel: 4, DState: 2, SeqLen: 2})
	if err != nil {
		return err
	}
	l.Exec.Backend = core.BackendNativeCPU
	x := fill(1, 2, 4)
	_, err = mamba.TrainNativeMSE(l, x, tgtLike(x), 0.05)
	if err != nil {
		return err
	}
	return l.SyncWeightsFromNative()
}

func demoGDN() error {
	l, err := gdn.New(gdn.Config{HiddenSize: 4, NumKeyHeads: 1, NumValueHeads: 1, KeyHeadDim: 2, ValueHeadDim: 2, ConvKernel: 2})
	if err != nil {
		return err
	}
	l.Exec.Backend = core.BackendNativeCPU
	x := fill(1, 2, 4)
	_, err = gdn.TrainNativeMSE(l, x, tgtLike(x), 0.05)
	if err != nil {
		return err
	}
	return l.SyncWeightsFromNative()
}

func demoMeta() error {
	l, err := metacognition.New(metacognition.Config{Dim: 4})
	if err != nil {
		return err
	}
	l.Exec.Backend = core.BackendNativeCPU
	x := fill(2, 4)
	_, err = metacognition.TrainNativeMSE(l, x, tgtLike(x), 0.05)
	if err != nil {
		return err
	}
	return l.SyncWeightsFromNative()
}
