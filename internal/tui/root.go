package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type MultiModel interface {
	tea.Model

	setScreen(tea.Msg) error
	AddNode(name string, screen tea.Model) error
	AddEdge(incomingEdge string, outgoingEdge string, action tea.Msg) error
}

type fsaEdge struct {
	edgeName string
	msg      tea.Msg
}

type RootModel struct {
	screens          map[string]tea.Model
	fsa              map[fsaEdge]string
	activeScreenName string
}

func (r *RootModel) AddNode(name string, screen tea.Model) error {
	if r.screens == nil {
		return fmt.Errorf("unallocated screens. Ensure model data members have been created")
	}
	r.screens[name] = screen
	return nil
}

func (r *RootModel) AddEdge(incomingEdge string, outgoingEdge string, action tea.Msg) error {
	if r.fsa == nil {
		return fmt.Errorf("unallocated fsa. Ensure model data members have been created")
	}
	r.fsa[fsaEdge{edgeName: outgoingEdge, msg: action}] = incomingEdge
	return nil
}

func NewRootModel() *RootModel {
	return &RootModel{
		screens: make(map[string]tea.Model),
		fsa:     make(map[fsaEdge]string),
	}
}

func (r *RootModel) setScreen(action tea.Msg) error {
	if r.fsa == nil {
		return fmt.Errorf("unallocated fsa. Ensure model data members have been created")
	}
	nextScreen, ok := r.fsa[fsaEdge{edgeName: r.activeScreenName, msg: action}]
	if ok {
		r.activeScreenName = nextScreen
	}
	return nil
}

func (r *RootModel) Init() tea.Cmd {
	for _, value := range r.screens {
		value.Init()
	}
	return nil
}

func (r *RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	r.setScreen(msg)
	return r.screens[r.activeScreenName].Update(msg)

}

func (r *RootModel) View() string {
	return r.screens[r.activeScreenName].View()
}
