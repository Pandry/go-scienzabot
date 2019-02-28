package main

import (
	"fmt"
)

//
// logLevel 	0: trace
//				1: info
//				2: warn
//				3: error
//				4: fatal

func logMessage(logLevel uint8, messages ...interface{}) {
	if *verboseMode {
		prependingString := ""
		switch logLevel {
		case 0:
			prependingString = "[ \033[44mTRACE\033[49m ]"
			break
		case 1:
			prependingString = "[ \033[94mINFO\033[49m ]"
			break
		case 2:
			prependingString = "[ \033[93mWARN\033[49m ]"
			break
		case 3:
			prependingString = "[ \033[91mERROR\033[49m ]"
			break
		case 4:
			prependingString = "[ \033[41mFATAL\033[49m ]"
			break
		}
		fmt.Println(prependingString, messages)
	}
}
