// Copyright 2013 Beego Samples authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package models

import (
	"container/list"
	"time"

	"github.com/astaxie/beego"
	"github.com/beego/i18n"
)

type EventType int

const (
	EVENT_JOIN = iota
	EVENT_LEAVE
	EVENT_MESSAGE
)

type BaseController struct {
	beego.Controller // Embed struct that has stub implementation of the interface.
	i18n.Locale      // For i18n usage when process data and render template.
}
type Event struct {
	TypeEvent       EventType // JOIN, LEAVE, MESSAGE
	User            string
	Profiles        []string
	Timestamp       int // Unix timestamp (secs)
	Content         map[string]interface{}
	FechaCreacion   time.Time
	UserDestination []string
	Alias           string
	EstiloIcono     string
	Estado          string
}

const archiveSize = 20

// Event archives.
var archive = list.New()

// NewArchive saves new event to archive list.
func NewArchive(event Event) {
	if archive.Len() >= archiveSize {
		archive.Remove(archive.Front())
	}
	// beego.Info("archivo", archive)
	archive.PushBack(event)
}

// GetEvents returns all events after lastReceived.
func GetEvents(lastReceived int) []Event {
	events := make([]Event, 0, archive.Len())
	for event := archive.Front(); event != nil; event = event.Next() {
		e := event.Value.(Event)
		if e.Timestamp > int(lastReceived) {
			events = append(events, e)
		}
	}
	return events
}
