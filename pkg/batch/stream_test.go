// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package batch

import (
	"time"

	ginkgov2 "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgov2.Describe("Stream", func() {

	type addTestArgs struct {
		entries        []Entry
		expectedStream Stream
	}

	type sortTestArgs struct {
		stream         Stream
		expectedStream Stream
	}

	timeStamp1 := time.Now()
	timeStamp2 := timeStamp1.Add(time.Second)
	timeStamp3 := timeStamp2.Add(time.Second)

	ginkgov2.DescribeTable("#add",
		func(args addTestArgs) {
			stream := Stream{}
			for _, entry := range args.entries {
				stream.add(entry.Timestamp, entry.Line)
			}

			gomega.Expect(stream).To(gomega.Equal(args.expectedStream))
		},
		ginkgov2.Entry("add one entry", addTestArgs{
			entries: []Entry{
				{
					Timestamp: timeStamp1,
					Line:      "Line1",
				},
			},
			expectedStream: Stream{
				Entries: []Entry{
					{
						Timestamp: timeStamp1,
						Line:      "Line1",
					},
				},
				lastTimestamp:     timeStamp1,
				isEntryOutOfOrder: false,
			},
		}),
		ginkgov2.Entry("add two entries", addTestArgs{
			entries: []Entry{
				{
					Timestamp: timeStamp1,
					Line:      "Line1",
				},
				{
					Timestamp: timeStamp2,
					Line:      "Line2",
				},
			},
			expectedStream: Stream{
				Entries: []Entry{
					{
						Timestamp: timeStamp1,
						Line:      "Line1",
					},
					{
						Timestamp: timeStamp2,
						Line:      "Line2",
					},
				},
				lastTimestamp:     timeStamp2,
				isEntryOutOfOrder: false,
			},
		}),
		ginkgov2.Entry("add two entries without order", addTestArgs{
			entries: []Entry{
				{
					Timestamp: timeStamp2,
					Line:      "Line2",
				},
				{
					Timestamp: timeStamp1,
					Line:      "Line1",
				},
			},
			expectedStream: Stream{
				Entries: []Entry{
					{
						Timestamp: timeStamp2,
						Line:      "Line2",
					},
					{
						Timestamp: timeStamp1,
						Line:      "Line1",
					},
				},
				lastTimestamp:     timeStamp2,
				isEntryOutOfOrder: true,
			},
		}),
	)

	ginkgov2.DescribeTable("#sort",
		func(args sortTestArgs) {
			args.stream.sort()
			gomega.Expect(args.stream).To(gomega.Equal(args.expectedStream))
		},
		ginkgov2.Entry("sort stream with two out of order entries", sortTestArgs{
			stream: Stream{
				Entries: []Entry{
					{
						Timestamp: timeStamp2,
						Line:      "Line2",
					},
					{
						Timestamp: timeStamp1,
						Line:      "Line1",
					},
				},
				lastTimestamp:     timeStamp2,
				isEntryOutOfOrder: true,
			},
			expectedStream: Stream{
				Entries: []Entry{
					{
						Timestamp: timeStamp1,
						Line:      "Line1",
					},
					{
						Timestamp: timeStamp2,
						Line:      "Line2",
					},
				},
				lastTimestamp:     timeStamp2,
				isEntryOutOfOrder: false,
			},
		}),
		ginkgov2.Entry("sort stream with three out of order entries", sortTestArgs{
			stream: Stream{
				Entries: []Entry{
					{
						Timestamp: timeStamp3,
						Line:      "Line3",
					},
					{
						Timestamp: timeStamp2,
						Line:      "Line2",
					},
					{
						Timestamp: timeStamp1,
						Line:      "Line1",
					},
				},
				lastTimestamp:     timeStamp3,
				isEntryOutOfOrder: true,
			},
			expectedStream: Stream{
				Entries: []Entry{
					{
						Timestamp: timeStamp1,
						Line:      "Line1",
					},
					{
						Timestamp: timeStamp2,
						Line:      "Line2",
					},
					{
						Timestamp: timeStamp3,
						Line:      "Line3",
					},
				},
				lastTimestamp:     timeStamp3,
				isEntryOutOfOrder: false,
			},
		}),
		ginkgov2.Entry("sort stream with no out of order entries", sortTestArgs{
			stream: Stream{
				Entries: []Entry{
					{
						Timestamp: timeStamp1,
						Line:      "Line1",
					},
					{
						Timestamp: timeStamp2,
						Line:      "Line2",
					},
					{
						Timestamp: timeStamp3,
						Line:      "Line3",
					},
				},
				lastTimestamp:     timeStamp3,
				isEntryOutOfOrder: false,
			},
			expectedStream: Stream{
				Entries: []Entry{
					{
						Timestamp: timeStamp1,
						Line:      "Line1",
					},
					{
						Timestamp: timeStamp2,
						Line:      "Line2",
					},
					{
						Timestamp: timeStamp3,
						Line:      "Line3",
					},
				},
				lastTimestamp:     timeStamp3,
				isEntryOutOfOrder: false,
			},
		}),
	)
})
