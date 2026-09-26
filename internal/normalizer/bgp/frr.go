// orac-data/internal/normalizer/bgp
package bgp

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"orac-data/internal/models"
)

type FRRNormalizer struct {
	// We will probably add configuration here later.
}

func NewFRRNormalizer() *FRRNormalizer {
	return &FRRNormalizer{}
}

func (n *FRRNormalizer) Normalize(rawJSON []byte) (*models.BGPData, error) {
	var frrResp models.FRRBGPRoutesResponse
	if err := json.Unmarshal(rawJSON, &frrResp); err != nil {
		return nil, fmt.Errorf("failed to parse FRR JSON %w", err)
	}

	return n.convertFRRToUnified(&frrResp), nil
}

func (n *FRRNormalizer) convertFRRToUnified(frrResp *models.FRRBGPRoutesResponse) *models.BGPData {
	// Prepare data structures.
	allRoutes := make([]models.UnifiedRoute, 0, frrResp.TotalPaths)
	routesByPrefix := make(map[string][]models.UnifiedRoute)
	bestRoutes := make(map[string]models.UnifiedRoute)

	for cidr, frrPaths := range frrResp.Routes {
		unifiedPaths := make([]models.UnifiedRoute, 0, len(frrPaths))

		for _, frrPath := range frrPaths {
			unified := n.convertFRRPathToUnified(cidr, &frrPath)
			unifiedPaths = append(unifiedPaths, unified)
			allRoutes = append(allRoutes, unified)

			if frrPath.Bestpath {
				bestRoutes[cidr] = unified
			}
		}

		routesByPrefix[cidr] = unifiedPaths
	}

	data := &models.BGPData{
		Vendor:         "frr",
		TotalRoutes:    frrResp.TotalRoutes,
		TotalPaths:     frrResp.TotalPaths,
		Routes:         allRoutes,
		RoutesByPrefix: routesByPrefix,
		BestRoutes:     bestRoutes,
	}

	data.Stats = n.calculateStats(data)
	return data
}

func (n *FRRNormalizer) convertFRRPathToUnified(cidr string, frrPath *models.FRRBGPPath) models.UnifiedRoute {
	asPath, asPathStr, asPathLen := n.parseASPath(frrPath.Path)

	nextHop, nextHopHostname, addressFamily := n.extractNextHop(frrPath)

	unified := models.UnifiedRoute{
		// Network information
		Network:      frrPath.Prefix,
		PrefixLength: frrPath.PrefixLen,
		CIDR:         cidr,

		// Next hop information
		NextHop:         nextHop,
		NextHopHostname: nextHopHostname,

		// Peer information
		PeerID: frrPath.PeerID,

		// BGP Attributes
		Metric: frrPath.Metric,
		Weight: frrPath.Weight,

		// AS path
		ASPath:    asPath,
		ASPathStr: asPathStr,
		ASPathLen: asPathLen,

		// Route status
		Origin:          frrPath.Origin,
		Best:            frrPath.Bestpath,
		Valid:           frrPath.Valid,
		SelectionReason: frrPath.SelectionReason,

		// Address family
		AddressFamily: addressFamily,

		// Preserve vendor-specific metadata
		VendorMetadata: n.extractVendorMetadata(frrPath),
	}

	return unified
}

func (n *FRRNormalizer) parseASPath(pathStr string) ([]int, string, int) {
	if pathStr == "" {
		return []int{}, "", 0
	}

	// Split by spaces
	parts := strings.Fields(pathStr)
	asPath := make([]int, 0, len(parts))

	for _, part := range parts {
		asn, err := strconv.Atoi(part)
		if err != nil {
			// If parsing fails, skip this ASN (could be a special value)
			continue
		}
		asPath = append(asPath, asn)
	}

	return asPath, pathStr, len(asPath)
}

func (n *FRRNormalizer) extractNextHop(frrPath *models.FRRBGPPath) (ip, hostname, afi string) {
	if len(frrPath.Nexthops) == 0 {
		return "", "", ""
	}

	// Use the first next-hop
	nh := frrPath.Nexthops[0]
	return nh.IP, nh.Hostname, nh.AFI
}

func (n *FRRNormalizer) extractVendorMetadata(frrPath *models.FRRBGPPath) map[string]interface{} {
	metadata := make(map[string]interface{})

	if frrPath.SelectionReason != "" {
		metadata["selectionReason"] = frrPath.SelectionReason
	}
	if frrPath.PathFrom != "" {
		metadata["pathFrom"] = frrPath.PathFrom
	}
	if frrPath.Version != 0 {
		metadata["version"] = frrPath.Version
	}
	if len(frrPath.Nexthops) > 1 {
		// Store additional next-hops if multiple
		additionalNextHops := make([]map[string]interface{}, 0)
		for i, nh := range frrPath.Nexthops {
			if i == 0 {
				continue
			}

			additionalNextHops = append(additionalNextHops, map[string]interface{}{
				"ip":       nh.IP,
				"hostname": nh.Hostname,
				"afi":      nh.AFI,
				"used":     nh.Used,
			})
		}

		if len(additionalNextHops) > 0 {
			metadata["additionalNextHops"] = additionalNextHops
		}
	}
	return metadata
}

func (n *FRRNormalizer) calculateStats(data *models.BGPData) models.BGPStats {
	stats := models.BGPStats{
		UniquePrefixes:     len(data.RoutesByPrefix),
		TotalPaths:         len(data.Routes),
		PrefixDistribution: make(map[int]int),
	}

	if len(data.Routes) == 0 {
		return stats
	}

	// Track ASNs and path lengths.
	asnSet := make(map[int]bool)
	totalASPathLen := 0
	maxASPathLen := 0

	for _, route := range data.Routes {
		stats.PrefixDistribution[route.PrefixLength]++

		// Track ASNs
		for _, asn := range route.ASPath {
			asnSet[asn] = true
		}

		asLen := route.ASPathLen
		totalASPathLen += asLen
		if asLen > maxASPathLen {
			maxASPathLen = asLen
		}
	}

	stats.MaxASPathLen = maxASPathLen
	if len(data.Routes) > 0 {
		stats.AverageASPathLen = float64(totalASPathLen) / float64(len(data.Routes))
	}

	// Convert ASN set to slice
	stats.UniqueASNs = make([]int, 0, len(asnSet))
	for asn := range asnSet {
		stats.UniqueASNs = append(stats.UniqueASNs, asn)
	}

	return stats
}
