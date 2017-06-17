package statsq

import (
	"bytes"
	"github.com/qnib/qframe-types"
	"io"
	"log"
	"strconv"
	"strings"
)

type MsgParser struct {
	reader           io.Reader
	buffer           []byte
	partialReads     bool
	done             bool
	debug            bool
	maxUdpPacketSize int
	prefix           string
	postfix          string
}

func NewParser(reader io.Reader, partialReads, debug bool, maxUdpPacketSize int, prefix, postfix string) *MsgParser {
	return &MsgParser{
		reader, []byte{},
		partialReads, false, debug,
		maxUdpPacketSize,
		prefix, postfix}
}

func (mp *MsgParser) Next() (*qtypes.StatsdPacket, bool) {
	buf := mp.buffer

	for {
		line, rest := mp.lineFrom(buf)

		if line != nil {
			mp.buffer = rest
			return mp.parseLine(line), true
		}

		if mp.done {
			return mp.parseLine(rest), false
		}

		idx := len(buf)
		end := idx
		if mp.partialReads {
			end += TCP_READ_SIZE
		} else {
			end += int(mp.maxUdpPacketSize)
		}
		if cap(buf) >= end {
			buf = buf[:end]
		} else {
			tmp := buf
			buf = make([]byte, end)
			copy(buf, tmp)
		}

		n, err := mp.reader.Read(buf[idx:])
		buf = buf[:idx+n]
		if err != nil {
			if err != io.EOF {
				log.Printf("ERROR: %s", err)
			}

			mp.done = true

			line, rest = mp.lineFrom(buf)
			if line != nil {
				mp.buffer = rest
				return mp.parseLine(line), len(rest) > 0
			}

			if len(rest) > 0 {
				return mp.parseLine(rest), false
			}

			return nil, false
		}
	}
}

func (mp *MsgParser) lineFrom(input []byte) ([]byte, []byte) {
	split := bytes.SplitAfterN(input, []byte("\n"), 2)
	if len(split) == 2 {
		return split[0][:len(split[0])-1], split[1]
	}

	if !mp.partialReads {
		if len(input) == 0 {
			input = nil
		}
		return input, []byte{}
	}

	if bytes.HasSuffix(input, []byte("\n")) {
		return input[:len(input)-1], []byte{}
	}

	return nil, input
}
/*
func (mp *MsgParser) parseLine(line []byte) *Packet {
	split := bytes.SplitN(line, []byte{'|'}, 3)
	if len(split) < 2 {
		mp.logParseFail(line)
		return nil
	}

	keyval := split[0]
	typeCode := string(split[1])

	sampling := float32(1)
	if strings.HasPrefix(typeCode, "c") || strings.HasPrefix(typeCode, "ms") {
		if len(split) == 3 && len(split[2]) > 0 {
			switch {
			case split[2][0] == '@':
				f64, err := strconv.ParseFloat(string(split[2][1:]), 32)
				if err != nil {
					log.Printf(
						"ERROR: failed to ParseFloat %s - %s",
						string(split[2][1:]),
						err,
					)
					return nil
				}
				sampling = float32(f64)
			}
		}
	}

	split = bytes.SplitN(keyval, []byte{':'}, 2)
	if len(split) < 2 {
		mp.logParseFail(line)
		return nil
	}
	name := string(split[0])
	val := split[1]
	if len(val) == 0 {
		mp.logParseFail(line)
		return nil
	}

	var (
		err      error
		floatval float64
		strval   string
	)

	switch typeCode {
	case "c":
		floatval, err = strconv.ParseFloat(string(val), 64)
		if err != nil {
			log.Printf("ERROR: failed to ParseFloat %s - %s", string(val), err)
			return nil
		}
	case "g":
		var s string

		if val[0] == '+' || val[0] == '-' {
			strval = string(val[0])
			s = string(val[1:])
		} else {
			s = string(val)
		}
		floatval, err = strconv.ParseFloat(s, 64)
		if err != nil {
			log.Printf("ERROR: failed to ParseFloat %s - %s", string(val), err)
			return nil
		}
	case "s":
		strval = string(val)
	case "ms":
		floatval, err = strconv.ParseFloat(string(val), 64)
		if err != nil {
			log.Printf("ERROR: failed to ParseFloat %s - %s", string(val), err)
			return nil
		}
	default:
		log.Printf("ERROR: unrecognized type code %q", typeCode)
		return nil
	}

	return &Packet{
		Bucket:   sanitizeBucket(mp.prefix + string(name) + mp.postfix),
		ValFlt:   floatval,
		ValStr:   strval,
		Modifier: typeCode,
		Sampling: sampling,
	}
}

*/
func (mp *MsgParser) parseLine(line []byte) *qtypes.StatsdPacket {
	splitDim := bytes.SplitN(line, []byte{' '}, 3)
	dims := qtypes.NewDimensions()
	switch len(splitDim) {
	case 2:
		dims = qtypes.NewDimensionsFromBytes(splitDim[1])
	}
	line = splitDim[0]
	split := bytes.SplitN(line, []byte{'|'}, 3)
	if len(split) < 2 {
		mp.logParseFail(line)
		return nil
	}

	keyval := split[0]
	typeCode := string(split[1])

	sampling := float32(1)
	if strings.HasPrefix(typeCode, "c") || strings.HasPrefix(typeCode, "ms") {
		if len(split) == 3 && len(split[2]) > 0 {
			switch {
			case split[2][0] == '@':
				f64, err := strconv.ParseFloat(string(split[2][1:]), 32)
				if err != nil {
					log.Printf("ERROR: failed to ParseFloat %s - %s", string(split[2][1:]), err)
					return nil
				}
				sampling = float32(f64)
			}
		}
	}
	split = bytes.SplitN(keyval, []byte{':'}, 2)
	if len(split) < 2 {
		mp.logParseFail(line)
		return nil
	}
	name := string(split[0])
	val := split[1]
	if len(val) == 0 {
		mp.logParseFail(line)
		return nil
	}

	var (
		err      error
		floatval float64
		strval   string
	)

	switch typeCode {
	case "c":
		floatval, err = strconv.ParseFloat(string(val), 64)
		if err != nil {
			log.Printf("ERROR: failed to ParseFloat %s - %s", string(val), err)
			return nil
		}
	case "g":
		var s string

		if val[0] == '+' || val[0] == '-' {
			strval = string(val[0])
			s = string(val[1:])
		} else {
			s = string(val)
		}
		floatval, err = strconv.ParseFloat(s, 64)
		if err != nil {
			log.Printf("ERROR: failed to ParseFloat %s - %s", string(val), err)
			return nil
		}
	case "s":
		strval = string(val)
	case "ms":
		floatval, err = strconv.ParseFloat(string(val), 64)
		if err != nil {
			log.Printf("ERROR: failed to ParseFloat %s - %s", string(val), err)
			return nil
		}
	default:
		log.Printf("ERROR: unrecognized type code %q", typeCode)
		return nil
	}

	return &qtypes.StatsdPacket{
		Bucket:     sanitizeBucket(mp.prefix + string(name) + mp.postfix),
		ValFlt:     floatval,
		ValStr:     strval,
		Modifier:   typeCode,
		Sampling:   sampling,
		Dimensions: dims,
	}
}

func (mp *MsgParser) logParseFail(line []byte) {
	if mp.debug {
		log.Printf("ERROR: failed to parse line: %q\n", string(line))
	}
}
