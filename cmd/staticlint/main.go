package main

import (
	"github.com/kisielk/errcheck/errcheck"

	"github.com/charithe/durationcheck"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/pkgfact"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/timeformat"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"
	"honnef.co/go/tools/staticcheck"
)

func main() {
	checks := []*analysis.Analyzer{
		durationcheck.Analyzer,
		errcheck.Analyzer,

		ifaceassert.Analyzer,
		bools.Analyzer,
		sigchanyzer.Analyzer,
		atomicalign.Analyzer,
		inspect.Analyzer,
		stdmethods.Analyzer,
		framepointer.Analyzer,
		sortslice.Analyzer,
		findcall.Analyzer,
		testinggoroutine.Analyzer,
		shadow.Analyzer,
		unsafeptr.Analyzer,
		lostcancel.Analyzer,
		shift.Analyzer,
		tests.Analyzer,
		unusedwrite.Analyzer,
		printf.Analyzer,
		structtag.Analyzer,
		atomic.Analyzer,
		composite.Analyzer,
		buildtag.Analyzer,
		ctrlflow.Analyzer,
		unmarshal.Analyzer,
		timeformat.Analyzer,
		assign.Analyzer,
		asmdecl.Analyzer,
		usesgenerics.Analyzer,
		deepequalerrors.Analyzer,
		httpresponse.Analyzer,
		copylock.Analyzer,
		stringintconv.Analyzer,
		pkgfact.Analyzer,
		unreachable.Analyzer,
		errorsas.Analyzer,
		nilness.Analyzer,
		unusedresult.Analyzer,
		fieldalignment.Analyzer,
		nilfunc.Analyzer,
		reflectvaluecompare.Analyzer,
		buildssa.Analyzer,
		loopclosure.Analyzer,
		cgocall.Analyzer,
	}
	for _, v := range staticcheck.Analyzers {
		checks = append(checks, v.Analyzer)
	}

	multichecker.Main(
		checks...,
	)

}
