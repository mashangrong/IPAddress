package ipaddr

import (
	"math/big"
)

type MACSegInt uint8

const useMACSegmentCache = true

type macSegmentValues struct {
	value      MACSegInt
	upperValue MACSegInt
	cache      divCache
}

func (seg *macSegmentValues) getAddrType() addrType {
	return macType
}

func (seg *macSegmentValues) includesZero() bool {
	return seg.value == 0
}

func (seg *macSegmentValues) includesMax() bool {
	return seg.upperValue == 0xff
}

func (seg *macSegmentValues) isMultiple() bool {
	return seg.value != seg.upperValue
}

func (seg *macSegmentValues) getCount() *big.Int {
	return big.NewInt(int64((seg.upperValue - seg.value)) + 1)
}

func (seg *macSegmentValues) getBitCount() BitCount {
	return MACBitsPerSegment
}

func (seg *macSegmentValues) getByteCount() int {
	return MACBytesPerSegment
}

func (seg *macSegmentValues) getValue() *BigDivInt {
	return big.NewInt(int64(seg.value))
}

func (seg *macSegmentValues) getUpperValue() *BigDivInt {
	return big.NewInt(int64(seg.upperValue))
}

func (seg *macSegmentValues) getDivisionValue() DivInt {
	return DivInt(seg.value)
}

func (seg *macSegmentValues) getUpperDivisionValue() DivInt {
	return DivInt(seg.upperValue)
}

func (seg *macSegmentValues) getDivisionPrefixLength() PrefixLen {
	return nil
}

func (seg *macSegmentValues) getSegmentValue() SegInt {
	return SegInt(seg.value)
}

func (seg *macSegmentValues) getUpperSegmentValue() SegInt {
	return SegInt(seg.upperValue)
}

func (seg *macSegmentValues) calcBytesInternal() (bytes, upperBytes []byte) {
	bytes = []byte{byte(seg.value)}
	if seg.isMultiple() {
		upperBytes = []byte{byte(seg.upperValue)}
	} else {
		upperBytes = bytes
	}
	return
}

func (seg *macSegmentValues) deriveNew(val, upperVal DivInt, _ PrefixLen) divisionValues {
	return newMACSegmentValues(MACSegInt(val), MACSegInt(upperVal))
}

func (seg *macSegmentValues) deriveNewSeg(val SegInt, _ PrefixLen) divisionValues {
	return newMACSegmentVal(MACSegInt(val))
}

func (seg *macSegmentValues) deriveNewMultiSeg(val, upperVal SegInt, _ PrefixLen) divisionValues {
	return newMACSegmentValues(MACSegInt(val), MACSegInt(upperVal))
}

func (seg *macSegmentValues) getCache() *divCache {
	return &seg.cache
}

var _ divisionValues = &macSegmentValues{}

var zeroMACSeg = NewMACSegment(0)

type MACAddressSegment struct {
	addressSegmentInternal
}

func (seg *MACAddressSegment) init() *MACAddressSegment {
	if seg.divisionValues == nil {
		return zeroMACSeg
	}
	return seg
}

// We must override getBitCount, getByteCount and others for the case when we construct as the zero value

func (seg *MACAddressSegment) GetBitCount() BitCount {
	return IPv4BitsPerSegment
}

func (seg *MACAddressSegment) GetByteCount() int {
	return IPv4BytesPerSegment
}

func (seg *MACAddressSegment) GetMaxValue() MACSegInt {
	return 0xff
}

func (seg *MACAddressSegment) setString(
	addressStr string,
	isStandardString bool,
	lowerStringStartIndex,
	lowerStringEndIndex int,
	originalLowerValue SegInt) {
	if cache := seg.getCache(); cache != nil {
		//TODO atomic writes only, the caches are shared, use cacheStr
		//TODO also, should I think about whether I should also set cached.string here?
		// The lower level types might be using it, though I'm not sure, depends on string generation and those params and whatnot
		// Actually, seems GetString in ipsegment.go uses it, so no prob here
		// BUT these two use it
		//func (div *addressDivisionInternal) getStringAsLower() string {
		//func (div *addressDivisionInternal) getString() string {
		// which are called from:
		// func (div *addressDivisionInternal) GetString() string
		// func (div *addressDivisionInternal) GetWildcardString() string {
		// So yeah, because you are not really overriding, you can convert to AddressDivision and then use that field
		// It seems this differs from Java, only IPAddressDivision uses cachedString
		// So either you make func (div *addressDivisionInternal) GetString() check for type and scale upwards, or you write to cacherString here too
		// The latter is  enough because we have code that always scales up to IPAddressSegment.
		// In go, cachedWildcardString is used ONLY by ip segments in GetWildcardString
		// in java, it used there too, but also by segments and divisions in general, as cachedString is present only in IPAddressDivision
		// So here just change this to use cachedString
		if cache.cachedWildcardString == nil && isStandardString && originalLowerValue == seg.getSegmentValue() {
			str := addressStr[lowerStringStartIndex:lowerStringEndIndex]
			cache.cachedWildcardString = &str
		}
	}
}

func (seg *MACAddressSegment) setRangeString(
	addressStr string,
	isStandardRangeString bool,
	lowerStringStartIndex,
	upperStringEndIndex int,
	rangeLower,
	rangeUpper SegInt) {
	if cache := seg.getCache(); cache != nil {
		//TODO atomic writes only, the caches are shared, use cacheStr
		//
		// TODO see above about cachedWildcardString.  Use only cachedString here.
		if cache.cachedWildcardString == nil {
			if seg.IsFullRange() {
				cache.cachedWildcardString = &segmentWildcardStr
			} else if isStandardRangeString && rangeLower == seg.getSegmentValue() && rangeUpper == seg.getUpperSegmentValue() {
				str := addressStr[lowerStringStartIndex:upperStringEndIndex]
				cache.cachedWildcardString = &str
			}
		}
	}
}

func (seg *MACAddressSegment) Iterator() MACSegmentIterator {
	return macSegmentIterator{seg.iterator()}
}

func (seg *MACAddressSegment) prefixBlockIterator(segmentPrefixLen BitCount) MACSegmentIterator {
	return macSegmentIterator{seg.prefixedBlockIterator(segmentPrefixLen)}
}

func (seg *MACAddressSegment) prefixIterator(segmentPrefixLen BitCount) MACSegmentIterator {
	return macSegmentIterator{seg.prefixedIterator(segmentPrefixLen)}
}

func (seg *MACAddressSegment) ReverseBits(_ bool) (res *MACAddressSegment, err IncompatibleAddressError) {
	if seg.divisionValues == nil {
		res = seg
		return
	}
	if seg.IsMultiple() {
		if isReversible, _ := seg.isReversibleRange(false); isReversible {
			res = seg
			return
		}
		err = &incompatibleAddressError{addressError{key: "ipaddress.error.reverseRange"}}
		return
	}
	oldVal := MACSegInt(seg.GetSegmentValue())
	val := MACSegInt(reverseUint8(uint8(oldVal)))
	if oldVal == val {
		res = seg
	} else {
		res = NewMACSegment(val)
	}
	return
}

func (seg *MACAddressSegment) ReverseBytes() (*MACAddressSegment, IncompatibleAddressError) {
	return seg, nil
}

func (seg *MACAddressSegment) ToAddressSegment() *AddressSegment {
	if seg == nil {
		return nil
	}
	return (*AddressSegment)(seg.init())
}

func NewMACSegment(val MACSegInt) *MACAddressSegment {
	return newMACSegment(newMACSegmentVal(val))
}

func NewMACRangeSegment(val, upperVal MACSegInt) *MACAddressSegment {
	return newMACSegment(newMACSegmentValues(val, upperVal))
}

func newMACSegment(vals *macSegmentValues) *MACAddressSegment {
	return &MACAddressSegment{
		addressSegmentInternal{
			addressDivisionInternal{
				addressDivisionBase{vals},
			},
		},
	}
}

var (
	allRangeValsMAC = &macSegmentValues{
		upperValue: MACMaxValuePerSegment,
	}
	segmentCacheMAC = makeSegmentCacheMAC()
)

func checkValuesMAC(value, upperValue MACSegInt, result *macSegmentValues) { //TODO remove eventually
	if result.value != value || result.upperValue != upperValue {
		panic("huh")
	}
}

func makeSegmentCacheMAC() (segmentCacheMAC []macSegmentValues) {
	if useMACSegmentCache {
		segmentCacheMAC = make([]macSegmentValues, MACMaxValuePerSegment+1)
		for i := range segmentCacheMAC {
			vals := &segmentCacheMAC[i]
			segi := MACSegInt(i)
			vals.value = segi
			vals.upperValue = segi
		}
	}
	return
}

func newMACSegmentVal(value MACSegInt) *macSegmentValues {
	if useMACSegmentCache {
		result := &segmentCacheMAC[value]
		checkValuesMAC(value, value, result)
		return result
	}
	return &macSegmentValues{value: value, upperValue: value}
}

func newMACSegmentValues(value, upperValue MACSegInt) *macSegmentValues {
	if value == upperValue {
		return newMACSegmentVal(value)
	}
	if useMACSegmentCache && value == 0 && upperValue == MACMaxValuePerSegment {
		return allRangeValsMAC
	}
	return &macSegmentValues{value: value, upperValue: upperValue}
}
