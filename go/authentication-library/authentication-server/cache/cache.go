package cache

import (
	"log"

	"github.com/markstanden/authentication"
)

// Heavily based on Ben Johnson's user cache gist
// https://gist.github.com/benbjohnson/ffed98c3be896af58c5d74dd52cf0234#file-cache-go

// UserCache
type userCache struct {
	emailCache map[string]*authentication.User
	tokenCache map[string]*authentication.User
}

// UserCache wraps a UserService to provide an in-memory cache.
type CachedStore struct {
	cache userCache
	store authentication.UserService
}

// NewUserCache returns a new read-through cache for service.
func NewUserCache(us authentication.UserService) *CachedStore {
	cache := userCache{
		emailCache: make(map[string]*authentication.User),
		tokenCache: make(map[string]*authentication.User),
	}
	return &CachedStore{
		cache: cache,
		store: us,
	}
}

// FindByEmail returns a user for a given email.
// Returns the cached instance if available.
func (c CachedStore) FindByEmail(email string) (*authentication.User, error) {
	// Check the local cache first.
	if u := c.cache.emailCache[email]; u != nil {
		log.Printf("authentication/cache: user (%d) read from emailCache, current size: %v Users", u.UniqueID, len(c.cache.emailCache))
		return u, nil
	}

	// Otherwise fetch from the underlying service.
	u, err := c.store.FindByEmail(email)
	if err != nil {
		return nil, err
	} else if u != nil {
		c.cache.emailCache[email] = u
	}
	return u, err
}

// FindByToken returns a user for a given token.
// Returns the cached instance if available.
func (c CachedStore) FindByToken(email string) (*authentication.User, error) {
	// Check the local cache first. ````
	if u := c.cache.tokenCache[email]; u != nil {
		log.Printf("authentication/cache: user (%d) read from tokenCache, current size: %v Users", u.UniqueID, len(c.cache.tokenCache))
		return u, nil
	}

	// Otherwise fetch from the underlying service.
	u, err := c.store.FindByToken(email)
	if err != nil {
		return nil, err
	} else if u != nil {
		c.cache.tokenCache[email] = u
	}
	return u, err
}

// FindByToken returns a user for a given token.
// Returns the cached instance if available.
func (c CachedStore) Add(u *authentication.User) (err error) {
	// add the user to the main store
	err = c.store.Add(u)
	if err != nil {
		return err
	}

	// no errors
	return nil
}

// Create forwards the error (if any) created by the main datastore
// following a call to drop the existing user table and rebuild.
// obviously this is for use in development only, to allow quick changes to
// the database structure / authentication.User struct object.
func (c CachedStore) FullReset() (err error) {

	// Reset the current user cache, otherwise previous user entries will persist
	// this won't be an issue for the current use case of the function, but may
	// prevent future bugs
	c.cache.emailCache = make(map[string]*authentication.User)
	c.cache.tokenCache = make(map[string]*authentication.User)

	// reset the main store and return any errors.
	return c.store.FullReset()
}
