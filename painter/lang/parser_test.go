package lang_test

import (
	"strings"
	"testing"

	"github.com/roman-mazur/architecture-lab-3/painter"
	"github.com/roman-mazur/architecture-lab-3/painter/lang"
	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	t.Run("parse valid commands", func(t *testing.T) {
		p := &lang.Parser{}
		ops, err := p.Parse(strings.NewReader(
			"white\nbgrect 0.1 0.1 0.9 0.9\nfigure 0.5 0.5\nmove 0.1 0.1\nupdate\n"))
		assert.NoError(t, err)
		assert.NotNil(t, ops)
		assert.Equal(t, 5, len(ops))
		assert.IsType(t, painter.OperationFunc(painter.WhiteFill), ops[0])
		assert.IsType(t, &painter.BgRect{}, ops[1])
		assert.IsType(t, &painter.MoveFigures{}, ops[2])
		assert.IsType(t, &painter.TFigure{}, ops[3])
		assert.IsType(t, painter.UpdateOp, ops[4])
	})

	t.Run("parse invalid command", func(t *testing.T) {
		p := &lang.Parser{}
		_, err := p.Parse(strings.NewReader("invalidComand\n"))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unknown command")
	})

	t.Run("parse invalid arguments", func(t *testing.T) {
		p := &lang.Parser{}
		_, err := p.Parse(strings.NewReader("bgrect 0.1 0.1 0.9\n"))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid number of arguments")
	})
}
