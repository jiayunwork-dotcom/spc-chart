package spc

import (
	"fmt"
	"io"
	"strings"
)

func FormatResult(w io.Writer, res Result) error {
	var b strings.Builder
	fmt.Fprintf(&b, "Individuals control chart (%d points, %.0f-sigma limits)\n",
		len(res.Points), res.Sigma)
	fmt.Fprintf(&b, "  center line (mean): %.4f\n", res.Mean)
	fmt.Fprintf(&b, "  std dev (pop)     : %.4f\n", res.StdDev)
	fmt.Fprintf(&b, "  UCL               : %.4f\n", res.UCL)
	fmt.Fprintf(&b, "  LCL               : %.4f\n", res.LCL)
	fmt.Fprintf(&b, "  out-of-control    : %d\n", res.Outlier)

	if res.Outlier > 0 {
		fmt.Fprintf(&b, "\nFlagged points (index:value):\n")
		for _, p := range res.Points {
			if p.Outlier {
				fmt.Fprintf(&b, "  [%d] %.4f\n", p.Index, p.Value)
			}
		}
	}

	_, err := io.WriteString(w, b.String())
	return err
}

func FormatCSV(w io.Writer, res Result) error {
	var b strings.Builder
	b.WriteString("index,value,outlier\n")
	for _, p := range res.Points {
		flag := "false"
		if p.Outlier {
			flag = "true"
		}
		fmt.Fprintf(&b, "%d,%.6f,%s\n", p.Index, p.Value, flag)
	}
	_, err := io.WriteString(w, b.String())
	return err
}

func Summary(res Result) string {
	return fmt.Sprintf("n=%d mean=%.4f sigma=%.4f UCL=%.4f LCL=%.4f OOC=%d",
		len(res.Points), res.Mean, res.StdDev, res.UCL, res.LCL, res.Outlier)
}
