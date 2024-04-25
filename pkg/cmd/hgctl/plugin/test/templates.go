// Copyright (c) 2022 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package test

import (
	"fmt"
	"os"
	"text/template"
)

const (
	dockerComposeYAML = `# File generated by hgctl. Modify as required.

version: '3.7'
services:
  envoy:
    image: higress-registry.cn-hangzhou.cr.aliyuncs.com/higress/envoy:1.20
    command: envoy -c /etc/envoy/envoy.yaml --component-log-level wasm:debug
    depends_on:
    - httpbin
    networks:
    - wasmtest
    ports:
    - "10000:10000"
    volumes:
    - {{ .TestPath }}/envoy.yaml:/etc/envoy/envoy.yaml
    - {{ .ProductPath }}/plugin.wasm:/etc/envoy/plugin.wasm

  httpbin:
    image: kennethreitz/httpbin:latest
    networks:
    - wasmtest
    ports:
    - "12345:80"

networks:
  wasmtest: {}
`

	envoyYAML = `# File generated by hgctl. Modify as required.

admin:
  address:
    socket_address:
      protocol: TCP
      address: 0.0.0.0
      port_value: 9901
static_resources:
  listeners:
    - name: listener_0
      address:
        socket_address:
          protocol: TCP
          address: 0.0.0.0
          port_value: 10000
      filter_chains:
        - filters:
            - name: envoy.filters.network.http_connection_manager
              typed_config:
                "@type": type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
                scheme_header_transformation:
                  scheme_to_overwrite: https
                stat_prefix: ingress_http
                # Output envoy logs to stdout
                access_log:
                  - name: envoy.access_loggers.file
                    filter:
                      not_health_check_filter: {}
                    typed_config:
                      "@type": type.googleapis.com/envoy.extensions.access_loggers.file.v3.FileAccessLog
                      path: /dev/stdout
                      log_format:
                        text_format_source:
                          inline_string: "{\"authority\":\"%REQ(X-ENVOY-ORIGINAL-HOST?:AUTHORITY)%\",\"bytes_received\":\"%BYTES_RECEIVED%\",\"bytes_sent\":\"%BYTES_SENT%\",\"downstream_local_address\":\"%DOWNSTREAM_LOCAL_ADDRESS%\",\"downstream_remote_address\":\"%DOWNSTREAM_REMOTE_ADDRESS%\",\"duration\":\"%DURATION%\",\"istio_policy_status\":\"%DYNAMIC_METADATA(istio.mixer:status)%\",\"method\":\"%REQ(:METHOD)%\",\"path\":\"%REQ(X-ENVOY-ORIGINAL-PATH?:PATH)%\",\"protocol\":\"%PROTOCOL%\",\"request_id\":\"%REQ(X-REQUEST-ID)%\",\"requested_server_name\":\"%REQUESTED_SERVER_NAME%\",\"response_code\":\"%RESPONSE_CODE%\",\"response_flags\":\"%RESPONSE_FLAGS%\",\"route_name\":\"%ROUTE_NAME%\",\"start_time\":\"%START_TIME%\",\"trace_id\":\"%REQ(X-B3-TRACEID)%\",\"upstream_cluster\":\"%UPSTREAM_CLUSTER%\",\"upstream_host\":\"%UPSTREAM_HOST%\",\"upstream_local_address\":\"%UPSTREAM_LOCAL_ADDRESS%\",\"upstream_service_time\":\"%RESP(X-ENVOY-UPSTREAM-SERVICE-TIME)%\",\"upstream_transport_failure_reason\":\"%UPSTREAM_TRANSPORT_FAILURE_REASON%\",\"user_agent\":\"%REQ(USER-AGENT)%\",\"x_forwarded_for\":\"%REQ(X-FORWARDED-FOR)%\"}\n"
                # Modify as required
                route_config:
                  name: local_route
                  virtual_hosts:
                    - name: local_service
                      domains: ["*"]
                      routes:
                        - match:
                            prefix: "/"
                          route:
                            cluster: httpbin
                http_filters:
                  - name: wasmtest
                    typed_config:
                      "@type": type.googleapis.com/udpa.type.v1.TypedStruct
                      type_url: type.googleapis.com/envoy.extensions.filters.http.wasm.v3.Wasm
                      value:
                        config:
                          name: wasmtest
                          vm_config:
                            runtime: envoy.wasm.runtime.v8
                            code:
                              local:
                                filename: /etc/envoy/plugin.wasm
                          configuration:
                            "@type": "type.googleapis.com/google.protobuf.StringValue"
                            # Modify as required
                            value: |
{{ .JSONExample }}
                  - name: envoy.filters.http.router
  clusters:
    - name: httpbin
      connect_timeout: 30s
      type: LOGICAL_DNS
      # Comment out the following line to test on v6 networks
      dns_lookup_family: V4_ONLY
      lb_policy: ROUND_ROBIN
      load_assignment:
        cluster_name: httpbin
        endpoints:
          - lb_endpoints:
              - endpoint:
                  address:
                    socket_address:
                      address: httpbin
                      port_value: 80
`
)

type DockerCompose struct {
	TestPath    string
	ProductPath string
}

type Envoy struct {
	JSONExample string
}

func genDockerComposeYAML(d *DockerCompose, dir string) error {
	path := fmt.Sprintf("%s/docker-compose.yaml", dir)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if err = template.Must(template.New("DockerComposeYAML").Parse(dockerComposeYAML)).Execute(f, d); err != nil {
		return err
	}

	return nil
}

func genEnvoyYAML(e *Envoy, dir string) error {
	path := fmt.Sprintf("%s/envoy.yaml", dir)
	f, err := os.Create(path)
	if err != nil {
		panic(fmt.Sprintf("failed to create %q: %v\n", path, err))
	}
	defer f.Close()

	if err = template.Must(template.New("EnvoyYAML").Parse(envoyYAML)).Execute(f, e); err != nil {
		return err
	}

	return nil
}
