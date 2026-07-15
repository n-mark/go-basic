package repository

import (
	"fmt"
	"example.com/go-basic/packages/internal/models/chess"
	"example.com/go-basic/packages/internal/models/player"
)

func (r *Repo) AddPlayer(p player.Player) (int, error) {
	if p.ID == 0 {
		p.ID = r.reservePlayerID()
	}
	p.ID = normalizeID(p.ID)
	tmp := p
	r.AppendPlayer(tmp)
	if err := r.PersistPlayers(); err != nil {
		r.RemovePlayer(r.PlayerIndexByID(tmp.ID))
		return 0, err
	}
	return p.ID, nil
}

func (r *Repo) UpdatePlayer(id int, p player.Player) error {
	idx := r.PlayerIndexByID(id)
	if idx < 0 {
		return fmt.Errorf("player с id=%d не найден", id)
	}
	p.ID = id
	p.ID = normalizeID(p.ID)
	old := r.Players[idx]
	r.ReplacePlayer(idx, p)
	if err := r.PersistPlayers(); err != nil {
		r.ReplacePlayer(idx, old)
		return err
	}
	return nil
}

func (r *Repo) DeletePlayer(id int) error {
	idx := r.PlayerIndexByID(id)
	if idx < 0 {
		return fmt.Errorf("player с id=%d не найден", id)
	}
	backup := r.Players[idx]
	r.RemovePlayer(idx)
	if err := r.PersistPlayers(); err != nil {
		r.AppendPlayer(backup)
		return err
	}
	return nil
}

func (r *Repo) GetPlayer(id int) (player.Player, bool) {
	idx := r.PlayerIndexByID(id)
	if idx < 0 {
		return player.Player{}, false
	}
	r.muPlayers.Lock()
	p := r.Players[idx]
	r.muPlayers.Unlock()
	return p, true
}

func (r *Repo) GetAllPlayers() []player.Player {
	return r.PlayersSnapshot()
}

func (r *Repo) AddFigure(f chess.Figure) (int, error) {
	if f.ID == 0 {
		f.ID = r.reserveFigureID()
	}
	f.ID = normalizeID(f.ID)
	tmp := f
	r.AppendFigure(tmp)
	if err := r.PersistFigures(); err != nil {
		r.RemoveFigure(r.FigureIndexByID(tmp.ID))
		return 0, err
	}
	return f.ID, nil
}

func (r *Repo) UpdateFigure(id int, f chess.Figure) error {
	idx := r.FigureIndexByID(id)
	if idx < 0 {
		return fmt.Errorf("figure с id=%d не найдена", id)
	}
	f.ID = id
	f.ID = normalizeID(f.ID)
	old := r.Figures[idx]
	r.ReplaceFigure(idx, f)
	if err := r.PersistFigures(); err != nil {
		r.ReplaceFigure(idx, old)
		return err
	}
	return nil
}

func (r *Repo) DeleteFigure(id int) error {
	idx := r.FigureIndexByID(id)
	if idx < 0 {
		return fmt.Errorf("figure с id=%d не найдена", id)
	}
	backup := r.Figures[idx]
	r.RemoveFigure(idx)
	if err := r.PersistFigures(); err != nil {
		r.AppendFigure(backup)
		return err
	}
	return nil
}

func (r *Repo) GetFigure(id int) (chess.Figure, bool) {
	idx := r.FigureIndexByID(id)
	if idx < 0 {
		return chess.Figure{}, false
	}
	r.muFigures.Lock()
	f := r.Figures[idx]
	r.muFigures.Unlock()
	return f, true
}

func (r *Repo) GetAllFigures() []chess.Figure {
	return r.FiguresSnapshot()
}

func (r *Repo) AddMove(m player.Move) (int, error) {
	if m.ID == 0 {
		m.ID = r.reserveMoveID()
	}
	m.ID = normalizeID(m.ID)
	tmp := m
	r.AppendMove(tmp)
	if err := r.PersistMoves(); err != nil {
		r.RemoveMove(r.MoveIndexByID(tmp.ID))
		return 0, err
	}
	return m.ID, nil
}

func (r *Repo) UpdateMove(id int, m player.Move) error {
	idx := r.MoveIndexByID(id)
	if idx < 0 {
		return fmt.Errorf("move с id=%d не найден", id)
	}
	m.ID = id
	m.ID = normalizeID(m.ID)
	old := r.Moves[idx]
	r.ReplaceMove(idx, m)
	if err := r.PersistMoves(); err != nil {
		r.ReplaceMove(idx, old)
		return err
	}
	return nil
}

func (r *Repo) DeleteMove(id int) error {
	idx := r.MoveIndexByID(id)
	if idx < 0 {
		return fmt.Errorf("move с id=%d не найден", id)
	}
	backup := r.Moves[idx]
	r.RemoveMove(idx)
	if err := r.PersistMoves(); err != nil {
		r.AppendMove(backup)
		return err
	}
	return nil
}

func (r *Repo) GetMove(id int) (player.Move, bool) {
	idx := r.MoveIndexByID(id)
	if idx < 0 {
		return player.Move{}, false
	}
	r.muMoves.Lock()
	m := r.Moves[idx]
	r.muMoves.Unlock()
	return m, true
}

func (r *Repo) GetAllMoves() []player.Move {
	return r.MovesSnapshot()
}

func normalizeID(id int) int {
	if id < 0 {
		return 0
	}
	return id
}
