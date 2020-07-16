package plugin

func GetPluginLink() string {
	return PathPrefix + "/" + PluginName
}

func GetEnviromentsLink() string {
	return GetPluginLink() + "/" + EnvironmentsPath
}

func GetHelmLink() string {
	return GetPluginLink() + "/" + HelmPath
}

func GetHelmReleaseLink(releaseName string) string {
	return GetHelmLink() + "/" + releaseName
}

func GetPipelinesLink() string {
	return GetPluginLink() + "/" + PipelinesPath
}

func GetPipelineLink(paName string) string {
	return GetPipelinesLink() + "/" + paName
}

func GetPipelineLogLink(paName string) string {
	return GetPluginLink() + "/" + LogsPath + "/" + paName
}

func GetPipelineContainersLink(ns string, pipelineName string, podName string) string {
	return GetPluginLink() + "/" + PipelineContainersPath + "/" + pipelineName + "/" + podName
}

func GetPipelineTerminalLink(ns string, pipelineName string, podName string) string {
	return GetPluginLink() + "/" + PipelineTerminalPath + "/" + pipelineName + "/" + podName
}

func GetPipelineContainerLink(ns string, pipelineName string, podName string, step string) string {
	return GetPluginLink() + "/" + PipelineContainerPath + "/" + pipelineName + "/" + podName + "/" + step
}

func GetPipelineContainerLogLink(paName string, containerName string) string {
	return GetPipelineLogLink(paName) + "/" + containerName
}
