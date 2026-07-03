package app

import (
	"luminrate/internal/engine"
	gamescene "luminrate/internal/scene/game"
	"luminrate/internal/scene/main_menu"

	tea "charm.land/bubbletea/v2"
)

const (
	defaultWidth  = 80
	defaultHeight = 24
)

type scene int

const (
	mainMenuScene scene = iota
	gameScene
)

type Model struct {
	width  int
	height int
	scene  scene

	mainMenu tea.Model
	game     tea.Model
}

func NewModel() Model {
	return Model{
		width:    defaultWidth,
		height:   defaultHeight,
		scene:    mainMenuScene,
		mainMenu: mainmenu.NewModel(),
	}
}

func NewModelWithInitialChapter(initial engine.ChapterID) (Model, bool) {
	game, ok := gamescene.NewModelWithSizeAndInitialChapter(defaultWidth, defaultHeight, initial)
	if !ok {
		return Model{}, false
	}

	return Model{
		width:    defaultWidth,
		height:   defaultHeight,
		scene:    gameScene,
		mainMenu: mainmenu.NewModel(),
		game:     game,
	}, true
}

func (m Model) Init() tea.Cmd {
	if m.scene == gameScene && m.game != nil {
		return m.game.Init()
	}

	return tea.RequestWindowSize
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case mainmenu.StartGameMsg:
		m.scene = gameScene
		m.game = gamescene.NewModelWithSize(m.width, m.height)
		return m, m.game.Init()
	}

	switch m.scene {
	case gameScene:
		return m.updateGame(msg)
	default:
		return m.updateMainMenu(msg)
	}
}

func (m Model) updateMainMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	next, cmd := m.mainMenu.Update(msg)
	m.mainMenu = next
	return m, cmd
}

func (m Model) updateGame(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.game == nil {
		m.game = gamescene.NewModelWithSize(m.width, m.height)
	}

	next, cmd := m.game.Update(msg)
	m.game = next
	return m, cmd
}

func (m Model) View() tea.View {
	if m.scene == gameScene {
		if m.game == nil {
			m.game = gamescene.NewModelWithSize(m.width, m.height)
		}
		return m.game.View()
	}

	return m.mainMenu.View()
}
