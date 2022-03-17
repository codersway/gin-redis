package chapter7

import (
	"github.com/stretchr/testify/assert"
	"redis-in-action/conf"
	"testing"
)

func TestJob(t *testing.T) {

	conn := conf.Conn()
	client := NewClient(conn)

	t.Run("test is qualified for job", func(t *testing.T) {
		client.AddJob("test", []string{"q1", "q2", "q3"})
		res := client.IsQualified("test", []string{"q1", "q2", "q3"})
		assert.Equal(t, 0, len(res))

		res = client.IsQualified("test", []string{"q1", "q2"})
		assert.Equal(t, 0, len(res))

		client.Conn.FlushDB(ctx)
	})

	t.Run("test index and find jobs", func(t *testing.T) {
		client.IndexJob("test1", []string{"q1", "q2", "q3"})
		client.IndexJob("test2", []string{"q1", "q2", "q3"})
		client.IndexJob("test3", []string{"q1", "q2", "q3"})

	})
}
