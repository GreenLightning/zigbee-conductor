package zdp

const (
	ClusterNWKAddrReq               = 0x0000
	ClusterIEEEAddrReq              = 0x0001
	ClusterNodeDescReq              = 0x0002
	ClusterPowerDescReq             = 0x0003
	ClusterSimpleDescReq            = 0x0004
	ClusterActiveEPReq              = 0x0005
	ClusterMatchDescReq             = 0x0006
	ClusterComplexDescReq           = 0x0010
	ClusterUserDescReq              = 0x0011
	ClusterDiscoveryCacheReq        = 0x0012
	ClusterDeviceAnnce              = 0x0013
	ClusterParentAnnce              = 0x001F
	ClusterUserDescSet              = 0x0014
	ClusterSystemServerDiscoveryReq = 0x0015
	ClusterDiscoveryStoreReq        = 0x0016
	ClusterNodeDescStoreReq         = 0x0017
	ClusterPowerDescStoreReq        = 0x0018
	ClusterActiveEPStoreReq         = 0x0019
	ClusterSimpleDescStoreReq       = 0x001a
	ClusterRemoveNodeCacheReq       = 0x001b
	ClusterFindNodeCacheReq         = 0x001c
	ClusterExtendedSimpleDescReq    = 0x001d
	ClusterExtendedActiveEPReq      = 0x001e

	ClusterEndDeviceBindReq       = 0x0020
	ClusterBindReq                = 0x0021
	ClusterUnbindReq              = 0x0022
	ClusterBindRegisterReq        = 0x0023
	ClusterReplaceDeviceReq       = 0x0024
	ClusterStoreBkupBindEntryReq  = 0x0025
	ClusterRemoveBkupBindEntryReq = 0x0026
	ClusterBackupBindTableReq     = 0x0027
	ClusterRecoverBindTableReq    = 0x0028
	ClusterBackupSourceBindReq    = 0x0029
	ClusterRecoverSourceBindReq   = 0x002a

	ClusterMgmtNWKDiscReq            = 0x0030
	ClusterMgmtLqiReq                = 0x0031
	ClusterMgmtRtgReq                = 0x0032
	ClusterMgmtBindReq               = 0x0033
	ClusterMgmtLeaveReq              = 0x0034
	ClusterMgmtDirectJoinReq         = 0x0035
	ClusterMgmtPermitJoiningReq      = 0x0036
	ClusterMgmtCacheReq              = 0x0037
	ClusterMgmtNWKUpdateReq          = 0x0038
	ClusterMgmtNWKEnhancedUpdateReq  = 0x0039
	ClusterMgmtNWKIEEEJoiningListReq = 0x003a
)

const (
	ClusterNWKAddrRsp               = 0x8000
	ClusterIEEEAddrRsp              = 0x8001
	ClusterNodeDescRsp              = 0x8002
	ClusterPowerDescRsp             = 0x8003
	ClusterSimpleDescRsp            = 0x8004
	ClusterActiveEPRsp              = 0x8005
	ClusterMatchDescRsp             = 0x8006
	ClusterComplexDescRsp           = 0x8010
	ClusterUserDescRsp              = 0x8011
	ClusterUserDescConf             = 0x8014
	ClusterParentAnnceRsp           = 0x801F
	ClusterSystemServerDiscoveryRsp = 0x8015
	ClusterDiscoveryStoreRsp        = 0x8016
	ClusterNodeDescStoreRsp         = 0x8017
	ClusterPowerDescStoreRsp        = 0x8018
	ClusterActiveEPStoreRsp         = 0x8019
	ClusterSimpleDescStoreRsp       = 0x801a
	ClusterRemoveNodeCacheRsp       = 0x801b
	ClusterFindNodeCacheRsp         = 0x801c
	ClusterExtendedSimpleDescRsp    = 0x801d
	ClusterExtendedActiveEPRsp      = 0x801e

	ClusterEndDeviceBindRsp       = 0x8020
	ClusterBindRsp                = 0x8021
	ClusterUnbindRsp              = 0x8022
	ClusterBindRegisterRsp        = 0x8023
	ClusterReplaceDeviceRsp       = 0x8024
	ClusterStoreBkupBindEntryRsp  = 0x8025
	ClusterRemoveBkupBindEntryRsp = 0x8026
	ClusterBackupBindTableRsp     = 0x8027
	ClusterRecoverBindTableRsp    = 0x8028
	ClusterBackupSourceBindRsp    = 0x8029
	ClusterRecoverSourceBindRsp   = 0x802a

	ClusterMgmtNWKDiscRsp                         = 0x8030
	ClusterMgmtLqiRsp                             = 0x8031
	ClusterMgmtRtgRsp                             = 0x8032
	ClusterMgmtBindRsp                            = 0x8033
	ClusterMgmtLeaveRsp                           = 0x8034
	ClusterMgmtDirectJoinRsp                      = 0x8035
	ClusterMgmtPermitJoiningRsp                   = 0x8036
	ClusterMgmtCacheRsp                           = 0x8037
	ClusterMgmtNWKUpdateNotify                    = 0x8038
	ClusterMgmtNWKEnhancedUpdateNotify            = 0x8039
	ClusterMgmtNWKIEEEJoiningListRsp              = 0x803a
	ClusterMgmtNWKUnsolicitedEnhancedUpdateNotify = 0x803b
)
