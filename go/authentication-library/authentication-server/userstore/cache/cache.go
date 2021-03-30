package cache

import (
	"log"

	"github.com/markstanden/authentication"
)

// initially based on Ben Johnson's user cache gist
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

func (c CachedStore) Find(key, value string) (*authentication.User, error) {

	var cache = make(map[string]*authentication.User)

	switch key {
	case "email":
		cache = c.cache.emailCache
	case "token":
		cache = c.cache.tokenCache
	}

	// Check the local cache first.
	if u, ok := cache[value]; ok {
		log.Printf(
			"authentication/cache: user (%d) read from %#v, current size: %v Users",
			u.UniqueID, cache, len(cache))

		// We have found a user in the cache
		// return early, no need to query main store.
		return u, nil
	}

	// User not found in the cache, so check in the wrapped service.
	u, err := c.store.Find(key, value)

	// If the user is not found return nil user, error
	if err != nil {
		return nil, err
	}

	// If the user is located - add to the correct cache
	if u != nil {
		cache[value] = u
	}

	// Return the found user
	return u, err

}

// Add passes the Add request to the wrapped store
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
