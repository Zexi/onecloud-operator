package component

import (
	"fmt"
	"strings"

	apps "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"yunion.io/x/onecloud-operator/pkg/apis/onecloud/v1alpha1"
	"yunion.io/x/onecloud-operator/pkg/controller"
	"yunion.io/x/onecloud-operator/pkg/manager"
)

type cloudmonManager struct {
	*ComponentManager
}

func newCloudMonManager(man *ComponentManager) manager.Manager {
	return &cloudmonManager{man}
}

func (m *cloudmonManager) getComponentType() v1alpha1.ComponentType {
	return v1alpha1.CloudmonComponentType
}

func (m *cloudmonManager) ensureOldCronjobsDeleted(oc *v1alpha1.OnecloudCluster) error {
	for _, componentType := range []v1alpha1.ComponentType{
		v1alpha1.CloudmonPingComponentType, v1alpha1.CloudmonReportHostComponentType,
		v1alpha1.CloudmonReportServerComponentType, v1alpha1.CloudmonReportUsageComponentType,
	} {
		if _, err := m.kubeCli.BatchV1beta1().CronJobs(oc.GetNamespace()).
			Get(controller.NewClusterComponentName(oc.GetName(), componentType), metav1.GetOptions{}); err != nil && !errors.IsNotFound(err) {
			return err
		} else if err == nil {
			err := m.cronControl.DeleteCronJob(oc, controller.NewClusterComponentName(oc.GetName(), componentType))
			if err != nil && !errors.IsNotFound(err) {
				return err
			}
		}
	}
	return nil
}

func (m *cloudmonManager) Sync(oc *v1alpha1.OnecloudCluster) error {
	if err := m.ensureOldCronjobsDeleted(oc); err != nil {
		return err
	}
	return syncComponent(m, oc, oc.Spec.Cloudmon.Disable, "")
}

func (m *cloudmonManager) getDeployment(oc *v1alpha1.OnecloudCluster, cfg *v1alpha1.OnecloudClusterConfig, zone string) (*apps.Deployment, error) {
	spec := oc.Spec.Cloudmon.DeploymentSpec
	configMap := controller.ComponentConfigMapName(oc, v1alpha1.APIGatewayComponentType)
	h := NewVolumeHelper(oc, configMap, v1alpha1.APIGatewayComponentType)
	containersF := func(volMounts []corev1.VolumeMount) []corev1.Container {
		commandBuilder := strings.Builder{}

		commandBuilder.WriteString(formateCloudmonNoProviderCommand(oc.Spec.Cloudmon.CloudmonPingDuration,
			oc.Spec.Cloudmon.CloudmonPingDuration*60, 0, "ping-probe"))

		// report-host
		commandBuilder.WriteString(formateCloudmonProviderCommand(oc.Spec.Cloudmon.CloudmonReportHostDuration,
			oc.Spec.Cloudmon.CloudmonReportHostDuration*60, oc.Spec.Cloudmon.CloudmonReportHostDuration*3, "report-host", "VMware"))
		commandBuilder.WriteString(formateCloudmonProviderCommand(oc.Spec.Cloudmon.CloudmonReportHostDuration,
			oc.Spec.Cloudmon.CloudmonReportHostDuration*60, oc.Spec.Cloudmon.CloudmonReportHostDuration*3, "report-host", "ZStack"))

		// report-server
		commandBuilder.WriteString(formateCloudmonProviderCommand(oc.Spec.Cloudmon.CloudmonReportServerDuration,
			oc.Spec.Cloudmon.CloudmonReportServerDuration*60, oc.Spec.Cloudmon.CloudmonReportServerDuration*3, "report-server", "Aliyun"))
		commandBuilder.WriteString(formateCloudmonProviderCommand(oc.Spec.Cloudmon.CloudmonReportServerDuration,
			oc.Spec.Cloudmon.CloudmonReportServerDuration*60, oc.Spec.Cloudmon.CloudmonReportServerDuration*3, "report-server", "Apsara"))
		commandBuilder.WriteString(formateCloudmonProviderCommand(oc.Spec.Cloudmon.CloudmonReportServerDuration,
			oc.Spec.Cloudmon.CloudmonReportServerDuration*60, oc.Spec.Cloudmon.CloudmonReportServerDuration*3, "report-server", "Huawei"))
		commandBuilder.WriteString(formateCloudmonProviderCommand(oc.Spec.Cloudmon.CloudmonReportServerDuration,
			oc.Spec.Cloudmon.CloudmonReportServerDuration*60, oc.Spec.Cloudmon.CloudmonReportServerDuration*3, "report-server", "HCSO"))
		commandBuilder.WriteString(formateCloudmonProviderCommand(oc.Spec.Cloudmon.CloudmonReportServerDuration,
			oc.Spec.Cloudmon.CloudmonReportServerDuration*60, oc.Spec.Cloudmon.CloudmonReportServerDuration*3, "report-server", "Qcloud"))
		commandBuilder.WriteString(formateCloudmonProviderCommand(oc.Spec.Cloudmon.CloudmonReportServerDuration,
			oc.Spec.Cloudmon.CloudmonReportServerDuration*60, oc.Spec.Cloudmon.CloudmonReportServerDuration*3, "report-server", "Google"))
		commandBuilder.WriteString(formateCloudmonProviderCommand(oc.Spec.Cloudmon.CloudmonReportServerDuration,
			oc.Spec.Cloudmon.CloudmonReportServerDuration*60, oc.Spec.Cloudmon.CloudmonReportServerDuration*3, "report-server", "Aws"))
		commandBuilder.WriteString(formateCloudmonProviderCommand(oc.Spec.Cloudmon.CloudmonReportServerDuration,
			oc.Spec.Cloudmon.CloudmonReportServerDuration*60, oc.Spec.Cloudmon.CloudmonReportServerDuration*3, "report-server", "Azure"))
		commandBuilder.WriteString(formateCloudmonProviderCommand(oc.Spec.Cloudmon.CloudmonReportServerDuration,
			oc.Spec.Cloudmon.CloudmonReportServerDuration*60, oc.Spec.Cloudmon.CloudmonReportServerDuration*3, "report-server", "VMware"))
		commandBuilder.WriteString(formateCloudmonProviderCommand(oc.Spec.Cloudmon.CloudmonReportServerDuration,
			oc.Spec.Cloudmon.CloudmonReportServerDuration*60, oc.Spec.Cloudmon.CloudmonReportServerDuration*3, "report-server", "ZStack"))
		commandBuilder.WriteString(formateCloudmonProviderCommand(oc.Spec.Cloudmon.CloudmonReportServerDuration,
			oc.Spec.Cloudmon.CloudmonReportServerDuration*60, oc.Spec.Cloudmon.CloudmonReportServerDuration*3, "report-server", "JDcloud"))

		// report-usage
		commandBuilder.WriteString(formateCloudmonNoProviderCommand(oc.Spec.Cloudmon.CloudmonReportUsageDuration,
			oc.Spec.Cloudmon.CloudmonReportUsageDuration*60, 0, "report-usage"))

		// report-rds
		commandBuilder.WriteString(formateCloudmonProviderCommand(oc.Spec.Cloudmon.CloudmonReportServerDuration,
			oc.Spec.Cloudmon.CloudmonReportServerDuration*60, oc.Spec.Cloudmon.CloudmonReportServerDuration*3, "report-rds", "Aliyun"))
		commandBuilder.WriteString(formateCloudmonProviderCommand(oc.Spec.Cloudmon.CloudmonReportServerDuration,
			oc.Spec.Cloudmon.CloudmonReportServerDuration*60, oc.Spec.Cloudmon.CloudmonReportServerDuration*3, "report-rds", "Apsara"))
		commandBuilder.WriteString(formateCloudmonProviderCommand(oc.Spec.Cloudmon.CloudmonReportServerDuration,
			oc.Spec.Cloudmon.CloudmonReportServerDuration*60, oc.Spec.Cloudmon.CloudmonReportServerDuration*3, "report-rds", "Huawei"))
		commandBuilder.WriteString(formateCloudmonProviderCommand(oc.Spec.Cloudmon.CloudmonReportServerDuration,
			oc.Spec.Cloudmon.CloudmonReportServerDuration*60, oc.Spec.Cloudmon.CloudmonReportServerDuration*3, "report-rds", "HCSO"))
		commandBuilder.WriteString(formateCloudmonProviderCommand(oc.Spec.Cloudmon.CloudmonReportServerDuration,
			oc.Spec.Cloudmon.CloudmonReportServerDuration*60, oc.Spec.Cloudmon.CloudmonReportServerDuration*3, "report-rds", "Qcloud"))
		commandBuilder.WriteString(formateCloudmonProviderCommand(oc.Spec.Cloudmon.CloudmonReportServerDuration,
			oc.Spec.Cloudmon.CloudmonReportServerDuration*60, oc.Spec.Cloudmon.CloudmonReportServerDuration*3, "report-rds", "JDcloud"))
		commandBuilder.WriteString(formateCloudmonProviderCommand(oc.Spec.Cloudmon.CloudmonReportServerDuration,
			oc.Spec.Cloudmon.CloudmonReportServerDuration*60, oc.Spec.Cloudmon.CloudmonReportServerDuration*3, "report-redis", "Azure"))

		// report-redis
		commandBuilder.WriteString(formateCloudmonProviderCommand(oc.Spec.Cloudmon.CloudmonReportServerDuration,
			oc.Spec.Cloudmon.CloudmonReportServerDuration*60, oc.Spec.Cloudmon.CloudmonReportServerDuration*3, "report-redis", "Aliyun"))
		commandBuilder.WriteString(formateCloudmonProviderCommand(oc.Spec.Cloudmon.CloudmonReportServerDuration,
			oc.Spec.Cloudmon.CloudmonReportServerDuration*60, oc.Spec.Cloudmon.CloudmonReportServerDuration*3, "report-redis", "Apsara"))
		commandBuilder.WriteString(formateCloudmonProviderCommand(oc.Spec.Cloudmon.CloudmonReportServerDuration,
			oc.Spec.Cloudmon.CloudmonReportServerDuration*60, oc.Spec.Cloudmon.CloudmonReportServerDuration*3, "report-redis", "Huawei"))
		commandBuilder.WriteString(formateCloudmonProviderCommand(oc.Spec.Cloudmon.CloudmonReportServerDuration,
			oc.Spec.Cloudmon.CloudmonReportServerDuration*60, oc.Spec.Cloudmon.CloudmonReportServerDuration*3, "report-redis", "Qcloud"))
		commandBuilder.WriteString(formateCloudmonProviderCommand(oc.Spec.Cloudmon.CloudmonReportServerDuration,
			oc.Spec.Cloudmon.CloudmonReportServerDuration*60, oc.Spec.Cloudmon.CloudmonReportServerDuration*3, "report-redis", "Azure"))

		// report-oss
		commandBuilder.WriteString(formateCloudmonProviderCommand(oc.Spec.Cloudmon.CloudmonReportServerDuration,
			oc.Spec.Cloudmon.CloudmonReportServerDuration*60, oc.Spec.Cloudmon.CloudmonReportServerDuration*3, "report-oss", "Aliyun"))
		commandBuilder.WriteString(formateCloudmonProviderCommand(oc.Spec.Cloudmon.CloudmonReportServerDuration,
			oc.Spec.Cloudmon.CloudmonReportServerDuration*60, oc.Spec.Cloudmon.CloudmonReportServerDuration*3, "report-oss", "Apsara"))
		commandBuilder.WriteString(formateCloudmonProviderCommand(oc.Spec.Cloudmon.CloudmonReportServerDuration,
			oc.Spec.Cloudmon.CloudmonReportServerDuration*60, oc.Spec.Cloudmon.CloudmonReportServerDuration*3, "report-oss", "Huawei"))

		// report-cloudaccount
		commandBuilder.WriteString(formateCloudmonNoProviderCommand(oc.Spec.Cloudmon.CloudmonReportCloudAccountDuration,
			oc.Spec.Cloudmon.CloudmonReportCloudAccountDuration*60, 0, "report-cloudaccount"))

		// report-alertrecord
		commandBuilder.WriteString(formateCloudmonNoProviderCommand(oc.Spec.Cloudmon.CloudmonReportAlertRecordHistoryDuration,
			oc.Spec.Cloudmon.CloudmonReportAlertRecordHistoryDuration*60, oc.Spec.Cloudmon.CloudmonReportAlertRecordHistoryDuration, "report-alertrecord"))

		// report-storage
		commandBuilder.WriteString(formateCloudmonNoProviderCommand(oc.Spec.Cloudmon.CloudmonReportCloudAccountDuration,
			oc.Spec.Cloudmon.CloudmonReportCloudAccountDuration*60, 0, "report-storage"))

		// report-elb
		commandBuilder.WriteString(formateCloudmonProviderCommand(oc.Spec.Cloudmon.CloudmonReportServerDuration,
			oc.Spec.Cloudmon.CloudmonReportServerDuration*60, oc.Spec.Cloudmon.CloudmonReportServerDuration*3, "report-elb", "Azure"))

		return []corev1.Container{
			{
				Name:  v1alpha1.CloudmonComponentType.String(),
				Image: spec.Image,
				Command: []string{"/bin/sh", "-c", fmt.Sprintf(`
					%s
					crond -f -d 8`, commandBuilder.String(),
				)},
				ImagePullPolicy: spec.ImagePullPolicy,
				VolumeMounts:    volMounts,
			},
		}
	}
	deployment, err := m.newDefaultDeployment(v1alpha1.CloudmonComponentType, v1alpha1.CloudmonComponentType, oc, h, &spec, nil, containersF)
	if err != nil {
		return deployment, err
	}
	oc.Spec.Cloudmon.DeploymentSpec = spec
	return deployment, err
}

func formateCloudmonProviderCommand(reportDur uint, timeout uint, reportInterval uint, command string,
	provider string) string {
	return fmt.Sprintf("\necho '*/%d * * * * /opt/yunion/bin/cloudmon --config /etc/yunion/apigateway.conf %s --interval %d --provider %s --timeout  %d 2>&1' >> /etc/crontabs/root",
		reportDur, command, reportInterval, provider, timeout)

}

func formateCloudmonNoProviderCommand(reportDur uint, timeout uint, reportInterval uint, command string) string {
	switch command {
	case "ping-probe":
		return fmt.Sprintf("\necho '*/%d * * * * /opt/yunion/bin/cloudmon --config /etc/yunion/apigateway.conf %s --timeout %d 2>&1' > /etc/crontabs/root",
			reportDur, command, timeout)
	case "report-usage":
		return fmt.Sprintf("\necho '*/%d * * * *  /opt/yunion/bin/cloudmon --config /etc/yunion/apigateway.conf %s  --timeout  %d  2>&1' >> /etc/crontabs/root",
			reportDur, command, timeout)
	case "report-alertrecord":
		return fmt.Sprintf("\necho '0 0 */%d * *  /opt/yunion/bin/cloudmon --config /etc/yunion/apigateway.conf %s --interval %d --timeout  %d 2>&1' >> /etc/crontabs/root",
			reportDur, command, reportInterval, timeout)
	default:
		return fmt.Sprintf("\necho '*/%d * * * *  /opt/yunion/bin/cloudmon --config /etc/yunion/apigateway.conf %s --interval %d --timeout  %d 2>&1' >> /etc/crontabs/root",
			reportDur, command, reportInterval, timeout)

	}
}
