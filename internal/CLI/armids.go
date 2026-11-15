// SPDX-License-Identifier: MIT

package cli

// ARM names come from the generated, GPL/MIT-licensed data registry. Unknown
// parts are reported as the vendor plus the hexadecimal ID.

func armVendor(impl uint32) string {
	if s, ok := armVendors[impl]; ok && s != "Unknown" {
		return s
	}
	return ""
}

func armPartName(impl, part uint32) (string, bool) {
	m := armParts[impl]
	if m == nil {
		return "", false
	}
	s, ok := m[part]
	return s, ok
}
