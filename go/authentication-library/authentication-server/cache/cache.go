package cache

import (
	"fmt"

	"github.com/markstanden/authentication"
)

// Heavily based on Ben Johnson's user cache gist
// https://gist.github.com/benbjohnson/ffed98c3be896af58c5d74dd52cf0234#file-cache-go

// UserCache wraps a UserService to provide an in-memory cache.
type UserCache struct {
	emailCache map[string]*authentication.User
	tokenCache map[string]*authentication.User
	store      authentication.UserService
}

// NewUserCache returns a new read-through cache for service.
func NewUserCache(us authentication.UserService) *UserCache {
	return &UserCache{
		emailCache: make(map[string]*authentication.User),
		tokenCache: make(map[string]*authentication.User),
		store:      us,
	}
}

// FindByEmail returns a user for a given email.
// Returns the cached instance if available.
func (c *UserCache) FindByEmail(email string) (*authentication.User, error) {
	// Check the local cache first.
	if u := c.emailCache[email]; u != nil {
		fmt.Println("From Cache")
		return u, nil
	}

	// Otherwise fetch from the underlying service.
	u, err := c.store.FindByEmail(email)
	if err != nil {
		return nil, err
	} else if u != nil {
		c.emailCache[email] = u
	}
	return u, err
}

// FindByToken returns a user for a given token.
// Returns the cached instance if available.
func (c *UserCache) FindByToken(email string) (*authentication.User, error) {
	// Check the local cache first. ````
	if u := c.tokenCache[email]; u != nil {
		fmt.Println("From Cache")
		return u, nil
	}

	// Otherwise fetch from the underlying service.
	u, err := c.store.FindByToken(email)
	if err != nil {
		return nil, err
	} else if u != nil {
		c.tokenCache[email] = u
	}
	return u, err
}

// FindByToken returns a user for a given token.
// Returns the cached instance if available.
func (c *UserCache) Add(u authentication.User) (err error) {
	// add the user to the main store
	err = c.store.Add(u)
	if err != nil {
		return err
	}

	// may as well add to each cache too
	c.emailCache[u.Email] = &u
	c.tokenCache[u.Token] = &u

	// no errors
	return nil
}
