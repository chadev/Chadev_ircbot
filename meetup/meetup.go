package meetup

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const baseURL = "https://api.meetup.com/2/events?&sign=true&photo-host=secure&group_urlname=%s&page=20&key=%s"

// Meetup contains the ruturn value from the Meetup API.
type Meetup struct {
	Results []MeetupEvents `json:"results"`
}

// MeetupEvents contains details for each event in the return value.
type MeetupEvents struct {
	Venue    MeetupVenue `json:"venue"`
	EventURL string      `json:"event_url"`
	Name     string      `json:"name"`
	Time     int64       `json:"time"`
}

// MeetupVenue contains details about the venue for each event.
type MeetupVenue struct {
	Name string `json:"name"`
}

func (e *MeetupEvents) parseDateTime() time.Time {
	// set the timezone, otherwise UTC will be used
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		// couldn't set timezone, falling back to UTC
		loc, _ = time.LoadLocation("UTC")
	}

	// parse the unix timestamp, and apply our timezone
	t := time.Unix(0, e.Time*int64(time.Millisecond)).In(loc)

	return t
}

func formatDateTime(dt time.Time) string {
	today := isToday(dt)

	var fTime string
	if today {
		fTime = dt.Format("3:04 pm")
	} else {
		fTime = dt.Format("01/02 3:04 pm")
	}

	return fTime
}

func isToday(dt time.Time) bool {
	now := time.Now().Format("2006-01-02")
	event := dt.Format("2006-01-02")

	if event != now {
		return false
	}

	return true
}

func (m *Meetup) string() string {
	dt := m.Results[0].parseDateTime()
	fTime := formatDateTime(dt)
	lunchDay := isToday(dt)

	if !lunchDay {
		return fmt.Sprintf("The next meetup is \"%s\", you can join us at %s on %s.  If you plan to come please make sure you have RSVPed at %s",
			m.Results[0].Name,
			m.Results[0].Venue.Name,
			fTime,
			m.Results[0].EventURL)
	}

	return fmt.Sprintf("The meetup today is \"%s\", you can join us at %s on %s.  If you plan to come please make sure you have RSVPed at %s",
		m.Results[0].Name,
		m.Results[0].Venue.Name,
		fTime,
		m.Results[0].EventURL)
}

func GetNextMeetup(group string) (string, error) {
	URL := fmt.Sprintf(baseURL, group,
		os.Getenv("CHADEV_MEETUP"))
	resp, err := http.Get(URL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var Events Meetup
	err = json.Unmarshal(body, &Events)
	if err != nil {
		return "", err
	}

	return Events.string(), nil
}
