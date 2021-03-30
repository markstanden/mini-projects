package cache

import (
	"log"

	"github.com/markstanden/authentication"
)

/*
	The idea and initial implementation was based on Ben Johnson's user cache gist
	https://gist.github.com/benbjohnson/ffed98c3be896af58c5d74dd52cf0234#file-cache-go
*/

/*
	UserCache is the base struct for the cache,
*/
type userCache struct {
	emailCache map[string]*authentication.User
	tokenIDCache map[string]*authentication.User
}

// UserCache wraps a UserService to provide an in-memory cache.
type CachedStore struct {
	cache userCache
	store authentication.UserService
}

func (c CachedStore) LogStatus(message string) {
	log.Println(message + " - Current Cache Values:\nEmail Cache:\n\t", c.cache.emailCache, "\nTokenID Cache:\n\t", c.cache.tokenIDCache)
}

// NewUserCache returns a new read-through cache for service.
func NewUserCache(us authentication.UserService) *CachedStore {
	cache := userCache{
		emailCache: make(map[string]*authentication.User),
		tokenIDCache: make(map[string]*authentication.User),
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
	case "tokenid":
		cache = c.cache.tokenIDCache
	}

	c.LogStatus("cache/Find Called")

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
	c.LogStatus("cache/Add Called")
	/* 
		add the user to the main store as
		we only add to the cache on read.
	*/
	return c.store.Add(u)

}

/*
	Update passes the Update request to the wrapped store
	It is now possible that the cache will hold out of date information
	so we will need to delete the entry from the cache(s).
*/
func (c CachedStore) Update(u *authentication.User) (err error) {
	c.LogStatus("cache/Update Called")
	
	/* 	
		delete the user from all caches as
		information may now be out of date
	*/
	c.deleteFromAll(u)

	/* 
		update the user in the main store,
		return any errors directly
	*/
	c.LogStatus("cache/Update Complete, passing to datastore...")
	return c.store.Update(u)
}

/* 
	deleteFromAll is intended as a single function to delete the current
	user from the cache, useful for deleting and updating users.
*/
func (c CachedStore) deleteFromAll(u *authentication.User) {
	
	/*
		Built in function delete only deletes the record if it exists,
		so no requirement for a comma ok.
	*/
	delete(c.cache.tokenIDCache, u.TokenID)
	delete(c.cache.emailCache, u.Email)
}

func (c CachedStore) Delete(u *authentication.User) (err error) {
	c.LogStatus("cache/Delete Called")
	c.deleteFromAll(u)

	c.LogStatus("cache/Delete Complete passing to datastore...")
	return c.store.Delete(u)
}

/*
	FullReset resets the cache, and forwards the error (if any)
	created by the main datastore following a call to drop the existing user table and rebuild.
	This is intended for use in *development only*, to allow quick changes to
	the database structure / authentication.User struct object while still experimenting with the
	naming and number of required fields.
*/
func (c CachedStore) FullReset() (err error) {
	c.LogStatus("cache/FullReset Called")
	/*
		Reset the current user cache, otherwise previous user entries will persist.
	*/
	c.cache.emailCache = make(map[string]*authentication.User)
	c.cache.tokenIDCache = make(map[string]*authentication.User)

	/*
		reset the main DataStore and return any errors.
	*/
	c.LogStatus("cache/FullReset Complete, Passing to datastore")
	return c.store.FullReset()
}
