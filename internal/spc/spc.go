package spc

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

type Point struct {
	Index   int
	Value   float64
	Outlier bool
}

type Result struct {
	Mean    float64
	StdDev  float64
	UCL     float64
	LCL     float64
	Sigma   float64
	Points  []Point
	Outlier int
}

func Analyze(values []float64, sigma float64) (Result, error) {
	if len(values) == 0 {
		return Result{}, fmt.Errorf("no values to analyze")
	}
	if sigma <= 0 {
		return Result{}, fmt.Errorf("sigma must be positive, got %v", sigma)
	}
	var sum float64
	for _, v := range values {
		sum += v
	}
	mean := sum / float64(len(values))

	var sq float64
	for _, v := range values {
		d := v - mean
		sq += d * d
	}
	stddev := 0.0
	if len(values) > 1 {
		stddev = math.Sqrt(sq / float64(len(values)))
	}
	ucl := mean + sigma*stddev
	lcl := mean - sigma*stddev

	pts := make([]Point, len(values))
	out := 0
	for i, v := range values {
		isOut := v > ucl || v < lcl
		if isOut {
			out++
		}
		pts[i] = Point{Index: i, Value: v, Outlier: isOut}
	}
	return Result{
		Mean:    mean,
		StdDev:  stddev,
		UCL:     ucl,
		LCL:     lcl,
		Sigma:   sigma,
		Points:  pts,
		Outlier: out,
	}, nil
}

func ParseValues(r io.Reader) ([]float64, error) {
	var vals []float64
	sc := bufio.NewScanner(r)
	line := 0
	for sc.Scan() {
		line++
		field := strings.TrimSpace(sc.Text())
		if field == "" {
			continue
		}
		for _, tok := range strings.FieldsFunc(field, func(r rune) bool {
			return r == ',' || r == ' ' || r == '\t'
		}) {
			v, err := strconv.ParseFloat(tok, 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: cannot parse %q as number", line, tok)
			}
			vals = append(vals, v)
		}
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	if len(vals) == 0 {
		return nil, fmt.Errorf("no numeric values found")
	}
	return vals, nil
}
