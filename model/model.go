package model

import (
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/flevin58/versions/cfg"
)

type cmdData struct {
	command  string
	version  string
	where    string
	found    bool
	homebrew bool
}

func (c *cmdData) Record() []string {
	return []string{c.command, c.version, c.where}
}

var (
	cmdHeader cmdData = cmdData{"Command", "Version", "Where", false, false}
	cmdList   []cmdData
)

func getVersion(cmd cfg.Command) string {
	const notavailable = "not available"
	command := exec.Command(cmd.Name, cmd.VersionFlag)
	out, err := command.Output()
	if err != nil {
		return notavailable
	}
	version := strings.Split(string(out), "\n")[cmd.VersionLine-1]
	re := regexp.MustCompile(`^[^\d]*([0-9.]+)`)
	match := re.FindStringSubmatch(version)
	result := match[1]
	return result
}

func Add(cmd cfg.Command) {
	var data cmdData
	where, err := exec.LookPath(cmd.Name)
	if err != nil {
		data = cmdData{
			found:    false,
			homebrew: false,
			command:  cmd.Name,
			version:  "",
			where:    "",
		}
	} else {
		data = cmdData{
			found:    true,
			homebrew: strings.HasPrefix(where, "/opt/homebrew"),
			command:  cmd.Name,
			version:  getVersion(cmd),
			where:    where,
		}
	}
	cmdList = append(cmdList, data)
}

func ToCSV(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to open file %v", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()
	header := cmdHeader.Record()
	if err = w.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %v", err)
	}
	for i := range cmdList {
		record := cmdList[i].Record()
		if err = w.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %v", err)
		}
	}
	return nil
}

func ToTable() {
	// Prepare the table
	HeaderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).AlignHorizontal(lipgloss.Center)
	EvenRowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#fff200"))
	OddRowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff"))
	table := table.New().Border(lipgloss.NormalBorder())
	table.Headers("Command", "Version", "Where")
	table.StyleFunc(func(row, col int) lipgloss.Style {
		switch {
		case row == 0:
			return HeaderStyle
		case row%2 == 0:
			return EvenRowStyle
		default:
			return OddRowStyle
		}
	})

	var emoji string
	for _, row := range cmdList {
		switch {
		case row.homebrew == true:
			emoji = "🍺 "
		case row.homebrew == false && row.found == true:
			emoji = "🟢 "
		default:
			emoji = "🔴 "
		}

		table.Row(emoji+row.command, row.version, row.where)
	}
	fmt.Println(table.Render())
}

func ToText() {
}
