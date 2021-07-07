package ipaddr

import (
	"sync/atomic"
	"unsafe"
)

// All IP address strings corresponds to exactly one of these types.
// In cases where there is no corresponding default IPAddress value (INVALID, ALL, and possibly EMPTY), these types can be used for comparison.
// EMPTY means a zero-length string (useful for validation, we can set validation to allow empty strings) that has no corresponding IPAddress value (validation options allow you to map empty to the loopback)
// INVALID means it is known that it is not any of the other allowed types (validation options can restrict the allowed types)
// ALL means it is wildcard(s) with no separators, like "*", which represents all addresses, whether IPv4, IPv6 or other, and thus has no corresponding IPAddress value
// These constants are ordered by address space size, from smallest to largest, and the ordering affects comparisons
type IPType int

func fromVersion(version IPVersion) IPType {
	switch version {
	case IPv4:
		return IPV4
	case IPv6:
		return IPV6
	default:
	}
	return UNINITIALIZED_TYPE
}

func (t IPType) isUnknown() bool {
	return t == UNINITIALIZED_TYPE
}

const (
	UNINITIALIZED_TYPE IPType = iota
	INVALID
	EMPTY
	IPV4
	IPV6
	//PREFIX_ONLY
	ALL
)

//TODO rename later IPAddressProvider, IPType, and the IPType constants, all the creator classes, etc, so not public, also same for macAddressProvider

type IPAddressProvider interface {
	getType() IPType

	getProviderHostAddress() (*IPAddress, IncompatibleAddressError)

	getProviderAddress() (*IPAddress, IncompatibleAddressError)

	getVersionedAddress(version IPVersion) (*IPAddress, IncompatibleAddressError)

	isSequential() bool

	getProviderSeqRange() *IPAddressSeqRange

	getProviderMask() *IPAddress

	// TODO getDivisionGrouping
	//default IPAddressDivisionSeries getDivisionGrouping() throws IncompatibleAddressError {
	//	return getProviderAddress();
	//}

	providerCompare(IPAddressProvider) (int, IncompatibleAddressError)

	providerEquals(IPAddressProvider) (bool, IncompatibleAddressError)

	getProviderIPVersion() IPVersion

	isProvidingIPAddress() bool

	isProvidingIPv4() bool

	isProvidingIPv6() bool

	isProvidingAllAddresses() bool

	isProvidingEmpty() bool

	isProvidingMixedIPv6() bool

	isProvidingBase85IPv6() bool

	getProviderNetworkPrefixLength() PrefixLen

	isInvalid() bool

	//isUninitialized() bool

	// If the address was created by parsing, this provides the parameters used when creating the address,
	// otherwise nil
	getParameters() IPAddressStringParameters
}

// TODO optimized contains: add these later
//	/**
//	 * An optimized contains that does not need to create address objects to return an answer.
//	 * Unconventional addresses may require that the address objects are created, in such cases null is returned.
//	 *
//	 * Addresses constructed from canonical or normalized representations with no wildcards will not return null.
//	 *
//	 * @param other
//	 * @return
//	 */
//	default Boolean contains(IPAddressProvider other) {
//		return null;
//	}
//
//	/**
//	 * An optimized contains that does not need to fully parse the other address to return an answer.
//	 *
//	 * Unconventional addresses may require full parsing, in such cases null is returned.
//	 *
//	 * Addresses constructed from canonical or normalized representations with no wildcards will not return null.
//	 *
//	 * @param other
//	 * @return
//	 */
//	default Boolean contains(String other) {
//		return null;
//	}
//
//	/**
//	 * An optimized prefix comparison that does not need to fully parse the other address to return an answer.
//	 *
//	 * Unconventional addresses may require full parsing, in such cases null is returned.
//	 *
//	 * Addresses constructed from canonical or normalized representations with no wildcards will not return null.
//	 *
//	 * @param other
//	 * @return
//	 */
//	default Boolean prefixEquals(String other) {
//		return null;
//	}
//
//	/**
//	 * An optimized prefix comparison that does not need to create addresses to return an answer.
//	 *
//	 * Unconventional addresses may require the address objects, in such cases null is returned.
//	 *
//	 * @param other
//	 * @return
//	 */
//	default Boolean prefixEquals(IPAddressProvider other) {
//		return null;
//	}
//
//	/**
//	 * An optimized prefix comparison that does not need to create addresses to return an answer.
//	 *
//	 * Unconventional addresses may require the address objects, in such cases null is returned.
//	 *
//	 * @param other
//	 * @return
//	 */
//	default Boolean prefixContains(String other) {
//		return null;
//	}
//
//	/**
//	 * An optimized prefix comparison that does not need to create addresses to return an answer.
//	 *
//	 * Unconventional addresses may require the address objects, in such cases null is returned.
//	 *
//	 * @param other
//	 * @return
//	 */
//	default Boolean prefixContains(IPAddressProvider other) {
//		return null;
//	}
//
//	/**
//	 * An optimized equality comparison that does not need to create addresses to return an answer.
//	 *
//	 * Unconventional addresses may require the address objects, in such cases null is returned.
//	 *
//	 * @param other
//	 * @return
//	 */
//	default Boolean parsedEquals(IPAddressProvider other) {
//		return null;
//	}
//
//	default boolean hasPrefixSeparator() {
//		return getProviderNetworkPrefixLength() != null;
//	}

type ipAddrProvider struct{}

func (p *ipAddrProvider) getType() IPType {
	return UNINITIALIZED_TYPE
}

func (p *ipAddrProvider) isSequential() bool {
	return false
}

func (p *ipAddrProvider) getProviderHostAddress() (*IPAddress, IncompatibleAddressError) {
	return nil, nil
}

func (p *ipAddrProvider) getProviderAddress() (*IPAddress, IncompatibleAddressError) {
	return nil, nil
}

func (p *ipAddrProvider) getProviderSeqRange() *IPAddressSeqRange {
	return nil
}

func (p *ipAddrProvider) getVersionedAddress(version IPVersion) (*IPAddress, IncompatibleAddressError) {
	return nil, nil
}

func (p *ipAddrProvider) getProviderMask() *IPAddress {
	return nil
}

func (p *ipAddrProvider) getProviderIPVersion() IPVersion {
	return INDETERMINATE_VERSION
}

func (p *ipAddrProvider) isProvidingIPAddress() bool {
	return false
}

func (p *ipAddrProvider) isProvidingIPv4() bool {
	return false
}

func (p *ipAddrProvider) isProvidingIPv6() bool {
	return false
}

func (p *ipAddrProvider) isProvidingAllAddresses() bool {
	return false
}

func (p *ipAddrProvider) isProvidingEmpty() bool {
	return false
}

func (p *ipAddrProvider) isInvalid() bool {
	return false
}

//func (p *ipAddrProvider) isUninitialized() bool {
//	return false
//}

func (p *ipAddrProvider) isProvidingMixedIPv6() bool {
	return false
}

func (p *ipAddrProvider) isProvidingBase85IPv6() bool {
	return false
}

func (p *ipAddrProvider) getProviderNetworkPrefixLength() PrefixLen {
	return nil
}

func (p *ipAddrProvider) getParameters() IPAddressStringParameters {
	return nil
}

func providerCompare(p, other IPAddressProvider) (res int, err IncompatibleAddressError) {
	if p == other {
		return
	}
	value, err := p.getProviderAddress()
	if err != nil {
		return
	}
	if value != nil {
		var otherValue *IPAddress
		otherValue, err = other.getProviderAddress()
		if err != nil {
			return
		}
		if otherValue != nil {
			//TODO compareTo on address
			//return value.compareTo(otherValue);
		}
	}
	var thisType, otherType IPType = p.getType(), other.getType()
	res = int(thisType - otherType)
	return
}

/**
* When a value provider produces no value, equality and comparison are based on the enum IPType,
* which can by null.
* @param o
* @return
 */
func providerEquals(p, other IPAddressProvider) (res bool, err IncompatibleAddressError) {
	if p == other {
		res = true
		return
	}
	value, err := p.getProviderAddress()
	if err != nil {
		return
	}
	if value != nil {
		var otherValue *IPAddress
		otherValue, err = other.getProviderAddress()
		if err != nil {
			return
		}
		if otherValue != nil {
			// TODO equals
			return
		} else {
			return
		}
	}
	res = p.getType() == other.getType()
	return
}

// if you have a type with 3 funcs, and 3 methods that defer to the funs
// then that is 4 decls, and then you can deine each of the 3 vars
// if you do a new type for each overridden method, that is 6 decls

type nullProvider struct {
	ipAddrProvider

	ipType                IPType
	isInvalidVal, isEmpty bool
	//isInvalidVal, isUninitializedVal, isEmpty bool
}

func (p *nullProvider) isInvalid() bool {
	return p.isInvalidVal
}

//func (p *nullProvider) isUninitialized() bool {
//	return p.isUninitializedVal
//}

func (p *nullProvider) isProvidingEmpty() bool {
	return p.isEmpty
}

func (p *nullProvider) getType() IPType {
	return p.ipType
}

func (p *nullProvider) providerCompare(other IPAddressProvider) (int, IncompatibleAddressError) {
	return providerCompare(p, other)
}

func (p *nullProvider) providerEquals(other IPAddressProvider) (bool, IncompatibleAddressError) {
	return providerEquals(p, other)
}

var (
	INVALID_PROVIDER = &nullProvider{isInvalidVal: true, ipType: INVALID}
	//NO_TYPE_PROVIDER = &nullProvider{isUninitializedVal: true, ipType: UNINITIALIZED_TYPE}
	EMPTY_PROVIDER = &nullProvider{isEmpty: true, ipType: EMPTY}
)

//type CachedIPAddresses struct {
//
//	//address is 1.2.0.0/16 and hostAddress is 1.2.3.4 for the string 1.2.3.4/16
//	address, hostAddress *IPAddress
//}
//
//func (cached *CachedIPAddresses) getAddress() *IPAddress {
//	return cached.address
//}
//
//func (cached *CachedIPAddresses) getHostAddress() *IPAddress {
//	return cached.hostAddress
//}

///**
//	 * Wraps an IPAddress for IPAddressString in the cases where no parsing is provided, the address exists already
//	 * @param value
//	 * @return
//	 */
func getProviderFor(address, hostAddress *IPAddress) IPAddressProvider {
	return &cachedAddressProvider{address: address, hostAddress: hostAddress}
}

type cachedAddressProvider struct {
	ipAddrProvider

	// addressCreator creates two addresses, the host address and address with prefix/mask, at the same time
	//addressCreator func() CachedIPAddresses
	addressCreator func() (address, hostAddress *IPAddress)

	//xxx must make this a pointer I think even though in some cases we jsut provide it off the bar, no sync required xxx
	//xxx but in the cases we do not, need to synhronize atomically on a ptr to the data we create
	//xxxx

	//xxx
	//cachedValues  CachedIPAddresses
	//createdValues *CachedIPAddresses

	address, hostAddress *IPAddress

	//CreationLock
}

//TODO do not forget you also need these two in all top level classes, including ParsedIPAddress, the mask, all and empty providers
// they are needed becaue of virtual calls to getType() and getProviderAddress()

func (cached *cachedAddressProvider) providerCompare(other IPAddressProvider) (int, IncompatibleAddressError) {
	return providerCompare(cached, other)
}

func (cached *cachedAddressProvider) providerEquals(other IPAddressProvider) (bool, IncompatibleAddressError) {
	return providerEquals(cached, other)
}

func (cached *cachedAddressProvider) isProvidingIPAddress() bool {
	return true
}

func (cached *cachedAddressProvider) getVersionedAddress(version IPVersion) (*IPAddress, IncompatibleAddressError) {
	thisVersion := cached.getProviderIPVersion()
	if version != thisVersion {
		return nil, nil
	}
	return cached.getProviderAddress()
}

func (cached *cachedAddressProvider) getProviderSeqRange() *IPAddressSeqRange {
	addr, _ := cached.getProviderAddress()
	if addr != nil {
		return addr.ToSequentialRange()
	}
	return nil
}

func (cached *cachedAddressProvider) isSequential() bool {
	addr, _ := cached.getProviderAddress()
	if addr != nil {
		return addr.IsSequential()
	}
	return false
}

//func (cached *cachedAddressProvider) hasCachedAddresses() bool {
//	return cached.addressCreator == nil || cached.isItemCreated()
//}

func (cached *cachedAddressProvider) getProviderHostAddress() (*IPAddress, IncompatibleAddressError) {
	res := cached.hostAddress
	if res == nil {
		_, res = cached.getCachedAddresses()
	}
	return res, nil
	//return cached.getCachedAddresses().getHostAddress(), nil
}

func (cached *cachedAddressProvider) getProviderAddress() (*IPAddress, IncompatibleAddressError) {
	res := cached.address
	if res == nil {
		res, _ = cached.getCachedAddresses()
	}
	return res, nil
	//return cached.getCachedAddresses().getAddress(), nil
}

func (cached *cachedAddressProvider) getCachedAddresses() (address, hostAddress *IPAddress) {
	if cached.addressCreator == nil {
		address, hostAddress = cached.addressCreator()
		dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&cached.address))
		atomic.StorePointer(dataLoc, unsafe.Pointer(address))
		dataLoc = (*unsafe.Pointer)(unsafe.Pointer(&cached.hostAddress))
		atomic.StorePointer(dataLoc, unsafe.Pointer(hostAddress))
	}
	return
	//xxx
	/*
		networkMaskLen, hostMaskLen := section.checkForPrefixMask()
			res := &maskLenSetting{networkMaskLen, hostMaskLen}
			dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&cache.cachedMaskLens))
			atomic.StorePointer(dataLoc, unsafe.Pointer(res))
	*/
	//if cached.addressCreator != nil && !cached.isItemCreated() {
	//	cached.create(func() {
	//		cached.values = cached.addressCreator()
	//	})
	//}
	//return &cached.values
}

func (cached *cachedAddressProvider) getProviderNetworkPrefixLength() (p PrefixLen) {
	addr, _ := cached.getProviderAddress()
	return addr.GetNetworkPrefixLength()
}

func (cached *cachedAddressProvider) getProviderIPVersion() IPVersion {
	addr, _ := cached.getProviderAddress()
	return addr.getIPVersion()
}

func (cached *cachedAddressProvider) getType() IPType {
	return fromVersion(cached.getProviderIPVersion())
}

func (cached *cachedAddressProvider) isProvidingIPv4() bool {
	addr, _ := cached.getProviderAddress()
	return addr.IsIPv4()
}

func (cached *cachedAddressProvider) isProvidingIPv6() bool {
	addr, _ := cached.getProviderAddress()
	return addr.IsIPv6()
}

type VersionedAddressCreator struct {
	cachedAddressProvider

	adjustedVersion IPVersion

	versionedAddressCreator func(IPVersion) *IPAddress

	//createdVersioned [2]CreationLock
	versionedValues [2]*IPAddress

	parameters IPAddressStringParameters
}

func (versioned *VersionedAddressCreator) getParameters() IPAddressStringParameters {
	return versioned.parameters
}

func (versioned *VersionedAddressCreator) isProvidingIPAddress() bool {
	return versioned.adjustedVersion != INDETERMINATE_VERSION
}

func (versioned *VersionedAddressCreator) isProvidingIPv4() bool {
	return versioned.adjustedVersion == IPv4
}

func (versioned *VersionedAddressCreator) isProvidingIPv6() bool {
	return versioned.adjustedVersion == IPv6
}

func (versioned *VersionedAddressCreator) getProviderIPVersion() IPVersion {
	return versioned.adjustedVersion
}

func (versioned *VersionedAddressCreator) getType() IPType {
	return fromVersion(versioned.adjustedVersion)
}

func (versioned *VersionedAddressCreator) getVersionedAddress(version IPVersion) (addr *IPAddress, err IncompatibleAddressError) {
	index := version.index()
	if index >= INDETERMINATE_VERSION.index() {
		return
	}
	if versioned.versionedAddressCreator != nil {
		cached := versioned.versionedValues[index]
		if cached == nil {
			res := versioned.versionedAddressCreator(version)
			dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&versioned.versionedValues[index]))
			atomic.StorePointer(dataLoc, unsafe.Pointer(res))
		}
	}

	//xxxx
	//if versioned.versionedAddressCreator != nil && !versioned.createdVersioned[index].isItemCreated() {
	//	versioned.createdVersioned[index].create(func() {
	//		xxxx
	//		versioned.versionedValues[index] = versioned.versionedAddressCreator(version)
	//	})
	//}

	addr = versioned.versionedValues[index]
	return
}

func newLoopbackCreator(options IPAddressStringParameters, zone string) *LoopbackCreator {
	// TODO an option to set preferred loopback here in IPAddressStringParameters, do the same in Java
	// the option will set one of three options, IPv4, IPv6, or INDETERMINATE_VERSION which is the default
	// In Go the default will be IPv4
	// There is another option I wanted to add, was in the validator code, I think allow empty zone with prefix like %/
	var preferIPv6 bool
	ipv6WithZoneLoop := func() *IPAddress {
		network := options.GetIPv6Parameters().GetNetwork()
		creator := network.GetIPv6AddressCreator()
		return creator.createAddressInternalFromBytes(network.GetLoopback().GetBytes(), zone)
	}
	ipv6Loop := func() *IPAddress {
		return options.GetIPv6Parameters().GetNetwork().GetLoopback()
	}
	ipv4Loop := func() *IPAddress {
		return options.GetIPv4Parameters().GetNetwork().GetLoopback()
	}
	double := func(one *IPAddress) (address, hostAddress *IPAddress) {
		return one, one
	}
	var addrCreator func() (address, hostAddress *IPAddress)
	var version IPVersion
	if len(zone) > 0 && preferIPv6 {
		addrCreator = func() (*IPAddress, *IPAddress) { return double(ipv6WithZoneLoop()) }
		version = IPv6
	} else if preferIPv6 {
		addrCreator = func() (*IPAddress, *IPAddress) { return double(ipv6Loop()) }
		version = IPv6
	} else {
		addrCreator = func() (*IPAddress, *IPAddress) { return double(ipv4Loop()) }
		version = IPv4
	}
	cached := cachedAddressProvider{addressCreator: addrCreator}
	versionedAddressCreator := func(version IPVersion) *IPAddress {
		addr := cached.address
		if addr != nil {
			if version == addr.GetIPVersion() {
				return addr
			}
		}
		if version.isIPv4() {
			return ipv4Loop()
		} else if version.isIPv6() {
			if len(zone) > 0 {
				return ipv6WithZoneLoop()
			}
			return ipv6Loop()
		}
		return nil
	}
	return &LoopbackCreator{
		VersionedAddressCreator: VersionedAddressCreator{
			adjustedVersion:         version,
			parameters:              options,
			cachedAddressProvider:   cached,
			versionedAddressCreator: versionedAddressCreator,
		},
		zone: zone,
	}
}

type LoopbackCreator struct {
	VersionedAddressCreator

	zone string
}

func (loop *LoopbackCreator) providerCompare(other IPAddressProvider) (int, IncompatibleAddressError) {
	return providerCompare(loop, other)
}

func (loop *LoopbackCreator) providerEquals(other IPAddressProvider) (bool, IncompatibleAddressError) {
	return providerEquals(loop, other)
}

func (loop *LoopbackCreator) getProviderNetworkPrefixLength() PrefixLen {
	return nil
}

type AdjustedAddressCreator struct {
	VersionedAddressCreator

	networkPrefixLength PrefixLen
}

func (adjusted *AdjustedAddressCreator) getProviderNetworkPrefixLength() PrefixLen {
	return adjusted.networkPrefixLength
}

func (adjusted *AdjustedAddressCreator) getProviderAddress() (*IPAddress, IncompatibleAddressError) {
	if !adjusted.isProvidingIPAddress() {
		return nil, nil
	}
	return adjusted.VersionedAddressCreator.getProviderAddress()
}

func (adjusted *AdjustedAddressCreator) getProviderHostAddress() (*IPAddress, IncompatibleAddressError) {
	if !adjusted.isProvidingIPAddress() {
		return nil, nil
	}
	return adjusted.VersionedAddressCreator.getProviderHostAddress()
}

// TODO the adjusted version passed in is the one adjusted due to zone %, or mask version, or prefix len >= 32
// INside this function we will handle the cases where it is still not determined, and that will be based on our new rules
// involving (a) maybe when < 32 we default to IPv4, otherwise IPv6
//			(b) this behaviour can be overridden by a string parameters option

func newMaskCreator(options IPAddressStringParameters, adjustedVersion IPVersion, networkPrefixLength PrefixLen) *MaskCreator {
	// TODO use the option for  preferred loopback also for preferred mask, do the same in Java
	// TODO also, consider the idea of preflen < 32 defaulting to IPv4, >= 32 IPv6
	// Drop "prefix only" type - it was never a good idea anyway!  Better to prefer one over the other.

	var preferIPv6 bool

	if adjustedVersion == INDETERMINATE_VERSION {
		if preferIPv6 {
			adjustedVersion = IPv6
		} else {
			adjustedVersion = IPv4
		}
	}
	createVersionedMask := func(version IPVersion, prefLen PrefixLen, withPrefixLength bool) *IPAddress {
		if version == IPv4 {
			network := options.GetIPv4Parameters().GetNetwork()
			//if withPrefixLength {
			//	return network.GetNetworkIPAddress(prefLen)
			//}
			//return network.GetNetworkMask(prefLen, false)
			return network.GetNetworkMask(*prefLen)
		} else if version == IPv6 {
			network := options.GetIPv6Parameters().GetNetwork()
			//if withPrefixLength {
			//	return network.GetNetworkIPAddress(prefLen)
			//}
			//return network.GetNetworkMask(prefLen, false)
			return network.GetNetworkMask(*prefLen)
		}
		return nil
	}
	versionedAddressCreator := func(version IPVersion) *IPAddress {
		return createVersionedMask(version, networkPrefixLength, true)
	}
	addrCreator := func() (address, hostAddress *IPAddress) {
		prefLen := networkPrefixLength
		return createVersionedMask(adjustedVersion, prefLen, true),
			createVersionedMask(adjustedVersion, prefLen, false)
	}
	cached := cachedAddressProvider{addressCreator: addrCreator}
	return &MaskCreator{
		AdjustedAddressCreator{
			networkPrefixLength: networkPrefixLength,
			VersionedAddressCreator: VersionedAddressCreator{
				adjustedVersion:         adjustedVersion,
				parameters:              options,
				cachedAddressProvider:   cached,
				versionedAddressCreator: versionedAddressCreator,
			},
		},
	}
}

type MaskCreator struct {
	AdjustedAddressCreator
}

// TODO the adjusted version passed in is the one adjusted due to zone %, or mask version, or prefix len >= 32
// INside this function we will handle the cases where it is still not determined, and that will be based on our new rules
// involving (a) maybe when < 32 we default to IPv4, otherwise IPv6
//			(b) this behaviour can be overridden by a string parameters option

func newAllCreator(qualifier *ParsedHostIdentifierStringQualifier, adjustedVersion IPVersion, originator HostIdentifierString, options IPAddressStringParameters) *AllCreator {
	// TODO use the option for  preferred loopback also for preferred mask, do the same in Java
	// consider using zero value instead of loopback - zero string becomes zero value
	// TODO also, consider the idea of preflen < 32 defaulting to IPv4, >= 32 IPv6
	// Drop "prefix only" type - it was never a good idea anyway!  Better to prefer one over the other.

	var preferIPv6 bool
	if adjustedVersion == INDETERMINATE_VERSION {
		// TODO do we defer to a version for "*"?  I prefer not.  I like it as is.
		// But we do use the adjusting rules (not this block) and we use the prefix length rules (this block).
		if preferIPv6 { // TODO this amounts to checkign the prefix length
			adjustedVersion = IPv6
		} else {
			adjustedVersion = IPv4
		}
	}
	var addrCreator func() (address, hostAddress *IPAddress)
	if *qualifier == *NO_QUALIFIER {
		addrCreator = func() (*IPAddress, *IPAddress) { //(address, hostAddress *IPAddress)
			addr := createAllAddress(adjustedVersion, NO_QUALIFIER, originator, options)
			return addr, addr
		}
	} else {
		addrCreator = func() (address, hostAddress *IPAddress) {
			addr := createAllAddress(adjustedVersion, qualifier, originator, options)
			if qualifier.zone == "" {
				hostAddr := createAllAddress(adjustedVersion, NO_QUALIFIER, originator, options)
				return addr, hostAddr
			}
			qualifier2 := ParsedHostIdentifierStringQualifier{zone: qualifier.zone}
			hostAddr := createAllAddress(adjustedVersion, &qualifier2, originator, options)
			return addr, hostAddr
		}
	}
	cached := cachedAddressProvider{addressCreator: addrCreator}
	return &AllCreator{
		AdjustedAddressCreator: AdjustedAddressCreator{
			networkPrefixLength: qualifier.getEquivalentPrefixLength(),
			VersionedAddressCreator: VersionedAddressCreator{
				adjustedVersion:       adjustedVersion,
				parameters:            options,
				cachedAddressProvider: cached,
				versionedAddressCreator: func(version IPVersion) *IPAddress {
					return createAllAddress(version, qualifier, originator, options)
				},
			},
		},
		originator: originator,
		qualifier:  *qualifier,
	}
}

type AllCreator struct {
	AdjustedAddressCreator

	originator HostIdentifierString
	qualifier  ParsedHostIdentifierStringQualifier //TODO copy the original to here
}

func (all *AllCreator) getType() IPType {
	if !all.adjustedVersion.isIndeterminate() {
		return fromVersion(all.adjustedVersion)
	}
	return ALL
}

func (all *AllCreator) isProvidingAllAddresses() bool {
	return all.adjustedVersion == INDETERMINATE_VERSION
}

func (all *AllCreator) getProviderNetworkPrefixLength() PrefixLen {
	return all.qualifier.getEquivalentPrefixLength()
}

func (all *AllCreator) getProviderMask() *IPAddress {
	return all.qualifier.getMaskLower()
}

func (all *AllCreator) getProviderSeqRange() *IPAddressSeqRange {
	if all.isProvidingAllAddresses() {
		return nil
	}
	mask := all.getProviderMask()
	if mask != nil && mask.GetBlockMaskPrefixLength(true) == nil {
		// we must apply the mask
		all := createAllAddress(all.adjustedVersion, NO_QUALIFIER, nil, all.parameters)
		upper, _ := all.GetUpper().Mask(mask)
		lower := all.GetLower() //TODO apply the mask? maybe I have this wrong in Java too
		rge, _ := lower.SpanWithRange(upper)
		return rge
	}
	return all.cachedAddressProvider.getProviderSeqRange()
}

// TODO the ones below later
//		@Override
//		public Boolean contains(IPAddressProvider otherProvider) {
//			if(otherProvider.isInvalid()) {
//				return Boolean.FALSE;
//			} else if(adjustedVersion == null) {
//				return Boolean.TRUE;
//			}
//			return adjustedVersion == otherProvider.getProviderIPVersion();
//		}
//
//
//		@Override
//		public boolean isSequential() {
//			return !isProvidingAllAddresses();
//		}
//
//		@Override
//		public IPAddressDivisionSeries getDivisionGrouping() throws IncompatibleAddressError {
//			if(isProvidingAllAddresses()) {
//				return null;
//			}
//			IPAddressNetwork<?, ?, ?, ?, ?> network = adjustedVersion.isIPv4() ?
//					options.getIPv4Parameters().getNetwork() : options.getIPv6Parameters().getNetwork();
//			IPAddress mask = getProviderMask();
//			if(mask != null && mask.getBlockMaskPrefixLength(true) == null) {
//				// there is a mask
//				Integer hostMaskPrefixLen = mask.getBlockMaskPrefixLength(false);
//				if(hostMaskPrefixLen == null) { // not a host mask
//					throw new IncompatibleAddressError(getProviderAddress(), mask, "ipaddress.error.maskMismatch");
//				}
//				IPAddress hostMask = network.getHostMask(hostMaskPrefixLen);
//				return hostMask.toPrefixBlock();
//			}
//			IPAddressDivisionSeries grouping;
//			if(adjustedVersion.isIPv4()) {
//				grouping = new IPAddressDivisionGrouping(new IPAddressBitsDivision[] {
//							new IPAddressBitsDivision(0, IPv4Address.MAX_VALUE, IPv4Address.BIT_COUNT, IPv4Address.DEFAULT_TEXTUAL_RADIX, network, qualifier.getEquivalentPrefixLength())
//						}, network);
//			} else if(adjustedVersion.isIPv6()) {
//				byte upperBytes[] = new byte[16];
//				Arrays.fill(upperBytes, (byte) 0xff);
//				grouping = new IPAddressLargeDivisionGrouping(new IPAddressLargeDivision[] {new IPAddressLargeDivision(new byte[IPv6Address.BYTE_COUNT], upperBytes, IPv6Address.BIT_COUNT, IPv6Address.DEFAULT_TEXTUAL_RADIX, network, qualifier.getEquivalentPrefixLength())}, network);
//			} else {
//				grouping = null;
//			}
//			return grouping;
//		}
//	}

// TODO NOW progress
// TODO NEXT NOW progress
//
// - also segment prefixContains and prefixEquals
// - you might take the approach of implementing the use-cases (excluding streams and tries) from the wiki to get the important stuff in, then fill in the gaps later
// - finish off the ip address creator interfaces (not sure if you need much here)
// - finish HostName (now it's mostly done, just a few methods left) <---
// - try to create the right set of constructors for sections and addresses, hopefully straightforward
// - check notes.txt in Java for functionality table
// - go over the java to-dos as some might make sense in golang too
// - did we do mac <-> ipv6?  Or ipv4 <-> ipv6?
// - finish the list of methods in ExtendedIPSegmentSeries - almost there <---
// ---> - we need to circle back to the parsing code and do all teh things we deferred, such as the locking, such as the optimized contains and equals, etc
//
// Still a lot of work, BUT, you are clearly past the bug hump, way past halfway, on the home stretch

// more: subtract, intersect, BitwiseOr

// I think you want to fix up the locking in the parsing, then move on to the prefixEquals and prefixContains stuff, but first do the reverse methods
