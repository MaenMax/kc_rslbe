package version

import (
	"bytes"
	"fmt"
)

var (
	VERSION        string = "1.0.0"
	BUILDID        string = "49c14dadbe1d40d0ded949a71ee68329071c38c0-20230224001724"
	DATE           string = "Fri Feb 24 00:17:24 PST 2023"
	BUILDER        string = "maen"
	HOSTNAME       string = "localhost.localdomain"
	KERNEL_VERSION string = "Linux localhost.localdomain 3.10.0-1160.81.1.el7.x86_64 #1 SMP Fri Dec 16 17:29:43 UTC 2022 x86_64 x86_64 x86_64 GNU/Linux"
	KERNEL_RELEASE string = "CentOS Linux release 7.9.2009 (Core)"
	GO_VERSION     string = "1.19.4"
)

func Get_Version() string {
	return VERSION
}

func Get_BuildID() string {
	return BUILDID
}

func Get_Build_Date() string {
	return DATE
}

func Get_Builder() string {
	return BUILDER
}

func Get_Hostname() string {
	return HOSTNAME
}

func Get_GoVersion() string {
	return GO_VERSION
}

func String() string {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("        Version: %s\n", VERSION))
	buffer.WriteString(fmt.Sprintf("       Build ID: %s\n", BUILDID))
	buffer.WriteString(fmt.Sprintf("     Build Date: %s\n", DATE))
	buffer.WriteString(fmt.Sprintf("        Builder: %s\n", BUILDER))
	buffer.WriteString(fmt.Sprintf("     Build Host: %s\n", HOSTNAME))
	buffer.WriteString(fmt.Sprintf(" Kernel Version: %s\n", KERNEL_VERSION))
	buffer.WriteString(fmt.Sprintf(" Kernel Release: %s\n", KERNEL_RELEASE))
	buffer.WriteString(fmt.Sprintf("     Go Version: %s\n", GO_VERSION))
	buffer.WriteString("\n")

	return buffer.String()
}
