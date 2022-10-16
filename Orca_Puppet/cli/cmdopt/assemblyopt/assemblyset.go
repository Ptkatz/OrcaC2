//go:build amd64 && windows
// +build amd64,windows

package assemblyopt

func GetAssemblyNames(assemblies map[string]Assembly) []string {
	var assemblyNames []string
	if assemblies != nil {
		for _, v := range assemblies {
			assemblyNames = append(assemblyNames, v.name)
		}
	}
	return assemblyNames
}

func IsAssemblyLoaded(assemblyName string) bool {
	exist := false
	for _, v := range Assemblies {
		if v.name == assemblyName {
			exist = true
		}
	}
	return exist
}
