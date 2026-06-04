package repository

import (
	"example.com/go-basic/packages/internal/models/chess"
	"example.com/go-basic/packages/internal/models/player"
)

func (r *Repo) recalculateCounters() {
	maxPlayerID := 0
	for _, p := range r.Players {
		if p.ID > maxPlayerID {
			maxPlayerID = p.ID
		}
	}
	r.playerIDCounter = maxPlayerID + 1

	maxFigureID := 0
	for _, f := range r.Figures {
		if f.ID > maxFigureID {
			maxFigureID = f.ID
		}
	}
	r.figureIDCounter = maxFigureID + 1

	maxMoveID := 0
	for _, m := range r.Moves {
		if m.ID > maxMoveID {
			maxMoveID = m.ID
		}
	}
	r.moveIDCounter = maxMoveID + 1
}

func (r *Repo) reservePlayerID() int {
	return r.reservePlayerIDWithLock()
}

func (r *Repo) reserveFigureID() int {
	return r.reserveFigureIDWithLock()
}

func (r *Repo) reserveMoveID() int {
	return r.reserveMoveIDWithLock()
}

func (r *Repo) reservePlayerIDWithLock() int {
	r.muPlayers.Lock()
	defer r.muPlayers.Unlock()
	id := r.playerIDCounter
	r.playerIDCounter++
	return id
}

func (r *Repo) reserveFigureIDWithLock() int {
	r.muFigures.Lock()
	defer r.muFigures.Unlock()
	id := r.figureIDCounter
	r.figureIDCounter++
	return id
}

func (r *Repo) reserveMoveIDWithLock() int {
	r.muMoves.Lock()
	defer r.muMoves.Unlock()
	id := r.moveIDCounter
	r.moveIDCounter++
	return id
}

func (r *Repo) ReservePlayerID() int {
	return r.reservePlayerID()
}

func (r *Repo) ReserveFigureID() int {
	return r.reserveFigureID()
}

func (r *Repo) ReserveMoveID() int {
	return r.reserveMoveID()
}

func (r *Repo) AppendPlayer(p player.Player) {
	r.muPlayers.Lock()
	defer r.muPlayers.Unlock()
	r.Players = append(r.Players, p)
}

func (r *Repo) ReplacePlayer(idx int, p player.Player) {
	r.muPlayers.Lock()
	defer r.muPlayers.Unlock()
	if idx >= 0 && idx < len(r.Players) {
		r.Players[idx] = p
	}
}

func (r *Repo) RemovePlayer(idx int) {
	r.muPlayers.Lock()
	defer r.muPlayers.Unlock()
	if idx >= 0 && idx < len(r.Players) {
		r.Players = append(r.Players[:idx], r.Players[idx+1:]...)
	}
}

func (r *Repo) PlayersSnapshot() []player.Player {
	r.muPlayers.Lock()
	defer r.muPlayers.Unlock()
	out := make([]player.Player, len(r.Players))
	copy(out, r.Players)
	return out
}

func (r *Repo) PlayerIndexByID(id int) int {
	r.muPlayers.Lock()
	defer r.muPlayers.Unlock()
	for i, p := range r.Players {
		if p.ID == id {
			return i
		}
	}
	return -1
}

func (r *Repo) PersistPlayers() error {
	if r.storage == nil {
		return nil
	}
	return r.storage.SavePlayers(r.Players)
}

func (r *Repo) AppendFigure(f chess.Figure) {
	r.muFigures.Lock()
	defer r.muFigures.Unlock()
	r.Figures = append(r.Figures, f)
}

func (r *Repo) ReplaceFigure(idx int, f chess.Figure) {
	r.muFigures.Lock()
	defer r.muFigures.Unlock()
	if idx >= 0 && idx < len(r.Figures) {
		r.Figures[idx] = f
	}
}

func (r *Repo) RemoveFigure(idx int) {
	r.muFigures.Lock()
	defer r.muFigures.Unlock()
	if idx >= 0 && idx < len(r.Figures) {
		r.Figures = append(r.Figures[:idx], r.Figures[idx+1:]...)
	}
}

func (r *Repo) FiguresSnapshot() []chess.Figure {
	r.muFigures.Lock()
	defer r.muFigures.Unlock()
	out := make([]chess.Figure, len(r.Figures))
	copy(out, r.Figures)
	return out
}

func (r *Repo) FigureIndexByID(id int) int {
	r.muFigures.Lock()
	defer r.muFigures.Unlock()
	for i, f := range r.Figures {
		if f.ID == id {
			return i
		}
	}
	return -1
}

func (r *Repo) PersistFigures() error {
	if r.storage == nil {
		return nil
	}
	return r.storage.SaveFigures(r.Figures)
}

func (r *Repo) AppendMove(m player.Move) {
	r.muMoves.Lock()
	defer r.muMoves.Unlock()
	r.Moves = append(r.Moves, m)
}

func (r *Repo) ReplaceMove(idx int, m player.Move) {
	r.muMoves.Lock()
	defer r.muMoves.Unlock()
	if idx >= 0 && idx < len(r.Moves) {
		r.Moves[idx] = m
	}
}

func (r *Repo) RemoveMove(idx int) {
	r.muMoves.Lock()
	defer r.muMoves.Unlock()
	if idx >= 0 && idx < len(r.Moves) {
		r.Moves = append(r.Moves[:idx], r.Moves[idx+1:]...)
	}
}

func (r *Repo) MovesSnapshot() []player.Move {
	r.muMoves.Lock()
	defer r.muMoves.Unlock()
	out := make([]player.Move, len(r.Moves))
	copy(out, r.Moves)
	return out
}

func (r *Repo) MoveIndexByID(id int) int {
	r.muMoves.Lock()
	defer r.muMoves.Unlock()
	for i, m := range r.Moves {
		if m.ID == id {
			return i
		}
	}
	return -1
}

func (r *Repo) PersistMoves() error {
	if r.storage == nil {
		return nil
	}
	return r.storage.SaveMoves(r.Moves)
}
