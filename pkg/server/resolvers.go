package server

import (
	"context"

	"github.com/ryicoh/apery-graphql/pkg"
)

type (
	Resolvers struct {
	}
	queryResolver struct {
		resolvers *Resolvers
	}
	mutationResolver struct {
		resolvers *Resolvers
	}
)

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
	panic("not implemented") // TODO: Implement
}
