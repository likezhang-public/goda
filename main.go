package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/google/subcommands"

	"github.com/likezhang-public/goda/internal/cut"
	"github.com/likezhang-public/goda/internal/exec"
	"github.com/likezhang-public/goda/internal/graph"
	"github.com/likezhang-public/goda/internal/list"
	"github.com/likezhang-public/goda/internal/pkgset"
	"github.com/likezhang-public/goda/internal/tree"
	"github.com/likezhang-public/goda/internal/weight"
)

func main() {
	cmds := subcommands.NewCommander(flag.CommandLine, path.Base(os.Args[0]))
	cmds.Register(cmds.HelpCommand(), "")

	cmds.Register(&list.Command{}, "")
	cmds.Register(&tree.Command{}, "")
	cmds.Register(&exec.Command{}, "")
	cmds.Register(&weight.Command{}, "")
	cmds.Register(&graph.Command{}, "")
	cmds.Register(&cut.Command{}, "")
	cmds.Register(&ExprHelp{}, "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(cmds.Execute(ctx)))
}

type ExprHelp struct{}

func (*ExprHelp) Name() string     { return "expr" }
func (*ExprHelp) Synopsis() string { return "Help about package expressions" }
func (*ExprHelp) Usage() string {
	return `Package expressions allow to specify calculations with dependencies.

The examples use X, Y and Z as placeholders for packages, packages paths or
package expressions.

# Basic operations
	There are a few basic oprations specified for manipulating sets of packages.

	X Y Z;
	X + Y + Z;  add(X, Y, Z);  or(X, Y, Z)
		all packages that match X, Y and Z

	X - Y - Z;  subtract(X, Y, Z);  exclude(X, Y, Z)
		packages that are used by X and not used by Y and Z

	shared(X, Y, Z);  intersect(X, Y, Z)
		packages that exist in all of X, Y and Z

	xor(X, Y);
		packages that are different between X and Y

# Selectors
	Selectors allow selecting parts of the dependency tree

	X:all
		select X and all of its direct and indirect dependencies
	X:import, X:imp
		select direct import of X
	X:import:all, X:imp:all
		select direct and indirect dependencies of X; X not included

	X:source
		packages that have no other package importing them
	X:-source
		shorthand for (X - X:source)

	X:main
		select packages named main

	X:test
		select test packages of X

# Functions:

	reach(X, Y);
		packages from X that can reach a package in Y

	incoming(X, Y);
		packages from X that directly import a package in Y, including Y

	transitive(X);
		a transitive reduction in package dependencies

# Tags and OS:

	test=1(X);
		include tests when resolving X

	goos=linux(X):
		set goos to "linux" tests when resolving X

	purego=1(X):
		add tag "purego" for resolving X

# Example expressions

	github.com/loov/goda:import
		all direct dependencies for "github.com/loov/goda" package

	shared(github.com/loov/goda/pkgset:all, github.com/loov/goda/templates:all)
		packages shared by "github.com/loov/goda/pkgset" and "github.com/loov/goda/templates"

	github.com/loov/goda/... - golang.org/x/tools/...
		all dependencies excluding golang.org/x/tools

	reach(github.com/loov/goda/...:all, golang.org/x/tools/go/packages)
		packages in github.com/loov/goda/ that use golang.org/x/tools/go/packages
`
}
func (*ExprHelp) SetFlags(f *flag.FlagSet) {}

func (cmd *ExprHelp) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if f.NArg() == 0 {
		return subcommands.ExitUsageError
	}

	result, err := pkgset.Parse(ctx, f.Args())
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return subcommands.ExitFailure
	}
	fmt.Fprintln(os.Stdout, result.Tree(0))

	return subcommands.ExitSuccess
}
