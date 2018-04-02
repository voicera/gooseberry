package polling

import (
	"math/rand"
	"time"

	"github.com/voicera/gooseberry/validate"
)

func init() {
	rand.Seed(time.Now().UnixNano()) // for Bernoulli sampling
}

// relaxationCondition represents a condition that can be composed in runtime;
// for example, it can be config-drived or created via dependency injection.
type relaxationCondition interface {
	shouldRelax(lastReceiptSucceeded bool) bool
}

// andRelaxationCondition is true if all conditions are true; order matters.
type and struct {
	conditions []relaxationCondition
}

func (and *and) shouldRelax(lastReceiptSucceeded bool) bool {
	for _, c := range and.conditions {
		if !c.shouldRelax(lastReceiptSucceeded) {
			return false
		}
	}
	return true
}

// A relaxationCondition that evaluates to false if the last call was
// empty-handed (no payload). The rationale is: if the last call was
// empty-handed, the next one is most probably alike.
type emptyHanded struct{}

func (*emptyHanded) shouldRelax(lastReceiptSucceeded bool) bool {
	return !lastReceiptSucceeded
}

// A relaxationCondition that runs a Bernoulli trial:
//     rand.Float64() < value
// to determine whether or not to relax.
type bernoulliSampler struct {
	float64
}

// newBernoulliSampler creates a probabilistic relaxer that relaxes
// if the following Bernoulli trial is true:
//     rand.Float64() < probability
func newBernoulliSampler(probability float64) (relaxationCondition, error) {
	if err := validate.InRange(probability, 0.0, 1.0, "probability"); err != nil {
		return nil, err
	}
	return &bernoulliSampler{probability}, nil
}

func (sampler *bernoulliSampler) shouldRelax(lastReceiptSucceeded bool) bool {
	return rand.Float64() < sampler.float64
}
