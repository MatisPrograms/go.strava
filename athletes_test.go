package strava

import (
	"reflect"
	"testing"
	"time"
)

func TestAthletesGet(t *testing.T) {
	client := newCassetteClient(testToken, "athlete_get")
	athlete, err := NewAthletesService(client).Get(3545423).Do()

	if err != nil {
		t.Fatalf("service error: %v", err)
	}

	expected := &AthleteSummary{}

	expected.Id = 3545423
	expected.FirstName = "Strava"
	expected.LastName = "Testing"
	expected.Friend = "accepted"
	expected.Follower = "accepted"
	expected.Profile = "avatar/athlete/large.png"
	expected.ProfileMedium = "avatar/athlete/medium.png"
	expected.City = "Palo Alto"
	expected.State = "CA"
	expected.Country = "United States"
	expected.Gender = "M"
	expected.CreatedAtString = "2013-12-26T19:19:36Z"
	expected.UpdatedAtString = "2014-01-12T00:20:58Z"
	expected.CreatedAt, _ = time.Parse(timeFormat, expected.CreatedAtString)
	expected.UpdatedAt, _ = time.Parse(timeFormat, expected.UpdatedAtString)

	if !reflect.DeepEqual(athlete, expected) {
		t.Errorf("should match\n%v\n%v", athlete, expected)
	}

	// from here on out just check the request parameters
	s := NewAthletesService(newStoreRequestClient())

	// path
	s.Get(111).Do()

	transport := s.client.httpClient.Transport.(*storeRequestTransport)
	if transport.request.URL.Path != "/api/v3/athletes/111" {
		t.Errorf("request path incorrect, got %v", transport.request.URL.Path)
	}

	if transport.request.URL.RawQuery != "" {
		t.Errorf("request query incorrect, got %v", transport.request.URL.RawQuery)
	}
}

func TestAthletesListFriends(t *testing.T) {
	client := newCassetteClient(testToken, "athlete_list_friends")
	friends, err := NewAthletesService(client).ListFriends(3545423).Do()

	if err != nil {
		t.Fatalf("service error: %v", err)
	}

	if len(friends) == 0 {
		t.Fatal("friends not parsed")
	}

	if friends[0].CreatedAt.IsZero() || friends[0].UpdatedAt.IsZero() {
		t.Error("dates not parsed")
	}

	// from here on out just check the request parameters
	s := NewAthletesService(newStoreRequestClient())

	// parameters
	s.ListFriends(123).Page(2).PerPage(3).Do()

	transport := s.client.httpClient.Transport.(*storeRequestTransport)
	if transport.request.URL.Path != "/api/v3/athletes/123/friends" {
		t.Errorf("request path incorrect, got %v", transport.request.URL.Path)
	}

	if transport.request.URL.RawQuery != "page=2&per_page=3" {
		t.Errorf("request query incorrect, got %v", transport.request.URL.RawQuery)
	}
}

func TestAthletesListFollowers(t *testing.T) {
	client := newCassetteClient(testToken, "athlete_list_followers")
	followers, err := NewAthletesService(client).ListFollowers(3545423).Do()

	if err != nil {
		t.Fatalf("service error: %v", err)
	}

	if len(followers) == 0 {
		t.Fatal("followers not parsed")
	}

	if followers[0].CreatedAt.IsZero() || followers[0].UpdatedAt.IsZero() {
		t.Error("dates not parsed")
	}

	// from here on out just check the request parameters
	s := NewAthletesService(newStoreRequestClient())

	// parameters
	s.ListFollowers(123).Page(2).PerPage(3).Do()

	transport := s.client.httpClient.Transport.(*storeRequestTransport)
	if transport.request.URL.Path != "/api/v3/athletes/123/followers" {
		t.Errorf("request path incorrect, got %v", transport.request.URL.Path)
	}

	if transport.request.URL.RawQuery != "page=2&per_page=3" {
		t.Errorf("request query incorrect, got %v", transport.request.URL.RawQuery)
	}
}

func TestAthletesListBothFollowing(t *testing.T) {
	client := newCassetteClient(testToken, "athlete_list_both_following")
	followers, err := NewAthletesService(client).ListBothFollowing(3545423).Do()

	if err != nil {
		t.Fatalf("service error: %v", err)
	}

	if len(followers) == 0 {
		t.Fatal("followers not parsed")
	}

	if followers[0].CreatedAt.IsZero() || followers[0].UpdatedAt.IsZero() {
		t.Error("dates not parsed")
	}

	// from here on out just check the request parameters
	s := NewAthletesService(newStoreRequestClient())

	// parameters
	s.ListBothFollowing(123).PerPage(7).Page(8).Do()

	transport := s.client.httpClient.Transport.(*storeRequestTransport)
	if transport.request.URL.Path != "/api/v3/athletes/123/both-following" {
		t.Errorf("request path incorrect, got %v", transport.request.URL.Path)
	}

	if transport.request.URL.RawQuery != "page=8&per_page=7" {
		t.Errorf("request query incorrect, got %v", transport.request.URL.RawQuery)
	}
}

func TestAthletesListKOMs(t *testing.T) {
	client := newCassetteClient(testToken, "athlete_list_koms")
	efforts, err := NewAthletesService(client).ListKOMs(3776).Do()

	if err != nil {
		t.Fatalf("service error: %v", err)
	}

	if len(efforts) == 0 {
		t.Fatal("efforts not parsed")
	}

	if efforts[0].StartDate.IsZero() || efforts[0].StartDateLocal.IsZero() {
		t.Error("dates not parsed")
	}

	// from here on out just check the request parameters
	s := NewAthletesService(newStoreRequestClient())

	// parameters
	s.ListKOMs(123).PerPage(9).Page(8).Do()

	transport := s.client.httpClient.Transport.(*storeRequestTransport)
	if transport.request.URL.Path != "/api/v3/athletes/123/koms" {
		t.Errorf("request path incorrect, got %v", transport.request.URL.Path)
	}

	if transport.request.URL.RawQuery != "page=8&per_page=9" {
		t.Errorf("request query incorrect, got %v", transport.request.URL.RawQuery)
	}
}

func TestAthletesListActivities(t *testing.T) {
	client := newCassetteClient(testToken, "athlete_list_activies")
	activities, err := NewAthletesService(client).ListActivities(14507).Do()

	if err != nil {
		t.Fatalf("service error: %v", err)
	}

	if len(activities) == 0 {
		t.Fatal("efforts not parsed")
	}

	if activities[0].StartDate.IsZero() || activities[0].StartDateLocal.IsZero() {
		t.Error("dates not parsed")
	}

	// from here on out just check the request parameters
	s := NewAthletesService(newStoreRequestClient())

	// path
	s.ListActivities(123).Do()

	transport := s.client.httpClient.Transport.(*storeRequestTransport)
	if transport.request.URL.Path != "/api/v3/athletes/123/activities" {
		t.Errorf("request path incorrect, got %v", transport.request.URL.Path)
	}

	if transport.request.URL.RawQuery != "" {
		t.Errorf("request query incorrect, got %v", transport.request.URL.RawQuery)
	}

	// parameters
	s.ListActivities(123).PerPage(9).Page(8).Do()

	transport = s.client.httpClient.Transport.(*storeRequestTransport)
	if transport.request.URL.RawQuery != "page=8&per_page=9" {
		t.Errorf("request query incorrect, got %v", transport.request.URL.RawQuery)
	}

	// parameters2
	s.ListActivities(123).Before(1391020072).Do()

	transport = s.client.httpClient.Transport.(*storeRequestTransport)

	if transport.request.URL.RawQuery != "before=1391020072" {
		t.Errorf("request query incorrect, got %v", transport.request.URL.RawQuery)
	}

	// parameters3
	s.ListActivities(123).After(0).Do()

	transport = s.client.httpClient.Transport.(*storeRequestTransport)
	if transport.request.URL.RawQuery != "after=0" {
		t.Errorf("request query incorrect, got %v", transport.request.URL.RawQuery)
	}
}

func TestAthletesBadJSON(t *testing.T) {
	var err error
	s := NewAthletesService(NewStubResponseClient("bad json"))

	_, err = s.Get(123).Do()
	if err == nil {
		t.Error("should return a bad json error")
	}

	_, err = s.ListFriends(123).Do()
	if err == nil {
		t.Error("should return a bad json error")
	}

	_, err = s.ListFollowers(123).Do()
	if err == nil {
		t.Error("should return a bad json error")
	}

	_, err = s.ListBothFollowing(123).Do()
	if err == nil {
		t.Error("should return a bad json error")
	}

	_, err = s.ListKOMs(123).Do()
	if err == nil {
		t.Error("should return a bad json error")
	}

	_, err = s.ListActivities(123).Do()
	if err == nil {
		t.Error("should return a bad json error")
	}
}
