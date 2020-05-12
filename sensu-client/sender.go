package main

import (
	"context"
	"fmt"

	"pack.ag/amqp"
)

type sender struct {
	s *amqp.Sender
}

func newSender(address, target string) (*sender, error) {
	client, err := amqp.Dial(address)
	if err != nil {
		return nil, fmt.Errorf("dialing AMQP server: %s", err)
	}

	session, err := client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("creating AMQP session: %s", err)
	}

	s, err := session.NewSender(amqp.LinkTargetAddress(target))
	if err != nil {
		return nil, fmt.Errorf("creating sender link: %s", err)
	}

	return &sender{
		s: s,
	}, nil
}

func (s *sender) send(ctx context.Context, message []byte) error {
	err := s.s.Send(ctx, amqp.NewMessage(message))
	if err != nil {
		return fmt.Errorf("sending AMQP message: %s", err)
	}
	return nil
}
