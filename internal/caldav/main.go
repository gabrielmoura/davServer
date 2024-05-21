package caldav

import (
	"context"
	"fmt"
	"github.com/emersion/go-ical"
	"github.com/emersion/go-webdav/caldav"
	"github.com/gabrielmoura/davServer/config"
	"github.com/gabrielmoura/davServer/internal/data"
	"log"
	"path/filepath"
)

type Backend struct {
	calendars []caldav.Calendar
	objectMap map[string][]caldav.CalendarObject
}

func (t Backend) ListCalendars(ctx context.Context) ([]caldav.Calendar, error) {
	return t.calendars, nil
}

func (t Backend) GetCalendar(ctx context.Context, path string) (*caldav.Calendar, error) {
	for _, cal := range t.calendars {
		if cal.Path == path {
			return &cal, nil
		}
	}
	return nil, fmt.Errorf("Calendar for path: %s not found", path)
}

func (t Backend) CalendarHomeSetPath(ctx context.Context) (string, error) {
	user, ok := ctx.Value("user").(data.User)
	if !ok {
		return "", fmt.Errorf("Usuário não encontrado")
	}
	fullPath := filepath.Join(config.Conf.ShareRootDir, user.Username, "calendars")
	return fullPath, nil
}

func (t Backend) CurrentUserPrincipal(ctx context.Context) (string, error) {
	user, ok := ctx.Value("user").(data.User)
	if !ok {
		return "", fmt.Errorf("Usuário não encontrado")
	}
	fullPath := filepath.Join(config.Conf.ShareRootDir, user.Username, "/")
	return fullPath, nil
}

func (t Backend) DeleteCalendarObject(ctx context.Context, path string) error {
	return nil
}

func (t Backend) GetCalendarObject(ctx context.Context, path string, req *caldav.CalendarCompRequest) (*caldav.CalendarObject, error) {
	log.Printf("GetCalendarObject Path: %s", path)
	log.Printf("GetCalendarObject Request: %v", req)
	log.Printf("GetCalendarObject ObjectMap: %v", t.objectMap)
	log.Printf("GetCalendarObject Calendars: %v", t.calendars)
	log.Printf("GetCalendarObject User: %v", ctx.Value("user"))
	for _, objs := range t.objectMap {
		for _, obj := range objs {
			if obj.Path == path {
				return &obj, nil
			}
		}
	}

	return nil, fmt.Errorf("Couldn't find calendar object at: %s", path)
}

func (t Backend) PutCalendarObject(ctx context.Context, path string, calendar *ical.Calendar, opts *caldav.PutCalendarObjectOptions) (string, error) {
	return "", nil
}

func (t Backend) ListCalendarObjects(ctx context.Context, path string, req *caldav.CalendarCompRequest) ([]caldav.CalendarObject, error) {
	return t.objectMap[path], nil
}

func (t Backend) QueryCalendarObjects(ctx context.Context, query *caldav.CalendarQuery) ([]caldav.CalendarObject, error) {
	return nil, nil
}
