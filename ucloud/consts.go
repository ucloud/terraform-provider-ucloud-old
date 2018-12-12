package ucloud

import (
	"time"
)

const (
	// DefaultMaxRetries is default max retry attempts number
	DefaultMaxRetries = 3

	// DefaultInSecure is a default value to enable https
	DefaultInSecure = false

	// DefaultWaitInterval is the inteval to wait for state changed after resource is created
	DefaultWaitInterval = 10 * time.Second

	// DefaultWaitMaxAttempts is the max attempts number to wait for state changed after resource is created
	DefaultWaitMaxAttempts = 10

	// DefaultWaitIgnoreError is if it will ignore error during wait for state changed after resource is created
	DefaultWaitIgnoreError = false
)

const (
	// StatusPending is status defined by provider, only use to wrap remote resource status as string representation for state waiter
	StatusPending = "pending"
)

// listenerStatusCvt is used to covert int to string for status after read lb listener
var listenerStatusCvt = newIntConverter(map[int]string{
	0: "allNormal",
	1: "partNormal",
	2: "allException",
})

// lbAttachmentStatusCvt is used to covert int to string for status after read lb attachment
var lbAttachmentStatusCvt = newIntConverter(map[int]string{
	0: "normalRunning",
	1: "exceptionRunning",
})

// lowerCaseProdCvt is used to covert one lower string to another lower string
var lowerCaseProdCvt = newStringConverter(map[string]string{
	"instance": "uhost",
	"lb":       "ulb",
})

// titleCaseProdCvt is used to covert one lower string to another string begin with uppercase letters
var titleCaseProdCvt = newStringConverter(map[string]string{
	"instance": "UHost",
	"lb":       "ULB",
})

// dbModeCvt is used to covert basic to Normal and convert ha to HA
var dbModeCvt = newStringConverter(map[string]string{
	"basic": "Normal",
	"ha":    "HA",
})

// backupTypeCvt is used to transform string to int for backup type when read db backups
var backupTypeCvt = newIntConverter(map[int]string{
	0: "automatic",
	1: "manual",
})

// pgValueTypeCvt is used to transform int to string for value type after read parameter groups
var pgValueTypeCvt = newIntConverter(map[int]string{
	0:  "unknown",
	10: "int",
	20: "string",
	30: "bool",
})

// upperCvt is used to transform uppercase with underscore to lowercase with underscore. eg. LOCAL_SSD -> local_ssd
var upperCvt = newUpperConverter(nil)

// lowerCamelCvt is used to transform lower camel case to lowercase with underscore. eg. localSSD -> local_ssd
var lowerCamelCvt = newLowerCamelConverter(nil)

// upperCamelCvt is used to transform uppercamel case to lowercase with underscore. eg. LocalSSD -> local_ssd
var upperCamelCvt = newUpperCamelConverter(nil)
