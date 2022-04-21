package web

import (
	"encoding/binary"
	"math"

	"github.com/nsqio/go-nsq"
)

const (
	N = 300
)

type RowData struct {
	Datas    []float64 `json:"data"`
	Min      float64   `json:"min"`
	Max      float64   `json:"max"`
	Variance float64   `json:"variance"`
	Mean     float64   `json:"mean"`
}

var Data RowData

type MyMessageHandler struct{}
type MeanHandler struct{}
type VarianceHandler struct{}
type MaxAndminHandler struct{}

//ByteToFloat64 byte转Float64
func ByteToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	return math.Float64frombits(bits)
}

func (h *MyMessageHandler) processMessage(m []byte) error {
	num := ByteToFloat64(m)
	// fmt.Println(num)
	if len(Data.Datas) > N {
		Data.Datas = Data.Datas[1:]
	}
	Data.Datas = append(Data.Datas, num)

	return nil
}

func (h *VarianceHandler) processMessage(m []byte) error {
	Data.Variance = ByteToFloat64(m)
	return nil
}

func (h *MeanHandler) processMessage(m []byte) error {
	Data.Mean = ByteToFloat64(m)

	return nil
}

func (h *MaxAndminHandler) processMessage(m []byte) error {

	byte := m[len(m)-1]

	m = m[:len(m)-1]

	if byte == 0 {
		Data.Max = ByteToFloat64(m)
	} else {
		Data.Min = ByteToFloat64(m)
	}
	return nil
}

// HandleMessage implements the Handler interface.
func (h *MyMessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}

	err := h.processMessage(m.Body)
	return err
}

// HandleMessage implements the Handler interface.
func (h *MeanHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}

	err := h.processMessage(m.Body)
	return err
}

// HandleMessage implements the Handler interface.
func (h *VarianceHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}

	err := h.processMessage(m.Body)
	return err
}

// HandleMessage implements the Handler interface.
func (h *MaxAndminHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}

	err := h.processMessage(m.Body)
	return err
}
