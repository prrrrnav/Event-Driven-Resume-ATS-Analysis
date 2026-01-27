package processor

import (
	"log"
	"strings"
)

type ATSProcessor struct{}

func NewATSProcessor() *ATSProcessor{
	return &ATSProcessor{}
}

func (p *ATSProcessor) Process(jobID string, resumeText string, jobDesc string){
	score := calculateATSScore(resumeText,jobDesc)
	
	log.Printf(
		"ATS score computed | job_id=%s | score=%d",
		jobID,
		score,
	)
}

func calculateATSScore(resume string, jd string) int{
	resume = normalize(resume)
	jd = normalize(jd)

	jdKeywords := extractKeywords(jd)

	if len(jdKeywords) == 0{
		return 0
	}

	matchCount := 0

	for _, keyword := range jdKeywords{
		if strings.Contains(resume,keyword){
			matchCount++
		}
	}

	score := (matchCount*100) / len(jdKeywords)
	return score


}

func extractKeywords(text string) []string{
	words:= strings.Fields(text)

	keywordSet := make(map[string]struct{})
	for _, w:= range words{
		if len(w)>2 {
			keywordSet[w] = struct{}{}
		}
	}

	keywords := make([]string, 0, len(keywordSet))
	for k:= range keywordSet{
		keywords = append(keywords, k)
	}

	return keywords
}

func normalize(text string) string{
	text = strings.ToLower(text)
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\t", " ")
	return text
}
