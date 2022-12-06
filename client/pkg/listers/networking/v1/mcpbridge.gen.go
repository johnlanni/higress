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

// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/alibaba/higress/client/pkg/apis/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// McpBridgeLister helps list McpBridges.
type McpBridgeLister interface {
	// List lists all McpBridges in the indexer.
	List(selector labels.Selector) (ret []*v1.McpBridge, err error)
	// McpBridges returns an object that can list and get McpBridges.
	McpBridges(namespace string) McpBridgeNamespaceLister
	McpBridgeListerExpansion
}

// mcpBridgeLister implements the McpBridgeLister interface.
type mcpBridgeLister struct {
	indexer cache.Indexer
}

// NewMcpBridgeLister returns a new McpBridgeLister.
func NewMcpBridgeLister(indexer cache.Indexer) McpBridgeLister {
	return &mcpBridgeLister{indexer: indexer}
}

// List lists all McpBridges in the indexer.
func (s *mcpBridgeLister) List(selector labels.Selector) (ret []*v1.McpBridge, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.McpBridge))
	})
	return ret, err
}

// McpBridges returns an object that can list and get McpBridges.
func (s *mcpBridgeLister) McpBridges(namespace string) McpBridgeNamespaceLister {
	return mcpBridgeNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// McpBridgeNamespaceLister helps list and get McpBridges.
type McpBridgeNamespaceLister interface {
	// List lists all McpBridges in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1.McpBridge, err error)
	// Get retrieves the McpBridge from the indexer for a given namespace and name.
	Get(name string) (*v1.McpBridge, error)
	McpBridgeNamespaceListerExpansion
}

// mcpBridgeNamespaceLister implements the McpBridgeNamespaceLister
// interface.
type mcpBridgeNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all McpBridges in the indexer for a given namespace.
func (s mcpBridgeNamespaceLister) List(selector labels.Selector) (ret []*v1.McpBridge, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.McpBridge))
	})
	return ret, err
}

// Get retrieves the McpBridge from the indexer for a given namespace and name.
func (s mcpBridgeNamespaceLister) Get(name string) (*v1.McpBridge, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("mcpbridge"), name)
	}
	return obj.(*v1.McpBridge), nil
}
