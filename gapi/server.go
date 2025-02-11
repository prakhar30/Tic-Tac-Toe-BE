package gapi

import (
	db "main/db/sqlc"
	"main/pb"
	"main/token"
	"main/utils"
)

type Server struct {
	pb.UnimplementedTicTacToeServer
	config     utils.Config
	store      db.Store
	tokenMaker token.Maker
}

func NewServer(config utils.Config, store db.Store, tokenMaker token.Maker) (*Server, error) {
	// tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	// if err != nil {
	// 	return nil, fmt.Errorf("cannot create token maker %w", err)
	// }

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
