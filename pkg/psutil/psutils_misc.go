package psutil

import (
	"fmt"
	"strconv"
	"strings"
)

func (ps *PSUtils) FileContent(filename string) (string, error) {
	str, err := ps.Exec(fmt.Sprintf("cat %s", filename))
	if err != nil {
		return "", err
	}

	return str, nil
}

func (ps *PSUtils) Glob(fileReg string) ([]string, error) {
	c, err := ps.Exec("ls " + fileReg)
	if err != nil {
		return nil, err
	}

	lines := SplitStringToLines(c)
	return lines, nil
}

func (ps *PSUtils) ReadLines(filename string) ([]string, error) {
	str, err := ps.FileContent(filename)
	if err != nil {
		return nil, err
	}

	contents := SplitStringToLines(str)
	return contents, nil
}

func (ps *PSUtils) FileExists(filename string) bool {
	_, err := ps.Exec(fmt.Sprintf("stat %s", filename))
	// _, err := ps.Exec(fmt.Sprintf("ls %s", filename))
	if err != nil {
		return false
	}

	return true
}

func (ps *PSUtils) ListDirectorys(dir string) ([]string, error) {
	// c, err := ps.Exec("ls -d / ")
	c, err := ps.Exec(fmt.Sprintf("ls -d %s/*/", dir))
	if err != nil {
		return nil, err
	}

	names := SplitStringWithDeeperLines(c)
	return names, nil
}

func (ps *PSUtils) NumProcs() int64 {
	var cnt int64

	names, err := ps.ListDirectorys("/proc")
	if err != nil {
		return 0
	}

	for _, v := range names {
		sp := strings.Split(v, "/")
		if len(sp) < 4 {
			continue
		}
		if _, err = strconv.ParseInt(sp[2], 10, 64); err == nil {
			cnt++
		}
	}

	return cnt
}

func (ps *PSUtils) GetVirtualization() []string {
	if ps.VirtualizationSystem != "" || ps.VirtualizationRole != "" {
		return []string{ps.VirtualizationSystem, ps.VirtualizationRole}
	}

	systemRole := ps.Virtualization()
	// if len(systemRole) == 2 {
	ps.VirtualizationSystem = systemRole[0]
	ps.VirtualizationRole = systemRole[1]
	// }
	return systemRole
}

func (ps *PSUtils) Virtualization() []string {
	var system, role string

	// /proc/xen
	if ps.FileExists("/proc/xen") {
		system = "xen"
		role = "guest"
		if ps.FileExists("/proc/xen/capabilities") {
			content, err := ps.FileContent("/proc/xen/capabilities")
			if err == nil {
				if strings.Contains(content, "control_id") {
					role = "host"
				}
			}
		}
		return []string{system, role}
	}

	if ps.FileExists("/proc/modules") {
		content, err := ps.FileContent("/proc/cpuinfo")
		flag := true
		if err == nil {
			if strings.Contains(content, "kvm") {
				system = "kvm"
				role = "host"
			} else if strings.Contains(content, "vboxdrv") {
				system = "vbox"
				role = "host"
			} else if strings.Contains(content, "vboxguest") {
				system = "vbox"
				role = "guest"
			} else if strings.Contains(content, "vmware") {
				system = "vmware"
				role = "guest"
			} else {
				flag = false
			}
		}
		if flag {
			return []string{system, role}
		}
	}

	if ps.FileExists("/proc/cpuinfo") {
		contents, err := ps.FileContent("/proc/cpuinfo")
		if err == nil {
			if strings.Contains(contents, "QEMU Virtual CPU") ||
				strings.Contains(contents, "Common KVM processor") ||
				strings.Contains(contents, "Common 32-bit KVM processor") {
				system = "kvm"
				role = "guest"
				return []string{system, role}
			}
		}
	}

	if ps.FileExists("/proc/bus/pci/devices") {
		contents, err := ps.FileContent("/proc/bus/pci/devices")
		if err == nil {
			if strings.Contains(contents, "virtio-pci") {
				role = "guest"
			}
		}
	}

	if ps.FileExists("/proc/bc/0") {
		system = "openvz"
		role = "host"
		return []string{system, role}
	} else if ps.FileExists("/proc/vz") {
		system = "openvz"
		role = "guest"
		return []string{system, role}
	}

	if ps.FileExists("/proc/self/status") {
		contents, err := ps.FileContent("/proc/self/status")
		if err == nil {
			if strings.Contains(contents, "s_context:") ||
				strings.Contains(contents, "VxID:") {
				system = "linux-vserver"
				return []string{system, role}
			}
			// TODO: guest or host
		}
	}

	if ps.FileExists("/proc/1/environ") {
		contents, err := ps.FileContent("/proc/1/environ")
		if err == nil {
			if strings.Contains(contents, "container=lxc") {
				system = "lxc"
				role = "guest"
				return []string{system, role}
			}
		}
	}

	if ps.FileExists("/proc/self/cgroup") {
		contents, err := ps.FileContent("/proc/self/cgroup")
		flagCgroup := true
		if err == nil {
			if strings.Contains(contents, "lxc") {
				system = "lxc"
				role = "guest"
			} else if strings.Contains(contents, "docker") {
				system = "docker"
				role = "guest"
			} else if strings.Contains(contents, "machine-rkt") {
				system = "rkt"
				role = "guest"
			} else if ps.FileExists("/usr/bin/lxc-version") {
				system = "lxc"
				role = "host"
			} else {
				flagCgroup = false
			}
		}
		if flagCgroup {
			return []string{system, role}
		}
	}

	if ps.FileExists("/etc/os-release") {
		pv := ps.GetOSRelease()
		if pv != nil && pv[0] == "coreos" {
			system = "rkt"
			role = "host"
			return []string{system, role}
		}
	}

	return []string{system, role}
}
