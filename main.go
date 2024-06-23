package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
)

type errMsg error

type model struct {
	quitting bool
	err      error
	canvas   [][]bool
	stable   bool
}

var printConfig = struct {
	emptyCell   string
	filledCell  string
	lineSpace   string
	refreshFreq int64
	noColor     bool
}{}

var (
	emptyCanvas string
)

func init() {
	flag.StringVar(&printConfig.filledCell, "filled", "■", "Filled cell character")
	flag.StringVar(&printConfig.emptyCell, "empty", " ", "Empty cell character")
	flag.StringVar(&printConfig.lineSpace, "line-space", "", "Line space character")
	flag.StringVar(&emptyCanvas, "canvas-only", "", "Empty canvas")
	flag.Int64Var(&printConfig.refreshFreq, "freq", 100, "Refresh frequency")
	flag.BoolVar(&printConfig.noColor, "no-color", false, "Disable colors")
}

func getTerminalSize() (int, int, error) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 0, 0, err
	}
	return height, width, nil
}

func initialModel() model {
	var canvas [][]bool
	for _, arg := range flag.Args() {
		f, err := os.Open(arg)
		if err != nil {
			fmt.Println("Error opening file: ", err)
			os.Exit(1)
		}

		conts, err := io.ReadAll(f)
		if err != nil {
			fmt.Println("Error reading file: ", err)
			os.Exit(1)
		}

		canvas, err = parseCanvas(conts)
		if err != nil {
			fmt.Println("Error parsing canvas: ", err)
			os.Exit(1)
		}

		rows, cols, err := getTerminalSize()
		if err != nil {
			fmt.Println("Error getting terminal size: ", err)
			os.Exit(1)
		}

		canvas = resizeCanvas(canvas, rows, cols) // -4 for the padding
		break
	}

	return model{canvas: canvas}
}

func (m model) Init() tea.Cmd {
	return tea.Tick(time.Millisecond*time.Duration(printConfig.refreshFreq), func(t time.Time) tea.Msg {
		return t
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}

	case errMsg:
		m.err = msg
		return m, nil

	case time.Time:
		if m.stable {
			return m, tea.Quit
		}

		canvas := canvasNext(m.canvas)
		if eq(canvas, m.canvas) {
			m.stable = true
		}

		m.canvas = canvas
		return m, m.Init()
	default:
		return m, nil
	}
}

func eq(a, b [][]bool) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}

		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}

	return true
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}

	color := yellow
	footMsg := "press q to quit"
	if m.stable {
		footMsg = "REACHED STABILITY!!"
		color = green
	}

	if printConfig.noColor {
		color = nil
	}

	return fmt.Sprintf("\n\n%s\n\n%s\n\n", printCanvas(m.canvas, color), footMsg)
}

func getNeighbors(canvas [][]bool, i, j int) int {
	neighbors := 0
	for x := i - 1; x <= i+1; x++ {
		for y := j - 1; y <= j+1; y++ {
			if x < 0 || y < 0 || x >= len(canvas) || y >= len(canvas[i]) {
				continue
			}

			if x == i && y == j {
				continue
			}

			if canvas[x][y] {
				neighbors++
			}
		}
	}

	return neighbors
}

func canvasNext(canvas [][]bool) [][]bool {
	next := make([][]bool, len(canvas))
	for i := range canvas {
		next[i] = make([]bool, len(canvas[i]))
	}

	for i := range canvas {
		for j := range canvas[i] {
			neighbors := getNeighbors(canvas, i, j)
			if canvas[i][j] {
				switch {
				case neighbors < 2, neighbors > 3:
					next[i][j] = false
				case neighbors >= 2, neighbors <= 3:
					next[i][j] = true
				}

				continue
			}

			if neighbors == 3 {
				next[i][j] = true
			}
		}
	}

	return next
}

func handleEmptyCanvas() {
	dimSplit := strings.Split(emptyCanvas, ",")
	if len(dimSplit) != 2 {
		fmt.Printf("Invalid dimensions: %s\n", emptyCanvas)
		os.Exit(1)
	}

	m, err := strconv.Atoi(dimSplit[0])
	if err != nil {
		fmt.Printf("Invalid dimensions: %s\n", emptyCanvas)
		os.Exit(1)
	}

	n, err := strconv.Atoi(dimSplit[1])
	if err != nil {
		fmt.Printf("Invalid dimensions: %s\n", emptyCanvas)
		os.Exit(1)
	}

	for _, row := range getEmptyCanvas(m, n) {
		for _, cell := range row {
			fmt.Print(cell)
		}
		fmt.Println()
	}
}

func parseCanvas(canvas []byte) ([][]bool, error) {
	parsedCanvas := [][]bool{}
	currRow := []bool{}

	for _, b := range bytes.TrimSpace(canvas) {
		switch b {
		case '1':
			currRow = append(currRow, true)
		case '0':
			currRow = append(currRow, false)
		case '\n':
			parsedCanvas = append(parsedCanvas, currRow)
			currRow = []bool{}
		case ' ':
			// ignore trailing spaces
			continue
		default:
			return nil, fmt.Errorf("unknown character '%v'", b)
		}
	}

	if len(currRow) > 0 {
		parsedCanvas = append(parsedCanvas, currRow)
	}

	return parsedCanvas, nil
}

func resizeCanvas(original [][]bool, x, y int) [][]bool {
	// Create a new 2D array with the new dimensions
	resized := make([][]bool, x)
	for i := range resized {
		resized[i] = make([]bool, y)
	}

	// Calculate the starting indices for centering the original array
	startRow := (x - len(original)) / 2
	startCol := (y - len(original[0])) / 2

	// Copy values from the original array to the centered position in the new array
	for i, row := range original {
		for j, cell := range row {
			resized[startRow+i][startCol+j] = cell
		}
	}

	return resized
}

func getEmptyCanvas(m, n int) [][]int {
	rows := make([][]int, m)
	for i := range rows {
		rows[i] = make([]int, n)
	}

	return rows
}

func printCanvas(canvas [][]bool, color colorFunc) string {
	var buf bytes.Buffer
	for _, row := range canvas {
		for _, cell := range row {
			if cell {
				if color != nil {
					buf.WriteString(color(printConfig.filledCell) + printConfig.lineSpace)
					continue
				}

				buf.WriteString(printConfig.filledCell + printConfig.lineSpace)
				continue
			}

			buf.WriteString(printConfig.emptyCell + printConfig.lineSpace)
		}

		buf.WriteString("\n")
	}

	return buf.String()
}

func main() {
	flag.Parse()

	if emptyCanvas != "" {
		handleEmptyCanvas()
		return
	}

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
