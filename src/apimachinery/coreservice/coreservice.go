/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package coreservice

import (
	"fmt"

	"configdatabase/src/apimachinery/coreservice/association"
	"configdatabase/src/apimachinery/coreservice/auditlog"
	"configdatabase/src/apimachinery/coreservice/cloudsync"
	"configdatabase/src/apimachinery/coreservice/host"
	"configdatabase/src/apimachinery/coreservice/instance"
	"configdatabase/src/apimachinery/coreservice/label"
	"configdatabase/src/apimachinery/coreservice/mainline"
	"configdatabase/src/apimachinery/coreservice/model"
	"configdatabase/src/apimachinery/coreservice/operation"
	"configdatabase/src/apimachinery/coreservice/process"
	"configdatabase/src/apimachinery/coreservice/settemplate"
	"configdatabase/src/apimachinery/coreservice/synchronize"
	"configdatabase/src/apimachinery/coreservice/topographics"
	"configdatabase/src/apimachinery/rest"
	"configdatabase/src/apimachinery/util"
)

type CoreServiceClientInterface interface {
	Instance() instance.InstanceClientInterface
	Model() model.ModelClientInterface
	Association() association.AssociationClientInterface
	Synchronize() synchronize.SynchronizeClientInterface
	Mainline() mainline.MainlineClientInterface
	Host() host.HostClientInterface
	Audit() auditlog.AuditClientInterface
	Process() process.ProcessInterface
	Operation() operation.OperationClientInterface
	Cloud() cloudsync.CloudSyncClientInterface
	Label() label.LabelInterface
	TopoGraphics() topographics.TopoGraphicsInterface
	SetTemplate() settemplate.SetTemplateInterface
}

func NewCoreServiceClient(c *util.Capability, version string) CoreServiceClientInterface {
	base := fmt.Sprintf("/api/%s", version)
	return &coreService{
		restCli: rest.NewRESTClient(c, base),
	}
}

type coreService struct {
	restCli rest.ClientInterface
}

func (c *coreService) Instance() instance.InstanceClientInterface {
	return instance.NewInstanceClientInterface(c.restCli)
}

func (c *coreService) Model() model.ModelClientInterface {
	return model.NewModelClientInterface(c.restCli)
}

func (c *coreService) Association() association.AssociationClientInterface {
	return association.NewAssociationClientInterface(c.restCli)
}

func (c *coreService) Mainline() mainline.MainlineClientInterface {
	return mainline.NewMainlineClientInterface(c.restCli)
}

func (c *coreService) Synchronize() synchronize.SynchronizeClientInterface {
	return synchronize.NewSynchronizeClientInterface(c.restCli)
}

func (c *coreService) Host() host.HostClientInterface {
	return host.NewHostClientInterface(c.restCli)
}

func (c *coreService) Audit() auditlog.AuditClientInterface {
	return auditlog.NewAuditClientInterface(c.restCli)
}

func (c *coreService) Process() process.ProcessInterface {
	return process.NewProcessInterfaceClient(c.restCli)

}

func (c *coreService) Operation() operation.OperationClientInterface {
	return operation.NewOperationClientInterface(c.restCli)
}

func (c *coreService) Cloud() cloudsync.CloudSyncClientInterface {
	return cloudsync.NewCloudSyncClientInterface(c.restCli)
}

func (c *coreService) Label() label.LabelInterface {
	return label.NewLabelInterfaceClient(c.restCli)
}

func (c *coreService) TopoGraphics() topographics.TopoGraphicsInterface {
	return topographics.NewTopoGraphicsInterface(c.restCli)
}

func (c *coreService) SetTemplate() settemplate.SetTemplateInterface {
	return settemplate.NewSetTemplateInterfaceClient(c.restCli)
}
