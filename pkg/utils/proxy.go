package utils

var KpandaEgressSVCName = GetEnvWithDefault("GATEWAY_EGRESS_SVC_NAME", "kpanda-egress.kpanda-system.svc.cluster.local") // L4 proxy, egress svc name

var KpandaEgressAssignPortMin = GetEnvWithDefault("GATEWAY_EGRESS_ASSIGN_PORT_MIN", "20000")

var KpandaEgressAssignPortMax = GetEnvWithDefault("GATEWAY_EGRESS_ASSIGN_PORT_MAX", "60000")

var KpandaAPIServerSVCName = GetEnvWithDefault("KPANDA_APISERVER_SVC_NAME", "kpanda-egress.kpanda-system.svc.cluster.local") // kpanda api svc name

var KpandaAPIServerSVCPort = GetEnvWithDefault("KPANDA_APISERVER_SVC_PORT", "80") // port which kpanda svc export to

var NginxPortForProxyKpandaAPIServer = GetEnvWithDefault("PORT_INGRESS_KPANDA", "8000") // ingress nginx  run in this port
