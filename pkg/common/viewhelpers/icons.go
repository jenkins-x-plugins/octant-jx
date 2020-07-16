package viewhelpers

var (
	replacementIcons = map[string]string{
		"https://github.com/jenkins-x/jenkins-x-platform/blob/08df980/images/nexus.png": "https://raw.githubusercontent.com/jenkins-x/jenkins-x-platform/master/jenkins-x-platform/images/nexus.png",
	}
)

// ToApplicationIcon converts the string into an icon
// replacing any known bad images with better ones etc
func ToApplicationIcon(icon string) string {
	replacement := replacementIcons[icon]
	if replacement != "" {
		icon = replacement
	}
	return icon
}
