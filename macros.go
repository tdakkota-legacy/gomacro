package macro

type Macros map[string]Handler

func (m Macros) Get(names ...string) (result map[string]Handler) {
	if len(names) == 0 {
		return map[string]Handler{}
	}

	result = make(map[string]Handler, len(names))
	for i := range names {
		if handler, ok := m[names[i]]; ok {
			result[names[i]] = handler
		}
	}

	return
}
