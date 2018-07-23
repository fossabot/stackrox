package index

import (
	"fmt"
	"testing"

	"bitbucket.org/stack-rox/apollo/central/globalindex"
	"bitbucket.org/stack-rox/apollo/generated/api/v1"
	"bitbucket.org/stack-rox/apollo/pkg/fixtures"
)

func getDeploymentIndex(b *testing.B) Indexer {
	tmpIndex, err := globalindex.TempInitializeIndices("")
	if err != nil {
		b.Fatal(err)
	}
	return New(tmpIndex)
}

func benchmarkAddDeploymentNumThen1(b *testing.B, numDeployments int) {
	indexer := getDeploymentIndex(b)
	deployment := fixtures.GetDeployment()
	addDeployments(indexer, deployment, numDeployments)
	deployment.Id = fmt.Sprintf("%d", numDeployments+1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		indexer.AddDeployment(deployment)
	}
}

func addDeployments(indexer Indexer, deployment *v1.Deployment, numDeployments int) {
	for i := 0; i < numDeployments; i++ {
		deployment.Id = fmt.Sprintf("%d", i)
		indexer.AddDeployment(deployment)
	}
}

func benchmarkAddDeployment(b *testing.B, numDeployments int) {
	indexer := getDeploymentIndex(b)
	deployment := fixtures.GetDeployment()
	for i := 0; i < b.N; i++ {
		addDeployments(indexer, deployment, numDeployments)
	}
}

func BenchmarkAddDeployments(b *testing.B) {
	for i := 1; i <= 1000; i *= 10 {
		b.Run(fmt.Sprintf("Add Deployments - %d", i), func(subB *testing.B) {
			benchmarkAddDeployment(subB, i)
		})
	}
}

func BenchmarkAddDeploymentssThen1(b *testing.B) {
	for i := 10; i <= 1000; i *= 10 {
		b.Run(fmt.Sprintf("Add Deployments %d then 1", i), func(subB *testing.B) {
			benchmarkAddDeploymentNumThen1(subB, i)
		})
	}
}
