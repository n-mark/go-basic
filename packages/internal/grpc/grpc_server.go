package grpc

import (
	"context"
	// "log/slog"
	"time"

	"example.com/go-basic/packages/internal/models/chess"
	"example.com/go-basic/packages/internal/models/player"
	"example.com/go-basic/packages/internal/service"
	pb "example.com/go-basic/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GrpcServer реализует gRPC API над EntityService.
type GrpcServer struct {
	pb.UnimplementedEntityServiceServer
	entityService *service.EntityService
}

// NewGrpcServer создает новый gRPC-сервер сущностей.
func NewGrpcServer(entityService *service.EntityService) *GrpcServer {
	return &GrpcServer{entityService: entityService}
}

// ----------------- Players -----------------

func (s *GrpcServer) CreatePlayer(_ context.Context, req *pb.CreatePlayerRequest) (*pb.Player, error) {
	// slog.Info("create player invoked via grpc", "payload", req)
	p := player.Player{
		Name:         req.GetName(),
		FiguresColor: req.GetFiguresColor(),
	}
	id, err := s.entityService.CreatePlayer(p)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	p.ID = id
	return playerToProto(p), nil
}

func (s *GrpcServer) UpdatePlayer(_ context.Context, req *pb.UpdatePlayerRequest) (*pb.Player, error) {
	// slog.Info("update player invoked via grpc", "payload", req)
	id := int(req.GetId())
	p := player.Player{
		ID:           id,
		Name:         req.GetName(),
		FiguresColor: req.GetFiguresColor(),
	}
	if err := s.entityService.UpdatePlayer(id, p); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if updated, ok := s.entityService.GetPlayer(id); ok {
		return playerToProto(updated), nil
	}
	return playerToProto(p), nil
}

func (s *GrpcServer) ListPlayers(_ context.Context, _ *pb.Empty) (*pb.PlayersList, error) {
	// slog.Info("list players invoked via grpc")
	players := s.entityService.ListPlayers()
	out := &pb.PlayersList{Players: make([]*pb.Player, 0, len(players))}
	for _, p := range players {
		out.Players = append(out.Players, playerToProto(p))
	}
	return out, nil
}

func (s *GrpcServer) GetPlayer(_ context.Context, req *pb.PlayerById) (*pb.Player, error) {
	// slog.Info("get player invoked via grpc", "payload", req)
	p, found := s.entityService.GetPlayer(int(req.GetId()))
	if !found {
		return nil, status.Error(codes.NotFound, "player не найден")
	}
	return playerToProto(p), nil
}

func (s *GrpcServer) DeletePlayer(_ context.Context, req *pb.PlayerById) (*pb.Response, error) {
	if err := s.entityService.DeletePlayer(int(req.GetId())); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &pb.Response{Message: "player удален"}, nil
}

// ----------------- Figures -----------------

func (s *GrpcServer) CreateFigure(_ context.Context, req *pb.CreateFigureRequest) (*pb.Figure, error) {
	f := chess.Figure{
		Symbol:     rune(req.GetSymbol()),
		PieceColor: req.GetPieceColor(),
		GameColor:  req.GetGameColor(),
	}
	id, err := s.entityService.CreateFigure(f)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	f.ID = id
	return figureToProto(f), nil
}

func (s *GrpcServer) UpdateFigure(_ context.Context, req *pb.UpdateFigureRequest) (*pb.Figure, error) {
	id := int(req.GetId())
	f := chess.Figure{
		ID:         id,
		Symbol:     rune(req.GetSymbol()),
		PieceColor: req.GetPieceColor(),
		GameColor:  req.GetGameColor(),
	}
	if err := s.entityService.UpdateFigure(id, f); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if updated, ok := s.entityService.GetFigure(id); ok {
		return figureToProto(updated), nil
	}
	return figureToProto(f), nil
}

func (s *GrpcServer) ListFigures(_ context.Context, _ *pb.Empty) (*pb.FiguresList, error) {
	figures := s.entityService.ListFigures()
	out := &pb.FiguresList{Figures: make([]*pb.Figure, 0, len(figures))}
	for _, f := range figures {
		out.Figures = append(out.Figures, figureToProto(f))
	}
	return out, nil
}

func (s *GrpcServer) GetFigure(_ context.Context, req *pb.FigureById) (*pb.Figure, error) {
	f, found := s.entityService.GetFigure(int(req.GetId()))
	if !found {
		return nil, status.Error(codes.NotFound, "figure не найдена")
	}
	return figureToProto(f), nil
}

func (s *GrpcServer) DeleteFigure(_ context.Context, req *pb.FigureById) (*pb.Response, error) {
	if err := s.entityService.DeleteFigure(int(req.GetId())); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &pb.Response{Message: "figure удалена"}, nil
}

// ----------------- Moves -----------------

func (s *GrpcServer) CreateMove(_ context.Context, req *pb.CreateMoveRequest) (*pb.Move, error) {
	m := player.Move{
		TimeTook:     time.Duration(req.GetTimeTook()),
		PositionFrom: req.GetPositionFrom(),
		PositionTo:   req.GetPositionTo(),
		Figure:       rune(req.GetFigure()),
	}
	id, err := s.entityService.CreateMove(m)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	m.ID = id
	return moveToProto(m), nil
}

func (s *GrpcServer) UpdateMove(_ context.Context, req *pb.UpdateMoveRequest) (*pb.Move, error) {
	id := int(req.GetId())
	m := player.Move{
		ID:           id,
		TimeTook:     time.Duration(req.GetTimeTook()),
		PositionFrom: req.GetPositionFrom(),
		PositionTo:   req.GetPositionTo(),
		Figure:       rune(req.GetFigure()),
	}
	if err := s.entityService.UpdateMove(id, m); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if updated, ok := s.entityService.GetMove(id); ok {
		return moveToProto(updated), nil
	}
	return moveToProto(m), nil
}

func (s *GrpcServer) ListMoves(_ context.Context, _ *pb.Empty) (*pb.MovesList, error) {
	moves := s.entityService.ListMoves()
	out := &pb.MovesList{Moves: make([]*pb.Move, 0, len(moves))}
	for _, m := range moves {
		out.Moves = append(out.Moves, moveToProto(m))
	}
	return out, nil
}

func (s *GrpcServer) GetMove(_ context.Context, req *pb.MoveById) (*pb.Move, error) {
	m, found := s.entityService.GetMove(int(req.GetId()))
	if !found {
		return nil, status.Error(codes.NotFound, "move не найден")
	}
	return moveToProto(m), nil
}

func (s *GrpcServer) DeleteMove(_ context.Context, req *pb.MoveById) (*pb.Response, error) {
	if err := s.entityService.DeleteMove(int(req.GetId())); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &pb.Response{Message: "move удален"}, nil
}

// ----------------- converters -----------------

func playerToProto(p player.Player) *pb.Player {
	figures := make([]*pb.Figure, 0, len(p.FiguresTook))
	for _, f := range p.FiguresTook {
		figures = append(figures, figureToProto(f))
	}
	moves := make([]*pb.Move, 0, len(p.Moves))
	for _, m := range p.Moves {
		moves = append(moves, moveToProto(m))
	}
	return &pb.Player{
		Id:           int64(p.ID),
		Name:         p.Name,
		FiguresColor: p.FiguresColor,
		FiguresTook:  figures,
		Moves:        moves,
	}
}

func figureToProto(f chess.Figure) *pb.Figure {
	return &pb.Figure{
		Id:         int64(f.ID),
		Symbol:     int32(f.Symbol),
		PieceColor: f.PieceColor,
		GameColor:  f.GameColor,
	}
}

func moveToProto(m player.Move) *pb.Move {
	return &pb.Move{
		Id:           int64(m.ID),
		TimeTook:     int64(m.TimeTook),
		PositionFrom: m.PositionFrom,
		PositionTo:   m.PositionTo,
		Figure:       int32(m.Figure),
	}
}
