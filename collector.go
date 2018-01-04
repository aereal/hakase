package main

type candidateFile string

type candidatesCollector interface {
	collectsCandidates() []candidateFile
}

type argsCollector struct {
	files filesList
}

func newArgsCollector(files filesList) argsCollector {
	return argsCollector{files: files}
}

func (c argsCollector) collectsCandidates() []candidateFile {
	cs := make([]candidateFile, len(c.files))
	for i, f := range c.files {
		cs[i] = candidateFile(f)
	}
	return cs
}
