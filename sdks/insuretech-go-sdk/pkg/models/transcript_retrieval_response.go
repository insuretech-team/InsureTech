package models


// TranscriptRetrievalResponse represents a transcript_retrieval_response
type TranscriptRetrievalResponse struct {
	Transcripts []*VoiceTranscript `json:"transcripts,omitempty"`
	Error *Error `json:"error,omitempty"`
}
