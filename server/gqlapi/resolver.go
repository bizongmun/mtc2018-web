//go:generate gqlgen

package gqlapi

import (
	"sync"

	"cloud.google.com/go/spanner"
	"github.com/mercari/mtc2018-web/server/domains"
	"go.uber.org/zap"
)

// NewResolver returns GraphQL root resolver.
func NewResolver(logger *zap.Logger, spannerClient *spanner.Client) (ResolverRoot, error) {

	sessionRepo, err := domains.NewSessionRepo()
	if err != nil {
		return nil, err
	}
	speakerRepo, err := domains.NewSpeakerRepo()
	if err != nil {
		return nil, err
	}

	newsRepo, err := domains.NewNewsRepo()
	if err != nil {
		return nil, err
	}

	var (
		likeRepo    domains.LikeRepo
		likeSumRepo domains.LikeSummaryRepo
	)
	if spannerClient != nil {
		var err error
		likeRepo, err = domains.NewLikeRepo(spannerClient)
		if err != nil {
			return nil, err
		}
		likeSumRepo, err = domains.NewLikeSummaryRepo(spannerClient)
		if err != nil {
			return nil, err
		}
	} else {
		var err error
		likeRepo, err = domains.NewFakeLikeRepo()
		if err != nil {
			return nil, err
		}
		likeSumRepo, err = domains.NewFakeLikeSummaryRepo()
		if err != nil {
			return nil, err
		}
	}

	storer := newStorer(logger, likeRepo, likeSumRepo)
	go storer.Run()

	eventCh := make(chan likeEvent, 128)

	observer := newObserver(logger, eventCh, likeSumRepo)
	go observer.Run()

	listener := newListener(logger, eventCh, sessionRepo)
	go listener.Run()

	r := &rootResolver{
		sessionRepo: sessionRepo,
		speakerRepo: speakerRepo,
		likeRepo:    likeRepo,
		newsRepo:    newsRepo,
		likeSumRepo: likeSumRepo,

		Logger:   logger,
		storer:   storer,
		observer: observer,
		listener: listener,
	}

	return r, nil
}

type rootResolver struct {
	sessionRepo domains.SessionRepo
	speakerRepo domains.SpeakerRepo
	likeRepo    domains.LikeRepo
	newsRepo    domains.NewsRepo
	likeSumRepo domains.LikeSummaryRepo

	Logger *zap.Logger
	mu     sync.Mutex

	storer   *storer
	observer *observer
	listener *listener
}

func (r *rootResolver) Query() QueryResolver {
	return &queryResolver{r}
}

func (r *rootResolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

func (r *rootResolver) Subscription() SubscriptionResolver {
	return &subscriptionResolver{r}
}

func (r *rootResolver) Session() SessionResolver {
	return &sessionResolver{r}
}

func (r *rootResolver) Speaker() SpeakerResolver {
	return &speakerResolver{r}
}

func (r *rootResolver) Like() LikeResolver {
	return &likeResolver{r}
}

func (r *rootResolver) News() NewsResolver {
	return &newsResolver{r}
}
