// Contains log related stuff.
package app

import (
	contour "github.com/mohae/contourp"
	jww "github.com/spf13/jwalterweatherman"
)

func init() {
	//Disable logger by default
	DisableLog()
}

// DisableLog disables all package output
func DisableLog() {
	jww.DiscardLogging()
}

// SetLog sets up logging, if it is enabled to stdout. At this point, the
// only overrides to logging will occur with CLI args. If the CLI args have any
// logging related flags, those will be processed and logging will be updated.
//
func SetLogging() error {
	if !contour.GetBool(CfgLog) {
		DisableLog()
		return nil
	}

	logFile := contour.GetString(CfgLog)
	if logFile == "" {
		logFile = DefaultLogFile
	}
	jww.SetLogFile(logFile)
	return nil
}
