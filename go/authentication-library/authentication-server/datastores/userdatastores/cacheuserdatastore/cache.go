package cache

import (
	"fmt"
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
	emailCache       map[string]*authentication.User
	tokenUserIDCache map[string]*authentication.User
}

/*
	UserCache wraps a UserService to provide an in-memory cache.
	The idea being to save database reads for basic	user authentication
*/
type UserServiceCache struct {
	cache userCache
	store authentication.UserDataStore
}

/*
	** LogStatus **
	takes a message and prefixes it to the start of a cache status log
	Helpful to ensure that the cache is acting as expected.
*/
func (usc UserServiceCache) LogStatus(message string) {
	log.Println(
		message+" - Current Cache Values:\nEmail Cache:\n\t",
		usc.cache.emailCache,
		"\nTokenUserID Cache:\n\t",
		usc.cache.tokenUserIDCache,
	)
}

/*
	** NewUserCache **
	returns an empty new read-through cache for the wrapped UserService.
*/
func NewUserCache(us authentication.UserDataStore) *UserServiceCache {
	uc := userCache{
		emailCache:       make(map[string]*authentication.User),
		tokenUserIDCache: make(map[string]*authentication.User),
	}
	return &UserServiceCache{
		cache: uc,
		store: us,
	}
}

/*
	** Find **
	looks up a user in the cache first,
	and if not present consults the wrapped UserService
*/
func (usc UserServiceCache) Find(key, value string) (*authentication.User, error) {

	var cache = make(map[string]*authentication.User)

	switch key {
	case "email":
		cache = usc.cache.emailCache
	case "tokenuserid":
		cache = usc.cache.tokenUserIDCache
	default:
		return nil, fmt.Errorf("authentication/cache:find - lookup key not found in switch statement")
	}

	// Check the local cache first.
	if u, ok := cache[value]; ok {
		log.Printf(
			"authentication/cache:find - user (%d) read from %#v, current size: %v Users",
			u.UniqueID, cache, len(cache))

		// We have found a user in the cache
		// return early, no need to query main store.
		return u, nil
	}

	// User not found in the cache, so check in the wrapped service.
	u, err := usc.store.Find(key, value)

	// If the user is not found return nil user, error
	if err != nil {
		return nil, err
	}

	// If the user is located - add to the correct cache
	if u != nil {
		cache[value] = u
	}

	// Return the found user
	return u, nil

}

/*
	** Add **
	passes the Add request to the wrapped store.
	Since the cache is only meant to speed up duplicate reads,
	and not replace the main UserService, this is a passthough method.
	This means the main UserService can perform any validation required
	(duplicate keys etc), and there is no need to duplicate here, or have
	potentially invald users stored in the cache.
*/
func (usc UserServiceCache) Add(u *authentication.User) (err error) {
	/*
		add the user to the main store as
		we only add to the cache on read.
	*/
	return usc.store.Add(u)

}

/*
	** Update **
	passes the Update request to the wrapped store
	It is now possible that the cache will hold out of date information
	so we will need to delete the entry from the cache(s).
*/
func (usc UserServiceCache) Update(u *authentication.User, fields authentication.User) (err error) {

	/*
		delete the user from all caches as
		information may now be out of date
	*/
	usc.deleteFromAll(u)

	/*
		update the user in the main store,
		return any errors directly
	*/
	return usc.store.Update(u, fields)
}

/*
	** deleteFromAll **
	is intended as a single function to delete the current
	user from the cache, useful for deleting and updating users.
*/
func (usc UserServiceCache) deleteFromAll(u *authentication.User) {

	/*
		Built in function delete only deletes the record if it exists,
		so no requirement for a comma ok.
	*/
	delete(usc.cache.tokenUserIDCache, u.TokenUserID)
	delete(usc.cache.emailCache, u.Email)
}

func (usc UserServiceCache) Delete(u *authentication.User) (err error) {
	usc.deleteFromAll(u)

	return usc.store.Delete(u)
}

/*
	** FullReset **
	resets the cache, and forwards the error (if any)
	created by the main datastore following a call to drop the existing user table and rebuild.
	This is intended for use in *development only*, to allow quick changes to
	the database structure / authentication.User struct object while still experimenting with the
	naming and number of required fields.
*/
func (usc UserServiceCache) FullReset() (err error) {
	/*
		Reset the current user cache, otherwise previous user entries will persist.
	*/
	usc.cache.emailCache = make(map[string]*authentication.User)
	usc.cache.tokenUserIDCache = make(map[string]*authentication.User)

	/*
		reset the main DataStore and return any errors.
	*/
	return usc.store.FullReset()
}
