# telnyx-api

A Go client library for the [Telnyx](https://telnyx.com) API covering:

- **Call Control** — full call command surface (dial, answer, transfer, gather, speak, record, stream, fork, AI, and more)
- **Webhooks** — ED25519 signature verification, typed dispatch for all 46 event types
- **Conversational AI** — Assistant CRUD, call attachment, and chat
- **Messaging** — SMS/MMS send, scheduling, retrieval, and messaging profile management
- **Conferences** — Full conference room management and participant control
- **Phone Numbers** — Search, purchase, update, and release numbers
- **Number Lookup** — Carrier, CNAM, and LERG/portability lookup
- **Recordings & Transcriptions** — List, retrieve, and delete call recordings
- **Media Storage** — Upload (URL or multipart), list, get, and delete media files
- **Fax** — Send faxes (URL or file), manage fax applications, list and cancel faxes
- **FQDN Connections** — Create and manage FQDN-based SIP connections and individual FQDNs
- **Verify / OTP** — Trigger SMS, call, and flashcall verifications; submit codes; manage verify profiles
- **Connections** — Credential and IP connection CRUD
- **Billing Groups** — Billing group CRUD
- **CDR Reports** — Request, poll, list, and delete call detail record reports
- **Audit Logs** — List account audit events with date and sort filters
- **Webhook Deliveries** — Query delivery history and retry status

Zero external dependencies — only the Go standard library.

## Installation

```bash
go get github.com/bnandaku/telnyx-api
```

Requires Go 1.21+.

## Quick start

```go
import telnyx "github.com/bnandaku/telnyx-api"

client := telnyx.NewClient("YOUR_API_KEY")
```

### Client options

```go
// Custom HTTP client
client := telnyx.NewClient("API_KEY",
    telnyx.WithHTTPClient(&http.Client{Timeout: 10 * time.Second}),
)

// Override base URL (useful for testing)
client := telnyx.NewClient("API_KEY",
    telnyx.WithBaseURL("http://localhost:9090"),
)
```

---

## Call Control

All call control commands follow the pattern: `client.CommandName(ctx, callControlID, Request{...})`.

### Dial

Initiate an outbound call.

```go
call, err := client.Dial(ctx, telnyx.DialRequest{
    ConnectionID: "your-connection-id",
    To:           "+15550002222",
    From:         "+15550001111",
    TimeoutSecs:  30,
    TimeLimitSecs: 3600,
    WebhookURL:   "https://yourserver.com/webhooks",
    ClientState:  base64.StdEncoding.EncodeToString([]byte(`{"order":"123"}`)),
})
fmt.Println(call.CallControlID)
```

**Key fields:** `ConnectionID` (required), `To` (string or `[]string`), `From`, `TimeoutSecs`, `TimeLimitSecs`, `AnsweringMachineDetection`, `Record`, `Transcription`, `StreamURL`, `Assistant`, `WebhookURL`.

### Retrieve a call

```go
call, err := client.GetCall(ctx, callControlID)
fmt.Println(call.IsAlive)
```

### Answer

Answer an inbound call.

```go
_, err := client.Answer(ctx, callControlID, telnyx.AnswerRequest{
    WebhookURL:    "https://yourserver.com/webhooks",
    Record:        "record-from-answer",
    RecordFormat:  "mp3",
    RecordChannels: "dual",
})
```

### Hangup

```go
_, err := client.Hangup(ctx, callControlID, telnyx.HangupRequest{})
```

### Reject

Reject an inbound call without answering.

```go
_, err := client.Reject(ctx, callControlID, telnyx.RejectRequest{
    Cause: "USER_BUSY", // or CALL_REJECTED
})
```

### Bridge

Connect two call legs together.

```go
_, err := client.Bridge(ctx, callControlID, telnyx.BridgeRequest{
    CallControlID: otherCallControlID,
    Record:        "record-from-answer",
})
```

### Transfer

Redirect a call to a new destination.

```go
_, err := client.Transfer(ctx, callControlID, telnyx.TransferRequest{
    To:          "+15550003333",
    TimeoutSecs: 30,
    WebhookURL:  "https://yourserver.com/webhooks",
})
```

### Speak (Text-to-Speech)

```go
_, err := client.Speak(ctx, callControlID, telnyx.SpeakRequest{
    Payload:     "Press 1 for sales, press 2 for support.",
    Voice:       "Telnyx.KokoroTTS.af",
    PayloadType: "text",     // text | ssml
    TargetLegs:  "self",     // self | opposite | both
})
```

### Play Audio

```go
_, err := client.PlayAudio(ctx, callControlID, telnyx.PlayAudioRequest{
    AudioURL:   "https://example.com/hold.mp3",
    Loop:       "infinity",
    TargetLegs: "self",
    Overlay:    false,
})

// Stop playback
_, err = client.PlaybackStop(ctx, callControlID, "current", "", "")
```

### Gather DTMF

Collect digit input from the caller.

```go
_, err := client.Gather(ctx, callControlID, telnyx.GatherRequest{
    MinimumDigits:           1,
    MaximumDigits:           1,
    TimeoutMillis:           10000,
    InterDigitTimeoutMillis: 3000,
    TerminatingDigit:        "#",
    ValidDigits:             "0123456789",
})
// Result arrives via EventCallGatherEnded webhook

// Cancel an in-progress gather
_, err = client.GatherStop(ctx, callControlID, "", "")
```

### Send DTMF

```go
_, err := client.SendDTMF(ctx, callControlID, telnyx.SendDTMFRequest{
    Digits:     "1234#",
    DurationMS: 100,
})
```

### Hold / Unhold

```go
_, err := client.Hold(ctx, callControlID, telnyx.HoldRequest{
    AudioURL: "https://example.com/hold.mp3",
})

_, err = client.Unhold(ctx, callControlID, "", "")
```

### Mute / Unmute

```go
_, err := client.Mute(ctx, callControlID, telnyx.MuteRequest{Mute: "self"})
_, err = client.Unmute(ctx, callControlID, telnyx.MuteRequest{Mute: "self"})
// Mute values: self | opposite | both
```

### Recording

```go
_, err := client.RecordStart(ctx, callControlID, telnyx.RecordStartRequest{
    Format:        "mp3",
    Channels:      "dual",
    MaxLength:     3600,
    Transcription: true,
    TranscriptionEngine: "B",
    TranscriptionLanguage: "en-US",
})

_, err = client.RecordPause(ctx, callControlID, "", "")
_, err = client.RecordResume(ctx, callControlID, "", "")
_, err = client.RecordStop(ctx, callControlID, "", "")
// Saved recording arrives via EventCallRecordingSaved
```

### Real-time Transcription

```go
_, err := client.TranscriptionStart(ctx, callControlID, telnyx.TranscriptionStartRequest{
    TranscriptionEngine:   "B",
    TranscriptionLanguage: "en-US",
    InterimResults:        true,
})

_, err = client.TranscriptionStop(ctx, callControlID, "", "")
// Segments arrive via EventCallTranscription
```

### Media Streaming (WebSocket)

```go
_, err := client.StreamingStart(ctx, callControlID, telnyx.StreamingStartRequest{
    StreamURL:   "wss://yourserver.com/media",
    StreamTrack: "both_tracks",
})

_, err = client.StreamingStop(ctx, callControlID, "", "")
```

### Media Forking

```go
_, err := client.ForkStart(ctx, callControlID, telnyx.ForkStartRequest{
    Target: "udp://10.0.0.1:8080",
})

_, err = client.ForkStop(ctx, callControlID, "", "")
```

### SIP REFER

```go
_, err := client.Refer(ctx, callControlID, telnyx.ReferRequest{
    SIPAddress: "sip:transfer@example.com",
})
// Events: EventCallReferStarted, EventCallReferCompleted, EventCallReferFailed
```

### SIPREC

```go
_, err := client.SIPRECStart(ctx, callControlID, telnyx.SIPRECStartRequest{
    SIPRECDestinationAddress: "sip:recorder@example.com",
    SIPRECTrack:              "both",
})

_, err = client.SIPRECStop(ctx, callControlID, "", "")
```

### Queue

```go
_, err := client.Enqueue(ctx, callControlID, telnyx.EnqueueRequest{
    QueueName:       "support",
    MaxWaitTimeSecs: 300,
})

_, err = client.LeaveQueue(ctx, callControlID, "", "")
```

### Answering Machine Detection

```go
call, err := client.Dial(ctx, telnyx.DialRequest{
    ConnectionID:              "conn-id",
    To:                        "+15550002222",
    From:                      "+15550001111",
    AnsweringMachineDetection: "premium", // premium | detect | detect_beep | detect_words | greeting_end | disabled
})
// Result via EventCallMachinePremiumDetectionEnded
```

### Update Client State

```go
_, err := client.UpdateClientState(ctx, callControlID, base64.StdEncoding.EncodeToString([]byte(`{"step":"2"}`)))
```

---

## Conversational AI

### Start AI Assistant on a call

```go
resp, err := client.StartAIAssistant(ctx, callControlID, telnyx.StartAIAssistantRequest{
    Assistant:   telnyx.CallAssistantRequest{ID: "asst_abc123"},
    Greeting:    "Hi! How can I help you today?",
    SendMessageHistoryUpdates: true,
    MessageHistory: []telnyx.MessageHistoryEntry{
        {Role: "system", Content: "The caller is a premium subscriber."},
    },
})
fmt.Println(resp.Data.ConversationID)
// Events: EventCallAIGatherEnded, EventCallConversationEnded, EventCallAIGatherPartialResults
```

### Stop AI Assistant

```go
_, err := client.StopAIAssistant(ctx, callControlID, telnyx.StopAIAssistantRequest{})
```

### Create an Assistant

```go
asst, err := client.CreateAssistant(ctx, telnyx.CreateAssistantRequest{
    Name:         "Support Bot",
    Instructions: "You are a helpful customer support agent. Be concise and friendly.",
    Model:        "gpt-4o",
    Greeting:     "Hello! How can I assist you today?",
    EnabledFeatures: []string{"telephony"},
    TelephonySettings: &telnyx.TelephonySettings{
        NoiseSuppression: true,
        SilenceTimeoutMS: 3000,
        MaxCallDuration:  1800,
    },
    VoiceSettings: &telnyx.VoiceSettings{
        Provider: "Telnyx",
        VoiceID:  "KokoroTTS.af",
    },
    Tags: []string{"support", "production"},
})
```

### Get an Assistant

```go
asst, err := client.GetAssistant(ctx, "asst_abc123")
```

### Update an Assistant

```go
asst, err := client.UpdateAssistant(ctx, "asst_abc123", telnyx.UpdateAssistantRequest{
    Instructions: "Updated instructions with new product information.",
    Model:        "gpt-4o-mini",
})
```

### List Assistants

```go
assistants, err := client.ListAssistants(ctx, telnyx.ListAssistantsParams{
    PageNumber: 1,
    PageSize:   20,
})
```

### Delete an Assistant

```go
err := client.DeleteAssistant(ctx, "asst_abc123")
```

### Chat (beta)

Send a message to an assistant and get a reply outside of a call.

```go
resp, err := client.Chat(ctx, "asst_abc123", telnyx.ChatRequest{
    Messages: []telnyx.ChatMessage{
        {Role: "user", Content: "What are your business hours?"},
    },
})
fmt.Println(resp.Data.Message.Content)
```

---

## Messaging

### Send an SMS

```go
msg, err := client.SendMessage(ctx, telnyx.SendMessageRequest{
    From: "+15550001111",
    To:   "+15550002222",
    Text: "Hello from Telnyx!",
})
fmt.Println(msg.ID, msg.To[0].Status)
```

### Send an MMS

```go
msg, err := client.SendMessage(ctx, telnyx.SendMessageRequest{
    From:      "+15550001111",
    To:        "+15550002222",
    Text:      "Check out this image!",
    MediaURLs: []string{"https://example.com/image.jpg"},
    Type:      "MMS",
})
```

### Send a Group MMS

```go
msg, err := client.SendMessage(ctx, telnyx.SendMessageRequest{
    From: "+15550001111",
    To:   []string{"+15550002222", "+15550003333"},
    Text: "Group message!",
    Type: "MMS",
})
```

### Get a Message

```go
msg, err := client.GetMessage(ctx, "msg-uuid")
fmt.Println(msg.To[0].Status) // queued | sent | delivered | delivery_failed
```

### List Messages

```go
messages, meta, err := client.ListMessages(ctx, telnyx.ListMessagesParams{
    PageNumber: 1,
    PageSize:   25,
    Direction:  "outbound",
    DateRangeStart: "2024-01-01T00:00:00Z",
    DateRangeEnd:   "2024-01-31T23:59:59Z",
})
fmt.Printf("%d total messages\n", meta.TotalResults)
```

### Schedule a Message

```go
msg, err := client.ScheduleMessage(ctx, telnyx.ScheduleMessageRequest{
    SendMessageRequest: telnyx.SendMessageRequest{
        From: "+15550001111",
        To:   "+15550002222",
        Text: "Scheduled reminder",
    },
    SendAt: "2024-06-01T09:00:00Z",
})
```

### Cancel a Scheduled Message

```go
err := client.CancelScheduledMessage(ctx, "msg-uuid")
```

### Messaging Profiles

```go
// Create
profile, err := client.CreateMessagingProfile(ctx, telnyx.CreateMessagingProfileRequest{
    Name:               "My Profile",
    WebhookURL:         "https://yourserver.com/webhooks/messaging",
    WebhookAPIVersion:  "2",
    WhitelistedDestinations: []string{"US", "CA"},
    MMSFallBackToSMS:   true,
    SmartEncoding:      true,
})

// Get
profile, err = client.GetMessagingProfile(ctx, profile.ID)

// Update
enabled := false
profile, err = client.UpdateMessagingProfile(ctx, profile.ID, telnyx.UpdateMessagingProfileRequest{
    Enabled: &enabled,
})

// List
profiles, meta, err := client.ListMessagingProfiles(ctx, 1, 20)

// Delete
err = client.DeleteMessagingProfile(ctx, profile.ID)
```

---

## Conferences

### Create a Conference

```go
conf, err := client.CreateConference(ctx, telnyx.CreateConferenceRequest{
    CallControlID:   callControlID,
    Name:            "Team Standup",
    BeepEnabled:     "on_enter",
    MaxParticipants: 50,
    DurationMinutes: 60,
    Region:          "US",
})
```

### Get a Conference

```go
conf, err := client.GetConference(ctx, conf.ID)
fmt.Println(conf.Status) // init | in_progress | completed
```

### List Conferences

```go
conferences, meta, err := client.ListConferences(ctx, 1, 20, "US")
```

### Join a Conference

```go
_, err := client.JoinConference(ctx, conf.ID, telnyx.JoinConferenceRequest{
    CallControlID:           callControlID,
    Mute:                    false,
    Hold:                    false,
    EndConferenceOnExit:     false,
    SoftEndConferenceOnExit: true,
    SupervisorRole:          "monitor", // barge | whisper | monitor
})
```

### List Participants

```go
participants, meta, err := client.ListConferenceParticipants(ctx, conf.ID, 1, 50)
for _, p := range participants {
    fmt.Printf("%s: muted=%v, on_hold=%v\n", p.CallControlID, p.Muted, p.OnHold)
}
```

### Mute / Unmute Participants

```go
// Mute specific participants
_, err := client.MuteConferenceParticipants(ctx, conf.ID, []string{callControlID})

// Unmute all participants
_, err = client.UnmuteConferenceParticipants(ctx, conf.ID, []string{})
```

### Hold / Unhold Participants

```go
_, err := client.HoldConferenceParticipants(ctx, conf.ID, telnyx.HoldConferenceParticipantsRequest{
    CallControlIDs: []string{callControlID},
    AudioURL:       "https://example.com/hold.mp3",
})

_, err = client.UnholdConferenceParticipants(ctx, conf.ID, []string{callControlID})
```

### Speak in a Conference

```go
_, err := client.ConferenceSpeak(ctx, conf.ID, telnyx.ConferenceSpeakRequest{
    Payload:  "The meeting will begin shortly.",
    Voice:    "Telnyx.KokoroTTS.af",
    Language: "en-US",
})
```

### Play Audio in a Conference

```go
_, err := client.ConferencePlayAudio(ctx, conf.ID, telnyx.ConferencePlayAudioRequest{
    AudioURL: "https://example.com/intro.mp3",
})

_, err = client.ConferenceStopAudio(ctx, conf.ID, []string{})
```

### Record a Conference

```go
_, err := client.ConferenceRecordStart(ctx, conf.ID, telnyx.ConferenceRecordStartRequest{
    Format:   "mp3",
    Channels: "single",
})

_, err = client.ConferenceRecordStop(ctx, conf.ID)
// Recording saved via EventConferenceRecordingSaved
```

### Kick Participants

```go
_, err := client.KickConferenceParticipants(ctx, conf.ID, []string{callControlID})
```

### Update a Participant

```go
t := true
_, err := client.UpdateConferenceParticipant(ctx, conf.ID, telnyx.UpdateConferenceParticipantRequest{
    CallControlID:           callControlID,
    SoftEndConferenceOnExit: &t,
})
```

---

## Phone Numbers

### Search Available Numbers

```go
numbers, _, err := client.SearchAvailableNumbers(ctx, telnyx.SearchAvailableNumbersParams{
    FilterCountryCode:    "US",
    FilterPhoneNumberType: "local",
    FilterFeatures:       "sms,voice",
    FilterAreaCode:       "+1212", // starts_with
    FilterLimit:          10,
    FilterQuickship:      true,
})
for _, n := range numbers {
    fmt.Println(n.PhoneNumber, n.Costs[0].Amount, n.Costs[0].Currency)
}
```

### Order Phone Numbers

```go
order, err := client.OrderPhoneNumbers(ctx, telnyx.NumberOrderRequest{
    PhoneNumbers: []telnyx.NumberOrderItem{
        {PhoneNumber: "+12125550100"},
        {PhoneNumber: "+12125550101"},
    },
    ConnectionID:       "conn-id",
    MessagingProfileID: "profile-id",
})
fmt.Println(order.Status) // pending | success | failure
```

### Get a Number Order

```go
order, err := client.GetNumberOrder(ctx, order.ID)
```

### List Your Numbers

```go
numbers, meta, err := client.ListPhoneNumbers(ctx, telnyx.ListPhoneNumbersParams{
    PageNumber:         1,
    PageSize:           50,
    FilterStatus:       "active",
    FilterPhoneNumber:  "+1212",
    Sort:               "purchased_at",
})
```

### Get a Number

```go
number, err := client.GetPhoneNumber(ctx, numberID)
fmt.Println(number.PhoneNumber, number.Status, number.PhoneNumberType)
```

### Update a Number

```go
number, err := client.UpdatePhoneNumber(ctx, numberID, telnyx.UpdatePhoneNumberRequest{
    ConnectionID:       "conn-id",
    MessagingProfileID: "profile-id",
    Tags:               []string{"production", "support"},
    CustomerReference:  "acct-456",
})
```

### Release a Number

```go
err := client.DeletePhoneNumber(ctx, numberID)
```

---

## Number Lookup

Perform a deep lookup on any E.164 phone number — carrier, CNAM, and LERG/portability data are fetched concurrently and merged.

```go
result, err := client.LookupNumber(ctx, "+15550001111")
if err != nil {
    log.Fatal(err)
}

fmt.Println(result.Carrier.Name)        // AT&T
fmt.Println(result.Carrier.Type)        // mobile | landline | voip | toll free
fmt.Println(result.CallerName.CallerName) // John Doe
fmt.Println(result.Portability.LRN)     // Local routing number
fmt.Println(result.Portability.PortedStatus) // Y | N
fmt.Println(result.Portability.City, result.Portability.State)
```

---

## Recordings

### List Recordings

```go
recordings, meta, err := client.ListRecordings(ctx, telnyx.ListRecordingsParams{
    PageNumber:          1,
    PageSize:            25,
    FilterCallLegID:     "leg-uuid",
    FilterStatus:        "completed",
    FilterSource:        "call_leg", // call_leg | conference
})
```

### Get a Recording

```go
recording, err := client.GetRecording(ctx, "rec-uuid")
fmt.Println(recording.Urls.MP3) // download URL
fmt.Println(recording.DurationMillis)
```

### Delete a Recording

```go
err := client.DeleteRecording(ctx, "rec-uuid")
```

### Get a Transcription

```go
transcript, err := client.GetTranscription(ctx, "transcription-uuid")
fmt.Println(transcript.TranscriptionText)
for _, speaker := range transcript.Speakers {
    for _, seg := range speaker.Segments {
        fmt.Printf("[%.2fs] %s\n", seg.StartTime, seg.Text)
    }
}
```

### List Transcriptions

```go
transcripts, meta, err := client.ListTranscriptions(ctx, "rec-uuid", 1, 20)
```

---

## Webhooks

### Setup

Get your public key from the Telnyx portal: **API Keys → Webhook signing secret** (base64-encoded ED25519 key).

```go
handler, err := telnyx.NewWebhookHandler("BASE64_PUBLIC_KEY")
if err != nil {
    log.Fatal(err)
}

http.Handle("/webhooks", handler)
log.Fatal(http.ListenAndServe(":8080", nil))
```

### Registering handlers

```go
handler.On(telnyx.EventCallInitiated, func(ctx context.Context, event telnyx.Event, payload any) error {
    p := payload.(*telnyx.CallInitiatedPayload)
    fmt.Printf("Incoming call: %s → %s (%s)\n", p.From, p.To, p.Direction)
    return nil
})

handler.On(telnyx.EventCallAnswered, func(ctx context.Context, event telnyx.Event, payload any) error {
    p := payload.(*telnyx.CallAnsweredPayload)
    // start IVR, record, etc.
    return nil
})

handler.On(telnyx.EventCallGatherEnded, func(ctx context.Context, event telnyx.Event, payload any) error {
    p := payload.(*telnyx.CallGatherEndedPayload)
    fmt.Printf("Gathered: %s (status: %s)\n", p.Digits, p.Status)
    return nil
})

handler.On(telnyx.EventCallHangup, func(ctx context.Context, event telnyx.Event, payload any) error {
    p := payload.(*telnyx.CallHangupPayload)
    fmt.Printf("Call ended: %s\n", p.HangupCause)
    return nil
})

handler.On(telnyx.EventCallRecordingSaved, func(ctx context.Context, event telnyx.Event, payload any) error {
    p := payload.(*telnyx.CallRecordingSavedPayload)
    fmt.Printf("Recording saved: %s (%dms)\n", p.RecordingID, p.DurationMillis)
    return nil
})

handler.On(telnyx.EventCallTranscription, func(ctx context.Context, event telnyx.Event, payload any) error {
    p := payload.(*telnyx.CallTranscriptionPayload)
    if p.TranscriptionData.IsFinal {
        fmt.Printf("[%s] %s\n", p.TranscriptionData.TranscriptFrom, p.TranscriptionData.Transcript)
    }
    return nil
})

handler.On(telnyx.EventCallConversationEnded, func(ctx context.Context, event telnyx.Event, payload any) error {
    p := payload.(*telnyx.CallConversationEndedPayload)
    fmt.Printf("AI conversation ended: %s (%d messages)\n", p.EndReason, len(p.MessageHistory))
    return nil
})

// Catch-all for unknown/future event types
handler.OnFallback(func(ctx context.Context, event telnyx.Event, payload any) error {
    fmt.Printf("Unhandled event: %s\n", event.Data.EventType)
    return nil
})
```

### All supported event types

| Constant | Event string | Payload type |
|---|---|---|
| `EventCallInitiated` | `call.initiated` | `*CallInitiatedPayload` |
| `EventCallAnswered` | `call.answered` | `*CallAnsweredPayload` |
| `EventCallHangup` | `call.hangup` | `*CallHangupPayload` |
| `EventCallHold` | `call.hold` | `*CallHoldPayload` |
| `EventCallUnhold` | `call.unhold` | `*CallUnholdPayload` |
| `EventCallBridged` | `call.bridged` | `*CallBridgedPayload` |
| `EventCallPlaybackStarted` | `call.playback.started` | `*CallPlaybackStartedPayload` |
| `EventCallPlaybackEnded` | `call.playback.ended` | `*CallPlaybackEndedPayload` |
| `EventCallSpeakStarted` | `call.speak.started` | `*CallSpeakStartedPayload` |
| `EventCallSpeakEnded` | `call.speak.ended` | `*CallSpeakEndedPayload` |
| `EventCallDTMFReceived` | `call.dtmf.received` | `*CallDTMFReceivedPayload` |
| `EventCallGatherEnded` | `call.gather.ended` | `*CallGatherEndedPayload` |
| `EventCallRecordingSaved` | `call.recording.saved` | `*CallRecordingSavedPayload` |
| `EventCallRecordingError` | `call.recording.error` | `*CallRecordingErrorPayload` |
| `EventCallRecordingTranscriptionSaved` | `call.recording.transcription.saved` | `*CallRecordingTranscriptionSavedPayload` |
| `EventCallMachineDetectionEnded` | `call.machine.detection.ended` | `*CallMachineDetectionEndedPayload` |
| `EventCallMachineGreetingEnded` | `call.machine.greeting.ended` | `*CallMachineGreetingEndedPayload` |
| `EventCallMachinePremiumDetectionEnded` | `call.machine.premium.detection.ended` | `*CallMachinePremiumDetectionEndedPayload` |
| `EventCallMachinePremiumGreetingEnded` | `call.machine.premium.greeting.ended` | `*CallMachinePremiumGreetingEndedPayload` |
| `EventCallForkStarted` | `call.fork.started` | `*CallForkStartedPayload` |
| `EventCallForkStopped` | `call.fork.stopped` | `*CallForkStoppedPayload` |
| `EventCallEnqueued` | `call.enqueued` | `*CallEnqueuedPayload` |
| `EventCallDequeued` | `call.dequeued` | `*CallDequeuedPayload` |
| `EventCallLeftQueue` | `call.left_queue` | `*CallLeftQueuePayload` |
| `EventCallTranscription` | `call.transcription` | `*CallTranscriptionPayload` |
| `EventStreamingStarted` | `streaming.started` | `*StreamingStartedPayload` |
| `EventStreamingStopped` | `streaming.stopped` | `*StreamingStoppedPayload` |
| `EventStreamingFailed` | `streaming.failed` | `*StreamingFailedPayload` |
| `EventCallReferStarted` | `call.refer.started` | `*CallReferStartedPayload` |
| `EventCallReferCompleted` | `call.refer.completed` | `*CallReferCompletedPayload` |
| `EventCallReferFailed` | `call.refer.failed` | `*CallReferFailedPayload` |
| `EventCallSIPRECStarted` | `call.siprec.started` | `*CallSIPRECStartedPayload` |
| `EventCallSIPRECStopped` | `call.siprec.stopped` | `*CallSIPRECStoppedPayload` |
| `EventCallSIPRECFailed` | `call.siprec.failed` | `*CallSIPRECFailedPayload` |
| `EventCallAIGatherEnded` | `call.ai_gather.ended` | `*CallAIGatherEndedPayload` |
| `EventCallAIGatherMessageHistoryUpdated` | `call.ai_gather.message_history.updated` | `*CallAIGatherMessageHistoryUpdatedPayload` |
| `EventCallAIGatherPartialResults` | `call.ai_gather.partial_results` | `*CallAIGatherPartialResultsPayload` |
| `EventCallConversationEnded` | `call.conversation.ended` | `*CallConversationEndedPayload` |
| `EventCallConversationInsightsGenerated` | `call.conversation.insights_generated` | `*CallConversationInsightsGeneratedPayload` |
| `EventCallDeepfakeDetectionResult` | `call.deepfake_detection.result` | `*DeepfakeDetectionResultPayload` |
| `EventCallDeepfakeDetectionError` | `call.deepfake_detection.error` | `*DeepfakeDetectionErrorPayload` |
| `EventConferenceCreated` | `conference.created` | `*ConferenceCreatedPayload` |
| `EventConferenceEnded` | `conference.ended` | `*ConferenceEndedPayload` |
| `EventConferenceFloorChanged` | `conference.floor_changed` | `*ConferenceFloorChangedPayload` |
| `EventConferenceParticipantJoined` | `conference.participant.joined` | `*ConferenceParticipantJoinedPayload` |
| `EventConferenceParticipantLeft` | `conference.participant.left` | `*ConferenceParticipantLeftPayload` |
| `EventConferencePlaybackStarted` | `conference.playback.started` | `*ConferencePlaybackStartedPayload` |
| `EventConferencePlaybackEnded` | `conference.playback.ended` | `*ConferencePlaybackEndedPayload` |
| `EventConferenceSpeakStarted` | `conference.speak.started` | `*ConferenceSpeakStartedPayload` |
| `EventConferenceSpeakEnded` | `conference.speak.ended` | `*ConferenceSpeakEndedPayload` |
| `EventConferenceRecordingSaved` | `conference.recording.saved` | `*ConferenceRecordingSavedPayload` |

### Signature verification

Every request is verified using ED25519. Requests with an invalid or missing `Telnyx-Signature-Ed25519` header receive `401 Unauthorized`. Webhooks older than 5 minutes are rejected to prevent replay attacks.

---

## Media Storage

```go
// Upload from URL
media, err := client.UploadMediaFromURL(ctx, telnyx.UploadMediaFromURLRequest{
    MediaURL:  "https://example.com/hold.mp3",
    MediaName: "hold-music",
})
fmt.Println(media.MediaName, media.ContentType)

// Upload a local file
f, _ := os.Open("audio.mp3")
defer f.Close()
media, err = client.UploadMediaFile(ctx, "audio.mp3", "audio/mpeg", f)

// List media
items, meta, err := client.ListMedia(ctx, telnyx.ListMediaParams{
    PageNumber:  1,
    PageSize:    25,
    ContentType: "audio/mpeg",
})

// Get and delete
media, err = client.GetMedia(ctx, media.MediaName)
err = client.DeleteMedia(ctx, media.MediaName)
```

---

## Fax

### Send a fax

```go
fax, err := client.SendFax(ctx, telnyx.SendFaxRequest{
    ConnectionID: "your-connection-id",
    To:           "+15550002222",
    From:         "+15550001111",
    MediaURL:     "https://example.com/document.pdf",
    Quality:      "high",    // normal | high | very-high
    StoreMedia:   true,
})
fmt.Println(fax.ID, fax.Status)
```

### Get / list / cancel / delete faxes

```go
fax, err := client.GetFax(ctx, fax.ID)

faxes, meta, err := client.ListFaxes(ctx, telnyx.ListFaxesParams{
    PageNumber: 1,
    PageSize:   25,
    FilterTo:   "+15550002222",
})

fax, err = client.CancelFax(ctx, fax.ID)
err = client.DeleteFax(ctx, fax.ID)
```

### Fax applications

```go
app, err := client.CreateFaxApplication(ctx, telnyx.CreateFaxApplicationRequest{
    ApplicationName: "My Fax App",
    WebhookEventURL: "https://yourserver.com/webhooks/fax",
})

app, err = client.GetFaxApplication(ctx, app.ID)
apps, meta, err := client.ListFaxApplications(ctx, 1, 20)
err = client.DeleteFaxApplication(ctx, app.ID)
```

---

## FQDN Connections

```go
conn, err := client.CreateFQDNConnection(ctx, telnyx.CreateFQDNConnectionRequest{
    ConnectionName:    "My FQDN Connection",
    TransportProtocol: "UDP",
    WebhookEventURL:   "https://yourserver.com/webhooks",
})

// Add an FQDN to the connection
fqdn, err := client.CreateFQDN(ctx, telnyx.CreateFQDNRequest{
    ConnectionID:  conn.ID,
    FQDN:          "sip.example.com",
    DNSRecordType: "a",
    Port:          5060,
})

conns, meta, err := client.ListFQDNConnections(ctx, telnyx.ListFQDNConnectionsParams{
    PageNumber: 1,
    PageSize:   20,
})

fqdns, meta, err := client.ListFQDNs(ctx, conn.ID, 1, 50)

err = client.DeleteFQDN(ctx, fqdn.ID)
err = client.DeleteFQDNConnection(ctx, conn.ID)
```

---

## Verify / OTP

### Send a verification

```go
// Via SMS
v, err := client.SendVerificationSMS(ctx, telnyx.CreateVerificationSMSRequest{
    PhoneNumber:     "+15550001111",
    VerifyProfileID: "profile-uuid",
})

// Via phone call
v, err = client.SendVerificationCall(ctx, telnyx.CreateVerificationCallRequest{
    PhoneNumber:     "+15550001111",
    VerifyProfileID: "profile-uuid",
})

// Via flashcall (missed call with code in caller ID)
v, err = client.SendVerificationFlashcall(ctx, telnyx.CreateVerificationFlashcallRequest{
    PhoneNumber:     "+15550001111",
    VerifyProfileID: "profile-uuid",
})

fmt.Println(v.ID, v.Status) // pending | accepted | invalid | expired | error
```

### Submit a code

```go
// By verification ID
resp, err := client.VerifyCodeByID(ctx, v.ID, telnyx.VerifyCodeByIDRequest{
    Code: "123456",
})
fmt.Println(resp.ResponseCode) // accepted | rejected

// By phone number
resp, err = client.VerifyCodeByPhone(ctx, "+15550001111", telnyx.VerifyCodeByPhoneRequest{
    Code:            "123456",
    VerifyProfileID: "profile-uuid",
})
```

### Verify profiles

```go
profile, err := client.CreateVerifyProfile(ctx, telnyx.CreateVerifyProfileRequest{
    Name:     "My App",
    Language: "en-US",
})

profile, err = client.GetVerifyProfile(ctx, profile.ID)
profiles, meta, err := client.ListVerifyProfiles(ctx, 1, 20)
err = client.DeleteVerifyProfile(ctx, profile.ID)
```

---

## Billing Groups

```go
group, err := client.CreateBillingGroup(ctx, telnyx.BillingGroupRequest{Name: "Team A"})
group, err = client.GetBillingGroup(ctx, group.ID)

groups, meta, err := client.ListBillingGroups(ctx, 1, 20)

group, err = client.UpdateBillingGroup(ctx, group.ID, telnyx.BillingGroupRequest{Name: "Team B"})
err = client.DeleteBillingGroup(ctx, group.ID)
```

---

## Connections

### Credential connections

```go
conn, err := client.CreateCredentialConnection(ctx, telnyx.CreateCredentialConnectionRequest{
    Name:     "My SIP Trunk",
    Username: "alice",
    Password: "s3cr3t",
})

conn, err = client.GetCredentialConnection(ctx, conn.ID)
conns, meta, err := client.ListCredentialConnections(ctx, 1, 20)
err = client.DeleteCredentialConnection(ctx, conn.ID)
```

### IP connections

```go
ipConn, err := client.CreateIPConnection(ctx, telnyx.CreateIPConnectionRequest{
    Name: "My IP Trunk",
})

ipConn, err = client.GetIPConnection(ctx, ipConn.ID)
ipConns, meta, err := client.ListIPConnections(ctx, 1, 20)
err = client.DeleteIPConnection(ctx, ipConn.ID)
```

---

## CDR Reports

```go
report, err := client.CreateCDRReport(ctx, telnyx.CDRFilter{
    StartTime: "2024-01-01T00:00:00Z",
    EndTime:   "2024-01-31T23:59:59Z",
    Timezone:  "UTC",
    CallTypes: []int{1, 2}, // 1=Inbound, 2=Outbound
})

// Poll until complete
for report.Status != telnyx.CDRStatusComplete {
    time.Sleep(5 * time.Second)
    report, err = client.GetCDRReport(ctx, report.ID)
}

reports, meta, err := client.ListCDRReports(ctx, 1, 20)
err = client.DeleteCDRReport(ctx, report.ID)
```

---

## Audit Logs

```go
logs, meta, err := client.ListAuditLogs(ctx, telnyx.ListAuditLogsParams{
    PageNumber:    1,
    PageSize:      50,
    CreatedAfter:  "2024-01-01T00:00:00Z",
    CreatedBefore: "2024-12-31T23:59:59Z",
    Sort:          "desc",
})
for _, l := range logs {
    fmt.Printf("%s: %s on %s\n", l.CreatedAt, l.ChangeType, l.ResourceID)
}
```

---

## Webhook Deliveries

```go
// List recent failed deliveries
deliveries, meta, err := client.ListWebhookDeliveries(ctx, telnyx.ListWebhookDeliveriesParams{
    FilterStatus:       "failed",
    FilterStartedAtGTE: "2024-01-01T00:00:00Z",
    PageNumber:         1,
    PageSize:           20,
})
for _, d := range deliveries {
    fmt.Printf("%s: %s → HTTP %d\n", d.Webhook.EventType, d.Status, d.Attempts[len(d.Attempts)-1].Response.Status)
}

// Get a single delivery with full attempt log
delivery, err := client.GetWebhookDelivery(ctx, "delivery-uuid")
```

---

## Error handling

All API errors return `*telnyx.APIError`:

```go
call, err := client.Dial(ctx, req)
if err != nil {
    var apiErr *telnyx.APIError
    if errors.As(err, &apiErr) {
        fmt.Printf("HTTP %d: %s — %s\n", apiErr.StatusCode, apiErr.Errors[0].Title, apiErr.Errors[0].Detail)
    }
    return err
}
```

Webhook signature failures return `*telnyx.ErrSignatureVerification`:

```go
var sigErr *telnyx.ErrSignatureVerification
if errors.As(err, &sigErr) {
    fmt.Println("Bad webhook signature:", sigErr.Reason)
}
```

---

## License

MIT
