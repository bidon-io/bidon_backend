package sdkapi

import "github.com/Masterminds/semver/v3"

var Version05GTEConstraint, _ = semver.NewConstraint(">= 0.5.0-next.1")
var Version05Constraint, _ = semver.NewConstraint(">= 0.5, < 0.6")
var VersionLessThan073Constraint, _ = semver.NewConstraint("< 0.7.3")
var Version07xConstraint, _ = semver.NewConstraint(">= 0.7, < 0.8")
var Version081Constraint, _ = semver.NewConstraint("= 0.8.1")
