// Code generated by protoc-gen-go.
// source: mfk.proto
// DO NOT EDIT!

package mfklib

import proto "code.google.com/p/goprotobuf/proto"
import json "encoding/json"
import math "math"

// Reference proto, json, and math imports to suppress error if they are not otherwise used.
var _ = proto.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf

type Vote_VoteType int32

const (
	Vote_UNSET Vote_VoteType = 0
	Vote_MFK   Vote_VoteType = 1
	Vote_MKF   Vote_VoteType = 2
	Vote_FMK   Vote_VoteType = 3
	Vote_FKM   Vote_VoteType = 4
	Vote_KMF   Vote_VoteType = 5
	Vote_KFM   Vote_VoteType = 6
	Vote_SKIP  Vote_VoteType = 7
)

var Vote_VoteType_name = map[int32]string{
	0: "UNSET",
	1: "MFK",
	2: "MKF",
	3: "FMK",
	4: "FKM",
	5: "KMF",
	6: "KFM",
	7: "SKIP",
}
var Vote_VoteType_value = map[string]int32{
	"UNSET": 0,
	"MFK":   1,
	"MKF":   2,
	"FMK":   3,
	"FKM":   4,
	"KMF":   5,
	"KFM":   6,
	"SKIP":  7,
}

func (x Vote_VoteType) Enum() *Vote_VoteType {
	p := new(Vote_VoteType)
	*p = x
	return p
}
func (x Vote_VoteType) String() string {
	return proto.EnumName(Vote_VoteType_name, int32(x))
}
func (x Vote_VoteType) MarshalJSON() ([]byte, error) {
	return json.Marshal(x.String())
}
func (x *Vote_VoteType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Vote_VoteType_value, data, "Vote_VoteType")
	if err != nil {
		return err
	}
	*x = Vote_VoteType(value)
	return nil
}

type ImageMetadata struct {
	Url              *string `protobuf:"bytes,1,opt,name=url" json:"url,omitempty"`
	Context          *string `protobuf:"bytes,2,opt,name=context" json:"context,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (this *ImageMetadata) Reset()         { *this = ImageMetadata{} }
func (this *ImageMetadata) String() string { return proto.CompactTextString(this) }
func (*ImageMetadata) ProtoMessage()       {}

func (this *ImageMetadata) GetUrl() string {
	if this != nil && this.Url != nil {
		return *this.Url
	}
	return ""
}

func (this *ImageMetadata) GetContext() string {
	if this != nil && this.Context != nil {
		return *this.Context
	}
	return ""
}

type Image struct {
	Metadata         *ImageMetadata `protobuf:"bytes,1,opt,name=metadata" json:"metadata,omitempty"`
	ContentType      *string        `protobuf:"bytes,2,opt,name=content_type" json:"content_type,omitempty"`
	Data             []byte         `protobuf:"bytes,3,opt,name=data" json:"data,omitempty"`
	XXX_unrecognized []byte         `json:"-"`
}

func (this *Image) Reset()         { *this = Image{} }
func (this *Image) String() string { return proto.CompactTextString(this) }
func (*Image) ProtoMessage()       {}

func (this *Image) GetMetadata() *ImageMetadata {
	if this != nil {
		return this.Metadata
	}
	return nil
}

func (this *Image) GetContentType() string {
	if this != nil && this.ContentType != nil {
		return *this.ContentType
	}
	return ""
}

func (this *Image) GetData() []byte {
	if this != nil {
		return this.Data
	}
	return nil
}

type Vote struct {
	TripleId         *int64         `protobuf:"varint,1,opt,name=triple_id" json:"triple_id,omitempty"`
	UserId           *string        `protobuf:"bytes,2,opt,name=user_id" json:"user_id,omitempty"`
	Vote             *Vote_VoteType `protobuf:"varint,3,opt,name=vote,enum=mfklib.Vote_VoteType" json:"vote,omitempty"`
	XXX_unrecognized []byte         `json:"-"`
}

func (this *Vote) Reset()         { *this = Vote{} }
func (this *Vote) String() string { return proto.CompactTextString(this) }
func (*Vote) ProtoMessage()       {}

func (this *Vote) GetTripleId() int64 {
	if this != nil && this.TripleId != nil {
		return *this.TripleId
	}
	return 0
}

func (this *Vote) GetUserId() string {
	if this != nil && this.UserId != nil {
		return *this.UserId
	}
	return ""
}

func (this *Vote) GetVote() Vote_VoteType {
	if this != nil && this.Vote != nil {
		return *this.Vote
	}
	return 0
}

type Triple struct {
	CreatorId        *string        `protobuf:"bytes,2,opt,name=creator_id" json:"creator_id,omitempty"`
	A                *Triple_Entity `protobuf:"bytes,3,opt,name=a" json:"a,omitempty"`
	B                *Triple_Entity `protobuf:"bytes,4,opt,name=b" json:"b,omitempty"`
	C                *Triple_Entity `protobuf:"bytes,5,opt,name=c" json:"c,omitempty"`
	XXX_unrecognized []byte         `json:"-"`
}

func (this *Triple) Reset()         { *this = Triple{} }
func (this *Triple) String() string { return proto.CompactTextString(this) }
func (*Triple) ProtoMessage()       {}

func (this *Triple) GetCreatorId() string {
	if this != nil && this.CreatorId != nil {
		return *this.CreatorId
	}
	return ""
}

func (this *Triple) GetA() *Triple_Entity {
	if this != nil {
		return this.A
	}
	return nil
}

func (this *Triple) GetB() *Triple_Entity {
	if this != nil {
		return this.B
	}
	return nil
}

func (this *Triple) GetC() *Triple_Entity {
	if this != nil {
		return this.C
	}
	return nil
}

type Triple_Entity struct {
	Name             *string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Image            *Image  `protobuf:"bytes,2,opt,name=image" json:"image,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (this *Triple_Entity) Reset()         { *this = Triple_Entity{} }
func (this *Triple_Entity) String() string { return proto.CompactTextString(this) }
func (*Triple_Entity) ProtoMessage()       {}

func (this *Triple_Entity) GetName() string {
	if this != nil && this.Name != nil {
		return *this.Name
	}
	return ""
}

func (this *Triple_Entity) GetImage() *Image {
	if this != nil {
		return this.Image
	}
	return nil
}

type TripleStats struct {
	Views            *uint64            `protobuf:"varint,1,opt,name=views" json:"views,omitempty"`
	Skips            *uint64            `protobuf:"varint,2,opt,name=skips" json:"skips,omitempty"`
	A                *TripleStats_Tally `protobuf:"bytes,3,opt,name=a" json:"a,omitempty"`
	B                *TripleStats_Tally `protobuf:"bytes,4,opt,name=b" json:"b,omitempty"`
	C                *TripleStats_Tally `protobuf:"bytes,5,opt,name=c" json:"c,omitempty"`
	XXX_unrecognized []byte             `json:"-"`
}

func (this *TripleStats) Reset()         { *this = TripleStats{} }
func (this *TripleStats) String() string { return proto.CompactTextString(this) }
func (*TripleStats) ProtoMessage()       {}

func (this *TripleStats) GetViews() uint64 {
	if this != nil && this.Views != nil {
		return *this.Views
	}
	return 0
}

func (this *TripleStats) GetSkips() uint64 {
	if this != nil && this.Skips != nil {
		return *this.Skips
	}
	return 0
}

func (this *TripleStats) GetA() *TripleStats_Tally {
	if this != nil {
		return this.A
	}
	return nil
}

func (this *TripleStats) GetB() *TripleStats_Tally {
	if this != nil {
		return this.B
	}
	return nil
}

func (this *TripleStats) GetC() *TripleStats_Tally {
	if this != nil {
		return this.C
	}
	return nil
}

type TripleStats_Tally struct {
	Marry            *uint64 `protobuf:"varint,1,opt,name=marry" json:"marry,omitempty"`
	Fuck             *uint64 `protobuf:"varint,2,opt,name=fuck" json:"fuck,omitempty"`
	Kill             *uint64 `protobuf:"varint,3,opt,name=kill" json:"kill,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (this *TripleStats_Tally) Reset()         { *this = TripleStats_Tally{} }
func (this *TripleStats_Tally) String() string { return proto.CompactTextString(this) }
func (*TripleStats_Tally) ProtoMessage()       {}

func (this *TripleStats_Tally) GetMarry() uint64 {
	if this != nil && this.Marry != nil {
		return *this.Marry
	}
	return 0
}

func (this *TripleStats_Tally) GetFuck() uint64 {
	if this != nil && this.Fuck != nil {
		return *this.Fuck
	}
	return 0
}

func (this *TripleStats_Tally) GetKill() uint64 {
	if this != nil && this.Kill != nil {
		return *this.Kill
	}
	return 0
}

type WrappedImageMetadata struct {
	Metadata         *ImageMetadata `protobuf:"bytes,1,opt,name=metadata" json:"metadata,omitempty"`
	XXX_unrecognized []byte         `json:"-"`
}

func (this *WrappedImageMetadata) Reset()         { *this = WrappedImageMetadata{} }
func (this *WrappedImageMetadata) String() string { return proto.CompactTextString(this) }
func (*WrappedImageMetadata) ProtoMessage()       {}

func (this *WrappedImageMetadata) GetMetadata() *ImageMetadata {
	if this != nil {
		return this.Metadata
	}
	return nil
}

type ImageSearchResponse struct {
	Image            []*WrappedImageMetadata `protobuf:"bytes,1,rep,name=image" json:"image,omitempty"`
	XXX_unrecognized []byte                  `json:"-"`
}

func (this *ImageSearchResponse) Reset()         { *this = ImageSearchResponse{} }
func (this *ImageSearchResponse) String() string { return proto.CompactTextString(this) }
func (*ImageSearchResponse) ProtoMessage()       {}

type MakeTripleRequest struct {
	A                *MakeTripleRequest_Entity `protobuf:"bytes,3,opt,name=a" json:"a,omitempty"`
	B                *MakeTripleRequest_Entity `protobuf:"bytes,4,opt,name=b" json:"b,omitempty"`
	C                *MakeTripleRequest_Entity `protobuf:"bytes,5,opt,name=c" json:"c,omitempty"`
	XXX_unrecognized []byte                    `json:"-"`
}

func (this *MakeTripleRequest) Reset()         { *this = MakeTripleRequest{} }
func (this *MakeTripleRequest) String() string { return proto.CompactTextString(this) }
func (*MakeTripleRequest) ProtoMessage()       {}

func (this *MakeTripleRequest) GetA() *MakeTripleRequest_Entity {
	if this != nil {
		return this.A
	}
	return nil
}

func (this *MakeTripleRequest) GetB() *MakeTripleRequest_Entity {
	if this != nil {
		return this.B
	}
	return nil
}

func (this *MakeTripleRequest) GetC() *MakeTripleRequest_Entity {
	if this != nil {
		return this.C
	}
	return nil
}

type MakeTripleRequest_Entity struct {
	Name             *string               `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Image            *WrappedImageMetadata `protobuf:"bytes,2,opt,name=image" json:"image,omitempty"`
	XXX_unrecognized []byte                `json:"-"`
}

func (this *MakeTripleRequest_Entity) Reset()         { *this = MakeTripleRequest_Entity{} }
func (this *MakeTripleRequest_Entity) String() string { return proto.CompactTextString(this) }
func (*MakeTripleRequest_Entity) ProtoMessage()       {}

func (this *MakeTripleRequest_Entity) GetName() string {
	if this != nil && this.Name != nil {
		return *this.Name
	}
	return ""
}

func (this *MakeTripleRequest_Entity) GetImage() *WrappedImageMetadata {
	if this != nil {
		return this.Image
	}
	return nil
}

type MakeTripleResponse struct {
	TripleId         *int64 `protobuf:"varint,1,opt,name=triple_id" json:"triple_id,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (this *MakeTripleResponse) Reset()         { *this = MakeTripleResponse{} }
func (this *MakeTripleResponse) String() string { return proto.CompactTextString(this) }
func (*MakeTripleResponse) ProtoMessage()       {}

func (this *MakeTripleResponse) GetTripleId() int64 {
	if this != nil && this.TripleId != nil {
		return *this.TripleId
	}
	return 0
}

func init() {
	proto.RegisterEnum("mfklib.Vote_VoteType", Vote_VoteType_name, Vote_VoteType_value)
}
