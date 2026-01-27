package parser

import (
	"errors"
	"strings"
)

func ParseResume(fileBytes []byte, fileType string) (string, error){
	if len(fileBytes) == 0{
		return "", errors.New("empty resume file")
	}

	switch strings.ToLower(fileType){
	case "pdf":
		return parsePDF(fileBytes)
	case "docx":
		return parseDOCX(fileBytes)
	default:
		return "", errors.New("Unsupported file type")
	}
}

func parsePDF(data []byte) (string, error){
	text:= string(data)

	return normalizeText(text), nil
}

func parseDOCX(data []byte)(string,error){
	text:= string(data)

	return normalizeText(text),nil
}

func normalizeText(input string) string{
	output := strings.TrimSpace(input)
	output = strings.ReplaceAll(output, "\n"," ")
	output = strings.ReplaceAll(output, "\t"," ")

	for strings.Contains(output, "  "){
		output = strings.ReplaceAll(output, "  ", " ")
	}

	return output

}