package mse6

import (
	"bytes"
	"compress/gzip"
	"github.com/rs/zerolog/log"
)

func gzipenc(input []byte) []byte {
wrt, _ := zipPool.Get().(*gzip.Writer)
buf := &bytes.Buffer{}
wrt.Reset(buf)

_, _ = wrt.Write(input)
_ = wrt.Close()
defer zipPool.Put(wrt)

enc := buf.Bytes()
log.Trace().Msgf("zipped byte buffer size %d", len(enc))
return enc
}