// orac-data/internal/normalizer/iface
package iface

import (
	"orac-data/internal/models"
	ifacemodel "orac-data/internal/models/iface"
)

// FRRInterfaceNormalizer converts FRR interface data into the unified format.
type FRRInterfaceNormalizer struct {
	// Future: configuration options.
}

func NewFRRInterfaceNormalizer() *FRRInterfaceNormalizer {
	return &FRRInterfaceNormalizer{}
}

// Converts an FRR interface response into unified InterfaceData.
func (n *FRRInterfaceNormalizer) Normalize(resp *ifacemodel.FRRInterfaceResponse) (*models.InterfaceData, error) {
	interfaces := make([]models.UnifiedInterface, 0, len(resp.Lib.Interface))
	interfacesByName := make(map[string]models.UnifiedInterface, len(resp.Lib.Interface))

	for _, frrIface := range resp.Lib.Interface {
		unified := n.convertFRRInterfaceToUnified(&frrIface)
		interfaces = append(interfaces, unified)
		interfacesByName[unified.Name] = unified
	}

	data := &models.InterfaceData{
		Vendor:           "frr",
		Interfaces:       interfaces,
		InterfacesByName: interfacesByName,
	}

	data.Stats = n.calculateStats(data)
	return data, nil
}

func (n *FRRInterfaceNormalizer) convertFRRInterfaceToUnified(frrIface *ifacemodel.FRRInterface) models.UnifiedInterface {
	return models.UnifiedInterface{
		Name:           frrIface.Name,
		VRF:            frrIface.VRF,
		Index:          frrIface.State.IfIndex,
		MTU:            frrIface.State.MTU,
		MTU6:           frrIface.State.MTU6,
		Speed:          frrIface.State.Speed,
		Metric:         frrIface.State.Metric,
		MACAddress:     frrIface.State.PhyAddress,
		VendorMetadata: n.extractVendorMetadata(frrIface),
	}
}

func (n *FRRInterfaceNormalizer) extractVendorMetadata(frrIface *ifacemodel.FRRInterface) map[string]interface{} {
	// Currently no extra fields beyond what the unified model captures.
	return nil
}

func (n *FRRInterfaceNormalizer) calculateStats(data *models.InterfaceData) models.InterfaceStats {
	stats := models.InterfaceStats{
		Total: len(data.Interfaces),
		ByVRF: make(map[string]int),
	}

	for _, ifc := range data.Interfaces {
		stats.ByVRF[ifc.VRF]++
	}

	return stats
}
