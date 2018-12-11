package ucloud

import (
	"time"
)

// DefaultMaxRetries is default max retry attempts number
const DefaultMaxRetries = 3

// DefaultInSecure is a default value to enable https
const DefaultInSecure = false

// DefaultWaitInterval is the inteval to wait for state changed after resource is created
const DefaultWaitInterval = 10 * time.Second

// DefaultWaitMaxAttempts is the max attempts number to wait for state changed after resource is created
const DefaultWaitMaxAttempts = 10

// DefaultWaitIgnoreError is if it will ignore error during wait for state changed after resource is created
const DefaultWaitIgnoreError = false

//listenerStatusCvt is used to covert int to string for status after read lb listener
var listenerStatusCvt = newIntConverter(map[int]string{
	0: "allNormal",
	1: "partNormal",
	2: "allException",
})

//lbAttachmentStatusCvt is used to covert int to string for status after read lb attachment
var lbAttachmentStatusCvt = newIntConverter(map[int]string{
	0: "normalRunning",
	1: "exceptionRunning",
})

//lowerCaseProdCvt is used to covert one lower string to another lower string
var lowerCaseProdCvt = newStringConverter(map[string]string{
	"instance": "uhost",
	"lb":       "ulb",
})

//titleCaseProdCvt is used to covert one lower string to another string begin with uppercase letters
var titleCaseProdCvt = newStringConverter(map[string]string{
	"instance": "UHost",
	"lb":       "ULB",
})

//dbModeCvt is used to covert basic to Normal and convert ha to HA
var dbModeCvt = newStringConverter(map[string]string{
	"basic": "Normal",
	"ha":    "HA",
})

//backupTypeCvt is used to transform string to int for backup type when read db backups
var backupTypeCvt = newIntConverter(map[int]string{
	0: "automatic",
	1: "manual",
})

//pgValueTypeCvt is used to transform int to string for value type after read parameter groups
var pgValueTypeCvt = newIntConverter(map[int]string{
	0:  "unknown",
	10: "int",
	20: "string",
	30: "bool",
})

// eg. LOCAL_SSD -> local_ssd
var upperCvt = newUpperConverter(nil)

// eg. localSSD -> local_ssd
var lowerCamelCvt = newLowerCamelConverter(nil)

// eg. LocalSSD -> local_ssd
var upperCamelCvt = newUpperCamelConverter(nil)
