package sdkapi

import "github.com/Masterminds/semver/v3"

var Version05GTEConstraint, _ = semver.NewConstraint(">= 0.5.0-next.1")
var Version05Constraint, _ = semver.NewConstraint(">= 0.5, < 0.6")
