package rule

type Config struct {
	EnabledRules []RuleID
}

func DefaultConfig() Config {
	return Config{
		EnabledRules: []RuleID{Rule1, Rule2, Rule3, Rule4, Rule5, Rule6, Rule7, Rule8},
	}
}

func WesternElectricConfig() Config {
	return Config{
		EnabledRules: []RuleID{Rule1, Rule2, Rule5, Rule6},
	}
}

type Engine struct {
	cfg      Config
	checkers map[RuleID]func(ChartZones) []Violation
}

func NewEngine(cfg Config) *Engine {
	e := &Engine{
		cfg: cfg,
		checkers: map[RuleID]func(ChartZones) []Violation{
			Rule1: CheckRule1,
			Rule2: CheckRule2,
			Rule3: CheckRule3,
			Rule4: CheckRule4,
			Rule5: CheckRule5,
			Rule6: CheckRule6,
			Rule7: CheckRule7,
			Rule8: CheckRule8,
		},
	}
	return e
}

type Result struct {
	Violations    []Violation
	RuleCounts    map[RuleID]int
	TotalViolated int
	PointFlags    []PointFlag
}

type PointFlag struct {
	Index int
	Rules []RuleID
}

func (e *Engine) Evaluate(cz ChartZones) Result {
	var allViolations []Violation
	ruleCounts := make(map[RuleID]int)

	for _, rid := range e.cfg.EnabledRules {
		checker, ok := e.checkers[rid]
		if !ok {
			continue
		}
		vs := checker(cz)
		allViolations = append(allViolations, vs...)
		ruleCounts[rid] = len(vs)
	}

	pointMap := make(map[int]map[RuleID]bool)
	for _, v := range allViolations {
		if _, ok := pointMap[v.Index]; !ok {
			pointMap[v.Index] = make(map[RuleID]bool)
		}
		pointMap[v.Index][v.Rule] = true
	}

	pointFlags := make([]PointFlag, 0, len(pointMap))
	for idx, rulesMap := range pointMap {
		var rules []RuleID
		for r := range rulesMap {
			rules = append(rules, r)
		}
		pointFlags = append(pointFlags, PointFlag{Index: idx, Rules: rules})
	}

	sortPointFlags(pointFlags)

	return Result{
		Violations:    allViolations,
		RuleCounts:    ruleCounts,
		TotalViolated: len(allViolations),
		PointFlags:    pointFlags,
	}
}

func (r *Result) ViolatedIndices() []int {
	seen := make(map[int]bool)
	for _, v := range r.Violations {
		seen[v.Index] = true
	}
	indices := make([]int, 0, len(seen))
	for idx := range seen {
		indices = append(indices, idx)
	}
	sortInts(indices)
	return indices
}

func (r *Result) HasViolation(idx int, rule RuleID) bool {
	for _, v := range r.Violations {
		if v.Index == idx && v.Rule == rule {
			return true
		}
	}
	return false
}

func sortPointFlags(pf []PointFlag) {
	for i := 1; i < len(pf); i++ {
		for j := i; j > 0 && pf[j].Index < pf[j-1].Index; j-- {
			pf[j], pf[j-1] = pf[j-1], pf[j]
		}
	}
}

func sortInts(a []int) {
	for i := 1; i < len(a); i++ {
		for j := i; j > 0 && a[j] < a[j-1]; j-- {
			a[j], a[j-1] = a[j-1], a[j]
		}
	}
}
