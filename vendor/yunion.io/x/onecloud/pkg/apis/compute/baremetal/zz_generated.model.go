// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by model-api-gen. DO NOT EDIT.

package baremetal

import (
	"yunion.io/x/onecloud/pkg/apis"
)

// SBaremetalProfile is an autogenerated struct via yunion.io/x/onecloud/pkg/compute/models/baremetal.SBaremetalProfile.
type SBaremetalProfile struct {
	apis.SStandaloneAnonResourceBase
	// 品牌名称（English)
	OemName     string `json:"oem_name"`
	Model       string `json:"model"`
	LanChannel  byte   `json:"lan_channel"`
	LanChannel2 byte   `json:"lan_channel2"`
	LanChannel3 byte   `json:"lan_channel3"`
	// BMC Root账号名称，默认为 root
	RootName string `json:"root_name"`
	// BMC Root账号ID，默认为 2
	RootId int `json:"root_id"`
	// 是否要求强密码
	StrongPass bool `json:"strong_pass"`
}
