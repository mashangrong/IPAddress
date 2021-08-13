package ipaddr

import (
	"math/big"
	"net"
	"unsafe"
)

const (
	IPv6SegmentSeparator                    = ':'
	IPv6ZoneSeparator                       = '%'
	IPv6AlternativeZoneSeparator            = '\u00a7'
	IPv6BitsPerSegment             BitCount = 16
	IPv6BytesPerSegment                     = 2
	IPv6SegmentCount                        = 8
	IPv6MixedReplacedSegmentCount           = 2
	IPv6MixedOriginalSegmentCount           = 6
	IPv6ByteCount                           = 16
	IPv6BitCount                   BitCount = 128
	IPv6DefaultTextualRadix                 = 16
	IPv6MaxValuePerSegment                  = 0xffff
	IPv6ReverseDnsSuffix                    = ".ip6.arpa"
	IPv6ReverseDnsSuffixDeprecated          = ".ip6.int"

	IPv6UncSegmentSeparator  = '-'
	IPv6UncZoneSeparator     = 's'
	IPv6UncRangeSeparator    = AlternativeRangeSeparator
	IPv6UncRangeSeparatorStr = string(AlternativeRangeSeparator)
	IPv6UncSuffix            = ".ipv6-literal.net"

	IPv6SegmentMaxChars             = 4
	IPv6SegmentBitsPerChar BitCount = 4
)

type Zone string

func (zone Zone) IsEmpty() bool {
	return zone == ""
}

const NoZone = ""

func NewIPv6Address(section *IPv6AddressSection) *IPv6Address {
	return createAddress(section.ToAddressSection(), NoZone).ToIPv6Address()
}

func NewIPv6AddressZoned(section *IPv6AddressSection, zone string) *IPv6Address {
	zoneVal := Zone(zone)
	result := createAddress(section.ToAddressSection(), zoneVal).ToIPv6Address()
	if zoneVal != NoZone {
		result.cache.stringCache = &stringCache{}
	}
	return result
}

func NewIPv6AddressFromIP(bytes net.IP) (addr *IPv6Address, err AddressValueError) {
	section, err := NewIPv6SectionFromSegmentedBytes(bytes, IPv6SegmentCount)
	if err == nil {
		addr = NewIPv6Address(section)
	}
	return
}

func NewIPv6AddressFromPrefixedIP(bytes net.IP, prefixLength PrefixLen) (addr *IPv6Address, err AddressValueError) {
	section, err := NewIPv6SectionFromPrefixedBytes(bytes, IPv6SegmentCount, prefixLength)
	if err == nil {
		addr = NewIPv6Address(section)
	}
	return
}

func newIPv6AddressFromZonedIP(bytes net.IP, zone string) (addr *IPv6Address, err AddressValueError) {
	addr, err = NewIPv6AddressFromIP(bytes)
	if err == nil {
		addr.zone = Zone(zone)
	}
	return
}

func NewIPv6AddressFromIPAddr(ipAddr *net.IPAddr) (addr *IPv6Address, err AddressValueError) {
	return newIPv6AddressFromZonedIP(ipAddr.IP, ipAddr.Zone)
}

func NewIPv6AddressFromPrefixedIPAddr(ipAddr *net.IPAddr, prefixLen PrefixLen) (addr *IPv6Address, err AddressValueError) {
	addr, err = NewIPv6AddressFromPrefixedIP(ipAddr.IP, prefixLen)
	if err == nil {
		addr.zone = Zone(ipAddr.Zone)
	}
	return
}

func NewIPv6AddressFromBigInt(bytes *big.Int) (addr *IPv6Address, err AddressValueError) {
	return NewIPv6AddressFromIP(bytes.Bytes())
}

func NewIPv6AddressFromPrefixedBigInt(bytes *big.Int, prefixLength PrefixLen) (addr *IPv6Address, err AddressValueError) {
	return NewIPv6AddressFromPrefixedIP(bytes.Bytes(), prefixLength)
}

func NewIPv6AddressFromZonedBigInt(bytes *big.Int, zone string) (addr *IPv6Address, err AddressValueError) {
	return newIPv6AddressFromZonedIP(bytes.Bytes(), zone)
}

func NewIPv6AddressFromPrefixedZonedBigInt(bytes *big.Int, prefixLength PrefixLen, zone string) (addr *IPv6Address, err AddressValueError) {
	return NewIPv6AddressFromPrefixedIPAddr(&net.IPAddr{IP: bytes.Bytes(), Zone: zone}, prefixLength)
}

func NewIPv6AddressFromSegments(segments []*IPv6AddressSegment) (addr *IPv6Address, err AddressValueError) {
	section, err := NewIPv6AddressFromSection(segments)
	if err == nil {
		addr = NewIPv6Address(section)
	}
	return
}

func NewIPv6AddressFromPrefixedSegments(segments []*IPv6AddressSegment, prefixLength PrefixLen) (addr *IPv6Address, err AddressValueError) {
	section, err := NewIPv6AddressFromPrefixedSection(segments, prefixLength)
	if err == nil {
		addr = NewIPv6Address(section)
	}
	return
}

func NewIPv6AddressFromZonedSegments(segments []*IPv6AddressSegment, zone string) (addr *IPv6Address, err AddressValueError) {
	section, err := NewIPv6AddressFromSection(segments)
	if err == nil {
		addr = NewIPv6AddressZoned(section, zone)
	}
	return
}

func NewIPv6AddressFromPrefixedZonedSegments(segments []*IPv6AddressSegment, prefixLength PrefixLen, zone string) (addr *IPv6Address, err AddressValueError) {
	section, err := NewIPv6AddressFromPrefixedSection(segments, prefixLength)
	if err == nil {
		addr = NewIPv6AddressZoned(section, zone)
	}
	return
}

func NewIPv6AddressFromUint64(highBytes, lowBytes uint64) *IPv6Address {
	section := NewIPv6SectionFromUint64(highBytes, lowBytes, IPv6SegmentCount)
	return NewIPv6Address(section)
}

func NewIPv6AddressFromPrefixedUint64(highBytes, lowBytes uint64, prefixLength PrefixLen) *IPv6Address {
	section := NewIPv6SectionFromPrefixedUint64(highBytes, lowBytes, IPv6SegmentCount, prefixLength)
	return NewIPv6Address(section)
}

func NewIPv6AddressFromZonedUint64(highBytes, lowBytes uint64, zone string) *IPv6Address {
	section := NewIPv6SectionFromUint64(highBytes, lowBytes, IPv6SegmentCount)
	return NewIPv6AddressZoned(section, zone)
}

func NewIPv6AddressFromPrefixedZonedUint64(highBytes, lowBytes uint64, prefixLength PrefixLen, zone string) *IPv6Address {
	section := NewIPv6SectionFromPrefixedUint64(highBytes, lowBytes, IPv6SegmentCount, prefixLength)
	return NewIPv6AddressZoned(section, zone)
}

func NewIPv6AddressFromVals(vals SegmentValueProvider) *IPv6Address {
	section := NewIPv6SectionFromValues(vals, IPv6SegmentCount)
	return NewIPv6Address(section)
}

func NewIPv6AddressFromPrefixedVals(vals SegmentValueProvider, prefixLength PrefixLen) *IPv6Address {
	section := NewIPv6SectionFromPrefixedValues(vals, IPv6SegmentCount, prefixLength)
	return NewIPv6Address(section)
}

func NewIPv6AddressFromRange(vals, upperVals SegmentValueProvider) *IPv6Address {
	section := NewIPv6SectionFromRangeValues(vals, upperVals, IPv6SegmentCount)
	return NewIPv6Address(section)
}

func NewIPv6AddressFromPrefixedRange(vals, upperVals SegmentValueProvider, prefixLength PrefixLen) *IPv6Address {
	section := NewIPv6SectionFromPrefixedRangeValues(vals, upperVals, IPv6SegmentCount, prefixLength)
	return NewIPv6Address(section)
}

func NewIPv6AddressFromZonedRange(vals, upperVals SegmentValueProvider, zone string) *IPv6Address {
	section := NewIPv6SectionFromRangeValues(vals, upperVals, IPv6SegmentCount)
	return NewIPv6AddressZoned(section, zone)
}

func NewIPv6AddressFromPrefixedZonedRange(vals, upperVals SegmentValueProvider, prefixLength PrefixLen, zone string) *IPv6Address {
	section := NewIPv6SectionFromPrefixedRangeValues(vals, upperVals, IPv6SegmentCount, prefixLength)
	return NewIPv6AddressZoned(section, zone)
}

var zeroIPv6 = initZeroIPv6()

func initZeroIPv6() *IPv6Address {
	div := NewIPv6Segment(0).ToAddressDivision()
	segs := []*AddressDivision{div, div, div, div, div, div, div, div}
	section, _ := newIPv6Section(segs, 0, false)
	return NewIPv6Address(section)
}

// TODO survey the IPv6 API (ie take a look at java to see what we are missing), I've already surveyed IPAddress API

//
//
// IPv6Address is an IPv6 address, or a subnet of multiple IPv6 addresses.  Each segment can represent a single value or a range of values.
// The zero value is ::
type IPv6Address struct {
	ipAddressInternal
}

func (addr *IPv6Address) GetBitCount() BitCount {
	return IPv6BitCount
}

func (addr *IPv6Address) GetByteCount() int {
	return IPv6ByteCount
}

func (addr *IPv6Address) GetBitsPerSegment() BitCount {
	return IPv6BitsPerSegment
}

func (addr *IPv6Address) GetBytesPerSegment() int {
	return IPv6BytesPerSegment
}

func (addr *IPv6Address) init() *IPv6Address {
	if addr.section == nil {
		return zeroIPv6
	}
	return addr
}

func (addr *IPv6Address) HasZone() bool {
	return addr.zone != NoZone
}

func (addr *IPv6Address) GetZone() Zone {
	return addr.zone
}

func (addr *IPv6Address) GetSection() *IPv6AddressSection {
	return addr.init().section.ToIPv6AddressSection()
}

// Gets the subsection from the series starting from the given index
// The first segment is at index 0.
func (addr *IPv6Address) GetTrailingSection(index int) *IPv6AddressSection {
	return addr.GetSection().GetTrailingSection(index)
}

//// Gets the subsection from the series starting from the given index and ending just before the give endIndex
//// The first segment is at index 0.
func (addr *IPv6Address) GetSubSection(index, endIndex int) *IPv6AddressSection {
	return addr.GetSection().GetSubSection(index, endIndex)
}

func (addr *IPv6Address) GetNetworkSection() *IPv6AddressSection {
	return addr.GetSection().GetNetworkSection()
}

func (addr *IPv6Address) GetNetworkSectionLen(prefLen BitCount) *IPv6AddressSection {
	return addr.GetSection().GetNetworkSectionLen(prefLen)
}

func (addr *IPv6Address) GetHostSection() *IPv6AddressSection {
	return addr.GetSection().GetHostSection()
}

func (addr *IPv6Address) GetHostSectionLen(prefLen BitCount) *IPv6AddressSection {
	return addr.GetSection().GetHostSectionLen(prefLen)
}

func (addr *IPv6Address) GetNetworkMask() *IPv6Address {
	return addr.getNetworkMask(DefaultIPv6Network).ToIPv6Address()
}

func (addr *IPv6Address) GetHostMask() *IPv6Address {
	return addr.getHostMask(DefaultIPv6Network).ToIPv6Address()
}

func (addr *IPv6Address) GetMixedAddressGrouping() (*IPv6v4MixedAddressGrouping, IncompatibleAddressError) {
	return addr.init().GetSection().getMixedAddressGrouping()
}

// GetEmbeddedIPv4AddressSection gets the IPv4 section corresponding to the lowest (least-significant) 4 bytes in the original address,
// which will correspond to between 0 and 4 bytes in this address.  Many IPv4 to IPv6 mapping schemes (but not all) use these 4 bytes for a mapped IPv4 address.
func (addr *IPv6Address) GetEmbeddedIPv4AddressSection() (*IPv4AddressSection, IncompatibleAddressError) {
	return addr.init().GetSection().getEmbeddedIPv4AddressSection()
}

// GetIPv4AddressSection produces an IPv4 address section from any sequence of bytes in this IPv6 address section
func (addr *IPv6Address) GetIPv4AddressSection(startIndex, endIndex int) (*IPv4AddressSection, IncompatibleAddressError) {
	return addr.init().GetSection().getIPv4AddressSection(startIndex, endIndex)
}

// CopySubSegments copies the existing segments from the given start index until but not including the segment at the given end index,
// into the given slice, as much as can be fit into the slice, returning the number of segments copied
func (addr *IPv6Address) CopySubSegments(start, end int, segs []*IPv6AddressSegment) (count int) {
	return addr.GetSection().CopySubSegments(start, end, segs)
}

// CopySubSegments copies the existing segments from the given start index until but not including the segment at the given end index,
// into the given slice, as much as can be fit into the slice, returning the number of segments copied
func (addr *IPv6Address) CopySegments(segs []*IPv6AddressSegment) (count int) {
	return addr.GetSection().CopySegments(segs)
}

// GetSegments returns a slice with the address segments.  The returned slice is not backed by the same array as this address.
func (addr *IPv6Address) GetSegments() []*IPv6AddressSegment {
	return addr.GetSection().GetSegments()
}

// GetSegment returns the segment at the given index
func (addr *IPv6Address) GetSegment(index int) *IPv6AddressSegment {
	return addr.init().getSegment(index).ToIPv6AddressSegment()
}

// GetSegmentCount returns the segment count
func (addr *IPv6Address) GetSegmentCount() int {
	return addr.GetDivisionCount()
}

// GetGenericDivision returns the segment at the given index as an DivisionType
func (addr *IPv6Address) GetGenericDivision(index int) DivisionType {
	return addr.init().getDivision(index)
}

// GetGenericSegment returns the segment at the given index as an AddressSegmentType
func (addr *IPv6Address) GetGenericSegment(index int) AddressSegmentType {
	return addr.init().getSegment(index)
}

// GetDivisionCount returns the segment count
func (addr *IPv6Address) GetDivisionCount() int {
	return addr.init().GetDivisionCount()
}

func (addr *IPv6Address) GetIPVersion() IPVersion {
	return IPv6
}

func (addr *IPv6Address) checkIdentity(section *IPv6AddressSection) *IPv6Address {
	sec := section.ToAddressSection()
	if sec == addr.section {
		return addr
	}
	return NewIPv6AddressZoned(section, string(addr.zone))
	//return &IPv6Address{ipAddressInternal{addressInternal{section: sec, zone: addr.zone, cache: &addressCache{}}}}
}

func (addr *IPv6Address) Mask(other *IPv6Address) (masked *IPv6Address, err IncompatibleAddressError) {
	return addr.maskPrefixed(other, false)
}

func (addr *IPv6Address) MaskPrefixed(other *IPv6Address) (masked *IPv6Address, err IncompatibleAddressError) {
	return addr.maskPrefixed(other, true)
}

func (addr *IPv6Address) maskPrefixed(other *IPv6Address, retainPrefix bool) (masked *IPv6Address, err IncompatibleAddressError) {
	addr = addr.init()
	sect, err := addr.GetSection().maskPrefixed(other.GetSection(), retainPrefix)
	if err == nil {
		masked = addr.checkIdentity(sect)
	}
	return
}

func (addr *IPv6Address) BitwiseOr(other *IPv6Address) (masked *IPv6Address, err IncompatibleAddressError) {
	return addr.bitwiseOrPrefixed(other, false)
}

func (addr *IPv6Address) BitwiseOrPrefixed(other *IPv6Address) (masked *IPv6Address, err IncompatibleAddressError) {
	return addr.bitwiseOrPrefixed(other, true)
}

func (addr *IPv6Address) bitwiseOrPrefixed(other *IPv6Address, retainPrefix bool) (masked *IPv6Address, err IncompatibleAddressError) {
	addr = addr.init()
	sect, err := addr.GetSection().bitwiseOrPrefixed(other.GetSection(), retainPrefix)
	if err == nil {
		masked = addr.checkIdentity(sect)
	}
	return
}

func (addr *IPv6Address) Subtract(other *IPv6Address) []*IPv6Address {
	addr = addr.init()
	sects, _ := addr.GetSection().Subtract(other.GetSection())
	sectLen := len(sects)
	if sectLen == 1 {
		sec := sects[0]
		if sec.ToAddressSection() == addr.section {
			return []*IPv6Address{addr}
		}
	}
	res := make([]*IPv6Address, sectLen)
	for i, sect := range sects {
		res[i] = NewIPv6AddressZoned(sect, string(addr.zone))
	}
	return res
}

func (addr *IPv6Address) Intersect(other *IPv6Address) *IPv6Address {
	addr = addr.init()
	section, _ := addr.GetSection().Intersect(other.GetSection())
	return addr.checkIdentity(section)
}

func (addr *IPv6Address) SpanWithRange(other *IPv6Address) *IPv6AddressSeqRange {
	return NewIPv6SeqRange(addr.init(), other.init())
}

func (addr *IPv6Address) GetLower() *IPv6Address {
	return addr.init().getLower().ToIPv6Address()
}

func (addr *IPv6Address) GetUpper() *IPv6Address {
	return addr.init().getUpper().ToIPv6Address()
}

// GetLowerIPAddress implements the IPAddressRange interface
func (addr *IPv6Address) GetLowerIPAddress() *IPAddress {
	return addr.GetLower().ToIPAddress()
}

// GetUpperIPAddress implements the IPAddressRange interface
func (addr *IPv6Address) GetUpperIPAddress() *IPAddress {
	return addr.GetUpper().ToIPAddress()
}

func (addr *IPv6Address) ToZeroHost() (*IPv6Address, IncompatibleAddressError) {
	res, err := addr.init().toZeroHost(false)
	return res.ToIPv6Address(), err
}

func (addr *IPv6Address) ToZeroHostLen(prefixLength BitCount) (*IPv6Address, IncompatibleAddressError) {
	res, err := addr.init().toZeroHostLen(prefixLength)
	return res.ToIPv6Address(), err
}

func (addr *IPv6Address) ToZeroNetwork() *IPv6Address {
	return addr.init().toZeroNetwork().ToIPv6Address()
}

func (addr *IPv6Address) ToMaxHost() (*IPv6Address, IncompatibleAddressError) {
	res, err := addr.init().toMaxHost()
	return res.ToIPv6Address(), err
}

func (addr *IPv6Address) ToMaxHostLen(prefixLength BitCount) (*IPv6Address, IncompatibleAddressError) {
	res, err := addr.init().toMaxHostLen(prefixLength)
	return res.ToIPv6Address(), err
}

func (addr *IPv6Address) ToPrefixBlock() *IPv6Address {
	return addr.init().toPrefixBlock().ToIPv6Address()
}

func (addr *IPv6Address) ToPrefixBlockLen(prefLen BitCount) *IPv6Address {
	return addr.init().toPrefixBlockLen(prefLen).ToIPv6Address()
}

func (addr *IPv6Address) ToBlock(segmentIndex int, lower, upper SegInt) *IPv6Address {
	return addr.init().toBlock(segmentIndex, lower, upper).ToIPv6Address()
}

func (addr *IPv6Address) WithoutPrefixLen() *IPv6Address {
	return addr.init().withoutPrefixLen().ToIPv6Address()
}

func (addr *IPv6Address) SetPrefixLen(prefixLen BitCount) *IPv6Address {
	return addr.init().setPrefixLen(prefixLen).ToIPv6Address()
}

func (addr *IPv6Address) SetPrefixLenZeroed(prefixLen BitCount) (*IPv6Address, IncompatibleAddressError) {
	res, err := addr.init().setPrefixLenZeroed(prefixLen)
	return res.ToIPv6Address(), err
}

func (addr *IPv6Address) AdjustPrefixLen(prefixLen BitCount) *IPv6Address {
	return addr.init().adjustPrefixLen(prefixLen).ToIPv6Address()
}

func (addr *IPv6Address) AdjustPrefixLenZeroed(prefixLen BitCount) (*IPv6Address, IncompatibleAddressError) {
	res, err := addr.init().adjustPrefixLenZeroed(prefixLen)
	return res.ToIPv6Address(), err
}

func (addr *IPv6Address) AssignPrefixForSingleBlock() *IPv6Address {
	return addr.init().assignPrefixForSingleBlock().ToIPv6Address()
}

func (addr *IPv6Address) AssignMinPrefixForBlock() *IPv6Address {
	return addr.init().assignMinPrefixForBlock().ToIPv6Address()
}

func (addr *IPv6Address) ContainsPrefixBlock(prefixLen BitCount) bool {
	return addr.init().ipAddressInternal.ContainsPrefixBlock(prefixLen)
}

func (addr *IPv6Address) ContainsSinglePrefixBlock(prefixLen BitCount) bool {
	return addr.init().ipAddressInternal.ContainsSinglePrefixBlock(prefixLen)
}

func (addr *IPv6Address) GetMinPrefixLengthForBlock() BitCount {
	return addr.init().ipAddressInternal.GetMinPrefixLengthForBlock()
}

func (addr *IPv6Address) GetPrefixLengthForSingleBlock() PrefixLen {
	return addr.init().ipAddressInternal.GetPrefixLengthForSingleBlock()
}

func (addr *IPv6Address) GetValue() *big.Int {
	return addr.init().section.GetValue()
}

func (addr *IPv6Address) GetUpperValue() *big.Int {
	return addr.init().section.GetUpperValue()
}

func (addr *IPv6Address) GetIPAddr() net.IPAddr {
	return net.IPAddr{
		IP:   addr.GetIP(),
		Zone: string(addr.GetZone()),
	}
}

func (addr *IPv6Address) GetIP() net.IP {
	return addr.GetBytes()
}

func (addr *IPv6Address) CopyIP(bytes net.IP) net.IP {
	return addr.CopyBytes(bytes)
}

func (addr *IPv6Address) GetUpperIP() net.IP {
	return addr.GetUpperBytes()
}

func (addr *IPv6Address) CopyUpperIP(bytes net.IP) net.IP {
	return addr.CopyUpperBytes(bytes)
}

func (addr *IPv6Address) GetBytes() []byte {
	return addr.init().section.GetBytes()
}

func (addr *IPv6Address) GetUpperBytes() []byte {
	return addr.init().section.GetUpperBytes()
}

func (addr *IPv6Address) CopyBytes(bytes []byte) []byte {
	return addr.init().section.CopyBytes(bytes)
}

func (addr *IPv6Address) CopyUpperBytes(bytes []byte) []byte {
	return addr.init().section.CopyUpperBytes(bytes)
}

func (addr *IPv6Address) IsMax() bool {
	return addr.init().section.IsMax()
}

func (addr *IPv6Address) IncludesMax() bool {
	return addr.init().section.IncludesMax()
}

// TestBit computes (this & (1 << n)) != 0), using the lower value of this segment.
func (addr *IPv6Address) TestBit(n BitCount) bool {
	return addr.init().testBit(n)
}

// Returns true if the bit in the lower value of this segment at the given index is 1, where index 0 is the most significant bit.
func (addr *IPv6Address) IsOneBit(bitIndex BitCount) bool {
	return addr.init().isOneBit(bitIndex)
}

func (addr *IPv6Address) CompareTo(item AddressItem) int {
	return CountComparator.Compare(addr.init(), item)
}

func (addr *IPv6Address) PrefixEquals(other AddressType) bool {
	return addr.init().prefixEquals(other)
}

func (addr *IPv6Address) PrefixContains(other AddressType) bool {
	return addr.init().prefixContains(other)
}

func (addr *IPv6Address) Contains(other AddressType) bool {
	return other.getAddrType() == ipv6Type && addr.init().section.sameCountTypeContains(other.ToAddress().GetSection()) &&
		addr.isSameZone(other.ToAddress())
}

func (addr *IPv6Address) Equals(other AddressType) bool {
	return other.getAddrType() == ipv6Type && addr.init().section.sameCountTypeEquals(other.ToAddress().GetSection()) &&
		addr.isSameZone(other.ToAddress())
}

func (addr *IPv6Address) GetMaxSegmentValue() SegInt {
	return addr.init().getMaxSegmentValue()
}

func (addr *IPv6Address) WithoutZone() *IPv6Address {
	if addr.HasZone() {
		return NewIPv6Address(addr.GetSection())
	}
	return addr
}

func (addr *IPv6Address) ToSequentialRange() *IPv6AddressSeqRange {
	if addr == nil {
		return nil
	}
	addr = addr.init().WithoutPrefixLen().WithoutZone()
	return newSeqRangeUnchecked(
		addr.GetLowerIPAddress(),
		addr.GetUpperIPAddress(),
		addr.IsMultiple()).ToIPv6SequentialRange()
}

func (addr *IPv6Address) ToAddressString() *IPAddressString {
	return addr.init().ToIPAddress().ToAddressString()
}

func (addr *IPv6Address) IncludesZeroHostLen(networkPrefixLength BitCount) bool {
	return addr.init().includesZeroHostLen(networkPrefixLength)
}

func (addr *IPv6Address) IncludesMaxHostLen(networkPrefixLength BitCount) bool {
	return addr.init().includesMaxHostLen(networkPrefixLength)
}

// IsLoopback returns whether this address is a loopback address, such as
// [::1] (aka [0:0:0:0:0:0:0:1]) or 127.0.0.1
func (addr *IPv6Address) IsLoopback() bool {
	if addr.section == nil {
		return false
	}
	//::1
	i := 0
	lim := addr.GetSegmentCount() - 1
	for ; i < lim; i++ {
		if !addr.GetSegment(i).IsZero() {
			return false
		}
	}
	return addr.GetSegment(i).Matches(1)
}

func (addr *IPv6Address) Iterator() IPv6AddressIterator {
	return ipv6AddressIterator{addr.init().addrIterator(nil)}
}

func (addr *IPv6Address) PrefixIterator() IPv6AddressIterator {
	return ipv6AddressIterator{addr.init().prefixIterator(false)}
}

func (addr *IPv6Address) PrefixBlockIterator() IPv6AddressIterator {
	return ipv6AddressIterator{addr.init().prefixIterator(true)}
}

func (addr *IPv6Address) BlockIterator(segmentCount int) IPv6AddressIterator {
	return ipv6AddressIterator{addr.init().blockIterator(segmentCount)}
}

func (addr *IPv6Address) SequentialBlockIterator() IPv6AddressIterator {
	return ipv6AddressIterator{addr.init().sequentialBlockIterator()}
}

func (addr *IPv6Address) GetSequentialBlockIndex() int {
	return addr.init().getSequentialBlockIndex()
}

func (addr *IPv6Address) GetSequentialBlockCount() *big.Int {
	return addr.getSequentialBlockCount()
}

func (addr *IPv6Address) IncrementBoundary(increment int64) *IPv6Address {
	return addr.init().incrementBoundary(increment).ToIPv6Address()
}

func (addr *IPv6Address) Increment(increment int64) *IPv6Address {
	return addr.init().increment(increment).ToIPv6Address()
}

func (addr *IPv6Address) SpanWithPrefixBlocks() []*IPv6Address {
	if addr.IsSequential() {
		if addr.IsSinglePrefixBlock() {
			return []*IPv6Address{addr}
		}
		wrapped := WrappedIPAddress{addr.ToIPAddress()}
		spanning := getSpanningPrefixBlocks(wrapped, wrapped)
		return cloneToIPv6Addrs(spanning)
	}
	wrapped := WrappedIPAddress{addr.ToIPAddress()}
	return cloneToIPv6Addrs(spanWithPrefixBlocks(wrapped))
}

func (addr *IPv6Address) SpanWithPrefixBlocksTo(other *IPv6Address) []*IPv6Address {
	return cloneToIPv6Addrs(
		getSpanningPrefixBlocks(
			WrappedIPAddress{addr.ToIPAddress()},
			WrappedIPAddress{other.ToIPAddress()},
		),
	)
}

func (addr *IPv6Address) SpanWithSequentialBlocks() []*IPv6Address {
	if addr.IsSequential() {
		return []*IPv6Address{addr}
	}
	wrapped := WrappedIPAddress{addr.ToIPAddress()}
	return cloneToIPv6Addrs(spanWithSequentialBlocks(wrapped))
}

func (addr *IPv6Address) SpanWithSequentialBlocksTo(other *IPv6Address) []*IPv6Address {
	return cloneToIPv6Addrs(
		getSpanningSequentialBlocks(
			WrappedIPAddress{addr.ToIPAddress()},
			WrappedIPAddress{other.ToIPAddress()},
		),
	)
}

func (addr *IPv6Address) CoverWithPrefixBlockTo(other *IPv6Address) *IPv6Address {
	return addr.init().coverWithPrefixBlockTo(other.ToIPAddress()).ToIPv6Address()
}

func (addr *IPv6Address) CoverWithPrefixBlock() *IPv6Address {
	return addr.init().coverWithPrefixBlock().ToIPv6Address()
}

//
// MergeToSequentialBlocks merges this with the list of addresses to produce the smallest array of blocks that are sequential
//
// The resulting array is sorted from lowest address value to highest, regardless of the size of each prefix block.
func (addr *IPv6Address) MergeToSequentialBlocks(addrs ...*IPv6Address) []*IPv6Address {
	series := cloneIPv6Addrs(addr, addrs)
	blocks := getMergedSequentialBlocks(series)
	return cloneToIPv6Addrs(blocks)
}

//
// MergeToPrefixBlocks merges this with the list of sections to produce the smallest array of prefix blocks.
//
// The resulting array is sorted from lowest address value to highest, regardless of the size of each prefix block.
func (addr *IPv6Address) MergeToPrefixBlocks(addrs ...*IPv6Address) []*IPv6Address {
	series := cloneIPv6Addrs(addr, addrs)
	blocks := getMergedPrefixBlocks(series)
	return cloneToIPv6Addrs(blocks)
}

func (addr *IPv6Address) ReverseBytes() (*IPv6Address, IncompatibleAddressError) {
	res, err := addr.GetSection().ReverseBytes()
	if err != nil {
		return nil, err
	}
	return addr.checkIdentity(res), nil
}

func (addr *IPv6Address) ReverseBits(perByte bool) (*IPv6Address, IncompatibleAddressError) {
	res, err := addr.GetSection().ReverseBits(perByte)
	if err != nil {
		return nil, err
	}
	return addr.checkIdentity(res), nil
}

func (addr *IPv6Address) ReverseSegments() *IPv6Address {
	return addr.checkIdentity(addr.GetSection().ReverseSegments())
}

func (addr *IPv6Address) GetLeadingBitCount(ones bool) BitCount {
	return addr.GetSection().GetLeadingBitCount(ones)
}

func (addr *IPv6Address) GetTrailingBitCount(ones bool) BitCount {
	return addr.GetSection().GetTrailingBitCount(ones)
}

func (addr IPv6Address) String() string {
	return addr.init().addressInternal.String()
}

func (addr *IPv6Address) ToCanonicalString() string {
	return addr.init().toCanonicalString()
}

func (addr *IPv6Address) ToNormalizedString() string {
	return addr.init().toNormalizedString()
}

func (addr *IPv6Address) ToCompressedString() string {
	return addr.init().toCompressedString()
}

func (addr *IPv6Address) ToCanonicalWildcardString() string {
	return addr.init().toCanonicalWildcardString()
}

func (addr *IPv6Address) ToNormalizedWildcardString() string {
	return addr.init().toNormalizedWildcardString()
}

func (addr *IPv6Address) ToSegmentedBinaryString() string {
	return addr.init().toSegmentedBinaryString()
}

func (addr *IPv6Address) ToSQLWildcardString() string {
	return addr.init().toSQLWildcardString()
}

func (addr *IPv6Address) ToFullString() string {
	return addr.init().toFullString()
}

func (addr *IPv6Address) ToPrefixLenString() string {
	return addr.init().toPrefixLenString()
}

func (addr *IPv6Address) ToSubnetString() string {
	return addr.init().toSubnetString()
}

func (addr *IPv6Address) ToCompressedWildcardString() string {
	return addr.init().toCompressedWildcardString()
}

func (addr *IPv6Address) ToReverseDNSString() (string, IncompatibleAddressError) {
	return addr.init().toReverseDNSString()
}

func (addr *IPv6Address) ToHexString(with0xPrefix bool) (string, IncompatibleAddressError) {
	return addr.init().toHexString(with0xPrefix)
}

func (addr *IPv6Address) ToOctalString(with0Prefix bool) (string, IncompatibleAddressError) {
	return addr.init().toOctalString(with0Prefix)
}

func (addr *IPv6Address) ToBinaryString(with0bPrefix bool) (string, IncompatibleAddressError) {
	return addr.init().toBinaryString(with0bPrefix)
}

// ToMixedString produces the mixed IPv6/IPv4 string.  It is the shortest such string (ie fully compressed).
// For some address sections with ranges of values in the IPv4 part of the address, there is not mixed string, and an error is returned.
func (addr *IPv6Address) ToMixedString() (string, IncompatibleAddressError) {
	if addr.hasZone() {
		var cacheField **string
		cacheField = &addr.getStringCache().mixedString
		return cacheStrErr(cacheField,
			func() (string, IncompatibleAddressError) {
				return addr.GetSection().toMixedStringZoned(addr.zone)
			})
	}
	return addr.GetSection().toMixedString()

}

func (addr *IPv6Address) ToAddress() *Address {
	if addr != nil {
		addr = addr.init()
	}
	return (*Address)(unsafe.Pointer(addr))
}

func (addr *IPv6Address) ToIPAddress() *IPAddress {
	if addr != nil {
		addr = addr.init()
	}
	return (*IPAddress)(unsafe.Pointer(addr))
}
