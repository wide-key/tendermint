package main

import (
	"encoding/json"
	"fmt"
	//"reflect"
	"os"

	"github.com/coinexchain/codon"
	"github.com/coinexchain/randsrc"
	"github.com/tendermint/tendermint/codec"
)

var Count = 100 * 10000

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("usage: %s filename\n", os.Args[0])
		return
	}
	r := randsrc.NewRandSrcFromFile(os.Args[1])
	runRandTest(r)
}

func runRandTest(r codec.RandSrc) {
	leafTypes := codec.GetLeafTypes()

	buf := make([]byte, 0, 4096)
	for i := 0; i < Count; i++ {
		if i%10000 == 0 {
			fmt.Printf("=== %d ===\n", i)
		}

		ifc := codec.RandAny(r)
		//codon.ShowInfoForVar(leafTypes, ifc)

		buf = buf[:0]
		codec.EncodeAny(&buf, ifc)
		ifcDec, _, err := codec.DecodeAny(buf)
		if err != nil {
			fmt.Printf("Now: %d\n", i)
			codon.ShowInfoForVar(leafTypes, ifc)
			panic(err)
		}
		cpDec := codec.DeepCopyAny(ifcDec)
		origS, _ := json.Marshal(ifc)
		decS, _ := json.Marshal(cpDec)
		if pos := findMismatch(origS, decS); pos!=-1 {
			fmt.Printf("Now: %d\n%s\n%s\n", i, string(origS), string(decS))
			fmt.Printf("Diff: %d\n%s\n%s\n", pos, string(origS[:pos]), string(decS[:pos]))
			codon.ShowInfoForVar(leafTypes, ifc)
			panic("Mismatch!")
		}
	}
}

func findMismatch(a, b []byte) int {
	length := len(a)
	if len(b) < len(a) {
		length = len(b)
	}
	for i := 0; i < length; i++ {
		if a[i] != b[i] {
			return i
		}
	}
	if len(b) != len(a) {
		return length
	}
	return -1
}

