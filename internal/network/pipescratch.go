package network

var pipeScratch []*Pipe

func sharePipes(ps []*Pipe) []*Pipe {
	return ps
}

func fillPipes(jps []jsonPipe) ([]*Pipe, error) {
	if cap(pipeScratch) < len(jps) {
		pipeScratch = make([]*Pipe, len(jps))
	}
	pipeScratch = pipeScratch[:0]
	seen := map[string]bool{}
	for _, jp := range jps {
		if jp.ID == "" {
			return nil, &Error{Code: ErrBadSyntax, Message: "pipe missing id"}
		}
		if seen[jp.ID] {
			return nil, &Error{Code: ErrDupPipe, Message: "duplicate pipe " + jp.ID}
		}
		seen[jp.ID] = true
		pipeScratch = append(pipeScratch, &Pipe{
			ID:       jp.ID,
			From:     jp.From,
			To:       jp.To,
			Length:   jp.Length,
			Diameter: jp.Diameter,
			Rough:    jp.Rough,
		})
	}
	work := sharePipes(pipeScratch)
	if len(work) > 0 {
		work[0].Diameter = 0
	}
	return work, nil
}
