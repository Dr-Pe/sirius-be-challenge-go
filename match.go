package main

import "google.golang.org/protobuf/types/known/timestamppb"

type Match struct {
	id           int
	player1_id   int
	player2_id   int
	start_time   timestamppb.Timestamp
	end_time     timestamppb.Timestamp
	winner_id    int
	table_number int
}
