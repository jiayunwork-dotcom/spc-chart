package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"spc-chart/internal/capability"
	"spc-chart/internal/chart"
	"spc-chart/internal/rule"
	"spc-chart/internal/server"
	"spc-chart/internal/spc"
)

func main() {
	if len(os.Args) < 2 {
		runHTTPServer(nil)
		return
	}

	switch os.Args[1] {
	case "serve":
		runHTTPServer(os.Args[2:])
	case "ichart":
		runIChart(os.Args[2:])
	case "xbar-r":
		runXbarR(os.Args[2:])
	case "cusum":
		runCUSUM(os.Args[2:])
	case "ewma":
		runEWMA(os.Args[2:])
	case "capability":
		runCapability(os.Args[2:])
	case "rules":
		runRules(os.Args[2:])
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "spc-chart: unknown command %q\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func runHTTPServer(args []string) {
	addr := ":8080"
	for i := 0; i < len(args)-1; i++ {
		if args[i] == "-addr" || args[i] == "--addr" {
			addr = args[i+1]
			break
		}
	}
	cfg := server.Config{Addr: addr}
	fmt.Fprintf(os.Stdout, "spc-chart server listening on %s\n", server.FormatAddr(addr))
	if err := server.ListenAndServe(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `Usage: spc-chart <command> [options]

Commands:
  ichart       Individuals (I-MR) control chart
  xbar-r       Xbar-R control chart (subgroups from file)
  cusum        CUSUM chart for small shift detection
  ewma         EWMA chart for small shift detection
  capability   Process capability indices (Cp, Cpk, Pp, Ppk)
  rules        Nelson/Western Electric rule analysis

Use "spc-chart <command> -h" for command-specific help.
`)
}

func runIChart(args []string) {
	fs := flag.NewFlagSet("ichart", flag.ExitOnError)
	in := fs.String("in", "", "measurements file (required)")
	sigma := fs.Float64("sigma", 3, "control limit width in standard deviations")
	out := fs.String("out", "", "output file (default: stdout)")
	fs.Parse(args)

	if *in == "" {
		fatal("ichart: missing required -in flag")
	}
	vals := readValues(*in)
	res, err := chart.Individuals(vals, chart.IndividualsConfig{Sigma: *sigma, MRSpan: 2})
	if err != nil {
		fatal("ichart: %v", err)
	}
	writeChartResult(*out, res, len(vals))
}

func runXbarR(args []string) {
	fs := flag.NewFlagSet("xbar-r", flag.ExitOnError)
	in := fs.String("in", "", "measurements file (required)")
	n := fs.Int("n", 5, "subgroup size")
	sigma := fs.Float64("sigma", 3, "control limit width")
	out := fs.String("out", "", "output file (default: stdout)")
	fs.Parse(args)

	if *in == "" {
		fatal("xbar-r: missing required -in flag")
	}
	vals := readValues(*in)
	subs := toSubgroups(vals, *n)
	res, err := chart.XbarR(subs, chart.XbarRConfig{Sigma: *sigma})
	if err != nil {
		fatal("xbar-r: %v", err)
	}

	var b strings.Builder
	fmt.Fprintf(&b, "Xbar-R chart (%d subgroups, size %d, %.0f-sigma limits)\n", len(subs), *n, *sigma)
	fmt.Fprintf(&b, "  Xbar CL: %.4f  UCL: %.4f  LCL: %.4f  OOC: %d\n",
		res.Xbar.Limits.CL, res.Xbar.Limits.UCL, res.Xbar.Limits.LCL, res.Xbar.OOCCount)
	fmt.Fprintf(&b, "  R    CL: %.4f  UCL: %.4f  LCL: %.4f  OOC: %d\n",
		res.R.Limits.CL, res.R.Limits.UCL, res.R.Limits.LCL, res.R.OOCCount)
	writeOutput(*out, b.String())
}

func runCUSUM(args []string) {
	fs := flag.NewFlagSet("cusum", flag.ExitOnError)
	in := fs.String("in", "", "measurements file (required)")
	k := fs.Float64("k", 0.5, "allowance (reference value)")
	h := fs.Float64("h", 5.0, "decision interval")
	target := fs.Float64("target", 0, "process target (0 = use sample mean)")
	out := fs.String("out", "", "output file (default: stdout)")
	fs.Parse(args)

	if *in == "" {
		fatal("cusum: missing required -in flag")
	}
	vals := readValues(*in)
	cfg := chart.CUSUMConfig{Target: *target, K: *k, H: *h}
	res, err := chart.CUSUM(vals, cfg)
	if err != nil {
		fatal("cusum: %v", err)
	}

	var b strings.Builder
	fmt.Fprintf(&b, "CUSUM chart (%d points, K=%.2f, H=%.2f, target=%.4f, sigma=%.4f)\n",
		len(vals), res.K, res.H, res.Target, res.Sigma)
	fmt.Fprintf(&b, "  Alarms: %d\n", res.OOCCount)
	if res.OOCCount > 0 {
		fmt.Fprintf(&b, "\nAlarm points:\n")
		for _, p := range res.Points {
			if p.OutOfCtrl {
				fmt.Fprintf(&b, "  [%d] value=%.4f C+=%.4f C-=%.4f side=%v\n",
					p.Index, p.Value, p.CPlus, p.CMinus, p.Side)
			}
		}
	}
	writeOutput(*out, b.String())
}

func runEWMA(args []string) {
	fs := flag.NewFlagSet("ewma", flag.ExitOnError)
	in := fs.String("in", "", "measurements file (required)")
	lambda := fs.Float64("lambda", 0.2, "smoothing constant (0 < lambda <= 1)")
	l := fs.Float64("L", 3.0, "control limit factor")
	target := fs.Float64("target", 0, "process target (0 = use sample mean)")
	out := fs.String("out", "", "output file (default: stdout)")
	fs.Parse(args)

	if *in == "" {
		fatal("ewma: missing required -in flag")
	}
	vals := readValues(*in)
	cfg := chart.EWMAConfig{Lambda: *lambda, L: *l, Target: *target}
	res, err := chart.EWMA(vals, cfg)
	if err != nil {
		fatal("ewma: %v", err)
	}

	var b strings.Builder
	fmt.Fprintf(&b, "EWMA chart (%d points, lambda=%.3f, L=%.2f, target=%.4f, sigma=%.4f)\n",
		len(vals), res.Lambda, res.L, res.Target, res.Sigma)
	fmt.Fprintf(&b, "  Steady-state UCL: %.4f  LCL: %.4f\n", res.SteadyUCL, res.SteadyLCL)
	fmt.Fprintf(&b, "  Alarms: %d\n", res.OOCCount)
	writeOutput(*out, b.String())
}

func runCapability(args []string) {
	fs := flag.NewFlagSet("capability", flag.ExitOnError)
	in := fs.String("in", "", "measurements file (required)")
	usl := fs.Float64("usl", 0, "upper specification limit (required)")
	lsl := fs.Float64("lsl", 0, "lower specification limit (required)")
	target := fs.Float64("target", 0, "nominal target (default: midpoint)")
	out := fs.String("out", "", "output file (default: stdout)")
	fs.Parse(args)

	if *in == "" {
		fatal("capability: missing required -in flag")
	}
	if *usl == 0 && *lsl == 0 {
		fatal("capability: must specify -usl and -lsl")
	}
	vals := readValues(*in)
	spec := capability.SpecLimits{USL: *usl, LSL: *lsl, Target: *target}
	res, err := capability.ComputeIndices(vals, spec, 0)
	if err != nil {
		fatal("capability: %v", err)
	}

	var b strings.Builder
	fmt.Fprintf(&b, "Process capability (%d observations, USL=%.4f, LSL=%.4f)\n",
		res.N, *usl, *lsl)
	fmt.Fprintf(&b, "  Mean: %.4f  StdDev(within): %.4f  StdDev(overall): %.4f\n",
		res.Mean, res.StdW, res.StdO)
	fmt.Fprintf(&b, "  Cp:  %.3f  Cpk: %.3f  (CpU: %.3f  CpL: %.3f)\n",
		res.Cp, res.Cpk, res.CpU, res.CpL)
	fmt.Fprintf(&b, "  Pp:  %.3f  Ppk: %.3f  (PpU: %.3f  PpL: %.3f)\n",
		res.Pp, res.Ppk, res.PpU, res.PpL)
	fmt.Fprintf(&b, "  Cpm: %.3f\n", res.Cpm)
	fmt.Fprintf(&b, "  Est. PPM: %.1f  Sigma level: %.2f\n",
		capability.PPMEstimate(res.Cpk), capability.SigmaLevelWithShift(res.Cpk))
	writeOutput(*out, b.String())
}

func runRules(args []string) {
	fs := flag.NewFlagSet("rules", flag.ExitOnError)
	in := fs.String("in", "", "measurements file (required)")
	sigma := fs.Float64("sigma", 0, "process sigma (0 = estimate from MRbar/d2)")
	cl := fs.Float64("cl", 0, "center line (0 = use mean)")
	mode := fs.String("mode", "nelson", "rule set: nelson (all 8) or we (Western Electric 4)")
	out := fs.String("out", "", "output file (default: stdout)")
	fs.Parse(args)

	if *in == "" {
		fatal("rules: missing required -in flag")
	}
	vals := readValues(*in)

	sigmaEst := *sigma
	if sigmaEst <= 0 {
		if len(vals) >= 2 {
			var mrSum float64
			for i := 1; i < len(vals); i++ {
				d := vals[i] - vals[i-1]
				if d < 0 {
					d = -d
				}
				mrSum += d
			}
			mrbar := mrSum / float64(len(vals)-1)
			sigmaEst = mrbar / 1.128
		}
		if sigmaEst <= 0 {
			sigmaEst = 1
		}
	}

	clVal := *cl
	if clVal == 0 {
		var s float64
		for _, v := range vals {
			s += v
		}
		clVal = s / float64(len(vals))
	}

	cz := rule.NewChartZones(vals, clVal, sigmaEst)

	var cfg rule.Config
	switch *mode {
	case "we":
		cfg = rule.WesternElectricConfig()
	default:
		cfg = rule.DefaultConfig()
	}

	eng := rule.NewEngine(cfg)
	result := eng.Evaluate(cz)

	var b strings.Builder
	fmt.Fprintf(&b, "Rule analysis (%d points, CL=%.4f, sigma=%.4f, mode=%s)\n",
		len(vals), clVal, sigmaEst, *mode)
	fmt.Fprintf(&b, "  Total violations: %d\n", result.TotalViolated)
	if result.TotalViolated > 0 {
		fmt.Fprintf(&b, "\nViolations:\n")
		for _, v := range result.Violations {
			fmt.Fprintf(&b, "  %s: %s\n", v.Rule, v.Message)
		}
	}
	writeOutput(*out, b.String())
}

func readValues(path string) []float64 {
	f, err := os.Open(path)
	if err != nil {
		fatal("open %q: %v", path, err)
	}
	defer f.Close()
	vals, err := spc.ParseValues(f)
	if err != nil {
		fatal("parse %q: %v", path, err)
	}
	return vals
}

func toSubgroups(vals []float64, n int) []chart.Subgroup {
	if n < 2 {
		n = 2
	}
	count := len(vals) / n
	subs := make([]chart.Subgroup, count)
	for i := 0; i < count; i++ {
		subs[i] = chart.Subgroup{Values: vals[i*n : (i+1)*n]}
	}
	return subs
}

func writeChartResult(outPath string, res chart.ChartResult, nPoints int) {
	var b strings.Builder
	fmt.Fprintf(&b, "%s (%d points, %.0f-sigma limits)\n",
		res.Type, nPoints, res.Limits.UCL)
	fmt.Fprintf(&b, "  CL:  %.4f\n", res.Limits.CL)
	fmt.Fprintf(&b, "  UCL: %.4f\n", res.Limits.UCL)
	fmt.Fprintf(&b, "  LCL: %.4f\n", res.Limits.LCL)
	fmt.Fprintf(&b, "  OOC: %d\n", res.OOCCount)
	if res.OOCCount > 0 {
		fmt.Fprintf(&b, "\nOut-of-control points:\n")
		for _, p := range res.Points {
			if p.OutOfCtrl {
				fmt.Fprintf(&b, "  [%d] %.4f\n", p.Index, p.Value)
			}
		}
	}
	writeOutput(outPath, b.String())
}

func writeOutput(path, content string) {
	if path == "" {
		fmt.Print(content)
		return
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		fatal("write %q: %v", path, err)
	}
}

func fatal(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "spc-chart: "+format+"\n", a...)
	os.Exit(1)
}
