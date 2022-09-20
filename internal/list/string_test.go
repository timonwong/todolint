package list

import (
	"flag"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	t.Parallel()

	set := NewString("TODO", "FIXME")
	assert.Equal(t, "TODO,FIXME", set.String())
}

func TestString_Flag(t *testing.T) {
	testCases := []struct {
		name      string
		flagValue string
		want      []string
	}{
		{
			name:      "empty",
			flagValue: "",
			want:      nil,
		},
		{
			name:      "todo",
			flagValue: "todo",
			want:      []string{"todo"},
		},
		{
			name:      "todo-and-fixme",
			flagValue: "todo,fixme",
			want:      []string{"todo", "fixme"},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			f := String{}
			fs := flag.NewFlagSet("test", flag.ContinueOnError)
			fs.SetOutput(io.Discard)
			fs.Var(&f, "set", "")

			err := fs.Parse([]string{"-set=" + tc.flagValue})
			require.NoError(t, err)
			assert.EqualValues(t, tc.want, f)
		})
	}
}
