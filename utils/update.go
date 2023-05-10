package utils

type Update struct {
	Type   string `json:"type"`
	Object struct {
		Message struct {
			Date                  int    `json:"date"`
			FromID                int    `json:"from_id"`
			ID                    int    `json:"id"`
			Out                   int    `json:"out"`
			PeerID                int    `json:"peer_id"`
			Text                  string `json:"text"`
			ConversationMessageID int    `json:"conversation_message_id"`
			Payload               string `json:"payload"`
		} `json:"message"`
	} `json:"object"`
}
