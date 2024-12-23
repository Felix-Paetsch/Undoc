package tokenizer

type TokenizerState string

const (
	StateAfterNewLine           TokenizerState = "AFTER_NEWLINE"             // We are directly after a new line
	StateMiddleOfLine           TokenizerState = "MIDDLE_OF_LINE"            // We are at a generic position after a new line
	StateSingleHashtagLine      TokenizerState = "SINGLE_HASHTAG_LINE"       // We are in the line after a single hashtag
	StateAfterSingleHashtagLine TokenizerState = "AFTER_SINGLE_HASHTAG_LINE" // We are in the line after a single hashtag, before getting a "{"	StateEOF               TokenizerState = "EOF"                 // File end has been reached
	StateStartOfFile            TokenizerState = "SOF"                       // We are at file start
	StateAfterTagStart          TokenizerState = "AFTER_TAG_START"           // After the { for tags
	StateAfterTagEnd            TokenizerState = "AFTER_TAG_END"             // After the } for tags
	StateError                  TokenizerState = "ERROR"
)

type TokenizerStateManager struct {
	State TokenizerState
}

func NewTokenizerStateManager() *TokenizerStateManager {
	return &TokenizerStateManager{StateStartOfFile}
}

type Transition struct {
	FromState TokenizerState
	TokenType TokenType
	NextState TokenizerState
}

var transitionTable = []Transition{}

func (sm *TokenizerStateManager) processToken(token Token) {
	for _, tr := range transitionTable {
		if tr.FromState == sm.State && tr.TokenType == token.Type {
			sm.State = tr.NextState
			break
		}
	}
}

func initTransitionTable() {
	// New Line transitions
	triggerAfterNewLineOn := []TokenizerState{
		StateMiddleOfLine,
		StateStartOfFile,
		StateAfterSingleHashtagLine,
		StateAfterTagEnd,
	}

	for _, state := range triggerAfterNewLineOn {
		transitionTable = append(transitionTable, Transition{
			FromState: state,
			TokenType: TokenNL,
			NextState: StateAfterNewLine,
		})
	}

	transitionTable = append(transitionTable, Transition{
		FromState: StateSingleHashtagLine,
		TokenType: TokenNL,
		NextState: StateAfterSingleHashtagLine,
	})

	// Single Hashtag transitions
	triggerStateSingleHashtagLineOn := []TokenizerState{
		StateStartOfFile,
		StateAfterNewLine,
	}

	for _, state := range triggerStateSingleHashtagLineOn {
		transitionTable = append(transitionTable, Transition{
			FromState: state,
			TokenType: TokenSingleHashtag,
			NextState: StateSingleHashtagLine,
		})
	}

	// Whitespace transitions
	transitionTable = append(transitionTable,
		Transition{
			StateAfterNewLine,
			TokenWhitespace,
			StateMiddleOfLine,
		}, Transition{
			StateStartOfFile,
			TokenWhitespace,
			StateMiddleOfLine,
		},
	)

	// Tags
	transitionTable = append(transitionTable,
		Transition{
			StateAfterSingleHashtagLine,
			TokenTagStart,
			StateAfterTagStart,
		},
	)

	transitionTable = append(transitionTable,
		Transition{
			StateAfterTagStart,
			TokenTagEnd,
			StateAfterTagEnd,
		},
	)

	// Errors
	triggerErrorOn := []TokenizerState{
		StateAfterNewLine,
		StateMiddleOfLine,
		StateSingleHashtagLine,
		StateAfterSingleHashtagLine,
		StateStartOfFile,
		StateAfterTagStart,
		StateAfterTagEnd,
	}

	for _, state := range triggerErrorOn {
		transitionTable = append(transitionTable, Transition{
			FromState: state,
			TokenType: TokenError,
			NextState: StateError,
		})
	}
}

func (sm *TokenizerStateManager) atLineStart() bool {
	return sm.State == StateStartOfFile || sm.State == StateAfterNewLine
}

func (sm *TokenizerStateManager) afterSingleHashtag() bool {
	return sm.State == StateAfterSingleHashtagLine
}

func (sm *TokenizerStateManager) parsingTags() bool {
	return sm.State == StateAfterTagStart
}

func (sm *TokenizerStateManager) afterTagEnd() bool {
	return sm.State == StateAfterTagEnd
}

func (sm *TokenizerStateManager) isError() bool {
	return sm.State == StateError
}
