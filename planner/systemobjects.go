// Copyright 2022 Molecula Corp. All rights reserved.

package planner

type viewSystemObject struct {
	name      string
	statement string
}

type functionSystemObject struct {
	name     string
	language string
	body     string
}

type modelSystemObject struct {
	name         string
	status       string
	modelType    string
	labels       []string
	inputColumns []string
}
