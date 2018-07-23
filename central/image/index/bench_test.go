package index

import (
	"fmt"
	"testing"

	"bitbucket.org/stack-rox/apollo/central/globalindex"
	"bitbucket.org/stack-rox/apollo/generated/api/v1"
	"bitbucket.org/stack-rox/apollo/pkg/fixtures"
)

func getImageIndex(b *testing.B) Indexer {
	tmpIndex, err := globalindex.TempInitializeIndices("")
	if err != nil {
		b.Fatal(err)
	}
	return New(tmpIndex)
}

func benchmarkAddImageNumThen1(b *testing.B, numImages int) {
	indexer := getImageIndex(b)
	image := fixtures.GetImage()
	addImages(indexer, image, numImages)
	image.Name.Sha = fmt.Sprintf("%d", numImages+1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		indexer.AddImage(image)
	}
}

func addImages(indexer Indexer, image *v1.Image, numImages int) {
	for i := 0; i < numImages; i++ {
		image.Name.Sha = fmt.Sprintf("%d", i)
		indexer.AddImage(image)
	}
}

func benchmarkAddImages(b *testing.B, numImages int) {
	indexer := getImageIndex(b)
	image := fixtures.GetImage()
	for i := 0; i < b.N; i++ {
		addImages(indexer, image, numImages)
	}
}

func BenchmarkAddImages(b *testing.B) {
	for i := 1; i <= 1000; i *= 10 {
		b.Run(fmt.Sprintf("Add Images - %d", i), func(subB *testing.B) {
			benchmarkAddImages(subB, i)
		})
	}
}

func BenchmarkAddImagesThen1(b *testing.B) {
	for i := 10; i <= 1000; i *= 10 {
		b.Run(fmt.Sprintf("Add Images %d then 1", i), func(subB *testing.B) {
			benchmarkAddImageNumThen1(subB, i)
		})
	}
}
