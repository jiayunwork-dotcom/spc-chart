package rule

type RuleID int

const (
	Rule1 RuleID = iota + 1
	Rule2
	Rule3
	Rule4
	Rule5
	Rule6
	Rule7
	Rule8
)

func (r RuleID) String() string {
	switch r {
	case Rule1:
		return "Rule 1 (beyond 3-sigma)"
	case Rule2:
		return "Rule 2 (9 same side)"
	case Rule3:
		return "Rule 3 (6 trending)"
	case Rule4:
		return "Rule 4 (14 alternating)"
	case Rule5:
		return "Rule 5 (2/3 beyond 2-sigma)"
	case Rule6:
		return "Rule 6 (4/5 beyond 1-sigma)"
	case Rule7:
		return "Rule 7 (15 within 1-sigma)"
	case Rule8:
		return "Rule 8 (8 beyond 1-sigma both sides)"
	default:
		return "Unknown"
	}
}

type Violation struct {
	Rule     RuleID
	Index    int
	StartIdx int
	EndIdx   int
	Message  string
}

type Zones struct {
	Value   float64
	ZScore  float64
	AboveCL bool
}

type ChartZones struct {
	CL    float64
	Sigma float64
	Zones []Zones
}

func NewChartZones(values []float64, cl, sigma float64) ChartZones {
	zones := make([]Zones, len(values))
	for i, v := range values {
		z := 0.0
		if sigma > 0 {
			z = (v - cl) / sigma
		}
		zones[i] = Zones{
			Value:   v,
			ZScore:  z,
			AboveCL: v > cl,
		}
	}
	return ChartZones{CL: cl, Sigma: sigma, Zones: zones}
}
