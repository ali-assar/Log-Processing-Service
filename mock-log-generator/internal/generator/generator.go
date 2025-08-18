package generator

import (
	crand "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/ali-assar/Log-Processing-Service/mock-log-generator/internal/types"
)


func NewRandomLog(r *rand.Rand, pickLevel func() string, pickService func() string) types.Log {
	return types.Log{
		Timestamp: time.Now().UnixMilli(),
		Level:     pickLevel(),
		Message:   PickFrom(r, types.Messages),
		Service:   pickService(),
		Component: PickFrom(r, types.Components),
		TraceID:   RandomID(),
		SpanID:    RandomID(),
		ParentID:  RandomID(),
	}
}

func PickFrom(r *rand.Rand, options []string) string {
	return options[r.Intn(len(options))]
}

func RandomID() string {
	b := make([]byte, 8)
	if _, err := crand.Read(b); err != nil {
		return fmt.Sprintf("%x", rand.Int63())
	}
	return hex.EncodeToString(b)
}

type WeightedPicker struct {
	items  []string
	weight []int
	total  int
	r      *rand.Rand
}

func NewWeightedPicker(weights map[string]int, r *rand.Rand) *WeightedPicker {
	p := &WeightedPicker{r: r}
	for _, lvl := range types.Levels {
		w := weights[strings.ToUpper(lvl)]
		if w <= 0 {
			continue
		}
		p.items = append(p.items, strings.ToUpper(lvl))
		p.weight = append(p.weight, w)
		p.total += w
	}
	if p.total == 0 {
		for _, lvl := range types.Levels {
			p.items = append(p.items, strings.ToUpper(lvl))
			p.weight = append(p.weight, 1)
			p.total++
		}
	}
	return p
}

func (p *WeightedPicker) Pick() string {
	if p.total <= 0 {
		return "INFO"
	}
	n := p.r.Intn(p.total)
	for i, w := range p.weight {
		if n < w {
			return p.items[i]
		}
		n -= w
	}
	return p.items[len(p.items)-1]
}

func ParseLevelWeights(raw string, defaults map[string]int) map[string]int {
	if strings.TrimSpace(raw) == "" {
		return defaults
	}
	out := map[string]int{}
	parts := strings.Split(raw, ",")
	for _, part := range parts {
		p := strings.SplitN(strings.TrimSpace(part), ":", 2)
		if len(p) != 2 {
			continue
		}
		lvl := strings.ToUpper(strings.TrimSpace(p[0]))
		val, err := strconv.Atoi(strings.TrimSpace(p[1]))
		if err != nil || val <= 0 {
			continue
		}
		out[lvl] = val
	}
	if len(out) == 0 {
		return defaults
	}
	return out
}

func Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}
