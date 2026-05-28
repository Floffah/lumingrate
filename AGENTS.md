# Lumingrate

Lumingrate is a terminal-based text adventure (interactive fiction) game. It uses an advanced TUI to provide an immersive user interface.

You should not need to care about the story, if you are asked to write story content, you should not write it yourself, but based on a consultation from the author.
The story is based on a world I call Luminous. It represents an alternate timeline branching in the year 2025, where a group of researchers discovered some technology that enables post-scarcity, post-capitalism, space travel, and more, but ethically could not be released to the general public or governments as it stands (too powerful), so the researchers set up a rudimentary research base on the planet Emergence which is hundreds of light years away from Earth. The research base grew into a full fledged eutopian community and nation, which is called Luminous. As this world is Asimov-like in size, the narrative of this story follows a specific person who is migrating to Luminous (no one knows about Luminous, so the government outreach team recruits people from earth), rather than the whole world.

Your main purpose is not to write the story, but to help build the harnesses that run the story: the tty and wasm harnesses.

The game is a text adventure, so the player will interact with the game by typing commands. The game will respond to the player's commands and provide feedback on the player's actions. The game will also provide a narrative that progresses as the player interacts with the game.

## Priorities

1. Maintainability: The code should be easy to read, understand, and modify. This includes organising the code in a logical manner.
2. Performance: The game should run smoothly and efficiently, without any noticeable lag or delays.
3. User Experience: The game should provide an immersive and engaging experience for the player.

## Structure

When there is a specific action you must take including linting, testing, and formatting, please use the justfile scripts (invoked in the command line with `just <name>`)

The codebase is organised with a separation between the harness and the content. The harness is the bubble tea TUI, responsible for rendering everything, while the chapters and engine are responsible for the story content and the game logic. 

## Technologies

- The codebase is entirely written in Go
- We primarily use the `github.com/charmbracelet/bubbletea` package for the TUI, and `github.com/charmbracelet/lipgloss` for styling the TUI.
- Some aspects of the `github.com/hajimehoshi/ebiten` package are pulled in for specific purposes like audio, but should never be used for rendering or engine content, as this part is custom.