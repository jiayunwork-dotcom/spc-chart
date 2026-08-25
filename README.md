# spc-chart

Statistical Process Control (SPC) toolkit for manufacturing quality analysis.
Implements control charts, process capability indices, Western Electric / Nelson
pattern-detection rules, and Gage R&R measurement system analysis.

## Features

- **Control charts**: I-chart (individuals), MR (moving range), Xbar-R, Xbar-S,
  CUSUM (cumulative sum), EWMA (exponentially weighted moving average)
- **Process capability**: Cp, Cpk, Pp, Ppk, Cpm with confidence intervals,
  Anderson-Darling and Shapiro-Wilk normality testing, PPM estimation
- **Rule engine**: All 8 Nelson rules plus Western Electric subset detection
- **Measurement System Analysis**: Gage R&R via ANOVA and Xbar-R methods,
  variance decomposition, %Study variation, NDC

## Build & Test

```bash
export GOTOOLCHAIN=local CGO_ENABLED=0
go build ./...
go test ./...
```

## Usage

```bash
# Individuals control chart
spc-chart ichart -in measurements.txt -sigma 3

# Xbar-R chart (subgroup size 5)
spc-chart xbar-r -in measurements.txt -n 5

# CUSUM for small shift detection
spc-chart cusum -in measurements.txt -k 0.5 -h 5

# EWMA chart
spc-chart ewma -in measurements.txt -lambda 0.2 -L 3

# Process capability
spc-chart capability -in measurements.txt -usl 12 -lsl 8

# Nelson/WE rule analysis
spc-chart rules -in measurements.txt -mode nelson
```

Input format: one number per line, or comma/space/tab-separated values.

## Package structure

```
internal/
├── chart/       Control chart algorithms (I, MR, Xbar-R, Xbar-S, CUSUM, EWMA)
├── capability/  Cp/Cpk/Pp/Ppk/Cpm indices, normality tests, confidence intervals
├── rule/        Nelson/Western Electric 8-rule pattern detection engine
├── msa/         Gage R&R measurement system analysis (ANOVA + Xbar-R)
└── spc/         Core parsing and basic I-chart (original module)
```

## License

MIT
