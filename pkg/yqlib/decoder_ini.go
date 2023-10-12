package yqlib

import (
	"io"

	ini "gopkg.in/ini.v1"
)

type iniDecoder struct {
	cfg      *ini.File
	finished bool
}

func NewIniDecoder() Decoder {
	return &iniDecoder{
		finished: false,
	}
}

func (dec *iniDecoder) Init(reader io.Reader) error {
	var err error
	dec.cfg, err = ini.Load(reader)
	if err != nil {
		return err
	}

	return nil
}

func (dec *iniDecoder) Decode() (*CandidateNode, error) {
	if dec.finished {
		return nil, io.EOF
	}

	log.Debug("ok here we go")
	return nil, nil
}
