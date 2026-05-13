# telnyx-api

A Go library for the [Telnyx](https://telnyx.com) Voice API covering call control, webhook handling, and conversational AI assistants.

- Zero external dependencies — only the Go standard library
- ED25519 webhook signature verification with replay protection
- Typed payload structs for all 46 webhook event types
- Full call control command surface
- Conversational AI assistant CRUD and chat

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

---

## Call control

### Dial an outbound call

```go
call, err := client.Dial(ctx, telnyx.DialRequest{
    ConnectionID: "your-connection-id",
    To:           "+15550002222",
    From:         "+15550001111",
    TimeoutSecs:  30,
    WebhookURL:   "https://yourserver.com/webhooks",
})
// call.CallControlID is used for all subsequent commands
```

### Answer an inbound call

```go
_, err := client.Answer(ctx, callControlID, telnyx.AnswerRequest{
    WebhookURL: "https://yourserver.com/webhooks",
})
```

### Hangup

```go
_, err := client.Hangup(ctx, callControlID, telnyx.HangupRequest{})
```

### Reject an inbound call

```go
_, err := client.Reject(ctx, callControlID, telnyx.RejectRequest{
    Cause: "USER_BUSY",
})
```

### Speak text-to-speech

```go
_, err := client.Speak(ctx, callControlID, telnyx.SpeakRequest{
    Payload:      "Hello, please press 1 to continue.",
    Voice:        "Telnyx.KokoroTTS.af",
    PayloadType:  "text",
    ClientState:  base64.StdEncoding.EncodeToString([]byte(`{"step":"greeting"}`)),
})
```

### Play audio

```go
_, err := client.PlayAudio(ctx, callControlID, telnyx.PlayAudioRequest{
    AudioURL:   "https://example.com/hold-music.mp3",
    Loop:       "infinity",
    TargetLegs: "self",
})

// Stop playback
_, err = client.PlaybackStop(ctx, callControlID, "current", "", "")
```

### Gather DTMF digits

```go
_, err := client.Gather(ctx, callControlID, telnyx.GatherRequest{
    MinimumDigits:           1,
    MaximumDigits:           1,
    TimeoutMillis:           10000,
    InterDigitTimeoutMillis: 3000,
    TerminatingDigit:        "#",
    ValidDigits:             "123456789",
})
// Result arrives via EventCallGatherEnded webhook
```

### Send DTMF

```go
_, err := client.SendDTMF(ctx, callControlID, telnyx.SendDTMFRequest{
    Digits:     "1234",
    DurationMS: 100,
})
```

### Bridge two legs

```go
_, err := client.Bridge(ctx, callControlID, telnyx.BridgeRequest{
    CallControlID: otherCallControlID,
})
```

### Transfer

```go
_, err := client.Transfer(ctx, callControlID, telnyx.TransferRequest{
    To:          "+15550003333",
    TimeoutSecs: 30,
})
```

### Hold / Unhold

```go
_, err := client.Hold(ctx, callControlID, telnyx.HoldRequest{
    AudioURL: "https://example.com/hold-music.mp3",
})

_, err = client.Unhold(ctx, callControlID, "", "")
```

### Mute / Unmute

```go
_, err := client.Mute(ctx, callControlID, telnyx.MuteRequest{Mute: "self"})
_, err = client.Unmute(ctx, callControlID, telnyx.MuteRequest{Mute: "self"})
```

### Recording

```go
_, err := client.RecordStart(ctx, callControlID, telnyx.RecordStartRequest{
    Format:   "mp3",
    Channels: "dual",
    MaxLength: 3600,
    Transcription: true,
    TranscriptionEngine: "B",
})

_, err = client.RecordStop(ctx, callControlID, "", "")
_, err = client.RecordPause(ctx, callControlID, "", "")
_, err = client.RecordResume(ctx, callControlID, "", "")
// Saved recording arrives via EventCallRecordingSaved webhook
```

### Real-time transcription

```go
_, err := client.TranscriptionStart(ctx, callControlID, telnyx.TranscriptionStartRequest{
    TranscriptionEngine:   "B",
    TranscriptionLanguage: "en-US",
    InterimResults:        true,
})

_, err = client.TranscriptionStop(ctx, callControlID, "", "")
// Segments arrive via EventCallTranscription webhooks
```

### Media streaming (WebSocket)

```go
_, err := client.StreamingStart(ctx, callControlID, telnyx.StreamingStartRequest{
    StreamURL:   "wss://yourserver.com/media",
    StreamTrack: "both_tracks",
})

_, err = client.StreamingStop(ctx, callControlID, "", "")
```

### Media forking

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

### Answering machine detection

Configure via `DialRequest.AnsweringMachineDetection`:

```go
call, err := client.Dial(ctx, telnyx.DialRequest{
    ConnectionID:              "conn-id",
    To:                        "+15550002222",
    From:                      "+15550001111",
    AnsweringMachineDetection: "premium",
})
// Result arrives via EventCallMachinePremiumDetectionEnded webhook
```

### Retrieve a call

```go
call, err := client.GetCall(ctx, callControlID)
fmt.Println(call.IsAlive)
```

---

## Conversational AI

### Start an AI assistant on a call

```go
resp, err := client.StartAIAssistant(ctx, callControlID, telnyx.StartAIAssistantRequest{
    Assistant: telnyx.CallAssistantRequest{ID: "asst_abc123"},
    Greeting:  "Hi! How can I help you today?",
    SendMessageHistoryUpdates: true,
})
fmt.Println(resp.Data.ConversationID)
```

### Stop the AI assistant

```go
_, err := client.StopAIAssistant(ctx, callControlID, telnyx.StopAIAssistantRequest{})
```

### Assistant CRUD

```go
// Create
asst, err := client.CreateAssistant(ctx, telnyx.CreateAssistantRequest{
    Name:         "Support Bot",
    Instructions: "You are a helpful customer support agent.",
    Model:        "gpt-4o",
    Greeting:     "Hello! How can I assist you?",
    EnabledFeatures: []string{"telephony"},
    TelephonySettings: &telnyx.TelephonySettings{
        NoiseSuppression: true,
        SilenceTimeoutMS: 3000,
    },
})

// Get
asst, err = client.GetAssistant(ctx, asst.ID)

// Update
asst, err = client.UpdateAssistant(ctx, asst.ID, telnyx.UpdateAssistantRequest{
    Instructions: "Updated instructions.",
})

// List
assistants, err := client.ListAssistants(ctx, telnyx.ListAssistantsParams{
    PageNumber: 1,
    PageSize:   20,
})

// Delete
err = client.DeleteAssistant(ctx, asst.ID)
```

### Chat with an assistant (beta)

```go
resp, err := client.Chat(ctx, asst.ID, telnyx.ChatRequest{
    Messages: []telnyx.ChatMessage{
        {Role: "user", Content: "What are your business hours?"},
    },
})
fmt.Println(resp.Data.Message.Content)
```

---

## Webhooks

### Setup

Get your public key from the Telnyx portal under **API Keys → Webhook signing secret** (base64-encoded ED25519 public key).

```go
handler, err := telnyx.NewWebhookHandler("BASE64_PUBLIC_KEY")
if err != nil {
    log.Fatal(err)
}
```

### Register event handlers

```go
handler.On(telnyx.EventCallInitiated, func(ctx context.Context, event telnyx.Event, payload any) error {
    p := payload.(*telnyx.CallInitiatedPayload)
    fmt.Printf("Incoming call from %s to %s (id: %s)\n", p.From, p.To, p.CallControlID)
    return nil
})

handler.On(telnyx.EventCallGatherEnded, func(ctx context.Context, event telnyx.Event, payload any) error {
    p := payload.(*telnyx.CallGatherEndedPayload)
    fmt.Printf("Gathered digits: %s (status: %s)\n", p.Digits, p.Status)
    return nil
})

handler.On(telnyx.EventCallHangup, func(ctx context.Context, event telnyx.Event, payload any) error {
    p := payload.(*telnyx.CallHangupPayload)
    fmt.Printf("Call ended: %s\n", p.HangupCause)
    return nil
})

// Catch-all for unregistered event types
handler.OnFallback(func(ctx context.Context, event telnyx.Event, payload any) error {
    fmt.Printf("Unhandled event: %s\n", event.Data.EventType)
    return nil
})

http.Handle("/webhooks", handler)
log.Fatal(http.ListenAndServe(":8080", nil))
```

### Supported event types

| Constant | Event string |
|---|---|
| `EventCallInitiated` | `call.initiated` |
| `EventCallAnswered` | `call.answered` |
| `EventCallHangup` | `call.hangup` |
| `EventCallHold` | `call.hold` |
| `EventCallUnhold` | `call.unhold` |
| `EventCallBridged` | `call.bridged` |
| `EventCallPlaybackStarted` | `call.playback.started` |
| `EventCallPlaybackEnded` | `call.playback.ended` |
| `EventCallSpeakStarted` | `call.speak.started` |
| `EventCallSpeakEnded` | `call.speak.ended` |
| `EventCallDTMFReceived` | `call.dtmf.received` |
| `EventCallGatherEnded` | `call.gather.ended` |
| `EventCallRecordingSaved` | `call.recording.saved` |
| `EventCallRecordingError` | `call.recording.error` |
| `EventCallRecordingTranscriptionSaved` | `call.recording.transcription.saved` |
| `EventCallMachineDetectionEnded` | `call.machine.detection.ended` |
| `EventCallMachineGreetingEnded` | `call.machine.greeting.ended` |
| `EventCallMachinePremiumDetectionEnded` | `call.machine.premium.detection.ended` |
| `EventCallMachinePremiumGreetingEnded` | `call.machine.premium.greeting.ended` |
| `EventCallForkStarted` | `call.fork.started` |
| `EventCallForkStopped` | `call.fork.stopped` |
| `EventCallEnqueued` | `call.enqueued` |
| `EventCallDequeued` | `call.dequeued` |
| `EventCallLeftQueue` | `call.left_queue` |
| `EventCallTranscription` | `call.transcription` |
| `EventStreamingStarted` | `streaming.started` |
| `EventStreamingStopped` | `streaming.stopped` |
| `EventStreamingFailed` | `streaming.failed` |
| `EventCallReferStarted` | `call.refer.started` |
| `EventCallReferCompleted` | `call.refer.completed` |
| `EventCallReferFailed` | `call.refer.failed` |
| `EventCallSIPRECStarted` | `call.siprec.started` |
| `EventCallSIPRECStopped` | `call.siprec.stopped` |
| `EventCallSIPRECFailed` | `call.siprec.failed` |
| `EventCallAIGatherEnded` | `call.ai_gather.ended` |
| `EventCallAIGatherMessageHistoryUpdated` | `call.ai_gather.message_history.updated` |
| `EventCallAIGatherPartialResults` | `call.ai_gather.partial_results` |
| `EventCallConversationEnded` | `call.conversation.ended` |
| `EventCallConversationInsightsGenerated` | `call.conversation.insights_generated` |
| `EventCallDeepfakeDetectionResult` | `call.deepfake_detection.result` |
| `EventCallDeepfakeDetectionError` | `call.deepfake_detection.error` |
| `EventConferenceCreated` | `conference.created` |
| `EventConferenceEnded` | `conference.ended` |
| `EventConferenceFloorChanged` | `conference.floor_changed` |
| `EventConferenceParticipantJoined` | `conference.participant.joined` |
| `EventConferenceParticipantLeft` | `conference.participant.left` |
| `EventConferencePlaybackStarted` | `conference.playback.started` |
| `EventConferencePlaybackEnded` | `conference.playback.ended` |
| `EventConferenceSpeakStarted` | `conference.speak.started` |
| `EventConferenceSpeakEnded` | `conference.speak.ended` |
| `EventConferenceRecordingSaved` | `conference.recording.saved` |

### Signature verification

The handler verifies every request using ED25519 — requests with an invalid or missing signature receive `401 Unauthorized`. Webhooks older than 5 minutes are also rejected to prevent replay attacks.

---

## Error handling

All API errors return `*telnyx.APIError`:

```go
call, err := client.Dial(ctx, req)
if err != nil {
    var apiErr *telnyx.APIError
    if errors.As(err, &apiErr) {
        fmt.Printf("HTTP %d: %s\n", apiErr.StatusCode, apiErr.Errors[0].Detail)
    }
    return err
}
```

Webhook signature failures return `*telnyx.ErrSignatureVerification`.

---

## Client options

```go
// Custom HTTP client (e.g. with proxy or custom timeout)
client := telnyx.NewClient("API_KEY",
    telnyx.WithHTTPClient(&http.Client{Timeout: 10 * time.Second}),
)

// Override base URL (useful for testing)
client := telnyx.NewClient("API_KEY",
    telnyx.WithBaseURL("http://localhost:9090"),
)
```

---

## License

MIT
