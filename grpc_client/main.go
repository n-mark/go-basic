package main

import (
	"context"
	"log"
	"time"

	pb "example.com/go-basic/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50001", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("cannot connect to grpc server: %v", err)
	}

	defer conn.Close()

	client := pb.NewEntityServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// ===== Player methods =====

	createdPlayer, err := client.CreatePlayer(ctx, &pb.CreatePlayerRequest{
		Name:         "Magnus",
		FiguresColor: "white",
	})
	if err != nil {
		log.Fatalf("CreatePlayer error: %v", err)
	}
	log.Printf("CreatePlayer response: %v", createdPlayer)

	updatedPlayer, err := client.UpdatePlayer(ctx, &pb.UpdatePlayerRequest{
		Id:           createdPlayer.Id,
		Name:         "Magnus Carlsen",
		FiguresColor: "black",
	})
	if err != nil {
		log.Fatalf("UpdatePlayer error: %v", err)
	}
	log.Printf("UpdatePlayer response: %v", updatedPlayer)

	playersList, err := client.ListPlayers(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("ListPlayers error: %v", err)
	}
	log.Printf("ListPlayers response: %v", playersList)

	player, err := client.GetPlayer(ctx, &pb.PlayerById{Id: createdPlayer.Id})
	if err != nil {
		log.Fatalf("GetPlayer error: %v", err)
	}
	log.Printf("GetPlayer response: %v", player)

	deletePlayerResp, err := client.DeletePlayer(ctx, &pb.PlayerById{Id: createdPlayer.Id})
	if err != nil {
		log.Fatalf("DeletePlayer error: %v", err)
	}
	log.Printf("DeletePlayer response: %v", deletePlayerResp)

	// ===== Figure methods =====

	createdFigure, err := client.CreateFigure(ctx, &pb.CreateFigureRequest{
		Symbol:     1,
		PieceColor: "white",
		GameColor:  "white",
	})
	if err != nil {
		log.Fatalf("CreateFigure error: %v", err)
	}
	log.Printf("CreateFigure response: %v", createdFigure)

	updatedFigure, err := client.UpdateFigure(ctx, &pb.UpdateFigureRequest{
		Id:         createdFigure.Id,
		Symbol:     2,
		PieceColor: "black",
		GameColor:  "black",
	})
	if err != nil {
		log.Fatalf("UpdateFigure error: %v", err)
	}
	log.Printf("UpdateFigure response: %v", updatedFigure)

	figuresList, err := client.ListFigures(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("ListFigures error: %v", err)
	}
	log.Printf("ListFigures response: %v", figuresList)

	figure, err := client.GetFigure(ctx, &pb.FigureById{Id: createdFigure.Id})
	if err != nil {
		log.Fatalf("GetFigure error: %v", err)
	}
	log.Printf("GetFigure response: %v", figure)

	deleteFigureResp, err := client.DeleteFigure(ctx, &pb.FigureById{Id: createdFigure.Id})
	if err != nil {
		log.Fatalf("DeleteFigure error: %v", err)
	}
	log.Printf("DeleteFigure response: %v", deleteFigureResp)

	// ===== Move methods =====

	createdMove, err := client.CreateMove(ctx, &pb.CreateMoveRequest{
		TimeTook:     1500,
		PositionFrom: "e2",
		PositionTo:   "e4",
		Figure:       1,
	})
	if err != nil {
		log.Fatalf("CreateMove error: %v", err)
	}
	log.Printf("CreateMove response: %v", createdMove)

	updatedMove, err := client.UpdateMove(ctx, &pb.UpdateMoveRequest{
		Id:           createdMove.Id,
		TimeTook:     2000,
		PositionFrom: "e7",
		PositionTo:   "e5",
		Figure:       2,
	})
	if err != nil {
		log.Fatalf("UpdateMove error: %v", err)
	}
	log.Printf("UpdateMove response: %v", updatedMove)

	movesList, err := client.ListMoves(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("ListMoves error: %v", err)
	}
	log.Printf("ListMoves response: %v", movesList)

	move, err := client.GetMove(ctx, &pb.MoveById{Id: createdMove.Id})
	if err != nil {
		log.Fatalf("GetMove error: %v", err)
	}
	log.Printf("GetMove response: %v", move)

	deleteMoveResp, err := client.DeleteMove(ctx, &pb.MoveById{Id: createdMove.Id})
	if err != nil {
		log.Fatalf("DeleteMove error: %v", err)
	}
	log.Printf("DeleteMove response: %v", deleteMoveResp)
}
