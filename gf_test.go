package apollogf

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/stretchr/testify/assert"
	_ "net/http/pprof"
	"strings"
	"testing"
)

func BenchmarkAdapter(b *testing.B) {
	namespaces := ""
	apolloArg := ApolloArg{
		Namespaces: strings.Split(namespaces, ","),
		AppId:      "",
		IP:         "",
		Cluster:    "",
	}
	SetAdapters(apolloArg)
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := g.Cfg().Get(ctx, "server.address")
		assert.NoError(b, err)
	}
}
