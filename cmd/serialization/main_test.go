package main

import "testing"

var movie = &Movie{
	Metadata{
		ID:          "123-456-789",
		Title:       "The Test",
		ReleaseYear: 2020,
		Description: "A test movie to benchmark various serialisation format",
	},
	7.5,
}

func BenchmarkJsonMarshalling(b *testing.B) {
	for b.Loop() {
		serializeToJson(movie)
	}
}

func BenchmarkXMLMarshalling(b *testing.B) {
	for b.Loop() {
		serializeToXML(movie)
	}
}

func BenchmarkProtobufMarshalling(b *testing.B) {
	pb := toProto(movie)
	for b.Loop() {
		serializeToProtobuf(pb)
	}
}
