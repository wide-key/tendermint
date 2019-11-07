package main

import (
	"os"

	"github.com/tendermint/tendermint/codec"
)

func main() {
	//codec.ShowInfo()
	genFile()
}

func genFile() {
	codec.GenerateCodecFile(os.Stdout)
}
