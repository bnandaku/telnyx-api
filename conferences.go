package telnyx

import (
	"context"
	"fmt"
)

// Conference is the conference object returned by the API.
type Conference struct {
	RecordType   string `json:"record_type"`
	ID           string `json:"id"`
	Name         string `json:"name"`
	CreatedAt    string `json:"created_at"`
	ExpiresAt    string `json:"expires_at"`
	UpdatedAt    string `json:"updated_at"`
	Region       string `json:"region"`
	Status       string `json:"status"`    // init | in_progress | completed
	EndReason    string `json:"end_reason"` // all_left | ended_via_api | host_left | time_exceeded
	ConnectionID string `json:"connection_id"`
}

type conferenceResponse struct {
	Data Conference `json:"data"`
}

type conferenceListResponse struct {
	Data []Conference   `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

// ConferenceCommandResponse is returned by most conference action endpoints.
type ConferenceCommandResponse struct {
	Data struct {
		Result string `json:"result"`
	} `json:"data"`
}

// CreateConferenceRequest creates a new conference via a call leg.
type CreateConferenceRequest struct {
	CallControlID          string `json:"call_control_id"`
	Name                   string `json:"name"`
	BeepEnabled            string `json:"beep_enabled,omitempty"` // always | never | on_enter | on_exit
	ClientState            string `json:"client_state,omitempty"`
	CommandID              string `json:"command_id,omitempty"`
	ComfortNoise           *bool  `json:"comfort_noise,omitempty"`
	DurationMinutes        int    `json:"duration_minutes,omitempty"`
	HoldAudioURL           string `json:"hold_audio_url,omitempty"`
	HoldMediaName          string `json:"hold_media_name,omitempty"`
	MaxParticipants        int    `json:"max_participants,omitempty"`
	StartConferenceOnCreate *bool `json:"start_conference_on_create,omitempty"`
	Region                 string `json:"region,omitempty"` // Australia | Europe | Middle East | US
}

// CreateConference creates a new conference room.
func (c *Client) CreateConference(ctx context.Context, req CreateConferenceRequest) (*Conference, error) {
	var resp conferenceResponse
	if err := c.post(ctx, "/conferences", req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// GetConference retrieves a conference by ID.
func (c *Client) GetConference(ctx context.Context, id string) (*Conference, error) {
	var resp conferenceResponse
	if err := c.get(ctx, fmt.Sprintf("/conferences/%s", id), &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// ListConferences returns all conferences with optional region filter.
func (c *Client) ListConferences(ctx context.Context, pageNumber, pageSize int, region string) ([]Conference, *PaginationMeta, error) {
	path := fmt.Sprintf("/conferences?page[number]=%d&page[size]=%d", pageNumber, pageSize)
	if region != "" {
		path += fmt.Sprintf("&filter[region]=%s", region)
	}
	var resp conferenceListResponse
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, nil, err
	}
	return resp.Data, &resp.Meta, nil
}

// ConferenceParticipant holds participant state within a conference.
type ConferenceParticipantDetail struct {
	RecordType              string     `json:"record_type"`
	ID                      string     `json:"id"`
	CallLegID               string     `json:"call_leg_id"`
	CallControlID           string     `json:"call_control_id"`
	Conference              Conference `json:"conference"`
	WhisperCallControlIDs   []string   `json:"whisper_call_control_ids"`
	CreatedAt               string     `json:"created_at"`
	UpdatedAt               string     `json:"updated_at"`
	EndConferenceOnExit     bool       `json:"end_conference_on_exit"`
	SoftEndConferenceOnExit bool       `json:"soft_end_conference_on_exit"`
	Status                  string     `json:"status"` // joining | joined | left
	Muted                   bool       `json:"muted"`
	OnHold                  bool       `json:"on_hold"`
}

type conferenceParticipantListResponse struct {
	Data []ConferenceParticipantDetail `json:"data"`
	Meta PaginationMeta                `json:"meta"`
}

// ListConferenceParticipants lists participants in a conference.
func (c *Client) ListConferenceParticipants(ctx context.Context, conferenceID string, pageNumber, pageSize int) ([]ConferenceParticipantDetail, *PaginationMeta, error) {
	path := fmt.Sprintf("/conferences/%s/participants?page[number]=%d&page[size]=%d", conferenceID, pageNumber, pageSize)
	var resp conferenceParticipantListResponse
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, nil, err
	}
	return resp.Data, &resp.Meta, nil
}

func (c *Client) conferenceAction(ctx context.Context, conferenceID, action string, body any) (*ConferenceCommandResponse, error) {
	var resp ConferenceCommandResponse
	path := fmt.Sprintf("/conferences/%s/actions/%s", conferenceID, action)
	if err := c.post(ctx, path, body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// JoinConferenceRequest adds a call leg to a conference.
type JoinConferenceRequest struct {
	CallControlID           string   `json:"call_control_id"`
	ClientState             string   `json:"client_state,omitempty"`
	CommandID               string   `json:"command_id,omitempty"`
	EndConferenceOnExit     bool     `json:"end_conference_on_exit,omitempty"`
	SoftEndConferenceOnExit bool     `json:"soft_end_conference_on_exit,omitempty"`
	Hold                    bool     `json:"hold,omitempty"`
	HoldAudioURL            string   `json:"hold_audio_url,omitempty"`
	HoldMediaName           string   `json:"hold_media_name,omitempty"`
	Mute                    bool     `json:"mute,omitempty"`
	StartConferenceonEnter  bool     `json:"start_conference_on_enter,omitempty"`
	SupervisorRole          string   `json:"supervisor_role,omitempty"` // barge | whisper | monitor
	WhisperCallControlIDs   []string `json:"whisper_call_control_ids,omitempty"`
	BeepEnabled             string   `json:"beep_enabled,omitempty"`
}

// JoinConference adds a call leg to an existing conference.
func (c *Client) JoinConference(ctx context.Context, conferenceID string, req JoinConferenceRequest) (*ConferenceCommandResponse, error) {
	return c.conferenceAction(ctx, conferenceID, "join", req)
}

// MuteConferenceParticipantsRequest mutes one or more participants.
type MuteConferenceParticipantsRequest struct {
	CallControlIDs []string `json:"call_control_ids,omitempty"`
}

// MuteConferenceParticipants mutes participants; empty list mutes all.
func (c *Client) MuteConferenceParticipants(ctx context.Context, conferenceID string, callControlIDs []string) (*ConferenceCommandResponse, error) {
	return c.conferenceAction(ctx, conferenceID, "mute", MuteConferenceParticipantsRequest{CallControlIDs: callControlIDs})
}

// UnmuteConferenceParticipants unmutes participants; empty list unmutes all.
func (c *Client) UnmuteConferenceParticipants(ctx context.Context, conferenceID string, callControlIDs []string) (*ConferenceCommandResponse, error) {
	return c.conferenceAction(ctx, conferenceID, "unmute", MuteConferenceParticipantsRequest{CallControlIDs: callControlIDs})
}

// HoldConferenceParticipantsRequest places participants on hold.
type HoldConferenceParticipantsRequest struct {
	CallControlIDs []string `json:"call_control_ids,omitempty"`
	AudioURL       string   `json:"audio_url,omitempty"`
	MediaName      string   `json:"media_name,omitempty"`
	Region         string   `json:"region,omitempty"`
}

// HoldConferenceParticipants places participants on hold; empty list holds all.
func (c *Client) HoldConferenceParticipants(ctx context.Context, conferenceID string, req HoldConferenceParticipantsRequest) (*ConferenceCommandResponse, error) {
	return c.conferenceAction(ctx, conferenceID, "hold", req)
}

// UnholdConferenceParticipants removes participants from hold; empty list unholds all.
func (c *Client) UnholdConferenceParticipants(ctx context.Context, conferenceID string, callControlIDs []string) (*ConferenceCommandResponse, error) {
	return c.conferenceAction(ctx, conferenceID, "unhold", map[string]any{
		"call_control_ids": callControlIDs,
	})
}

// ConferenceSpeakRequest sends TTS to conference participants.
type ConferenceSpeakRequest struct {
	Payload        string   `json:"payload"`
	Voice          string   `json:"voice"`
	Language       string   `json:"language,omitempty"`
	PayloadType    string   `json:"payload_type,omitempty"`
	ServiceLevel   string   `json:"service_level,omitempty"`
	CallControlIDs []string `json:"call_control_ids,omitempty"`
	CommandID      string   `json:"command_id,omitempty"`
}

// ConferenceSpeak plays TTS to conference participants.
func (c *Client) ConferenceSpeak(ctx context.Context, conferenceID string, req ConferenceSpeakRequest) (*ConferenceCommandResponse, error) {
	return c.conferenceAction(ctx, conferenceID, "speak", req)
}

// ConferencePlayAudioRequest plays audio to conference participants.
type ConferencePlayAudioRequest struct {
	AudioURL       string   `json:"audio_url,omitempty"`
	MediaName      string   `json:"media_name,omitempty"`
	Loop           any      `json:"loop,omitempty"`
	CallControlIDs []string `json:"call_control_ids,omitempty"`
	CommandID      string   `json:"command_id,omitempty"`
}

// ConferencePlayAudio plays an audio file to conference participants.
func (c *Client) ConferencePlayAudio(ctx context.Context, conferenceID string, req ConferencePlayAudioRequest) (*ConferenceCommandResponse, error) {
	return c.conferenceAction(ctx, conferenceID, "play_audio", req)
}

// ConferenceStopAudio stops audio playback in the conference.
func (c *Client) ConferenceStopAudio(ctx context.Context, conferenceID string, callControlIDs []string) (*ConferenceCommandResponse, error) {
	return c.conferenceAction(ctx, conferenceID, "stop_audio", map[string]any{
		"call_control_ids": callControlIDs,
	})
}

// ConferenceRecordStartRequest starts recording a conference.
type ConferenceRecordStartRequest struct {
	Format         string `json:"format"`   // wav | mp3
	Channels       string `json:"channels"` // single | dual
	ClientState    string `json:"client_state,omitempty"`
	CommandID      string `json:"command_id,omitempty"`
	PlayBeep       bool   `json:"play_beep,omitempty"`
	MaxLength      int    `json:"max_length,omitempty"`
	Trim           string `json:"trim,omitempty"`
	CustomFileName string `json:"custom_file_name,omitempty"`
}

// ConferenceRecordStart begins recording a conference.
func (c *Client) ConferenceRecordStart(ctx context.Context, conferenceID string, req ConferenceRecordStartRequest) (*ConferenceCommandResponse, error) {
	return c.conferenceAction(ctx, conferenceID, "record_start", req)
}

// ConferenceRecordStop stops recording a conference.
func (c *Client) ConferenceRecordStop(ctx context.Context, conferenceID string) (*ConferenceCommandResponse, error) {
	return c.conferenceAction(ctx, conferenceID, "record_stop", map[string]any{})
}

// KickConferenceParticipantRequest removes a participant from the conference.
type KickConferenceParticipantRequest struct {
	CallControlIDs []string `json:"call_control_ids"`
}

// KickConferenceParticipants removes participants from the conference.
func (c *Client) KickConferenceParticipants(ctx context.Context, conferenceID string, callControlIDs []string) (*ConferenceCommandResponse, error) {
	return c.conferenceAction(ctx, conferenceID, "kick", KickConferenceParticipantRequest{CallControlIDs: callControlIDs})
}

// UpdateConferenceParticipantRequest updates participant settings.
type UpdateConferenceParticipantRequest struct {
	CallControlID           string `json:"call_control_id"`
	EndConferenceOnExit     *bool  `json:"end_conference_on_exit,omitempty"`
	SoftEndConferenceOnExit *bool  `json:"soft_end_conference_on_exit,omitempty"`
}

// UpdateConferenceParticipant updates a participant's settings.
func (c *Client) UpdateConferenceParticipant(ctx context.Context, conferenceID string, req UpdateConferenceParticipantRequest) (*ConferenceCommandResponse, error) {
	return c.conferenceAction(ctx, conferenceID, "update", req)
}
