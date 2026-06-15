package main

import (
	"reflect"
	"testing"
)

func TestFlagsForPostAddsLinkedInOnlyForHashtags(t *testing.T) {
	tests := []struct {
		name  string
		flags []string
		text  string
		want  []string
	}{
		{
			name:  "defaults without hashtag",
			flags: []string{"-bmt"},
			text:  "plain post",
			want:  []string{"-bmt"},
		},
		{
			name:  "defaults with hashtag",
			flags: []string{"-bmt"},
			text:  "launching today #golang",
			want:  []string{"-bmt", "-l"},
		},
		{
			name:  "strips explicit linkedin without hashtag",
			flags: []string{"-bmt", "--linkedin"},
			text:  "plain post",
			want:  []string{"-bmt"},
		},
		{
			name:  "strips bundled linkedin without hashtag",
			flags: []string{"-bmtl"},
			text:  "plain post",
			want:  []string{"-bmt"},
		},
		{
			name:  "deduplicates bundled linkedin with hashtag",
			flags: []string{"-lbmt"},
			text:  "launching today #golang",
			want:  []string{"-bmt", "-l"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := flagsForPost(tt.flags, tt.text)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("flagsForPost(%v, %q) = %v, want %v", tt.flags, tt.text, got, tt.want)
			}
		})
	}
}

func TestHasHashtag(t *testing.T) {
	tests := []struct {
		name string
		text string
		want bool
	}{
		{name: "plain text", text: "no tags here", want: false},
		{name: "hash only", text: "this is not #", want: false},
		{name: "word hashtag", text: "post to LinkedIn #launch", want: true},
		{name: "numeric hashtag", text: "day #100", want: true},
		{name: "embedded hash ignored", text: "abc#notatag", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasHashtag(tt.text)
			if got != tt.want {
				t.Fatalf("hasHashtag(%q) = %v, want %v", tt.text, got, tt.want)
			}
		})
	}
}
