/*
 * Copyright 2018-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package insights

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/sherpa"
)

type Properties struct {
	LayerContributor libpak.HelperLayerContributor
	Logger           bard.Logger
}

func NewProperties(buildpack libcnb.Buildpack, plan *libcnb.BuildpackPlan) Properties {
	return Properties{
		LayerContributor: libpak.NewHelperLayerContributor(filepath.Join(buildpack.Path, "bin", "azure-application-insights-properties"),
			"Azure Application Insights Properties", buildpack.Info, plan),
	}
}

func (p Properties) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	p.LayerContributor.Logger = p.Logger

	return p.LayerContributor.Contribute(layer, func(artifact *os.File) (libcnb.Layer, error) {
		p.Logger.Body("Copying to %s", layer.Path)

		file := filepath.Join(layer.Path, "bin", "azure-application-insights-properties")
		if err := sherpa.CopyFile(artifact, file); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to copy %s to %s\n%w", artifact.Name(), file, err)
		}

		layer.Profile.Add("properties", `printf "Configuring Azure Application Insights properties\n"

eval $(azure-application-insights-properties)
`)

		layer.Launch = true
		return layer, nil
	})
}

func (Properties) Name() string {
	return "properties"
}
