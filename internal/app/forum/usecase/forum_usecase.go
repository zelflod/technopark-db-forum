package forumUsecase

import (
	"github.com/nozimy/technopark-db-forum/internal/app/forum"
	"github.com/nozimy/technopark-db-forum/internal/app/thread"
	"github.com/nozimy/technopark-db-forum/internal/app/user"
	"github.com/nozimy/technopark-db-forum/internal/model"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
)

type ForumUsecase struct {
	forumRep  forum.Repository
	userRep   user.Repository
	threadRep thread.Repository
	cache     *cache.Cache
}

func (f ForumUsecase) GetThreadsByForum(forumSlug string, params map[string][]string) (model.Threads, int, error) {
	forumObj, err := f.Find(forumSlug)
	if forumObj == nil || err != nil {
		return nil, 404, err
	}

	threads, err := f.forumRep.FindForumThreads(forumSlug, params)
	if err != nil {
		return nil, 404, err
	}

	return threads, 200, nil
}

func (f ForumUsecase) GetUsersByForum(forumSlug string, params map[string][]string) (model.Users, int, error) {
	forumObj, err := f.Find(forumSlug)
	if forumObj == nil || err != nil {
		return nil, 404, err
	}

	users, err := f.forumRep.FindForumUsers(forumObj, params)
	if err != nil {
		return nil, 404, err
	}

	return users, 200, nil
}

func (f ForumUsecase) CreateThread(forumSlug string, newThread *model.NewThread) (*model.Thread, int, error) {
	userObj, err := f.userRep.FindByNickname(newThread.Author)
	if userObj == nil || err != nil {
		return nil, 404, err
	}

	forumObj, err := f.Find(forumSlug)
	if forumObj == nil || err != nil {
		return nil, 404, err
	}

	newThread.Forum = forumObj.Slug

	threadObj, err := f.threadRep.FindByIdOrSlug(0, newThread.Slug)
	if threadObj != nil {
		return threadObj, 409, err
	}

	threadObj, err = f.threadRep.CreateThread(newThread)
	if err != nil {
		return nil, 409, err
	}

	return threadObj, 201, nil
}

func (f ForumUsecase) CreateForum(data *model.Forum) (*model.Forum, int, error) {
	userObj, err := f.userRep.FindByNickname(data.User)
	if userObj == nil || err != nil {
		return nil, 404, err
	}

	data.User = userObj.Nickname

	if err := f.forumRep.Create(data); err != nil {
		forumObj, err := f.forumRep.Find(data.Slug)
		if err != nil {
			return nil, 409, err
		}

		return forumObj, 409, err
	}

	return data, 201, nil
}

func (f ForumUsecase) Find(slug string) (*model.Forum, error) {
	var forumObj *model.Forum
	var err error

	fromCache, found := f.cache.Get(slug)
	if !found {
		forumObj, err = f.forumRep.Find(slug)
		if err != nil {
			return nil, errors.Wrap(err, "forumRep.Find()")
		}
		f.cache.Set(slug, forumObj, cache.DefaultExpiration)
	} else {
		forumObj = fromCache.(*model.Forum)
	}

	return forumObj, nil
}

func NewForumUsecase(f forum.Repository, u user.Repository, t thread.Repository, c *cache.Cache) forum.Usecase {
	return &ForumUsecase{
		forumRep:  f,
		userRep:   u,
		threadRep: t,
		cache:     c,
	}
}
