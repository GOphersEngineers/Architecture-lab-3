package lang

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/roman-mazur/architecture-lab-3/painter"
	"golang.org/x/exp/shiny/screen"
)

type Parser struct {
	Background painter.Operation
	Rectangle  *painter.BgRect
	Figures    []*painter.TFigure
	Movements  []painter.Operation
	Update     painter.Operation
}

func (cp *Parser) Parse(input io.Reader) ([]painter.Operation, error) {
	cp.initialize()
	scanner := bufio.NewScanner(input)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		if err := cp.parseLine(line); err != nil {
			return nil, err
		}
	}
	return cp.getFinalResult(), nil
}

func (cp *Parser) parseLine(line string) error {
	tokens := strings.Split(line, " ")
	if len(tokens) < 1 {
		return fmt.Errorf("invalid command format: %s", line)
	}

	command := tokens[0]
	args, err := convertToIntArgs(tokens, 400)
	if err != nil {
		return fmt.Errorf("invalid argument format for %s: %s", command, line)
	}

	switch command {
	case "white":
		cp.setBackground(painter.WhiteFill)
	case "green":
		cp.setBackground(painter.GreenFill)
	case "bgrect":
		if len(args) != 4 {
			return fmt.Errorf("invalid number of arguments for bgrect: %s", line)
		}
		cp.setRectangle(args)
	case "figure":
		if len(args) != 2 {
			return fmt.Errorf("invalid number of arguments for figure: %s", line)
		}
		cp.addFigure(args)
	case "move":
		if len(args) != 2 {
			return fmt.Errorf("invalid number of arguments for move: %s", line)
		}
		cp.addMovement(args)
	case "reset":
		cp.resetState()
		cp.setBackground(painter.Reset)
	case "update":
		cp.setUpdate(painter.UpdateOp)
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
	return nil
}

func convertToIntArgs(tokens []string, screenSize int) ([]int, error) {
	args := make([]int, 0, len(tokens)-1)
	for _, arg := range tokens[1:] {
		val, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			return nil, err
		}
		args = append(args, int(val*float64(screenSize)))
	}
	return args, nil
}

func (cp *Parser) getFinalResult() []painter.Operation {
	var result []painter.Operation
	result = append(result, cp.Background)
	if cp.Rectangle != nil {
		result = append(result, cp.Rectangle)
	}
	result = append(result, cp.Movements...)
	cp.Movements = nil // Clear accumulated move operations
	if len(cp.Figures) > 0 {
		for _, figure := range cp.Figures {
			result = append(result, figure)
		}
	}
	if cp.Update != nil {
		result = append(result, cp.Update)
	}
	return result
}

func (cp *Parser) initialize() {
	if cp.Background == nil {
		cp.setBackground(painter.Reset)
	}
	cp.Update = nil
}

func (cp *Parser) resetState() {
	cp.Rectangle = nil
	cp.Figures = nil
	cp.Movements = nil
	cp.Update = nil
}

func (cp *Parser) setBackground(opFunc func(screen.Texture)) {
	cp.Background = painter.OperationFunc(opFunc)
}

func (cp *Parser) setRectangle(args []int) {
	cp.Rectangle = &painter.BgRect{X1: args[0], Y1: args[1], X2: args[2], Y2: args[3]}
}

func (cp *Parser) addFigure(args []int) {
	cp.Figures = append(cp.Figures, &painter.TFigure{X: args[0], Y: args[1]})
}

func (cp *Parser) addMovement(args []int) {
	cp.Movements = append(cp.Movements, &painter.MoveFigures{X: args[0], Y: args[1], Figures: cp.Figures})
}

func (cp *Parser) setUpdate(op painter.Operation) {
	cp.Update = op
}
