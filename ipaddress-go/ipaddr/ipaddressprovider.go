package ipaddr

import (
	"sync/atomic"
	"unsafe"
)

// All IP address strings corresponds to exactly one of these types.
// In cases where there is no corresponding default IPAddress value (invalidType, allType, and possibly emptyType), these types can be used for comparison.
// emptyType means a zero-length string (useful for validation, we can set validation to allow empty strings) that has no corresponding IPAddress value (validation options allow you to map empty to the loopback)
// invalidType means it is known that it is not any of the other allowed types (validation options can restrict the allowed types)
// allType means it is wildcard(s) with no separators, like "*", which represents all addresses, whether IPv4, IPv6 or other, and thus has no corresponding IPAddress value
// These constants are ordered by address space size, from smallest to largest, and the ordering affects comparisons
type ipType int

func fromVersion(version IPVersion) ipType {
	switch version {
	case IPv4:
		return ipv4AddrType
	case IPv6:
		return ipv6AddrType
	default:
	}
	return uninitializedType
}

func (t ipType) isUnknown() bool {
	return t == uninitializedType
}

const (
	uninitializedType ipType = iota
	invalidType
	emptyType
	ipv4AddrType
	ipv6AddrType
	//PREFIX_ONLY
	allType
)

type ipAddressProvider interface {
	getType() ipType

	getProviderHostAddress() (*IPAddress, IncompatibleAddressError)

	getProviderAddress() (*IPAddress, IncompatibleAddressError)

	getVersionedAddress(version IPVersion) (*IPAddress, IncompatibleAddressError)

	isSequential() bool

	getProviderSeqRange() *IPAddressSeqRange

	getProviderMask() *IPAddress

	// TODO LATER getDivisionGrouping
	//default IPAddressDivisionSeries getDivisionGrouping() throws IncompatibleAddressError {
	//	return getProviderAddress();
	//}

	providerCompare(ipAddressProvider) (int, IncompatibleAddressError)

	providerEquals(ipAddressProvider) (bool, IncompatibleAddressError)

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

	// If the address was created by parsing, this provides the parameters used when creating the address,
	// otherwise nil
	getParameters() IPAddressStringParameters

	// containsProvider is an optimized contains that does not need to create address objects to return an answer.
	// Unconventional addresses may require that the address objects are created, in such cases null is returned.
	//
	// Addresses constructed from canonical or normalized representations with no wildcards will not return null.
	containsProvider(ipAddressProvider) boolSetting

	// contains is an optimized contains that does not need to fully parse the other address to return an answer.
	//
	// Unconventional addresses may require full parsing, in such cases null is returned.
	//
	// Addresses constructed from canonical or normalized representations with no wildcards will not return null.
	contains(string) boolSetting

	// prefixEquals is an optimized prefix comparison that does not need to fully parse the other address to return an answer.
	//
	// Unconventional addresses may require full parsing, in such cases null is returned.
	//
	// Addresses constructed from canonical or normalized representations with no wildcards will not return null.
	prefixEquals(string) boolSetting

	// prefixEqualsProvider is an optimized prefix comparison that does not need to create addresses to return an answer.
	//
	// Unconventional addresses may require the address objects, in such cases null is returned.
	prefixEqualsProvider(ipAddressProvider) boolSetting

	// prefixContains is an optimized prefix comparison that does not need to create addresses to return an answer.
	//
	// Unconventional addresses may require the address objects, in such cases null is returned.
	prefixContains(string) boolSetting

	// prefixContainsProvider is an optimized prefix comparison that does not need to create addresses to return an answer.
	//
	// Unconventional addresses may require the address objects, in such cases null is returned.
	prefixContainsProvider(ipAddressProvider) boolSetting

	// parsedEquals is an optimized equality comparison that does not need to create addresses to return an answer.
	//
	// Unconventional addresses may require the address objects, in such cases null is returned.
	parsedEquals(ipAddressProvider) boolSetting
}

type ipAddrProvider struct{}

func (p *ipAddrProvider) getType() ipType {
	return uninitializedType
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

func (p *ipAddrProvider) getVersionedAddress(_ IPVersion) (*IPAddress, IncompatibleAddressError) {
	return nil, nil
}

func (p *ipAddrProvider) getProviderMask() *IPAddress {
	return nil
}

func (p *ipAddrProvider) getProviderIPVersion() IPVersion {
	return IndeterminateIPVersion
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

func (p *ipAddrProvider) containsProvider(ipAddressProvider) boolSetting {
	return boolSetting{}
}

func (p *ipAddrProvider) contains(string) boolSetting {
	return boolSetting{}
}

func (p *ipAddrProvider) prefixEquals(string) boolSetting {
	return boolSetting{}
}

func (p *ipAddrProvider) prefixEqualsProvider(ipAddressProvider) boolSetting {
	return boolSetting{}
}

func (p *ipAddrProvider) prefixContains(string) boolSetting {
	return boolSetting{}
}

func (p *ipAddrProvider) prefixContainsProvider(ipAddressProvider) boolSetting {
	return boolSetting{}
}

func (p *ipAddrProvider) parsedEquals(ipAddressProvider) boolSetting {
	return boolSetting{}
}

func providerCompare(p, other ipAddressProvider) (res int, err IncompatibleAddressError) {
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
			res = value.CompareTo(otherValue)
			return
		}
	}
	var thisType, otherType = p.getType(), other.getType()
	res = int(thisType - otherType)
	return
}

/**
* When a value provider produces no value, equality and comparison are based on the enum ipType,
* which can by null.
* @param o
* @return
 */
func providerEquals(p, other ipAddressProvider) (res bool, err IncompatibleAddressError) {
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
			res = value.Equals(otherValue)
			return
		} else {
			return // returns false
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

	ipType                ipType
	isInvalidVal, isEmpty bool
	//isInvalidVal, isUninitializedVal, isEmpty bool
}

func (p *nullProvider) isInvalid() bool {
	return p.isInvalidVal
}

func (p *nullProvider) isProvidingEmpty() bool {
	return p.isEmpty
}

func (p *nullProvider) getType() ipType {
	return p.ipType
}

func (p *nullProvider) providerCompare(other ipAddressProvider) (int, IncompatibleAddressError) {
	return providerCompare(p, other)
}

func (p *nullProvider) providerEquals(other ipAddressProvider) (bool, IncompatibleAddressError) {
	return providerEquals(p, other)
}

var (
	invalidProvider = &nullProvider{isInvalidVal: true, ipType: invalidType}
	emptyProvider   = &nullProvider{isEmpty: true, ipType: emptyType}
)

///**
//	 * Wraps an IPAddress for IPAddressString in the cases where no parsing is provided, the address exists already
//	 * @param value
//	 * @return
//	 */
func getProviderFor(address, hostAddress *IPAddress) ipAddressProvider {
	return &cachedAddressProvider{addresses: &addressResult{address: address, hostAddress: hostAddress}}
}

type addressResult struct {
	address, hostAddress *IPAddress

	// addrErr applies to address, hostErr to hostAddress
	addrErr, hostErr IncompatibleAddressError
}

type cachedAddressProvider struct {
	ipAddrProvider

	// addressCreator creates two addresses, the host address and address with prefix/mask, at the same time
	addressCreator func() (address, hostAddress *IPAddress, addrErr, hostErr IncompatibleAddressError)

	addresses *addressResult
}

func (cached *cachedAddressProvider) providerCompare(other ipAddressProvider) (int, IncompatibleAddressError) {
	return providerCompare(cached, other)
}

func (cached *cachedAddressProvider) providerEquals(other ipAddressProvider) (bool, IncompatibleAddressError) {
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

func (cached *cachedAddressProvider) getProviderHostAddress() (res *IPAddress, err IncompatibleAddressError) {
	addrs := cached.addresses
	if addrs == nil {
		_, res, _, err = cached.getCachedAddresses() // sets cached.addresses
	} else {
		res, err = addrs.hostAddress, addrs.hostErr
	}
	return
}

func (cached *cachedAddressProvider) getProviderAddress() (res *IPAddress, err IncompatibleAddressError) {
	addrs := cached.addresses
	if addrs == nil {
		res, _, err, _ = cached.getCachedAddresses() // sets cached.addresses
	} else {
		res, err = addrs.address, addrs.addrErr
	}
	return
}

func (cached *cachedAddressProvider) getCachedAddresses() (address, hostAddress *IPAddress, addrErr, hostErr IncompatibleAddressError) {
	addrs := cached.addresses
	if addrs == nil {
		if cached.addressCreator != nil {
			address, hostAddress, addrErr, hostErr = cached.addressCreator()
			addresses := &addressResult{
				address:     address,
				hostAddress: hostAddress,
				addrErr:     addrErr,
				hostErr:     hostErr,
			}
			dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&cached.addresses))
			atomic.StorePointer(dataLoc, unsafe.Pointer(addresses))
		}
	} else {
		address, hostAddress, addrErr, hostErr = addrs.address, addrs.hostAddress, addrs.addrErr, addrs.hostErr
	}
	return
	//xxx
	/*
		networkMaskLen, hostMaskLen := section.checkForPrefixMask()
			res := &maskLenSetting{networkMaskLen, hostMaskLen}
			dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&cacheBitCountx.cachedMaskLens))
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
	if addr, _ := cached.getProviderAddress(); addr != nil {
		p = addr.GetNetworkPrefixLength()
	}
	return
}

func (cached *cachedAddressProvider) getProviderIPVersion() IPVersion {
	if addr, _ := cached.getProviderAddress(); addr != nil {
		return addr.getIPVersion()
	}
	return IndeterminateIPVersion
}

func (cached *cachedAddressProvider) getType() ipType {
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

type versionedAddressCreator struct {
	cachedAddressProvider

	adjustedVersion IPVersion

	versionedAddressCreatorFunc func(IPVersion) (*IPAddress, IncompatibleAddressError)

	versionedValues [2]*IPAddress

	parameters IPAddressStringParameters
}

func (versioned *versionedAddressCreator) getParameters() IPAddressStringParameters {
	return versioned.parameters
}

func (versioned *versionedAddressCreator) isProvidingIPAddress() bool {
	return versioned.adjustedVersion != IndeterminateIPVersion
}

func (versioned *versionedAddressCreator) isProvidingIPv4() bool {
	return versioned.adjustedVersion == IPv4
}

func (versioned *versionedAddressCreator) isProvidingIPv6() bool {
	return versioned.adjustedVersion == IPv6
}

func (versioned *versionedAddressCreator) getProviderIPVersion() IPVersion {
	return versioned.adjustedVersion
}

func (versioned *versionedAddressCreator) getType() ipType {
	return fromVersion(versioned.adjustedVersion)
}

func (versioned *versionedAddressCreator) getVersionedAddress(version IPVersion) (addr *IPAddress, err IncompatibleAddressError) {
	index := version.index()
	if index >= IndeterminateIPVersion.index() {
		return
	}
	if versioned.versionedAddressCreatorFunc != nil {
		addr = versioned.versionedValues[index]
		if addr == nil {
			addr, err = versioned.versionedAddressCreatorFunc(version)
			if err == nil {
				dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&versioned.versionedValues[index]))
				atomic.StorePointer(dataLoc, unsafe.Pointer(addr))
			}
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

func newLoopbackCreator(options IPAddressStringParameters, zone Zone) *loopbackCreator {
	// TODO an option to set preferred loopback here in IPAddressStringParameters, do the same in Java
	// the option will set one of three options, IPv4, IPv6, or IndeterminateIPVersion which is the default
	// In Go the default will be IPv4
	// There is another option I wanted to add, was in the validator code, I think allow empty zone with prefix like %/
	// ALSO, consider using zero value instead of loopback - zero string becomes zero value
	var preferIPv6 bool
	ipv6WithZoneLoop := func() *IPAddress {
		network := DefaultIPv6Network
		creator := network.getIPAddressCreator()
		return creator.createAddressInternalFromBytes(network.GetLoopback().GetBytes(), zone)
	}
	ipv6Loop := func() *IPAddress {
		return DefaultIPv6Network.GetLoopback()
	}
	ipv4Loop := func() *IPAddress {
		return DefaultIPv4Network.GetLoopback()
	}
	double := func(one *IPAddress) (address, hostAddress *IPAddress) {
		return one, one
	}
	var lbackCreator func() (address, hostAddress *IPAddress)
	var version IPVersion
	if preferIPv6 {
		if len(zone) > 0 {
			lbackCreator = func() (*IPAddress, *IPAddress) { return double(ipv6WithZoneLoop()) }
		} else {
			lbackCreator = func() (*IPAddress, *IPAddress) { return double(ipv6Loop()) }
		}
		version = IPv6
	} else {
		lbackCreator = func() (*IPAddress, *IPAddress) { return double(ipv4Loop()) }
		version = IPv4
	}
	cached := cachedAddressProvider{
		addressCreator: func() (address, hostAddress *IPAddress, addrErr, hostErr IncompatibleAddressError) {
			address, hostAddress = lbackCreator()
			return
		},
	}
	loopbackCreatorFunc := func(version IPVersion) *IPAddress {
		addresses := cached.addresses
		if addresses != nil {
			addr := addresses.address
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
	versionedAddressCreatorFunc := func(version IPVersion) (*IPAddress, IncompatibleAddressError) {
		return loopbackCreatorFunc(version), nil
	}
	return &loopbackCreator{
		versionedAddressCreator: versionedAddressCreator{
			adjustedVersion:             version,
			parameters:                  options,
			cachedAddressProvider:       cached,
			versionedAddressCreatorFunc: versionedAddressCreatorFunc,
		},
		zone: zone,
	}
}

type loopbackCreator struct {
	versionedAddressCreator

	zone Zone
}

func (loop *loopbackCreator) providerCompare(other ipAddressProvider) (int, IncompatibleAddressError) {
	return providerCompare(loop, other)
}

func (loop *loopbackCreator) providerEquals(other ipAddressProvider) (bool, IncompatibleAddressError) {
	return providerEquals(loop, other)
}

func (loop *loopbackCreator) getProviderNetworkPrefixLength() PrefixLen {
	return nil
}

type adjustedAddressCreator struct {
	versionedAddressCreator

	networkPrefixLength PrefixLen
}

func (adjusted *adjustedAddressCreator) getProviderNetworkPrefixLength() PrefixLen {
	return adjusted.networkPrefixLength
}

func (adjusted *adjustedAddressCreator) getProviderAddress() (*IPAddress, IncompatibleAddressError) {
	if !adjusted.isProvidingIPAddress() {
		return nil, nil
	}
	return adjusted.versionedAddressCreator.getProviderAddress()
}

func (adjusted *adjustedAddressCreator) getProviderHostAddress() (*IPAddress, IncompatibleAddressError) {
	if !adjusted.isProvidingIPAddress() {
		return nil, nil
	}
	return adjusted.versionedAddressCreator.getProviderHostAddress()
}

// TODO the adjusted version passed in is the one adjusted due to zone %, or mask version, or prefix len >= 32
// INside this function we will handle the cases where it is still not determined, and that will be based on our new rules
// involving (a) maybe when < 32 we default to IPv4, otherwise IPv6
//			(b) this behaviour can be overridden by a string parameters option

func newMaskCreator(options IPAddressStringParameters, adjustedVersion IPVersion, networkPrefixLength PrefixLen) *maskCreator {
	// TODO use the option for  preferred loopback also for preferred mask, do the same in Java
	// Drop "prefix only" type - it was never a good idea anyway!  Better to prefer one over the other.

	var preferIPv6 bool

	if adjustedVersion == IndeterminateIPVersion {
		if preferIPv6 {
			adjustedVersion = IPv6
		} else {
			adjustedVersion = IPv4
		}
	}
	createVersionedMask := func(version IPVersion, prefLen PrefixLen, withPrefixLength bool) *IPAddress {
		if version == IPv4 {
			network := DefaultIPv4Network
			return network.GetNetworkMask(*prefLen)
		} else if version == IPv6 {
			network := DefaultIPv6Network
			return network.GetNetworkMask(*prefLen)
		}
		return nil
	}
	versionedAddressCreatorFunc := func(version IPVersion) (*IPAddress, IncompatibleAddressError) {
		return createVersionedMask(version, networkPrefixLength, true), nil
	}
	maskCreatorFunc := func() (address, hostAddress *IPAddress) {
		prefLen := networkPrefixLength
		return createVersionedMask(adjustedVersion, prefLen, true),
			createVersionedMask(adjustedVersion, prefLen, false)
	}
	addrCreator := func() (address, hostAddress *IPAddress, addrErr, hostErr IncompatibleAddressError) {
		address, hostAddress = maskCreatorFunc()
		return
	}
	cached := cachedAddressProvider{addressCreator: addrCreator}
	return &maskCreator{
		adjustedAddressCreator{
			networkPrefixLength: networkPrefixLength,
			versionedAddressCreator: versionedAddressCreator{
				adjustedVersion:             adjustedVersion,
				parameters:                  options,
				cachedAddressProvider:       cached,
				versionedAddressCreatorFunc: versionedAddressCreatorFunc,
			},
		},
	}
}

type maskCreator struct {
	adjustedAddressCreator
}

// TODO the adjusted version passed in is the one adjusted due to zone %, or mask version, or prefix len >= 32
// INside this function we will handle the cases where it is still not determined, and that will be based on our new rules
// involving (a) maybe when < 32 we default to IPv4, otherwise IPv6
//			(b) this behaviour can be overridden by a string parameters option

func newAllCreator(qualifier *parsedHostIdentifierStringQualifier, adjustedVersion IPVersion, originator HostIdentifierString, options IPAddressStringParameters) ipAddressProvider {
	result := &allCreator{
		adjustedAddressCreator: adjustedAddressCreator{
			networkPrefixLength: qualifier.getEquivalentPrefixLength(),
			versionedAddressCreator: versionedAddressCreator{
				adjustedVersion: adjustedVersion,
				parameters:      options,
			},
		},
		originator: originator,
		qualifier:  *qualifier,
	}
	result.addressCreator = result.createAddrs
	result.versionedAddressCreatorFunc = result.versionedCreate
	return result
}

type allCreator struct {
	adjustedAddressCreator

	originator HostIdentifierString
	qualifier  parsedHostIdentifierStringQualifier

	rng *IPAddressSeqRange
}

func (all *allCreator) getType() ipType {
	if !all.adjustedVersion.isIndeterminate() {
		return fromVersion(all.adjustedVersion)
	}
	return allType
}

func (all *allCreator) providerCompare(other ipAddressProvider) (int, IncompatibleAddressError) {
	return providerCompare(all, other)
}

func (all *allCreator) providerEquals(other ipAddressProvider) (bool, IncompatibleAddressError) {
	return providerEquals(all, other)
}

func (all *allCreator) isProvidingAllAddresses() bool {
	return all.adjustedVersion == IndeterminateIPVersion
}

func (all *allCreator) getProviderNetworkPrefixLength() PrefixLen {
	return all.qualifier.getEquivalentPrefixLength()
}

func (all *allCreator) getProviderMask() *IPAddress {
	return all.qualifier.getMaskLower()
}

func (all *allCreator) createAll() (rng *IPAddressSeqRange, addr *IPAddress, hostAddr *IPAddress, addrErr IncompatibleAddressError, hostErr IncompatibleAddressError) {
	rng = all.rng
	addrs := all.addresses
	if rng == nil || addrs == nil {
		var lower, upper *IPAddress
		addr, hostAddr, lower, upper, addrErr = createAllAddress(
			all.adjustedVersion,
			&all.qualifier,
			all.originator)
		rng, _ = lower.SpanWithRange(upper)
		dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&all.rng))
		atomic.StorePointer(dataLoc, unsafe.Pointer(rng))
		addresses := &addressResult{
			address:     addr,
			hostAddress: hostAddr,
			addrErr:     addrErr,
			hostErr:     hostErr,
		}
		dataLoc = (*unsafe.Pointer)(unsafe.Pointer(&all.addresses))
		atomic.StorePointer(dataLoc, unsafe.Pointer(addresses))
	} else {
		addr, hostAddr, addrErr, hostErr = addrs.address, addrs.hostAddress, addrs.addrErr, addrs.hostErr
	}
	return
}

func (all *allCreator) createRange() (rng *IPAddressSeqRange) {
	rng, _, _, _, _ = all.createAll()
	return
}

func (all *allCreator) createAddrs() (addr *IPAddress, hostAddr *IPAddress, addrErr IncompatibleAddressError, hostErr IncompatibleAddressError) {
	_, addr, hostAddr, addrErr, hostErr = all.createAll()
	return
}

func (all *allCreator) versionedCreate(version IPVersion) (addr *IPAddress, addrErr IncompatibleAddressError) {
	if version == all.adjustedVersion {
		return all.getProviderAddress()
	} else if all.adjustedVersion != IndeterminateIPVersion {
		return nil, nil
	}
	addr, _, _, _, addrErr = createAllAddress(
		version,
		&all.qualifier,
		all.originator)
	return
}

func (all *allCreator) getProviderSeqRange() *IPAddressSeqRange {
	if all.isProvidingAllAddresses() {
		return nil
	}
	rng := all.rng
	if rng == nil {
		rng = all.createRange()
	}
	return rng
}

func (all *allCreator) prefixContainsProvider(otherProvider ipAddressProvider) boolSetting {
	return all.containsProviderFunc(otherProvider, (*IPAddress).prefixContains)
}

func (all *allCreator) containsProvider(otherProvider ipAddressProvider) (res boolSetting) {
	return all.containsProviderFunc(otherProvider, (*IPAddress).contains)
}

func (all *allCreator) containsProviderFunc(otherProvider ipAddressProvider, functor func(*IPAddress, AddressType) bool) (res boolSetting) {
	if otherProvider.isInvalid() {
		return boolSetting{true, false}
	} else if all.adjustedVersion == IndeterminateIPVersion {
		return boolSetting{true, true}
	} else if all.adjustedVersion != otherProvider.getProviderIPVersion() {
		return boolSetting{true, false}
	} else if all.qualifier.getMaskLower() == nil && all.qualifier.getZone() == NoZone {
		return boolSetting{true, true}
	} else if addr, err := all.getProviderAddress(); err != nil {
		return boolSetting{true, false}
	} else if otherAddr, err := all.getProviderAddress(); err != nil {
		return boolSetting{true, false}
	} else {
		return boolSetting{true, functor(addr, otherAddr)}
		//return boolSetting{true, addr.Contains(otherAddr)}
	}
}

// TODO later getDivisionGrouping()
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
// - you might take the approach of implementing the use-cases (excluding streams and tries) from the wiki to get the important stuff in, then fill in the gaps later
// - finish HostName (now it's mostly done, just a few methods left) <---
// - try to create the right set of constructors for sections and addresses, hopefully straightforward
// - check notes.txt in Java for functionality table
// - go over the java to-dos as some might make sense in golang too
// - did we do mac <-> ipv6?  Or ipv4 <-> ipv6? not yet
// ---> - we need to circle back to the parsing code and do all teh things we deferred, such as the locking, such as the optimized contains and equals, etc
//
// Still a lot of work, BUT, you are clearly past the bug hump, way past halfway, on the home stretch

// TODO append and replace in sections: we only allow at top-level.
// This ensures we do not have weirdness with IPv6v4MixedSection or whatnot.  Keeps ipv4 sections as ipv4.  Etc.
// Appending to IPv6v4MixedSection, what should happen?
// avoiding it at lower level prevents weirdness like ipv4 becoming not ipv4 unpredictably.
// Or appending to IPv4, we must ensure the division groupings are also ipv4.  I am inclined to (a) only maintain addrType when appending at highest level,
// (b) drop the addrType at lower levels.  It is possible you could check addrType and upscale, but this does not help with IPv6v4MixedSection.
// So, you could just upscale selectively.  I like that.
// But in Java, you do not allow append or replace at lower levels at all.  So, maybe you do that.  In fact, that alleviates confusion.
// And any grouping can simply be reconstitued from the divisions as desired, you don't need it at lower level.
// In Java, it is really problematic for (a) the type of the append or replace arg and (b) what to do when there is no match
