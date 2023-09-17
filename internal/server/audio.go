package server

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"poker"
	"poker/internal"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/polly"
	ptypes "github.com/aws/aws-sdk-go-v2/service/polly/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gorilla/mux"
)

func (s *server) handleGetDashboardTimerLevelAudio(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	entry := s.logger.WithContext(ctx)

	user := internal.UserFromContext(ctx)

	entry = entry.WithField("user_id", user.ID)

	vars := mux.Vars(r)

	timerID, ok := vars["timerID"]
	if !ok {
		entry.Error("var timerID missing from request context")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	entry = entry.WithField("timerID", timerID)

	levelID, ok := vars["levelID"]
	if !ok {
		entry.Error("var levelID missing from request context")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	entry = entry.WithField("levelID", levelID)

	actionStr, ok := vars["action"]
	if !ok {
		entry.Error("var action missing from request context")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	entry = entry.WithField("action", actionStr)

	action := _action(actionStr)
	if !action.valid() {
		entry.Error("action is invalid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	timer, err := s.timerRepo.Timer(ctx, timerID)
	if err != nil {
		entry.WithError(err).Error("failed to fetch timer")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if timer == nil {
		entry.Error("timer not found, returning not found page")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if timer.UserID != user.ID {
		entry.Error("timer is not owned by authenticated user")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var level *poker.TimerLevel
	for _, lvl := range timer.Levels {
		if lvl.ID != levelID {
			continue
		}
		level = lvl
		break
	}
	if level == nil {
		entry.Error("timer does not contain requested level")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	buffer, contentType, err := s.generateAndSaveAudio(ctx, level, action)
	if err != nil {
		entry.WithError(err).Error("failed to generate/save audio file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", contentType)
	_, _ = buffer.WriteTo(w)

}

func (s *server) generateAndSaveAudio(ctx context.Context, level *poker.TimerLevel, action _action) (io.WriterTo, string, error) {

	var objectKey = fmt.Sprintf("%s-%s.mp3", level.AudioS3Key(), action)

	entry := s.logger.WithField("objectKey", objectKey).WithContext(ctx)

	objectOutput, err := s.s3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.audioBucket),
		Key:    aws.String(objectKey),
	})
	if err == nil {
		entry.Info("cached audio file found, returning")
		defer objectOutput.Body.Close()
		var buffer = new(bytes.Buffer)
		_, _ = buffer.ReadFrom(objectOutput.Body)
		return buffer, aws.ToString(objectOutput.ContentType), nil
	}

	text := generateSpeechText(action, level)
	entry.WithField("text", text).Info("generating audio file for text")

	synthesizeSpeechOutput, err := s.polly.SynthesizeSpeech(ctx, &polly.SynthesizeSpeechInput{
		Engine:       ptypes.EngineNeural,
		OutputFormat: ptypes.OutputFormatMp3,
		LanguageCode: ptypes.LanguageCodeEnUs,
		Text:         aws.String(text),
		VoiceId:      ptypes.VoiceIdStephen,
	})
	if err != nil {
		entry.WithError(err).Error("failed to synthesize speech")
		return nil, "", fmt.Errorf("failed to synthesize speech: %w", err)
	}

	defer synthesizeSpeechOutput.AudioStream.Close()

	var buffer = new(bytes.Buffer)
	_, _ = buffer.ReadFrom(synthesizeSpeechOutput.AudioStream)

	var buffer2 = bytes.NewBuffer(bytes.Clone(buffer.Bytes()))

	_, err = s.s3.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.audioBucket),
		Key:         aws.String(objectKey),
		Body:        buffer,
		ContentType: synthesizeSpeechOutput.ContentType,
	})
	if err != nil {
		entry.WithError(err).Error("failed to put audio file in S3")
		return nil, "", fmt.Errorf("failed to put audio file in S3: %w", err)
	}

	fmt.Println(buffer.Len(), buffer2.Len())
	return buffer2, aws.ToString(synthesizeSpeechOutput.ContentType), nil
}

type _action string

func (a _action) String() string {
	return string(a)
}

const (
	continueAction _action = "continue"
	playAction     _action = "play"
)

var allActions = []_action{continueAction, playAction}

func (a _action) valid() bool {
	for _, _a := range allActions {
		if a == _a {
			return true
		}
	}

	return true

}

func generateSpeechText(action _action, level *poker.TimerLevel) string {

	levelType := level.Type

	var prefixMap = map[poker.LevelType]map[_action]string{
		poker.LevelTypeBlind: {
			continueAction: "Blinds Up.",
			playAction:     "Let's Play Poker.",
		},
		poker.LevelTypeBreak: {
			continueAction: "It's break time.",
			playAction:     "It's break time.",
		},
	}

	var blindMap = map[poker.LevelType]string{
		poker.LevelTypeBlind: "The blinds are now %.0f/%.0f.",
	}

	var durationMap = map[poker.LevelType]string{
		poker.LevelTypeBlind: "This level will last for %.0f minutes",
		poker.LevelTypeBreak: "This break will last for %.0f minutes",
	}

	var prefix string = prefixMap[levelType][action]

	var blinds string
	var blindFmt = blindMap[levelType]
	if blindFmt != "" {
		blinds = fmt.Sprintf(blindFmt, level.SmallBlind, level.BigBlind)
	}

	var duration string
	var durationFmt = durationMap[levelType]
	if durationFmt != "" {
		duration = fmt.Sprintf(durationFmt, level.DurationMin)
	}

	return strings.Join([]string{prefix, blinds, duration}, " ")

}
