package server

import (
	"context"
	"errors"
	"time"

	"github.com/ryicoh/apery-graphql/pkg"
	"github.com/ryicoh/apery-graphql/pkg/apery"
)

type (
	Resolvers struct {
		aperyClient apery.AperyClient
	}
	queryResolver struct {
		resolvers *Resolvers
	}
	mutationResolver struct {
		resolvers *Resolvers
	}
)

func NewResolvers(cli apery.AperyClient) pkg.ResolverRoot {
	return &Resolvers{cli}
}

func (r *Resolvers) Mutation() pkg.MutationResolver {
	return &mutationResolver{r}
}

func (m *mutationResolver) Foo(ctx context.Context) (bool, error) {
	panic("not implemented") // TODO: Implement
}

func (r *Resolvers) Query() pkg.QueryResolver {
	return &queryResolver{r}
}

func (q *queryResolver) Evaluate(ctx context.Context, input pkg.EvaluateInput) (*pkg.EvaluateOutput, error) {
	if input.TimeoutSecond > 30 {
		return nil, errors.New("timeoutSecond は30秒より少なくしてください")
	}

	result, err := q.resolvers.aperyClient.Evaluate(
		ctx, input.Sfen, input.Moves,
		time.Duration(input.TimeoutSecond)*time.Second)
	if err != nil {
		return nil, err
	}

	return &pkg.EvaluateOutput{
		Value:    result.Value,
		Bestmove: result.Bestmove,
		Pv:       result.Pv,
		Nodes:    result.Nodes,
		Depth:    result.Depth,
	}, nil
}
