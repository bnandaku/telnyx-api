# telnyx-api

A complete Go client library for the [Telnyx](https://telnyx.com) API.

Zero external dependencies — only the Go standard library.

## Contents

- [Installation](#installation)
- [Client](#client)
- [Call Control](#call-control)
- [Conversational AI](#conversational-ai)
- [Webhooks](#webhooks)
- [Messaging](#messaging)
- [Conferences](#conferences)
- [Phone Numbers](#phone-numbers)
- [Number Lookup](#number-lookup)
- [Recordings & Transcriptions](#recordings--transcriptions)
- [Media Storage](#media-storage)
- [Fax](#fax)
- [FQDN Connections](#fqdn-connections)
- [Verify / OTP](#verify--otp)
- [SIP Connections](#sip-connections)
- [Billing Groups](#billing-groups)
- [CDR Reports](#cdr-reports)
- [Audit Logs](#audit-logs)
- [Webhook Deliveries](#webhook-deliveries)
- [Error Handling](#error-handling)
- [Examples](#examples)

---

## Installation

```bash
go get github.com/bnandaku/telnyx-api
```

Requires Go 1.21+.

---

## Client

```go
import telnyx "github.com/bnandaku/telnyx-api"

client := telnyx.NewClient("YOUR_API_KEY")
```

### Options

```go
import "net/http"
import "time"

client := telnyx.NewClient("YOUR_API_KEY",
    // Custom HTTP client (e.g. to set a timeout)
    telnyx.WithHTTPClient(&http.Client{Timeout: 15 * time.Second}),

    // Override the base URL (useful for tests / proxies)
    telnyx.WithBaseURL("http://localhost:9090"),
)
```

All methods accept a `context.Context` as their first argument.

---

## Call Control

### Dial

```go
call, err := client.Dial(ctx, telnyx.DialRequest{
    ConnectionID: "your-connection-id", // required
    To:           "+15550002222",       // string or []string for multi-party
    From:         "+15550001111",

    // Optional timing
    TimeoutSecs:   30,
    TimeLimitSecs: 3600,

    // Optional answering machine detection
    AnsweringMachineDetection: "premium", // premium | detect | detect_beep | detect_words | greeting_end | disabled
    AnsweringMachineDetectionConfig: &telnyx.AMDConfig{
        TotalAnalysisTimeMillis: 5000,
        InitialSilenceMillis:    2000,
    },

    // Optional deepfake detection
    DeepfakeDetection: &telnyx.DeepfakeDetectionConfig{
        Enabled: true,
        Timeout: 30,
    },

    // Optional recording at dial time
    Record:         "record-from-answer", // record-from-answer | record-from-ringing
    RecordChannels: "dual",               // single | dual
    RecordFormat:   "mp3",               // mp3 | wav

    // Optional real-time transcription at dial time
    Transcription: true,
    TranscriptionConfig: &telnyx.TranscriptionConfig{
        TranscriptionEngine:   "B",
        TranscriptionLanguage: "en-US",
    },

    // Optional WebSocket streaming at dial time
    StreamURL:   "wss://yourserver.com/media",
    StreamTrack: "both_tracks", // both_tracks | inbound_track | outbound_track

    // Optional AI assistant at dial time
    Assistant: &telnyx.CallAssistantRequest{ID: "asst_abc123"},

    // Optional webhook override
    WebhookURL: "https://yourserver.com/webhooks",

    // Optional passthrough state (base64-encoded)
    ClientState: "base64data==",
})
fmt.Println(call.CallControlID, call.IsAlive)
```

### Get a call

```go
call, err := client.GetCall(ctx, callControlID)
fmt.Println(call.IsAlive, call.CallDuration)
```

### Answer

```go
_, err := client.Answer(ctx, callControlID, telnyx.AnswerRequest{
    // All fields optional
    Record:        "record-from-answer",
    RecordChannels: "dual",
    RecordFormat:  "mp3",
    Transcription: true,
    WebhookURL:    "https://yourserver.com/webhooks",
    ClientState:   "base64state==",
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
    Cause: "USER_BUSY", // USER_BUSY | CALL_REJECTED
})
```

### Bridge

Connect this call leg to another active call leg.

```go
_, err := client.Bridge(ctx, callControlID, telnyx.BridgeRequest{
    CallControlID:  otherCallControlID,
    Record:         "record-from-answer",
    RecordChannels: "dual",
})
```

### Transfer

Redirect the call to a new destination.

```go
_, err := client.Transfer(ctx, callControlID, telnyx.TransferRequest{
    To:          "+15550003333",
    TimeoutSecs: 30,
    WebhookURL:  "https://yourserver.com/webhooks",

    // Optional AMD on the transferred leg
    AnsweringMachineDetection: "detect",

    // Optional recording on the transferred leg
    Record:         "record-from-answer",
    RecordChannels: "dual",
})
```

### Speak (Text-to-Speech)

```go
_, err := client.Speak(ctx, callControlID, telnyx.SpeakRequest{
    Payload:     "Press 1 for sales, press 2 for support.",
    Voice:       "Telnyx.KokoroTTS.af",
    PayloadType: "text",  // text | ssml
    TargetLegs:  "self",  // self | opposite | both
    ClientState: "base64state==",
})
// Completion fires EventCallSpeakEnded
```

### Play Audio

```go
_, err := client.PlayAudio(ctx, callControlID, telnyx.PlayAudioRequest{
    AudioURL:   "https://example.com/hold.mp3",
    Loop:       "infinity", // or an integer
    TargetLegs: "self",
    Overlay:    false,
})

// Stop playback
_, err = client.PlaybackStop(ctx, callControlID,
    "current", // stop: current | all
    "",         // clientState (optional)
    "",         // commandID (optional)
)
```

### Gather DTMF

```go
_, err := client.Gather(ctx, callControlID, telnyx.GatherRequest{
    MinimumDigits:           1,
    MaximumDigits:           1,
    TimeoutMillis:           10000,
    InterDigitTimeoutMillis: 3000,
    TerminatingDigit:        "#",
    ValidDigits:             "0123456789*#",
    ClientState:             "base64state==",
})
// Result arrives via EventCallGatherEnded

// Cancel a running gather
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
// mute: self | opposite | both
_, err := client.Mute(ctx, callControlID, telnyx.MuteRequest{Mute: "self"})
_, err = client.Unmute(ctx, callControlID, telnyx.MuteRequest{Mute: "self"})
```

### Recording

```go
_, err := client.RecordStart(ctx, callControlID, telnyx.RecordStartRequest{
    Format:   "mp3",  // mp3 | wav
    Channels: "dual", // single | dual

    // Optional transcription of the recording
    Transcription:                   true,
    TranscriptionEngine:             "B",
    TranscriptionLanguage:           "en-US",
    TranscriptionSpeakerDiarization: true,
    TranscriptionMinSpeakerCount:    1,
    TranscriptionMaxSpeakerCount:    4,
})
// Recording URL arrives via EventCallRecordingSaved

_, err = client.RecordPause(ctx, callControlID, "", "")
_, err = client.RecordResume(ctx, callControlID, "", "")
_, err = client.RecordStop(ctx, callControlID, "", "")
```

### Real-time Transcription

```go
_, err := client.TranscriptionStart(ctx, callControlID, telnyx.TranscriptionStartRequest{
    TranscriptionEngine:   "B",
    TranscriptionLanguage: "en-US",
    InterimResults:        true,
})
// Partial + final transcripts arrive via EventCallTranscription

_, err = client.TranscriptionStop(ctx, callControlID, "", "")
```

### WebSocket Media Streaming

```go
_, err := client.StreamingStart(ctx, callControlID, telnyx.StreamingStartRequest{
    StreamURL:   "wss://yourserver.com/media",
    StreamTrack: "both_tracks", // both_tracks | inbound_track | outbound_track
    StreamBidirectionalMode: "mp3",
})
// Events: EventStreamingStarted, EventStreamingStopped, EventStreamingFailed

_, err = client.StreamingStop(ctx, callControlID, "", "")
```

### Media Forking

Fork RTP audio to a UDP target.

```go
_, err := client.ForkStart(ctx, callControlID, telnyx.ForkStartRequest{
    Target: "udp://10.0.0.1:8080",
    RxMode: "single",
    TxMode: "single",
})
// Events: EventCallForkStarted, EventCallForkStopped

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
    SIPRECTrack:              "both", // both | inbound | outbound
})
// Events: EventCallSIPRECStarted, EventCallSIPRECStopped, EventCallSIPRECFailed

_, err = client.SIPRECStop(ctx, callControlID, "", "")
```

### Queue

```go
_, err := client.Enqueue(ctx, callControlID, telnyx.EnqueueRequest{
    QueueName:       "support",
    MaxSize:         100,
    MaxWaitTimeSecs: 300,
})
// Events: EventCallEnqueued, EventCallDequeued, EventCallLeftQueue

_, err = client.LeaveQueue(ctx, callControlID, "", "")
```

### AI Assistant on a call

```go
resp, err := client.StartAIAssistant(ctx, callControlID, telnyx.StartAIAssistantRequest{
    Assistant: telnyx.CallAssistantRequest{ID: "asst_abc123"},
    Greeting:  "Hi, how can I help you today?",
    MessageHistory: []telnyx.MessageHistoryEntry{
        {Role: "system", Content: "The caller is a premium subscriber."},
    },
    SendMessageHistoryUpdates: true,
    InterruptionSettings: &telnyx.InterruptionSettings{
        InterruptionSensitivity: "high",
        HoldAfterInterrupt:      false,
    },
})
fmt.Println(resp.Data.ConversationID)
// Events: EventCallAIGatherEnded, EventCallConversationEnded,
//         EventCallAIGatherPartialResults, EventCallAIGatherMessageHistoryUpdated,
//         EventCallConversationInsightsGenerated

_, err = client.StopAIAssistant(ctx, callControlID, telnyx.StopAIAssistantRequest{})
```

### Update Client State

```go
_, err := client.UpdateClientState(ctx, callControlID,
    base64.StdEncoding.EncodeToString([]byte(`{"step":"2"}`)),
)
```

---

## Conversational AI

### Create an Assistant

```go
asst, err := client.CreateAssistant(ctx, telnyx.CreateAssistantRequest{
    Name:         "Support Bot",
    Instructions: "You are a helpful customer support agent.",
    Model:        "gpt-4o",
    Greeting:     "Hello! How can I assist you today?",

    EnabledFeatures: []string{"telephony", "messaging"},

    VoiceSettings: &telnyx.VoiceSettings{
        Provider: "Telnyx",
        VoiceID:  "KokoroTTS.af",
    },

    TelephonySettings: &telnyx.TelephonySettings{
        NoiseSuppression: true,
        SilenceTimeoutMS: 3000,
        MaxCallDuration:  1800,
    },

    MessagingSettings: &telnyx.MessagingSettings{
        MessagingProfileID:  "profile-uuid",
        InactivityTimeoutMS: 60000,
    },

    Transcription: &telnyx.AssistantTranscriptionConfig{
        Language: "en-US",
    },

    InterruptionSettings: &telnyx.InterruptionSettings{
        InterruptionSensitivity: "medium",
    },

    FallbackConfig: &telnyx.FallbackConfig{Model: "gpt-4o-mini"},

    Tags:             []string{"support", "production"},
    DynamicVariables: map[string]any{"company_name": "Acme Inc."},
})
fmt.Println(asst.ID, asst.VersionID)
```

### Get / Update / Delete

```go
asst, err := client.GetAssistant(ctx, "asst_abc123")

asst, err = client.UpdateAssistant(ctx, "asst_abc123", telnyx.UpdateAssistantRequest{
    Instructions: "Updated instructions with new product information.",
    Model:        "gpt-4o-mini",
    Tags:         []string{"updated"},
})

err = client.DeleteAssistant(ctx, "asst_abc123")
```

### List Assistants

```go
assistants, err := client.ListAssistants(ctx, telnyx.ListAssistantsParams{
    PageNumber: 1,
    PageSize:   20,
})
```

### Chat (beta)

Send a message to an assistant outside of a call.

```go
resp, err := client.Chat(ctx, "asst_abc123", telnyx.ChatRequest{
    Messages: []telnyx.ChatMessage{
        {Role: "user", Content: "What are your business hours?"},
    },
})
fmt.Println(resp.Data.Message.Content)
```

---

## Webhooks

### Setup

Get your public key from the Telnyx portal under **API Keys → Webhook signing secret** (base64-encoded ED25519 key).

```go
wh, err := telnyx.NewWebhookHandler("BASE64_PUBLIC_KEY")
if err != nil {
    log.Fatal(err)
}

http.Handle("/webhooks", wh)
log.Fatal(http.ListenAndServe(":8080", nil))
```

Every incoming request is verified with ED25519. Requests with an invalid or missing `Telnyx-Signature-Ed25519` header receive `401 Unauthorized`. Webhooks older than 5 minutes are rejected to prevent replay attacks.

### Registering handlers

```go
// Multiple handlers for the same event type are called in registration order.
wh.On(telnyx.EventCallInitiated, func(ctx context.Context, event telnyx.Event, payload any) error {
    p := payload.(*telnyx.CallInitiatedPayload)
    fmt.Printf("incoming call: %s → %s direction=%s\n", p.From, p.To, p.Direction)
    return nil
})

wh.On(telnyx.EventCallAnswered, func(ctx context.Context, event telnyx.Event, payload any) error {
    p := payload.(*telnyx.CallAnsweredPayload)
    fmt.Println("answered:", p.CallControlID)
    return nil
})

wh.On(telnyx.EventCallGatherEnded, func(ctx context.Context, event telnyx.Event, payload any) error {
    p := payload.(*telnyx.CallGatherEndedPayload)
    fmt.Printf("gathered: %q status=%s\n", p.Digits, p.Status)
    return nil
})

wh.On(telnyx.EventCallRecordingSaved, func(ctx context.Context, event telnyx.Event, payload any) error {
    p := payload.(*telnyx.CallRecordingSavedPayload)
    fmt.Printf("recording %s mp3=%s\n", p.RecordingID, p.RecordingURL.MP3)
    return nil
})

wh.On(telnyx.EventCallTranscription, func(ctx context.Context, event telnyx.Event, payload any) error {
    p := payload.(*telnyx.CallTranscriptionPayload)
    if p.TranscriptionData.IsFinal {
        fmt.Printf("[%s] %s\n", p.TranscriptionData.TranscriptFrom, p.TranscriptionData.Transcript)
    }
    return nil
})

wh.On(telnyx.EventCallConversationEnded, func(ctx context.Context, event telnyx.Event, payload any) error {
    p := payload.(*telnyx.CallConversationEndedPayload)
    fmt.Printf("AI conversation ended: reason=%s messages=%d\n", p.EndReason, len(p.MessageHistory))
    return nil
})

// Catch-all for any event without a registered handler
wh.OnFallback(func(ctx context.Context, event telnyx.Event, payload any) error {
    fmt.Printf("unhandled event: %s\n", event.Data.EventType)
    return nil
})
```

### All 46 event types

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

### Common payload fields

All call event payloads embed `BasePayload`:

```go
type BasePayload struct {
    CallControlID string // use for subsequent call commands
    CallLegID     string
    CallSessionID string
    ConnectionID  string
    ClientState   string // your passthrough state, base64-encoded
    From          string
    To            string
}
```

All conference event payloads embed `ConferenceBasePayload`:

```go
type ConferenceBasePayload struct {
    ConferenceID string
    ConnectionID string
    OccurredAt   string
    ClientState  string
}
```

---

## Messaging

### Send SMS

```go
msg, err := client.SendMessage(ctx, telnyx.SendMessageRequest{
    From: "+15550001111",
    To:   "+15550002222",
    Text: "Hello from Telnyx!",
})
fmt.Println(msg.ID, msg.To[0].Status) // queued | sent | delivered | delivery_failed
```

### Send MMS

```go
msg, err := client.SendMessage(ctx, telnyx.SendMessageRequest{
    From:      "+15550001111",
    To:        "+15550002222",
    Text:      "Check this out!",
    MediaURLs: []string{"https://example.com/image.jpg"},
    Type:      "MMS",
})
```

### Send Group MMS

```go
msg, err := client.SendMessage(ctx, telnyx.SendMessageRequest{
    From: "+15550001111",
    To:   []string{"+15550002222", "+15550003333"},
    Text: "Group message!",
    Type: "MMS",
})
```

### Get a message

```go
msg, err := client.GetMessage(ctx, "msg-uuid")
fmt.Println(msg.Direction, msg.To[0].Status)
```

### List messages

```go
messages, meta, err := client.ListMessages(ctx, telnyx.ListMessagesParams{
    PageNumber:     1,
    PageSize:       25,
    Direction:      "outbound",         // inbound | outbound
    DateRangeStart: "2024-01-01T00:00:00Z",
    DateRangeEnd:   "2024-01-31T23:59:59Z",
})
fmt.Printf("%d messages, %d total\n", len(messages), meta.TotalResults)
```

### Schedule a message

```go
msg, err := client.ScheduleMessage(ctx, telnyx.ScheduleMessageRequest{
    SendMessageRequest: telnyx.SendMessageRequest{
        From: "+15550001111",
        To:   "+15550002222",
        Text: "Scheduled reminder",
    },
    SendAt: "2024-06-01T09:00:00Z", // ISO 8601
})

// Cancel before it sends
err = client.CancelScheduledMessage(ctx, msg.ID)
```

### Messaging profiles

```go
profile, err := client.CreateMessagingProfile(ctx, telnyx.CreateMessagingProfileRequest{
    Name:                    "My Profile",
    WebhookURL:              "https://yourserver.com/webhooks/messaging",
    WebhookAPIVersion:       "2",
    WhitelistedDestinations: []string{"US", "CA"},
    MMSFallBackToSMS:        true,
    SmartEncoding:           true,
})

profile, err = client.GetMessagingProfile(ctx, profile.ID)

enabled := false
profile, err = client.UpdateMessagingProfile(ctx, profile.ID, telnyx.UpdateMessagingProfileRequest{
    Enabled: &enabled,
})

profiles, meta, err := client.ListMessagingProfiles(ctx, 1, 20)

err = client.DeleteMessagingProfile(ctx, profile.ID)
```

---

## Conferences

### Create

```go
conf, err := client.CreateConference(ctx, telnyx.CreateConferenceRequest{
    CallControlID:   callControlID, // required: call that creates the conf
    Name:            "Team Standup",
    BeepEnabled:     "on_enter", // always | never | on_enter | on_exit
    MaxParticipants: 50,
    DurationMinutes: 60,
    Region:          "US", // US | Europe | Australia | Middle East
})
fmt.Println(conf.ID, conf.Status) // init | in_progress | completed
```

### Get / list

```go
conf, err := client.GetConference(ctx, conf.ID)

conferences, meta, err := client.ListConferences(ctx, 1, 20, "US") // region filter optional
```

### Join

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

### List participants

```go
participants, meta, err := client.ListConferenceParticipants(ctx, conf.ID, 1, 50)
for _, p := range participants {
    fmt.Printf("%s muted=%v on_hold=%v status=%s\n",
        p.CallControlID, p.Muted, p.OnHold, p.Status)
}
```

### Mute / Unmute

```go
// Pass specific IDs, or empty slice to mute all
_, err := client.MuteConferenceParticipants(ctx, conf.ID, []string{callControlID})
_, err = client.UnmuteConferenceParticipants(ctx, conf.ID, []string{})
```

### Hold / Unhold

```go
_, err := client.HoldConferenceParticipants(ctx, conf.ID, telnyx.HoldConferenceParticipantsRequest{
    CallControlIDs: []string{callControlID},
    AudioURL:       "https://example.com/hold.mp3",
})

_, err = client.UnholdConferenceParticipants(ctx, conf.ID, []string{callControlID})
```

### Speak / play audio

```go
_, err := client.ConferenceSpeak(ctx, conf.ID, telnyx.ConferenceSpeakRequest{
    Payload:  "The meeting will begin shortly.",
    Voice:    "Telnyx.KokoroTTS.af",
    Language: "en-US",
})

_, err = client.ConferencePlayAudio(ctx, conf.ID, telnyx.ConferencePlayAudioRequest{
    AudioURL: "https://example.com/intro.mp3",
})

_, err = client.ConferenceStopAudio(ctx, conf.ID, []string{}) // empty = stop for all
```

### Record

```go
_, err := client.ConferenceRecordStart(ctx, conf.ID, telnyx.ConferenceRecordStartRequest{
    Format:   "mp3",   // mp3 | wav
    Channels: "single", // single | dual
})
// Saved recording fires EventConferenceRecordingSaved

_, err = client.ConferenceRecordStop(ctx, conf.ID)
```

### Kick / update participant

```go
_, err := client.KickConferenceParticipants(ctx, conf.ID, []string{callControlID})

t := true
_, err = client.UpdateConferenceParticipant(ctx, conf.ID, telnyx.UpdateConferenceParticipantRequest{
    CallControlID:           callControlID,
    SoftEndConferenceOnExit: &t,
})
```

---

## Phone Numbers

### Search available numbers

```go
numbers, meta, err := client.SearchAvailableNumbers(ctx, telnyx.SearchAvailableNumbersParams{
    FilterCountryCode:    "US",
    FilterPhoneNumberType: "local", // local | toll_free | mobile | national
    FilterFeatures:       "sms,voice",
    FilterAreaCode:       "+1212",
    FilterLimit:          10,
    FilterQuickship:      true,
})
for _, n := range numbers {
    fmt.Println(n.PhoneNumber, n.Costs[0].Amount, n.Costs[0].Currency)
}
```

### Order numbers

```go
order, err := client.OrderPhoneNumbers(ctx, telnyx.NumberOrderRequest{
    PhoneNumbers: []telnyx.NumberOrderItem{
        {PhoneNumber: "+12125550100"},
        {PhoneNumber: "+12125550101"},
    },
    ConnectionID:       "conn-id",
    MessagingProfileID: "profile-id",
    BillingGroupID:     "group-id",
})
fmt.Println(order.ID, order.Status) // pending | success | failure

order, err = client.GetNumberOrder(ctx, order.ID)
```

### List / get / update / release

```go
numbers, meta, err := client.ListPhoneNumbers(ctx, telnyx.ListPhoneNumbersParams{
    PageNumber:  1,
    PageSize:    50,
    FilterStatus: "active", // active | pending | inactive
    Sort:        "purchased_at",
})

number, err := client.GetPhoneNumber(ctx, numberID)
fmt.Println(number.PhoneNumber, number.PhoneNumberType, number.Status)

number, err = client.UpdatePhoneNumber(ctx, numberID, telnyx.UpdatePhoneNumberRequest{
    ConnectionID:       "conn-id",
    MessagingProfileID: "profile-id",
    BillingGroupID:     "group-id",
    Tags:               []string{"production"},
    CustomerReference:  "acct-456",
})

err = client.DeletePhoneNumber(ctx, numberID)
```

---

## Number Lookup

Carrier, CNAM, and LERG/portability are fetched concurrently and merged into one result.

```go
result, err := client.LookupNumber(ctx, "+15550001111")

fmt.Println(result.PhoneNumber)
fmt.Println(result.CountryCode, result.NationalFormat)

if result.Carrier != nil {
    fmt.Println(result.Carrier.Name)              // AT&T
    fmt.Println(result.Carrier.Type)              // mobile | voip | landline | toll free
    fmt.Println(result.Carrier.NormalizedCarrier)
}

if result.CallerName != nil {
    fmt.Println(result.CallerName.CallerName)     // John Doe
}

if result.Portability != nil {
    p := result.Portability
    fmt.Println(p.LRN)         // Local routing number
    fmt.Println(p.PortedStatus) // Y | N
    fmt.Println(p.OCN)
    fmt.Println(p.City, p.State)
}
```

---

## Recordings & Transcriptions

### List recordings

```go
recordings, meta, err := client.ListRecordings(ctx, telnyx.ListRecordingsParams{
    PageNumber:          1,
    PageSize:            25,
    FilterCallLegID:     "leg-uuid",
    FilterCallSessionID: "session-uuid",
    FilterStatus:        "completed", // processing | completed | failed
    FilterSource:        "call_leg",  // call_leg | conference
})
```

### Get / delete a recording

```go
rec, err := client.GetRecording(ctx, "rec-uuid")
fmt.Println(rec.Urls.MP3)        // download URL
fmt.Println(rec.DurationMillis)  // duration in ms
fmt.Println(rec.Format)          // mp3 | wav
fmt.Println(rec.Channels)        // single | dual

err = client.DeleteRecording(ctx, "rec-uuid")
```

### Transcriptions

```go
// Get a single transcription
t, err := client.GetTranscription(ctx, "transcription-uuid")
fmt.Println(t.TranscriptionText)
for _, speaker := range t.Speakers {
    for _, seg := range speaker.Segments {
        fmt.Printf("[%.2fs–%.2fs %.0f%%] %s\n",
            seg.StartTime, seg.EndTime, seg.Confidence*100, seg.Text)
    }
}

// List transcriptions, optionally filtered by recording
transcripts, meta, err := client.ListTranscriptions(ctx, "rec-uuid", 1, 20)
```

---

## Media Storage

Pre-upload audio files and reference them by `media_name` in call commands.

### Upload from URL

```go
media, err := client.UploadMediaFromURL(ctx,
    "https://example.com/hold.mp3", // mediaURL (≤20 MB)
    "hold-music",                    // mediaName (unique key)
    0,                               // ttlSecs (0 = default 2 days)
)
fmt.Println(media.MediaName, media.ContentType, media.ExpiresAt)
```

### Upload a local file

```go
f, _ := os.Open("audio.mp3")
defer f.Close()

media, err := client.UploadMediaFile(ctx,
    "my-audio",   // mediaName
    0,            // ttlSecs
    "audio.mp3",  // filename (used for the multipart field)
    f,            // io.Reader
)
```

### Use in call commands

```go
// Play the pre-uploaded file
_, err := client.PlayAudio(ctx, callControlID, telnyx.PlayAudioRequest{
    MediaName: "hold-music",
    Loop:      "infinity",
})

// Or use it in a dial request
call, err := client.Dial(ctx, telnyx.DialRequest{
    MediaName: "hold-music", // played while ringing
    // ...
})
```

### List / get / delete media

```go
items, meta, err := client.ListMedia(ctx, 1, 25, []string{"audio/mpeg", "audio/wav"})

item, err := client.GetMedia(ctx, "hold-music")

err = client.DeleteMedia(ctx, "hold-music")
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
    Quality:      "high",       // normal | high | very-high
    StoreMedia:   true,         // store a copy of the sent PDF
    StorePreview: true,         // generate a preview image
    T38Enabled:   nil,          // nil = Telnyx default
    Monochrome:   false,
    WebhookURL:   "https://yourserver.com/webhooks/fax",
    ClientState:  "base64state==",
})
fmt.Println(fax.ID, fax.Status)
// Webhook events: fax.queued, fax.media.processed,
//                 fax.sending.started, fax.delivered, fax.failed
```

### Get / list / cancel / delete

```go
fax, err := client.GetFax(ctx, fax.ID)

faxes, meta, err := client.ListFaxes(ctx, telnyx.ListFaxesParams{
    PageNumber: 1,
    PageSize:   25,
    FilterTo:   "+15550002222",
    FilterFrom: "+15550001111",
})

fax, err = client.CancelFax(ctx, fax.ID)
fax, err = client.RefreshFax(ctx, fax.ID)  // re-fetch updated URLs
err = client.DeleteFax(ctx, fax.ID)
```

### Fax applications

```go
app, err := client.CreateFaxApplication(ctx, telnyx.CreateFaxApplicationRequest{
    ApplicationName:         "My Fax App",
    WebhookEventURL:         "https://yourserver.com/webhooks/fax",
    WebhookEventFailoverURL: "https://backup.yourserver.com/webhooks/fax",
    Active:                  true,
})

app, err = client.GetFaxApplication(ctx, app.ID)

t := true
app, err = client.UpdateFaxApplication(ctx, app.ID, telnyx.UpdateFaxApplicationRequest{
    Active:            &t,
    FaxEmailRecipient: "faxes@example.com",
})

apps, meta, err := client.ListFaxApplications(ctx, 1, 20)
err = client.DeleteFaxApplication(ctx, app.ID)
```

---

## FQDN Connections

### Connections

```go
conn, err := client.CreateFQDNConnection(ctx, telnyx.CreateFQDNConnectionRequest{
    ConnectionName:    "My FQDN Connection",
    TransportProtocol: "UDP",            // UDP | TCP | TLS
    DTMFType:          "RFC 2833",       // RFC 2833 | Inband | SIP INFO
    NoiseSuppression:  "both",           // inbound | outbound | both | disabled
    WebhookEventURL:   "https://yourserver.com/webhooks",
    WebhookAPIVersion: "2",

    Inbound: &telnyx.FQDNInbound{
        SIPRegion:        "US",          // US | Europe | Australia
        ChannelLimit:     &limit,
        SHAKENSTIREnabled: true,
    },
    Outbound: &telnyx.FQDNOutbound{
        OutboundVoiceProfileID: "profile-id",
        T38ReinviteSource:      "telnyx",
    },
    RTCPSettings: &telnyx.RTCPSettings{
        Port:           "rtp+1",
        CaptureEnabled: true,
    },
})

conn, err = client.GetFQDNConnection(ctx, conn.ID)

conn, err = client.UpdateFQDNConnection(ctx, conn.ID, telnyx.CreateFQDNConnectionRequest{
    ConnectionName: "Renamed Connection",
})

conns, meta, err := client.ListFQDNConnections(ctx, telnyx.ListFQDNConnectionsParams{
    FilterNameContains: "my",
    PageNumber:         1,
    PageSize:           20,
    Sort:               "created_at", // created_at | connection_name | active; prefix - for desc
})

err = client.DeleteFQDNConnection(ctx, conn.ID)
```

### FQDNs

```go
fqdn, err := client.CreateFQDN(ctx, telnyx.CreateFQDNRequest{
    ConnectionID:  conn.ID,
    FQDN:          "sip.example.com",
    DNSRecordType: "a",  // a | cname
    Port:          5060,
})

fqdn, err = client.GetFQDN(ctx, fqdn.ID)
fqdn, err = client.UpdateFQDN(ctx, fqdn.ID, telnyx.CreateFQDNRequest{Port: 5061})
fqdns, meta, err := client.ListFQDNs(ctx, conn.ID, 1, 50)
err = client.DeleteFQDN(ctx, fqdn.ID)
```

---

## Verify / OTP

### Send a verification code

```go
// SMS
v, err := client.SendVerificationSMS(ctx, telnyx.CreateVerificationSMSRequest{
    PhoneNumber:     "+15550001111",
    VerifyProfileID: "profile-uuid",
    TimeoutSecs:     300, // code expires after this many seconds
})
fmt.Println(v.ID, v.Status) // pending | accepted | invalid | expired | error

// Voice call
v, err = client.SendVerificationCall(ctx, telnyx.CreateVerificationCallRequest{
    PhoneNumber:     "+15550001111",
    VerifyProfileID: "profile-uuid",
    Extension:       "123", // optional DTMF digits after answer
})

// Flashcall (missed call — code is in caller ID)
v, err = client.SendVerificationFlashcall(ctx, telnyx.CreateVerificationFlashcallRequest{
    PhoneNumber:     "+15550001111",
    VerifyProfileID: "profile-uuid",
})
```

### Submit a code

```go
// By verification ID
resp, err := client.VerifyCodeByID(ctx, v.ID, telnyx.VerifyCodeByIDRequest{
    Code: "123456",
})
fmt.Println(resp.ResponseCode) // accepted | rejected

// By phone number (simpler flow — no need to store the verification ID)
resp, err = client.VerifyCodeByPhone(ctx, "+15550001111", telnyx.VerifyCodeByPhoneRequest{
    Code:            "123456",
    VerifyProfileID: "profile-uuid",
})
```

### Get / list verifications

```go
v, err := client.GetVerification(ctx, verificationID)

verifications, meta, err := client.GetVerificationsByPhone(ctx, "+15550001111", 1, 20)
```

### Verify profiles

```go
profile, err := client.CreateVerifyProfile(ctx, telnyx.CreateVerifyProfileRequest{
    Name:                   "My App OTP",
    Language:               "en-US",
    WebhookURL:             "https://yourserver.com/webhooks/verify",
    DailySpendLimitEnabled: true,
    DailySpendLimit:        10.00,
})

profile, err = client.GetVerifyProfile(ctx, profile.ID)
profiles, meta, err := client.ListVerifyProfiles(ctx, 1, 20)

profile, err = client.UpdateVerifyProfile(ctx, profile.ID, telnyx.CreateVerifyProfileRequest{
    DailySpendLimit: 20.00,
})

err = client.DeleteVerifyProfile(ctx, profile.ID)
```

---

## SIP Connections

### Credential connections (username/password auth)

```go
conn, err := client.CreateCredentialConnection(ctx, telnyx.CreateCredentialConnectionRequest{
    UserName:       "alice",
    Password:       "s3cr3t",
    ConnectionName: "My SIP Trunk",

    DTMFType:                      "RFC 2833",
    NoiseSuppression:              "both",
    OnnetT38PassthroughEnabled:    true,
    EncodeContactHeaderEnabled:    true,
    WebhookEventURL:               "https://yourserver.com/webhooks",
    WebhookAPIVersion:             "2",

    Inbound: &telnyx.ConnectionInbound{
        Codecs:       []string{"G722", "PCMU", "PCMA"},
        ChannelLimit: 50,
    },
    Outbound: &telnyx.ConnectionOutbound{
        OutboundVoiceProfileID: "profile-id",
    },
    JitterBuffer: &telnyx.JitterBuffer{
        Enabled: true,
        MinMS:   40,
        MaxMS:   200,
    },
    RTCPSettings: &telnyx.RTCPSettings{
        CaptureEnabled:  true,
        ReportFrequency: 5,
    },
})

conn, err = client.GetCredentialConnection(ctx, conn.ID)

conn, err = client.UpdateCredentialConnection(ctx, conn.ID, telnyx.UpdateCredentialConnectionRequest{
    Password: "newpassword",
    Tags:     []string{"updated"},
})

conns, meta, err := client.ListCredentialConnections(ctx, 1, 20)
err = client.DeleteCredentialConnection(ctx, conn.ID)
```

### IP connections (IP-authenticated)

```go
ipConn, err := client.CreateIPConnection(ctx, telnyx.CreateIPConnectionRequest{
    ConnectionName:    "My IP Trunk",
    TransportProtocol: "UDP", // UDP | TCP | TLS
    NoiseSuppression:  "inbound",
    WebhookEventURL:   "https://yourserver.com/webhooks",
})

ipConn, err = client.GetIPConnection(ctx, ipConn.ID)

ipConn, err = client.UpdateIPConnection(ctx, ipConn.ID, telnyx.UpdateIPConnectionRequest{
    TransportProtocol: "TLS",
})

ipConns, meta, err := client.ListIPConnections(ctx, 1, 20)
err = client.DeleteIPConnection(ctx, ipConn.ID)
```

---

## Billing Groups

Organise numbers and usage under named cost centres.

```go
group, err := client.CreateBillingGroup(ctx, "Team A")
fmt.Println(group.ID, group.OrganizationID)

group, err = client.GetBillingGroup(ctx, group.ID)

group, err = client.UpdateBillingGroup(ctx, group.ID, "Team B")

groups, meta, err := client.ListBillingGroups(ctx, 1, 20)

err = client.DeleteBillingGroup(ctx, group.ID)
```

---

## CDR Reports

CDR reports are generated asynchronously. Create the report then poll until complete.

```go
report, err := client.CreateCDRReport(ctx, telnyx.CreateCDRReportRequest{
    StartTime:  "2024-01-01T00:00:00Z",
    EndTime:    "2024-01-31T23:59:59Z",
    Timezone:   "America/New_York",
    ReportName: "January 2024",

    // Filter by direction (1=Inbound, 2=Outbound)
    CallTypes: []int{telnyx.CDRCallTypeInbound, telnyx.CDRCallTypeOutbound},

    // Filter by completeness (1=Complete, 2=Incomplete, 3=Errors)
    RecordTypes: []int{telnyx.CDRRecordTypeComplete},

    // Filter by source
    Source: "call-control", // calls | call-control | fax-api | webrtc
})

// Poll until complete
for report.Status == telnyx.CDRStatusPending {
    time.Sleep(5 * time.Second)
    report, err = client.GetCDRReport(ctx, report.ID)
}
if report.Status == telnyx.CDRStatusComplete {
    fmt.Println("download:", report.ReportURL)
}

reports, err := client.ListCDRReports(ctx)
err = client.DeleteCDRReport(ctx, report.ID)
```

CDR status constants: `CDRStatusPending` (1), `CDRStatusComplete` (2), `CDRStatusFailed` (3), `CDRStatusExpired` (4).

---

## Audit Logs

```go
logs, meta, err := client.ListAuditLogs(ctx, telnyx.ListAuditLogsParams{
    PageNumber:    1,
    PageSize:      50,
    CreatedAfter:  "2024-01-01T00:00:00Z",
    CreatedBefore: "2024-12-31T23:59:59Z",
    Sort:          "desc", // asc | desc
})

for _, l := range logs {
    fmt.Printf("[%s] %s made %s change on %s\n",
        l.CreatedAt, l.ChangeMadeBy, l.ChangeType, l.ResourceID)
    for _, c := range l.Changes {
        fmt.Printf("  %s: %v → %v\n", c.Field, c.From, c.To)
    }
}
```

`ChangeMadeBy` values: `telnyx` | `account_manager` | `account_owner` | `organization_member`

---

## Webhook Deliveries

Query the delivery log for any webhook event sent by Telnyx.

```go
// List recent failed deliveries
deliveries, meta, err := client.ListWebhookDeliveries(ctx, telnyx.ListWebhookDeliveriesParams{
    FilterStatus:       "failed",               // delivered | failed
    FilterEventType:    "call.hangup",           // comma-separated for multiple
    FilterWebhook:      "yourserver.com",        // substring match on webhook URL
    FilterStartedAtGTE: "2024-01-01T00:00:00Z",
    FilterStartedAtLTE: "2024-01-31T23:59:59Z",
    PageNumber:         1,
    PageSize:           20,
})

for _, d := range deliveries {
    last := d.Attempts[len(d.Attempts)-1]
    fmt.Printf("[%s] %s → HTTP %d\n", d.Webhook.EventType, d.Status, last.Response.Status)
    fmt.Println("  request URL:", last.Request.URL)
    fmt.Println("  response body:", last.Response.Body)
}

// Full detail for a single delivery
delivery, err := client.GetWebhookDelivery(ctx, "delivery-uuid")
```

---

## Error Handling

All API errors return `*telnyx.APIError`:

```go
call, err := client.Dial(ctx, req)
if err != nil {
    var apiErr *telnyx.APIError
    if errors.As(err, &apiErr) {
        fmt.Printf("HTTP %d: %s – %s\n",
            apiErr.StatusCode,
            apiErr.Errors[0].Title,
            apiErr.Errors[0].Detail)
    }
    return err
}
```

Webhook signature failures return `*telnyx.ErrSignatureVerification`:

```go
// (returned internally by WebhookHandler; surface it from your own middleware if needed)
var sigErr *telnyx.ErrSignatureVerification
if errors.As(err, &sigErr) {
    fmt.Println("bad webhook:", sigErr.Reason)
}
```

---

## Examples

Working example programs are in the [`examples/`](./examples) directory.

| Example | Description |
|---|---|
| [`examples/ivr/`](./examples/ivr/) | Inbound IVR: answers calls, plays a menu, gathers DTMF, routes to sales/support/voicemail |
| [`examples/outbound_call/`](./examples/outbound_call/) | Outbound dialer with premium AMD: speaks to humans, leaves voicemail after machine beep |
| [`examples/conference/`](./examples/conference/) | Conference bridge: creates a room on the first caller, joins all subsequent callers |
| [`examples/send_sms/`](./examples/send_sms/) | Sends an SMS and listens for delivery webhooks |
| [`examples/number_lookup/`](./examples/number_lookup/) | CLI: looks up carrier, CNAM, and portability for a phone number |

### Run an example

```bash
TELNYX_API_KEY=... TELNYX_PUBLIC_KEY=... go run ./examples/ivr
```

---

## License

MIT
