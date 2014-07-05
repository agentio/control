#!/bin/sh
#
# builds necessary Go components
# must be run as control
#
go get control/...
go install control/...

