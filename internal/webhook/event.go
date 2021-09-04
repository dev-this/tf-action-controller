package webhook

type Event struct {
	sha   string
	ref   string
	owner string
	repo  string
}

func NewEvent(sha, ref, owner, repo string) Event {
	return Event{
		sha:   sha,
		ref:   ref,
		owner: owner,
		repo:  repo,
	}
}
