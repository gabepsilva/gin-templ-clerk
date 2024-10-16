package model

import "time"

type Event struct {
	ID                   uint64    `json:"id" gorm:"primaryKey"`                       // ID (uint): The unique identifier for the event, serves as the primary key.
	CreatedBy            string    `json:"createdBy" gorm:"createdBy;not null"`        // CreatedBy (string): The user ID of the event creator, linking to the User entity.
	User                 User      `json:"-" gorm:"foreignKey:CreatedBy" validate:"-"` // User (User): The user object associated with the event creator.
	Title                string    `json:"title" gorm:"not null"`                      // Title (string): The title of the event, required for easy identification.
	Description          string    `json:"description"`                                // Description (string): A brief explanation of what the event is about.
	Location             string    `json:"location"`                                   // Location (string): The physical or virtual location where the event will take place.
	Images               string    `json:"images" gorm:"type:json"`                    // Images ([]string): Array of image URLs associated with the event (e.g., event posters).
	StartTime            time.Time `json:"start_time" gorm:"default:null"`             // StartTime (time.Time): When the event is scheduled to begin.
	EndTime              time.Time `json:"end_time" gorm:"default:null"`               // EndTime (time.Time): When the event is scheduled to end.
	CreatedAt            time.Time `json:"created_at" gorm:"autoCreateTime"`           // CreatedAt (time.Time): The timestamp when the event was created, automatically set.
	UpdatedAt            time.Time `json:"updated_at" gorm:"autoUpdateTime"`           // UpdatedAt (time.Time): The timestamp when the event was last updated, automatically set.
	UpdatedBy            string    `json:"updated_by"`                                 // UpdatedBy (string): The user ID of the person who last updated the event.
	Status               string    `json:"status" gorm:"default:'draft'"`              // Status (string): The current state of the event (e.g., draft, published, cancelled).
	MaxAttendees         uint      `json:"max_attendees"`                              // MaxAttendees (uint): The maximum number of attendees allowed.
	AttendeesCount       uint      `json:"attendees_count"`                            // AttendeesCount (uint): The current number of registered attendees.
	IsPublic             bool      `json:"is_public" gorm:"default:true"`              // IsPublic (bool): Whether the event is public or private. Defaults to public.
	RSVPRequired         bool      `json:"rsvp_required" gorm:"default:false"`         // RSVPRequired (bool): Whether an RSVP is required to attend the event.
	Tags                 string    `json:"tags" gorm:"type:json"`                      // Tags ([]string): For categorizing events using tags like "conference", "workshop", etc.
	OrganizerContactInfo string    `json:"organizer_contact_info"`                     // OrganizerContactInfo (string): Contact details for the event organizer.
	ExternalLink         string    `json:"external_link"`                              // ExternalLink (string): Link to an external site related to the event (e.g., event registration page or official website).
	IsFeatured           bool      `json:"is_featured" gorm:"default:false"`           // IsFeatured (bool): Indicates whether this event is featured or highlighted on the platform.
	EventType            string    `json:"event_type"`                                 // EventType (string): The type or category of the event (e.g., webinar, in-person, hybrid).
}
