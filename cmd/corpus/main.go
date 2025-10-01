package main

import (
	"database/sql"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/corpus/internal/database"
	tui "github.com/corpus/internal/tui/blog"
)

func main() {
	db, err := sql.Open("sqlite", "./my.db?_pragma=foreign_keys(1)")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer db.Close()

	_, err = database.DropArticleTable(db)
	if err != nil {
		fmt.Println(err)
		return
	}

	// create the articles table
	_, err = database.CreateArticleTable(db)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Table articles was created successfully.")

	_, err = database.FillArticleTable(db)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Table articles filled with mock database")

	model := tui.InitialModel()
	article_repo := database.NewSQLArticleRepo(db)
	articles, err := article_repo.GetAll()
	if err != nil {
		fmt.Println(err)
		return
	}

	model.Articles = articles

	p := tea.NewProgram(
		&model,
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel)
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("There's been an error: %v", err)
		os.Exit(1)
	}
}
