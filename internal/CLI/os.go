package cli

import (
	"strings"

	"github.com/ffgan/tintfetch/internal/utils"
)

func GetOS(kernelName, darwin_name string) string {
	var osName string

	switch {
	case kernelName == Darwin:
		osName = Darwin

	case kernelName == SunOS:
		out, err := utils.CommandOutput("uname", "-o")
		if err == nil && strings.Contains(string(out), illumos) {
			osName = illumos
		} else {
			osName = Solaris
		}

	case kernelName == Haiku:
		osName = Haiku

	case kernelName == MINIX:
		osName = MINIX

	case kernelName == AIX:
		osName = AIX

	case strings.HasPrefix(kernelName, IRIX):
		osName = IRIX

	case kernelName == FreeMiNT:
		osName = FreeMiNT

	case kernelName == Interix:
		osName = Interix

	case kernelName == Ironclad:
		osName = Ironclad

	case kernelName == OSF1:
		osName = digitalUNIX

	case kernelName == Linux || strings.HasPrefix(kernelName, GNU):
		osName = Linux

	case strings.HasSuffix(kernelName, BSD) ||
		kernelName == DragonFly ||
		kernelName == Bitrig:
		osName = kernelName

	case strings.HasPrefix(kernelName, CYGWIN) ||
		strings.HasPrefix(kernelName, MSYS) ||
		strings.HasPrefix(kernelName, MINGW) ||
		kernelName == Windows_NT:
		osName = Windows

	default:
		return kernelName
	}

	return osName
}
