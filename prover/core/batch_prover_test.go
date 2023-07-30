//go:build ffi

package core_test

import (
	"encoding/json"
	"flag"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/scroll-tech/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"

	"scroll-tech/prover/config"
	"scroll-tech/prover/core"
)

var (
	paramsPath    = flag.String("params", "/assets/test_params", "params dir")
	tracesPath    = flag.String("traces", "/assets/traces", "traces dir")
	proofDumpPath = flag.String("dump", "/assets/proof_data", "the path proofs dump to")
)

func TestFFI(t *testing.T) {
	as := assert.New(t)
	cfg := &config.ProverCoreConfig{
		ParamsPath: *paramsPath,
	}
	proverCore, err := core.NewProverCore(cfg)
	as.NoError(err)

	files, err := os.ReadDir(*tracesPath)
	as.NoError(err)

	traces := make([]*types.BlockTrace, 0)
	for _, file := range files {
		var (
			f   *os.File
			byt []byte
		)
		f, err = os.Open(filepath.Join(*tracesPath, file.Name()))
		as.NoError(err)
		byt, err = io.ReadAll(f)
		as.NoError(err)
		trace := &types.BlockTrace{}
		as.NoError(json.Unmarshal(byt, trace))
		traces = append(traces, trace)
	}
	proof, err := proverCore.Prove("test", traces)
	as.NoError(err)
	t.Log("prove success")

	// dump the proof
	os.RemoveAll(*proofDumpPath)
	proofByt, err := json.Marshal(proof)
	as.NoError(err)
	proofFile, err := os.Create(*proofDumpPath)
	as.NoError(err)
	_, err = proofFile.Write(proofByt)
	as.NoError(err)
}
