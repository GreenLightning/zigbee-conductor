package zdp

import "fmt"

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

func ClusterName(clusterID uint16) string {
	switch clusterID {
	case ClusterNWKAddrReq:
		return "NWKAddrReq"
	case ClusterIEEEAddrReq:
		return "IEEEAddrReq"
	case ClusterNodeDescReq:
		return "NodeDescReq"
	case ClusterPowerDescReq:
		return "PowerDescReq"
	case ClusterSimpleDescReq:
		return "SimpleDescReq"
	case ClusterActiveEPReq:
		return "ActiveEPReq"
	case ClusterMatchDescReq:
		return "MatchDescReq"
	case ClusterComplexDescReq:
		return "ComplexDescReq"
	case ClusterUserDescReq:
		return "UserDescReq"
	case ClusterDiscoveryCacheReq:
		return "DiscoveryCacheReq"
	case ClusterDeviceAnnce:
		return "DeviceAnnce"
	case ClusterParentAnnce:
		return "ParentAnnce"
	case ClusterUserDescSet:
		return "UserDescSet"
	case ClusterSystemServerDiscoveryReq:
		return "SystemServerDiscoveryReq"
	case ClusterDiscoveryStoreReq:
		return "DiscoveryStoreReq"
	case ClusterNodeDescStoreReq:
		return "NodeDescStoreReq"
	case ClusterPowerDescStoreReq:
		return "PowerDescStoreReq"
	case ClusterActiveEPStoreReq:
		return "ActiveEPStoreReq"
	case ClusterSimpleDescStoreReq:
		return "SimpleDescStoreReq"
	case ClusterRemoveNodeCacheReq:
		return "RemoveNodeCacheReq"
	case ClusterFindNodeCacheReq:
		return "FindNodeCacheReq"
	case ClusterExtendedSimpleDescReq:
		return "ExtendedSimpleDescReq"
	case ClusterExtendedActiveEPReq:
		return "ExtendedActiveEPReq"
	case ClusterEndDeviceBindReq:
		return "EndDeviceBindReq"
	case ClusterBindReq:
		return "BindReq"
	case ClusterUnbindReq:
		return "UnbindReq"
	case ClusterBindRegisterReq:
		return "BindRegisterReq"
	case ClusterReplaceDeviceReq:
		return "ReplaceDeviceReq"
	case ClusterStoreBkupBindEntryReq:
		return "StoreBkupBindEntryReq"
	case ClusterRemoveBkupBindEntryReq:
		return "RemoveBkupBindEntryReq"
	case ClusterBackupBindTableReq:
		return "BackupBindTableReq"
	case ClusterRecoverBindTableReq:
		return "RecoverBindTableReq"
	case ClusterBackupSourceBindReq:
		return "BackupSourceBindReq"
	case ClusterRecoverSourceBindReq:
		return "RecoverSourceBindReq"
	case ClusterMgmtNWKDiscReq:
		return "MgmtNWKDiscReq"
	case ClusterMgmtLqiReq:
		return "MgmtLqiReq"
	case ClusterMgmtRtgReq:
		return "MgmtRtgReq"
	case ClusterMgmtBindReq:
		return "MgmtBindReq"
	case ClusterMgmtLeaveReq:
		return "MgmtLeaveReq"
	case ClusterMgmtDirectJoinReq:
		return "MgmtDirectJoinReq"
	case ClusterMgmtPermitJoiningReq:
		return "MgmtPermitJoiningReq"
	case ClusterMgmtCacheReq:
		return "MgmtCacheReq"
	case ClusterMgmtNWKUpdateReq:
		return "MgmtNWKUpdateReq"
	case ClusterMgmtNWKEnhancedUpdateReq:
		return "MgmtNWKEnhancedUpdateReq"
	case ClusterMgmtNWKIEEEJoiningListReq:
		return "MgmtNWKIEEEJoiningListReq"
	case ClusterNWKAddrRsp:
		return "NWKAddrRsp"
	case ClusterIEEEAddrRsp:
		return "IEEEAddrRsp"
	case ClusterNodeDescRsp:
		return "NodeDescRsp"
	case ClusterPowerDescRsp:
		return "PowerDescRsp"
	case ClusterSimpleDescRsp:
		return "SimpleDescRsp"
	case ClusterActiveEPRsp:
		return "ActiveEPRsp"
	case ClusterMatchDescRsp:
		return "MatchDescRsp"
	case ClusterComplexDescRsp:
		return "ComplexDescRsp"
	case ClusterUserDescRsp:
		return "UserDescRsp"
	case ClusterUserDescConf:
		return "UserDescConf"
	case ClusterParentAnnceRsp:
		return "ParentAnnceRsp"
	case ClusterSystemServerDiscoveryRsp:
		return "SystemServerDiscoveryRsp"
	case ClusterDiscoveryStoreRsp:
		return "DiscoveryStoreRsp"
	case ClusterNodeDescStoreRsp:
		return "NodeDescStoreRsp"
	case ClusterPowerDescStoreRsp:
		return "PowerDescStoreRsp"
	case ClusterActiveEPStoreRsp:
		return "ActiveEPStoreRsp"
	case ClusterSimpleDescStoreRsp:
		return "SimpleDescStoreRsp"
	case ClusterRemoveNodeCacheRsp:
		return "RemoveNodeCacheRsp"
	case ClusterFindNodeCacheRsp:
		return "FindNodeCacheRsp"
	case ClusterExtendedSimpleDescRsp:
		return "ExtendedSimpleDescRsp"
	case ClusterExtendedActiveEPRsp:
		return "ExtendedActiveEPRsp"
	case ClusterEndDeviceBindRsp:
		return "EndDeviceBindRsp"
	case ClusterBindRsp:
		return "BindRsp"
	case ClusterUnbindRsp:
		return "UnbindRsp"
	case ClusterBindRegisterRsp:
		return "BindRegisterRsp"
	case ClusterReplaceDeviceRsp:
		return "ReplaceDeviceRsp"
	case ClusterStoreBkupBindEntryRsp:
		return "StoreBkupBindEntryRsp"
	case ClusterRemoveBkupBindEntryRsp:
		return "RemoveBkupBindEntryRsp"
	case ClusterBackupBindTableRsp:
		return "BackupBindTableRsp"
	case ClusterRecoverBindTableRsp:
		return "RecoverBindTableRsp"
	case ClusterBackupSourceBindRsp:
		return "BackupSourceBindRsp"
	case ClusterRecoverSourceBindRsp:
		return "RecoverSourceBindRsp"
	case ClusterMgmtNWKDiscRsp:
		return "MgmtNWKDiscRsp"
	case ClusterMgmtLqiRsp:
		return "MgmtLqiRsp"
	case ClusterMgmtRtgRsp:
		return "MgmtRtgRsp"
	case ClusterMgmtBindRsp:
		return "MgmtBindRsp"
	case ClusterMgmtLeaveRsp:
		return "MgmtLeaveRsp"
	case ClusterMgmtDirectJoinRsp:
		return "MgmtDirectJoinRsp"
	case ClusterMgmtPermitJoiningRsp:
		return "MgmtPermitJoiningRsp"
	case ClusterMgmtCacheRsp:
		return "MgmtCacheRsp"
	case ClusterMgmtNWKUpdateNotify:
		return "MgmtNWKUpdateNotify"
	case ClusterMgmtNWKEnhancedUpdateNotify:
		return "MgmtNWKEnhancedUpdateNotify"
	case ClusterMgmtNWKIEEEJoiningListRsp:
		return "MgmtNWKIEEEJoiningListRsp"
	case ClusterMgmtNWKUnsolicitedEnhancedUpdateNotify:
		return "MgmtNWKUnsolicitedEnhancedUpdateNotify"
	default:
		return fmt.Sprintf("ClusterID(0x%04x)", uint16(clusterID))
	}
}
