package immich

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/jaybinks/immich-go/helpers/tzone"
)

type UnsupportedMedia struct {
	msg string
}

func (u UnsupportedMedia) Error() string {
	return u.msg
}

func (u UnsupportedMedia) Is(target error) bool {
	_, ok := target.(*UnsupportedMedia)
	return ok
}

type PingResponse struct {
	Res string `json:"res"`
}

type User struct {
	ID                   string    `json:"id"`
	Email                string    `json:"email"`
	FirstName            string    `json:"firstName"`
	LastName             string    `json:"lastName"`
	StorageLabel         string    `json:"storageLabel"`
	ExternalPath         string    `json:"externalPath"`
	ProfileImagePath     string    `json:"profileImagePath"`
	ShouldChangePassword bool      `json:"shouldChangePassword"`
	IsAdmin              bool      `json:"isAdmin"`
	CreatedAt            time.Time `json:"createdAt"`
	DeletedAt            time.Time `json:"deletedAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
	OauthID              string    `json:"oauthId"`
}

type List[T comparable] struct {
	list []T
	lock sync.RWMutex
}

func (l *List[T]) Includes(v T) bool {
	l.lock.RLock()
	defer l.lock.RUnlock()
	for i := range l.list {
		if l.list[i] == v {
			return true
		}
	}
	return false
}

func (l *List[T]) Push(v T) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.list = append(l.list, v)
}

func (l *List[T]) MarshalJSON() ([]byte, error) {
	return nil, errors.New("MarshalJSON not implemented for List[T]")
}

func (l *List[T]) UnmarshalJSON(data []byte) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.list == nil {
		l.list = []T{}
	}
	return json.Unmarshal(data, &l.list)
}

type myBool bool

func (b myBool) String() string {
	if b {
		return "true"
	}
	return "false"
}

// immich Asset simplified
type Asset struct {
	ID               string            `json:"id"`
	DeviceAssetID    string            `json:"deviceAssetId"`
	OwnerID          string            `json:"ownerId"`
	DeviceID         string            `json:"deviceId"`
	Type             string            `json:"type"`
	OriginalPath     string            `json:"originalPath"`
	OriginalFileName string            `json:"originalFileName"`
	Resized          bool              `json:"resized"`
	Thumbhash        string            `json:"thumbhash"`
	FileCreatedAt    ImmichTime        `json:"fileCreatedAt"`
	FileModifiedAt   ImmichTime        `json:"fileModifiedAt"`
	UpdatedAt        ImmichTime        `json:"updatedAt"`
	IsFavorite       bool              `json:"isFavorite"`
	IsArchived       bool              `json:"isArchived"`
	IsTrashed        bool              `json:"isTrashed"`
	Duration         string            `json:"duration"`
	ExifInfo         ExifInfo          `json:"exifInfo"`
	LivePhotoVideoID any               `json:"livePhotoVideoId"`
	Tags             []any             `json:"tags"`
	Checksum         string            `json:"checksum"`
	StackParentId    string            `json:"stackParentId"`
	JustUploaded     bool              `json:"-"`
	Albums           []AlbumSimplified `json:"-"` // Albums that asset belong to
}

type ExifInfo struct {
	Make             string     `json:"make"`
	Model            string     `json:"model"`
	ExifImageWidth   int        `json:"exifImageWidth"`
	ExifImageHeight  int        `json:"exifImageHeight"`
	FileSizeInByte   int        `json:"fileSizeInByte"`
	Orientation      string     `json:"orientation"`
	DateTimeOriginal ImmichTime `json:"dateTimeOriginal,omitempty"`
	// 	ModifyDate       time.Time `json:"modifyDate"`
	TimeZone string `json:"timeZone"`
	// LensModel        string    `json:"lensModel"`
	// 	FNumber          float64   `json:"fNumber"`
	// 	FocalLength      float64   `json:"focalLength"`
	// 	Iso              int       `json:"iso"`
	// 	ExposureTime     string    `json:"exposureTime"`
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
	// 	City             string    `json:"city"`
	// 	State            string    `json:"state"`
	// 	Country          string    `json:"country"`
	Description string `json:"description"`
}

type ImmichTime struct {
	time.Time
}

// ImmichTime.UnmarshalJSON read time from the JSON string.
// The json provides a time UTC, but the server and the images dates are given in local time.
// The get the correct time into the struct, we capture the UTC time and return it in the local zone.
//
// workaround for: error at connection to immich server: cannot parse "+174510-04-28T00:49:44.000Z" as "2006" #28
// capture the error

func (t *ImmichTime) UnmarshalJSON(b []byte) error {
	local, err := tzone.Local()
	if err != nil {
		return err
	}
	var ts time.Time
	if len(b) < 3 {
		t.Time = time.Time{}
		return nil
	}
	b = b[1 : len(b)-1]
	ts, err = time.ParseInLocation("2006-01-02T15:04:05.000Z", string(b), time.UTC)
	if err != nil {
		t.Time = time.Time{}
		return nil
	}
	t.Time = ts.In(local)
	return nil
}
