package manifest

const (
	ClusterlinkClusterNode = `---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.11.0
  creationTimestamp: null
  name: clusternodes.kosmos.io
spec:
  group: kosmos.io
  names:
    kind: ClusterNode
    listKind: ClusterNodeList
    plural: clusternodes
    singular: clusternode
  scope: Cluster
  versions:
  - additionalPrinterColumns:
    - jsonPath: .spec.roles
      name: ROLES
      type: string
    - jsonPath: .spec.interfaceName
      name: INTERFACE
      type: string
    - jsonPath: .spec.ip
      name: IP
      type: string
    name: v1alpha1
    schema:
      openAPIV3Schema:
        properties:
          apiVersion:
            description: 'APIVersion defines the versioned schema of this representation
              of an object. Servers should convert recognized schemas to the latest
              internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
            type: string
          kind:
            description: 'Kind is a string value representing the REST resource this
              object represents. Servers may infer this from the endpoint the client
              submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
            type: string
          metadata:
            type: object
          spec:
            properties:
              clusterName:
                type: string
              interfaceName:
                type: string
              ip:
                type: string
              ip6:
                type: string
              nodeName:
                type: string
              podCIDRs:
                items:
                  type: string
                type: array
              roles:
                items:
                  type: string
                type: array
            type: object
          status:
            type: object
        required:
        - spec
        type: object
    served: true
    storage: true
    subresources:
      status: {}
`

	ClusterlinkCluster = `---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.11.0
  creationTimestamp: null
  name: clusters.kosmos.io
spec:
  group: kosmos.io
  names:
    kind: Cluster
    listKind: ClusterList
    plural: clusters
    singular: cluster
  scope: Cluster
  versions:
  - additionalPrinterColumns:
    - jsonPath: .spec.networkType
      name: NETWORK_TYPE
      type: string
    - jsonPath: .spec.ipFamily
      name: IP_FAMILY
      type: string
    name: v1alpha1
    schema:
      openAPIV3Schema:
        properties:
          apiVersion:
            description: 'APIVersion defines the versioned schema of this representation
              of an object. Servers should convert recognized schemas to the latest
              internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
            type: string
          kind:
            description: 'Kind is a string value representing the REST resource this
              object represents. Servers may infer this from the endpoint the client
              submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
            type: string
          metadata:
            type: object
          spec:
            description: Spec is the specification for the behaviour of the cluster.
            properties:
              bridgeCIDRs:
                default:
                  ip: 220.0.0.0/8
                  ip6: 9470::/16
                properties:
                  ip:
                    type: string
                  ip6:
                    type: string
                required:
                - ip
                - ip6
                type: object
              cni:
                default: calico
                type: string
              defaultNICName:
                default: '*'
                type: string
              globalCIDRsMap:
                additionalProperties:
                  type: string
                type: object
              imageRepository:
                type: string
              ipFamily:
                default: all
                type: string
              kubeconfig:
                format: byte
                type: string
              localCIDRs:
                default:
                  ip: 210.0.0.0/8
                  ip6: 9480::/16
                properties:
                  ip:
                    type: string
                  ip6:
                    type: string
                required:
                - ip
                - ip6
                type: object
              namespace:
                default: clusterlink-system
                type: string
              networkType:
                default: p2p
                enum:
                - p2p
                - gateway
                type: string
              nicNodeNames:
                items:
                  properties:
                    interfaceName:
                      type: string
                    nodeName:
                      items:
                        type: string
                      type: array
                  required:
                  - interfaceName
                  - nodeName
                  type: object
                type: array
              useIPPool:
                default: false
                type: boolean
            type: object
          status:
            description: Status describes the current status of a cluster.
            properties:
              podCIDRs:
                items:
                  type: string
                type: array
              serviceCIDRs:
                items:
                  type: string
                type: array
            type: object
        required:
        - spec
        type: object
    served: true
    storage: true
    subresources: {}
`

	ClusterlinkNodeConfig = `---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.11.0
  creationTimestamp: null
  name: nodeconfigs.kosmos.io
spec:
  group: kosmos.io
  names:
    kind: NodeConfig
    listKind: NodeConfigList
    plural: nodeconfigs
    singular: nodeconfig
  scope: Cluster
  versions:
  - name: v1alpha1
    schema:
      openAPIV3Schema:
        properties:
          apiVersion:
            description: 'APIVersion defines the versioned schema of this representation
              of an object. Servers should convert recognized schemas to the latest
              internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
            type: string
          kind:
            description: 'Kind is a string value representing the REST resource this
              object represents. Servers may infer this from the endpoint the client
              submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
            type: string
          metadata:
            type: object
          spec:
            properties:
              arps:
                items:
                  properties:
                    dev:
                      type: string
                    ip:
                      type: string
                    mac:
                      type: string
                  required:
                  - dev
                  - ip
                  - mac
                  type: object
                type: array
              devices:
                items:
                  properties:
                    addr:
                      type: string
                    bindDev:
                      type: string
                    id:
                      format: int32
                      type: integer
                    mac:
                      type: string
                    name:
                      type: string
                    port:
                      format: int32
                      type: integer
                    type:
                      type: string
                  required:
                  - addr
                  - bindDev
                  - id
                  - mac
                  - name
                  - port
                  - type
                  type: object
                type: array
              fdbs:
                items:
                  properties:
                    dev:
                      type: string
                    ip:
                      type: string
                    mac:
                      type: string
                  required:
                  - dev
                  - ip
                  - mac
                  type: object
                type: array
              iptables:
                items:
                  properties:
                    chain:
                      type: string
                    rule:
                      type: string
                    table:
                      type: string
                  required:
                  - chain
                  - rule
                  - table
                  type: object
                type: array
              routes:
                items:
                  properties:
                    cidr:
                      type: string
                    dev:
                      type: string
                    gw:
                      type: string
                  required:
                  - cidr
                  - dev
                  - gw
                  type: object
                type: array
            type: object
          status:
            properties:
              lastChangeTime:
                format: date-time
                type: string
              lastSyncTime:
                format: date-time
                type: string
            type: object
        required:
        - spec
        type: object
    served: true
    storage: true
    subresources:
      status: {}
`
)
