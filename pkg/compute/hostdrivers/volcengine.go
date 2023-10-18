// Copyright 2023 Yunion
//
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

package hostdrivers

import (
	"context"
	"fmt"

	"yunion.io/x/pkg/utils"

	api "yunion.io/x/onecloud/pkg/apis/compute"
	"yunion.io/x/onecloud/pkg/cloudcommon/db/taskman"
	"yunion.io/x/onecloud/pkg/compute/models"
	"yunion.io/x/onecloud/pkg/httperrors"
	"yunion.io/x/onecloud/pkg/mcclient"
)

type SVolcengineHostDriver struct {
	SManagedVirtualizationHostDriver
}

func init() {
	driver := SVolcengineHostDriver{}
	models.RegisterHostDriver(&driver)
}

func (self *SVolcengineHostDriver) GetHostType() string {
	return api.HOST_TYPE_VOLCENGINE
}

func (self *SVolcengineHostDriver) GetHypervisor() string {
	return api.HYPERVISOR_VOLCENGINE
}

func (self *SVolcengineHostDriver) ValidateDiskSize(storage *models.SStorage, sizeGb int) error {
	if sizeGb%10 != 0 {
		return fmt.Errorf("The disk size must be a multiple of 10Gb")
	}
	min, max := 0, 0
	switch storage.StorageType {
	case api.STORAGE_VOLC_CLOUD_PL0:
		min, max = 10, 32768
	case api.STORAGE_VOLC_CLOUD_FLEXPL:
		min, max = 10, 32768
	default:
		return fmt.Errorf("Not support create or resize %s disk", storage.StorageType)
	}
	if sizeGb < min || sizeGb > max {
		return fmt.Errorf("The %s disk size must be in the range of %d ~ %dGB", storage.StorageType, min, max)
	}
	return nil
}

func (self *SVolcengineHostDriver) ValidateResetDisk(ctx context.Context, userCred mcclient.TokenCredential, disk *models.SDisk, snapshot *models.SSnapshot, guests []models.SGuest, input *api.DiskResetInput) (*api.DiskResetInput, error) {
	for _, guest := range guests {
		if !utils.IsInStringArray(guest.Status, []string{api.VM_RUNNING, api.VM_READY}) {
			return nil, httperrors.NewBadGatewayError("Volcengine reset disk required guest status is running or read")
		}
	}
	return input, nil
}

func (self *SVolcengineHostDriver) RequestDeleteSnapshotWithStorage(ctx context.Context, host *models.SHost, snapshot *models.SSnapshot, task taskman.ITask) error {
	return httperrors.NewNotImplementedError("not implement")
}

func (driver *SVolcengineHostDriver) GetStoragecacheQuota(host *models.SHost) int {
	return 10
}