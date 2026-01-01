/*
Package ui handles the interactive terminal interface using Bubble Tea.
*/
package ui

import (
	"fmt"
	"strings"

	"github.com/004Ongoro/swiftstack/internal/cache"
	"github.com/004Ongoro/swiftstack/internal/engine"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().MarginLeft(2).Bold(true).Foreground(lipgloss.Color("#00ADD8"))
	docStyle   = lipgloss.NewStyle().Margin(1, 2)
	checkStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
)

type item struct {
	id, title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type Step int

const (
	StepName Step = iota
	StepBase
	StepAddons
	StepConfirm
	StepProcessing
	StepDone
)

type WizardModel struct {
	step           Step
	projectName    textinput.Model
	baseList       list.Model
	addonList      list.Model
	selectedBase   string
	selectedAddons map[int]struct{} // Tracks indexes of checked addons
	err            error
	status         string
}

func InitialModel() WizardModel {
	// 1. Load the dynamic manifest
	manifest, _ := cache.LoadManifest()

	// 2. Initialize Project Name Input
	ti := textinput.New()
	ti.Placeholder = "my-awesome-app"
	ti.Focus()

	// 3. Initialize Base Selection List from manifest
	baseItems := make([]list.Item, len(manifest.Bases))
	for i, b := range manifest.Bases {
		baseItems[i] = item{id: b.ID, title: b.Title, desc: b.Description}
	}
	bl := list.New(baseItems, list.NewDefaultDelegate(), 0, 0)
	bl.Title = "Select Base Template"

	// 4. Initialize Addon Selection List from manifest
	addonItems := make([]list.Item, len(manifest.Addons))
	for i, a := range manifest.Addons {
		addonItems[i] = item{id: a.ID, title: a.Title, desc: a.Description}
	}
	al := list.New(addonItems, list.NewDefaultDelegate(), 0, 0)
	al.Title = "Select Addons (Space to toggle)"

	return WizardModel{
		step:           StepName,
		projectName:    ti,
		baseList:       bl,
		addonList:      al,
		selectedAddons: make(map[int]struct{}),
	}
}

func (m WizardModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m WizardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.baseList.SetSize(msg.Width-h, msg.Height-v)
		m.addonList.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case " ": // Toggle Addon
			if m.step == StepAddons {
				idx := m.addonList.Index()
				if _, ok := m.selectedAddons[idx]; ok {
					delete(m.selectedAddons, idx)
				} else {
					m.selectedAddons[idx] = struct{}{}
				}
				return m, nil
			}

		case "enter":
			switch m.step {
			case StepName:
				if m.projectName.Value() != "" {
					m.step = StepBase
				}
			case StepBase:
				i, ok := m.baseList.SelectedItem().(item)
				if ok {
					m.selectedBase = i.id
					m.step = StepAddons
				}
			case StepAddons:
				m.step = StepConfirm
			case StepConfirm:
				m.step = StepProcessing
				m.status = "Assembling project..."
				return m, executeGeneration(m)
			case StepDone:
				return m, tea.Quit
			}
		}

	case error:
		m.err = msg
		m.step = StepDone
		return m, nil
	}

	switch m.step {
	case StepName:
		m.projectName, cmd = m.projectName.Update(msg)
	case StepBase:
		m.baseList, cmd = m.baseList.Update(msg)
	case StepAddons:
		m.addonList, cmd = m.addonList.Update(msg)
	}

	return m, cmd
}

func (m WizardModel) View() string {
	switch m.step {
	case StepName:
		return fmt.Sprintf("\n%s\n\n%s\n\n(Enter to continue)", titleStyle.Render("Project Name?"), m.projectName.View())
	case StepBase:
		return docStyle.Render(m.baseList.View())
	case StepAddons:
		// We override the default view to show checkmarks
		view := m.addonList.View()
		return docStyle.Render(view + "\n [Space] Toggle | [Enter] Continue")
	case StepConfirm:
		addons := []string{}
		for idx := range m.selectedAddons {
			addons = append(addons, cache.GetAvailableAddons()[idx].Title)
		}
		return fmt.Sprintf("\n%s\n\nProject: %s\nBase: %s\nAddons: %s\n\n(Enter to Start)",
			titleStyle.Render("Final Check"), m.projectName.Value(), m.selectedBase, strings.Join(addons, ", "))
	case StepProcessing:
		return fmt.Sprintf("\n⏳ %s", m.status)
	case StepDone:
		if m.err != nil {
			return fmt.Sprintf("\n❌ Error: %v", m.err)
		}
		return titleStyle.Render("\n✨ Done! Project assembled. Press 'q' to exit.")
	}
	return ""
}

func executeGeneration(m WizardModel) tea.Cmd {
	return func() tea.Msg {
		addons := []string{}
		available := cache.GetAvailableAddons()
		for idx := range m.selectedAddons {
			addons = append(addons, available[idx].ID)
		}

		opts := engine.ProjectOptions{
			Name:        m.projectName.Value(),
			OutputPath:  ".",
			BaseSlice:   m.selectedBase,
			AddonSlices: addons,
		}
		return engine.GenerateProject(opts)
	}
}
