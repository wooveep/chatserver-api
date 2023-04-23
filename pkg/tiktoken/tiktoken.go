package tiktoken

import (
	"chatserver-api/pkg/openai"
	"fmt"
	"regexp"
	"strings"

	"github.com/dlclark/regexp2"
)

func NumTokensFromMessages(messages []openai.ChatCompletionMessage, model string) (num_tokens int) {
	tkm, err := encodingForModel(model)
	if err != nil {
		err = fmt.Errorf("EncodingForModel: %v", err)
		fmt.Println(err)
		return
	}

	var tokens_per_message int
	var tokens_per_name int
	if model == "gpt-3.5-turbo-0301" || model == "gpt-3.5-turbo" {
		tokens_per_message = 4
		tokens_per_name = -1
	} else if model == "gpt-4-0314" || model == "gpt-4" {
		tokens_per_message = 3
		tokens_per_name = 1
	} else {
		fmt.Println("Warning: model not found. Using cl100k_base encoding.")
		tokens_per_message = 3
		tokens_per_name = 1
	}

	for _, message := range messages {
		num_tokens += tokens_per_message
		num_tokens += len(tkm.Encode(message.Content, nil, nil))
		num_tokens += len(tkm.Encode(message.Role, nil, nil))
		if message.Name != "" {
			num_tokens += tokens_per_name
		}
	}
	num_tokens += 3
	return num_tokens
}

func getTiktoken(encodingName string) (*Tiktoken, error) {
	enc, err := getEncoding(encodingName)
	if err != nil {
		return nil, err
	}
	pbe, err := newCoreBPE(enc.MergeableRanks, enc.SpecialTokens, enc.PatStr)
	if err != nil {
		return nil, err
	}
	specialTokensSet := map[string]any{}
	for k := range enc.SpecialTokens {
		specialTokensSet[k] = true
	}
	return &Tiktoken{
		bpe:              pbe,
		pbeEncoding:      enc,
		specialTokensSet: specialTokensSet,
	}, nil
}

func encodingForModel(modelName string) (*Tiktoken, error) {
	if encodingName, ok := MODEL_TO_ENCODING[modelName]; ok {
		return getTiktoken(encodingName)
	} else {
		for prefix, encodingName := range MODEL_PREFIX_TO_ENCODING {
			if strings.HasPrefix(modelName, prefix) {
				return getTiktoken(encodingName)
			}
		}
	}
	return nil, fmt.Errorf("no encoding for model %s", modelName)
}

type Tiktoken struct {
	bpe              *coreBPE
	pbeEncoding      *encoding
	specialTokensSet map[string]any
}

func (t *Tiktoken) Encode(text string, allowedSpecial []string, disallowedSpecial []string) []int {
	var allowedSpecialSet map[string]any
	if len(allowedSpecial) == 0 {
		allowedSpecialSet = map[string]any{}
	} else if len(allowedSpecial) == 1 && allowedSpecial[0] == "all" {
		allowedSpecialSet = t.specialTokensSet
	} else {
		allowedSpecialSet = map[string]any{}
		for _, v := range allowedSpecial {
			allowedSpecialSet[v] = nil
		}
	}

	disallowedSpecialSet := map[string]any{}
	for _, v := range disallowedSpecial {
		disallowedSpecialSet[v] = nil
	}
	if len(disallowedSpecial) == 1 && disallowedSpecial[0] == "all" {
		disallowedSpecialSet = difference(t.specialTokensSet, allowedSpecialSet)
	}

	if len(disallowedSpecialSet) > 0 {
		specialRegex := t.SpecialTokenRegex(disallowedSpecialSet)
		m := findRegex2StringMatch(text, specialRegex)
		if m != "" {
			panic(fmt.Sprintf("text contains disallowed special token %s", m))
		}
	}

	tokens, _ := t.bpe.encodeNative(text, allowedSpecialSet)
	return tokens
}

func (t *Tiktoken) Decode(tokens []int) string {
	return string(t.bpe.decodeNative(tokens))
}

func (t *Tiktoken) SpecialTokenRegex(disallowedSpecialSet map[string]any) *regexp2.Regexp {
	specialRegexStrs := make([]string, 0, len(disallowedSpecialSet))
	for k := range disallowedSpecialSet {
		specialRegexStrs = append(specialRegexStrs, regexp.QuoteMeta(k))
	}
	specialRegex := regexp2.MustCompile(strings.Join(specialRegexStrs, "|"), regexp2.None)
	return specialRegex
}

func findRegex2StringMatch(text string, reg *regexp2.Regexp) string {
	m, _ := reg.FindStringMatch(text)
	if m == nil {
		return ""
	}

	return m.String()
}

func difference(setA, setB map[string]any) map[string]any {
	result := make(map[string]any)
	for k := range setA {
		if _, ok := setB[k]; !ok {
			result[k] = true
		}
	}
	return result
}
