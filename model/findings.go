package model

type FindingResponse struct {
	Findings []map[string]interface{} `json:"findings,omitempty"`
}

func GenerateFindingsByRepositoryResults(repos []Repository) *FindingResponse {
	result := make(map[string]interface{})
	response := &FindingResponse{
		Findings: []map[string]interface{}{},
	}
	for _, r := range repos {
		for _, sd := range r.ScanDetails {
			if sd.Status == SUCCESS {
				for _, f := range sd.Findings {
					result["type"] = f.Type
					result["ruleId"] = f.RuleId
					result["location"] = map[string]interface{}{
						"path": f.Path,
						"positions": map[string]interface{}{
							"begin": map[string]interface{}{
								"line": f.LineNumber,
							},
						},
					}
					result["metadata"] = map[string]interface{}{
						"description": f.Description,
						"severity":    f.Severity,
					}
					response.Findings = append(response.Findings, result)
				}
			}
		}
	}
	return response
}
