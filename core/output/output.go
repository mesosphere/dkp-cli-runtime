package output

import "io"

type Output interface {
	// Info displays informative output.
	//
	// Example:
	//  output.Info("namespace created")
	Info(msg string)

	// Infof displays informative output.
	//
	// Example:
	//  output.Infof("namespace %q created", namespace)
	Infof(format string, args ...interface{})

	// InfoWriter returns a writer for informative output.
	//
	// Example:
	//  io.WriteString(output.InfoWriter(), "namespace created")
	InfoWriter() io.Writer

	// Error communicates an error to the users.
	//
	// Example:
	//  output.Error(err, "namespace could not be created")
	Error(err error, msg string)

	// Errorf communicates an error to the users.
	//
	// Example:
	//  output.Errorf(err, "namespace %q could not be created", namespace)
	Errorf(err error, format string, args ...interface{})

	// ErrorWriter returns a writer for errors.
	//
	// Example:
	//  io.WriteString(output.ErrorWriter(), "namespace could not be created")
	ErrorWriter() io.Writer

	// StartOperation communicates the beginning of a long-running operation.
	// If running in a terminal, a progress animation will be shown. Starting a
	// new operation ends any previously running operation.
	//
	// Example:
	//  output.StartOperation("installing package")
	//  err := installPackage()
	//  if err != nil {
	//  	output.EndOperation(false)
	//  	output.Error(err, "")
	//  	return
	//  }
	//  output.EndOperation(true)
	StartOperation(status string)

	// EndOperation communicates the end of a long-running operation, either because
	// the operation completed successfully or failed (parameter success).
	//
	// Example:
	//  output.StartOperation("installing package")
	//  err := installPackage()
	//  if err != nil {
	//  	output.EndOperation(false)
	//  	output.Error(err, "")
	//  	return
	//  }
	//  output.EndOperation(true)
	EndOperation(success bool)

	// Result outputs the result of an operation, e.g. a "get <something>" command.
	//
	// Example:
	//  output.Result(pods.String())
	Result(result string)

	// ResultWriter returns a writer for command results.
	//
	// Example:
	//  encoder := json.NewEncoder(output.ResultWriter())
	//  encoder.Encode(object)
	ResultWriter() io.Writer

	// V returns an Output with a higher verbosity level (default: 0).
	// Info and Error output with a higher verbosity is only displayed if the
	// "--verbose" flag is set to an equal or higher value.
	//
	// Example:
	//  output.V(1).Info("verbose information")
	V(level int) Output

	// WithValues returns an Output with additional context in the form of
	// structured data (key-value pairs). Not displayed in interactive shells.
	//
	// Example:
	//  output.WithValues("cluster", clusterName).Info("namespace created")
	WithValues(keysAndValues ...interface{}) Output
}
