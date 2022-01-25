<!--
 Copyright 2022 D2iQ, Inc. All rights reserved.
 SPDX-License-Identifier: Apache-2.0
 -->

# Unified CLI output

Since the CLIs for Konvoy, Kommander, Diagnostics will be bundled together in a single `dkp` CLI starting with v2.2,
they must output information to the user in a consistent way.

We achieve this by offering a common package with a simple interface for user output in this package. By using this
interface for all output, the CLIs only communicate *what* they want to output, this package handles *how* the output
looks like. This also makes eventual future changes to the output formatting easier, because it only needs to happen
in one place.

All terminal output must use this interface, direct use of StdOut or StdErr is discouraged.

## Usage examples

```go
// output is preconfigured by the common root command (e.g. verbosity set from --verbose flag)
rootCmd, rootOptions := cmd.Root()
output := rootOptions.Output()
```

### Displaying information and errors

```go
// optional output only displayed with higher verbosity
output.V(1).Info("kubeconfig read from file")
err := createNameSpace(namespaceName)
if err != nil {
    output.Errorf(err, "failed to create namespace %q", namespaceName)
    os.Exit(1)
}
output.Infof("namespace %q created" namespaceName)
```

### Long-running operations

```go
output.StartOperation("installing packages")
for _, package := range packages {
    err := installPackage(package)
    if err != nil {
        output.EndOperation(false)
        output.Errorf(err, "failed to install package %q", package.Name)
        os.Exit(1)
    }
    output.V(1).Infof("package %q installed", package.Name)
}
output.EndOperation(true)
output.Info("All packages installed successfully")
```

### Output results

```go
pods, err := getPods(namespaceName)
if err != nil {
    output.Error(err, "failed to get pods")
    os.Exit(1)
}
if outputJSON {
    output.Result(pods.ToJSON())
} else {
    output.Result(pods.String())
}
```

## What's the difference between Info() and Result()?

`Result()` is meant to output the result of an operation. This might be clear text, but can also be e.g. JSON encoded
depending on the use case. A command's result is the only thing that's sent to StdOut and will always be output "as is".
All other output is sent to StdErr.

`Info()` is meant to communicate information (e.g. progress, successful execution) to a user. This output (together with
error messages and animations) is sent to StdErr and might be formatted in different ways, e.g: For an interactive
terminal the output can be colored, progress messages can be animated. If not running in a terminal, these messages
might be prefixed with a timestamp or formatted in a more machine readable way.

## Why using StdOut only for results?

This makes sure the result can be used directly, e.g. in scripts, piped to other tools, redirected to a file, etc.
Informative output only meant for a human user is kept out of StdOut.

`kubectl` uses the same convention, see e.g.
[here (result)](https://github.com/kubernetes/kubectl/blob/3f7abd9859b92958fcf8f6e5bb9d7b354aee4781/pkg/cmd/get/get.go#L824-L826)
vs. [here (informative)](https://github.com/kubernetes/kubectl/blob/3f7abd9859b92958fcf8f6e5bb9d7b354aee4781/pkg/cmd/get/get.go#L596-L603).

See also [this StackExchange discussion](https://unix.stackexchange.com/questions/331611/do-progress-reports-logging-information-belong-on-stderr-or-stdout).
