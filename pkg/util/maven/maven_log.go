/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package maven

import (
	"encoding/json"

	"github.com/apache/camel-k/pkg/util/log"
)

// nolint: stylecheck
type mavenLog struct {
	Level            string `json:"level"`
	Ts               string `json:"ts"`
	Logger           string `json:"logger"`
	Msg              string `json:"msg"`
	Class            string `json:"class"`
	CallerMethodName string `json:"caller_method_name"`
	CallerFileName   string `json:"caller_file_name"`
	CallerLineNumber int    `json:"caller_line_number"`
	Thread           string `json:"thread"`
}

const (
	TRACE = "TRACE"
	DEBUG = "DEBUG"
	INFO  = "INFO"
	WARN  = "WARN"
	ERROR = "ERROR"
	FATAL = "FATAL"
)

var mavenLogger = log.WithName("maven.build")

func mavenLogHandler(s string) {
	mavenLog, parseError := parseLog(s)
	if parseError == nil {
		normalizeLog(mavenLog)
	} else {
		// Why we are ignoring the parsing errors here: there are a few scenarios where this would likely occur.
		// For example, if something outside of Maven outputs something (i.e.: the JDK, a misbehaved plugin,
		// etc). The build may still have succeeded, though.
		nonNormalizedLog(s)
	}
}

func parseLog(line string) (mavenLog, error) {
	var l mavenLog
	err := json.Unmarshal([]byte(line), &l)
	return l, err
}

func normalizeLog(mavenLog mavenLog) {
	switch mavenLog.Level {
	case DEBUG, TRACE:
		mavenLogger.Debug(mavenLog.Msg)
	case INFO, WARN:
		mavenLogger.Info(mavenLog.Msg)
	case ERROR, FATAL:
		mavenLogger.Errorf(nil, mavenLog.Msg)
	}
}

func nonNormalizedLog(rawLog string) {
	mavenLogger.Info(rawLog)
}
