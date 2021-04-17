// 	usercache wraps a userstore returning recently looked up records, saving a db lookup
//	The idea and initial implementation was based on Ben Johnson's user cache gist
//	https://gist.github.com/benbjohnson/ffed98c3be896af58c5d74dd52cf0234#file-cache-go
package usercache

import (
	"fmt"
	"log"

	"github.com/markstanden/authentication"
)

//	UserCache is the base struct for the cache,
//	it wraps a UserStore to provide an in-memory cache.
//	The idea being to save database reads for basic	user authentication
type Cache struct {
	email  map[string]*authentication.User
	UserID map[string]*authentication.User
	store  authentication.UserDataStore
}

//	LogStatus takes a message and prefixes it to the start of a cache status log
//	Helpful to ensure that the cache is acting as expected.
func (c Cache) LogStatus(message string) {
	log.Println(
		message+" - Current Cache Values:\nEmail Cache:\n\t",
		c.email,
		"\nUserID Cache:\n\t",
		c.UserID,
	)
}

// 	New returns an empty new read-through cache for the wrapped UserService.
func New(us authentication.UserDataStore) *Cache {
	return &Cache{
		email:  make(map[string]*authentication.User),
		UserID: make(map[string]*authentication.User),
		store:  us,
	}
}

// 	Find looks up a user in the cache first,
// 	and if not present consults the wrapped UserService
func (c Cache) Find(key, value string) (*authentication.User, error) {

	var cache = make(map[string]*authentication.User)

	switch key {
	case "email":
		cache = c.email
	case "tokenuserid":
		cache = c.UserID
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
	return u, nil

}

// 	Add passes the Add request to the wrapped store.
//	Since the cache is only meant to speed up duplicate reads,
//	and not replace the main UserStore, this is a passthough method.
//	This means the main UserStore can perform any validation required
//	(duplicate keys etc), and there is no need to duplicate here, or have
//	potentially invald users stored in the cache.
func (c Cache) Add(u *authentication.User) (err error) {

	// 	add the user to the main store as we only add to the cache on read.
	return c.store.Add(u)

}

// 	Update passes the Update request to the wrapped store
//	It is now possible that the cache will hold out of date information
//	so we will need to delete the entry from the cache(s).
func (c Cache) Update(user *authentication.User, updatedUser authentication.User) (err error) {

	// 	delete the user from all caches
	// 	as information may now be out of date
	c.deleteFromAllCaches(user)

	//	update the user in the main store,
	//	return any errors directly
	return c.store.Update(user, updatedUser)
}

func (c Cache) UpdateRefreshToken(u *authentication.User, newRefreshToken string) (err error) {

	//	delete the user from all caches as
	//	information may now be out of date
	c.deleteFromAllCaches(u)

	//	update the user in the main store,
	//	return any errors directly
	return c.store.UpdateRefreshToken(u, newRefreshToken)
}

// 	deleteFromAllCaches is intended as a single function to delete the current
// 	user from the cache, useful for deleting and updating users.
// 	Will only delete caches that exist, and will not panic if a non existant
// 	cache is attempted to be deleted.
func (c Cache) deleteFromAllCaches(u *authentication.User) {

	// 	Built in function delete only deletes the record if it exists,
	// 	so no requirement for a comma ok.
	delete(c.UserID, u.TokenUserID)
	delete(c.email, u.Email)
}

// 	Delete removes a user from the caches and
// 	calls the delete function from the main store
func (c Cache) Delete(u *authentication.User) (err error) {
	c.deleteFromAllCaches(u)

	return c.store.Delete(u)
}

//	FullReset resets the cache, and forwards the error (if any)
//	created by the main datastore following a call to drop the existing user table and rebuild.
//	This is intended for use in *development only*, to allow quick changes to
//	the database structure / authentication.User struct object while still experimenting with the
//	naming and number of required fields.
func (c Cache) FullReset() (err error) {

	//	Reset the current user cache, otherwise previous user entries will persist.
	c.email = make(map[string]*authentication.User)
	c.UserID = make(map[string]*authentication.User)

	// reset the main DataStore and return any errors.
	return c.store.FullReset()
}
