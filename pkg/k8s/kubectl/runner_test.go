package kubectl

import (
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type mockCommandFactory struct {
	mock.Mock
}

func (m *mockCommandFactory) create(config *CommandConfig) (commands []*Command, err error) {
	args := m.Called(config)
	return args.Get(0).([]*Command), args.Error(1)
}

var _ CommandFactory = &mockCommandFactory{}

func TestRunCommand(t *testing.T) {

	command := Command{
		Command: &cobra.Command{
			Run: func(cmd *cobra.Command, args []string) {
				t.Log("running...")
			},
		},
		args: []string{"arg1", "arg2"},
	}
	command.Flags().String("flag", "test", "test")
	_ = command.Flags().Set("flag", "flag1-value")

	mock := new(mockCommandFactory)
	c := &CommandConfig{}
	mock.On("create", c).Return([]*Command{&command, &command}, nil)

	if err := RunCommand(mock, c); err != nil {
		assert.Error(t, err)
	}
}
