// Code generated by MockGen. DO NOT EDIT.
// Source: store/sqldb/database.go

// Package sqldb is a generated GoMock package.
package sqldb

import (
	reflect "reflect"

	types "github.com/bnb-chain/greenfield-storage-provider/service/types"
	types0 "github.com/bnb-chain/greenfield/x/sp/types"
	types1 "github.com/bnb-chain/greenfield/x/storage/types"
	gomock "github.com/golang/mock/gomock"
)

// MockJob is a mock of Job interface.
type MockJob struct {
	ctrl     *gomock.Controller
	recorder *MockJobMockRecorder
}

// MockJobMockRecorder is the mock recorder for MockJob.
type MockJobMockRecorder struct {
	mock *MockJob
}

// NewMockJob creates a new mock instance.
func NewMockJob(ctrl *gomock.Controller) *MockJob {
	mock := &MockJob{ctrl: ctrl}
	mock.recorder = &MockJobMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockJob) EXPECT() *MockJobMockRecorder {
	return m.recorder
}

// CreateUploadJob mocks base method.
func (m *MockJob) CreateUploadJob(objectInfo *types1.ObjectInfo) (*types.JobContext, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUploadJob", objectInfo)
	ret0, _ := ret[0].(*types.JobContext)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUploadJob indicates an expected call of CreateUploadJob.
func (mr *MockJobMockRecorder) CreateUploadJob(objectInfo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUploadJob", reflect.TypeOf((*MockJob)(nil).CreateUploadJob), objectInfo)
}

// GetJobByID mocks base method.
func (m *MockJob) GetJobByID(jobID uint64) (*types.JobContext, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetJobByID", jobID)
	ret0, _ := ret[0].(*types.JobContext)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetJobByID indicates an expected call of GetJobByID.
func (mr *MockJobMockRecorder) GetJobByID(jobID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetJobByID", reflect.TypeOf((*MockJob)(nil).GetJobByID), jobID)
}

// GetJobByObjectID mocks base method.
func (m *MockJob) GetJobByObjectID(objectID uint64) (*types.JobContext, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetJobByObjectID", objectID)
	ret0, _ := ret[0].(*types.JobContext)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetJobByObjectID indicates an expected call of GetJobByObjectID.
func (mr *MockJobMockRecorder) GetJobByObjectID(objectID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetJobByObjectID", reflect.TypeOf((*MockJob)(nil).GetJobByObjectID), objectID)
}

// UpdateJobState mocks base method.
func (m *MockJob) UpdateJobState(objectID uint64, state types.JobState) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateJobState", objectID, state)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateJobState indicates an expected call of UpdateJobState.
func (mr *MockJobMockRecorder) UpdateJobState(objectID, state interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateJobState", reflect.TypeOf((*MockJob)(nil).UpdateJobState), objectID, state)
}

// MockObject is a mock of Object interface.
type MockObject struct {
	ctrl     *gomock.Controller
	recorder *MockObjectMockRecorder
}

// MockObjectMockRecorder is the mock recorder for MockObject.
type MockObjectMockRecorder struct {
	mock *MockObject
}

// NewMockObject creates a new mock instance.
func NewMockObject(ctrl *gomock.Controller) *MockObject {
	mock := &MockObject{ctrl: ctrl}
	mock.recorder = &MockObjectMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockObject) EXPECT() *MockObjectMockRecorder {
	return m.recorder
}

// GetObjectInfo mocks base method.
func (m *MockObject) GetObjectInfo(objectID uint64) (*types1.ObjectInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetObjectInfo", objectID)
	ret0, _ := ret[0].(*types1.ObjectInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetObjectInfo indicates an expected call of GetObjectInfo.
func (mr *MockObjectMockRecorder) GetObjectInfo(objectID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetObjectInfo", reflect.TypeOf((*MockObject)(nil).GetObjectInfo), objectID)
}

// SetObjectInfo mocks base method.
func (m *MockObject) SetObjectInfo(objectID uint64, objectInfo *types1.ObjectInfo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetObjectInfo", objectID, objectInfo)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetObjectInfo indicates an expected call of SetObjectInfo.
func (mr *MockObjectMockRecorder) SetObjectInfo(objectID, objectInfo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetObjectInfo", reflect.TypeOf((*MockObject)(nil).SetObjectInfo), objectID, objectInfo)
}

// MockObjectIntegrity is a mock of ObjectIntegrity interface.
type MockObjectIntegrity struct {
	ctrl     *gomock.Controller
	recorder *MockObjectIntegrityMockRecorder
}

// MockObjectIntegrityMockRecorder is the mock recorder for MockObjectIntegrity.
type MockObjectIntegrityMockRecorder struct {
	mock *MockObjectIntegrity
}

// NewMockObjectIntegrity creates a new mock instance.
func NewMockObjectIntegrity(ctrl *gomock.Controller) *MockObjectIntegrity {
	mock := &MockObjectIntegrity{ctrl: ctrl}
	mock.recorder = &MockObjectIntegrityMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockObjectIntegrity) EXPECT() *MockObjectIntegrityMockRecorder {
	return m.recorder
}

// GetObjectIntegrity mocks base method.
func (m *MockObjectIntegrity) GetObjectIntegrity(objectID uint64) (*IntegrityMeta, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetObjectIntegrity", objectID)
	ret0, _ := ret[0].(*IntegrityMeta)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetObjectIntegrity indicates an expected call of GetObjectIntegrity.
func (mr *MockObjectIntegrityMockRecorder) GetObjectIntegrity(objectID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetObjectIntegrity", reflect.TypeOf((*MockObjectIntegrity)(nil).GetObjectIntegrity), objectID)
}

// SetObjectIntegrity mocks base method.
func (m *MockObjectIntegrity) SetObjectIntegrity(integrity *IntegrityMeta) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetObjectIntegrity", integrity)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetObjectIntegrity indicates an expected call of SetObjectIntegrity.
func (mr *MockObjectIntegrityMockRecorder) SetObjectIntegrity(integrity interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetObjectIntegrity", reflect.TypeOf((*MockObjectIntegrity)(nil).SetObjectIntegrity), integrity)
}

// MockSPInfo is a mock of SPInfo interface.
type MockSPInfo struct {
	ctrl     *gomock.Controller
	recorder *MockSPInfoMockRecorder
}

// MockSPInfoMockRecorder is the mock recorder for MockSPInfo.
type MockSPInfoMockRecorder struct {
	mock *MockSPInfo
}

// NewMockSPInfo creates a new mock instance.
func NewMockSPInfo(ctrl *gomock.Controller) *MockSPInfo {
	mock := &MockSPInfo{ctrl: ctrl}
	mock.recorder = &MockSPInfoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSPInfo) EXPECT() *MockSPInfoMockRecorder {
	return m.recorder
}

// FetchAllSp mocks base method.
func (m *MockSPInfo) FetchAllSp(status ...types0.Status) ([]*types0.StorageProvider, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range status {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "FetchAllSp", varargs...)
	ret0, _ := ret[0].([]*types0.StorageProvider)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchAllSp indicates an expected call of FetchAllSp.
func (mr *MockSPInfoMockRecorder) FetchAllSp(status ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchAllSp", reflect.TypeOf((*MockSPInfo)(nil).FetchAllSp), status...)
}

// FetchAllSpWithoutOwnSp mocks base method.
func (m *MockSPInfo) FetchAllSpWithoutOwnSp(status ...types0.Status) ([]*types0.StorageProvider, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range status {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "FetchAllSpWithoutOwnSp", varargs...)
	ret0, _ := ret[0].([]*types0.StorageProvider)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchAllSpWithoutOwnSp indicates an expected call of FetchAllSpWithoutOwnSp.
func (mr *MockSPInfoMockRecorder) FetchAllSpWithoutOwnSp(status ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchAllSpWithoutOwnSp", reflect.TypeOf((*MockSPInfo)(nil).FetchAllSpWithoutOwnSp), status...)
}

// GetOwnSpInfo mocks base method.
func (m *MockSPInfo) GetOwnSpInfo() (*types0.StorageProvider, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOwnSpInfo")
	ret0, _ := ret[0].(*types0.StorageProvider)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOwnSpInfo indicates an expected call of GetOwnSpInfo.
func (mr *MockSPInfoMockRecorder) GetOwnSpInfo() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOwnSpInfo", reflect.TypeOf((*MockSPInfo)(nil).GetOwnSpInfo))
}

// GetSpByAddress mocks base method.
func (m *MockSPInfo) GetSpByAddress(address string, addressType SpAddressType) (*types0.StorageProvider, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSpByAddress", address, addressType)
	ret0, _ := ret[0].(*types0.StorageProvider)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSpByAddress indicates an expected call of GetSpByAddress.
func (mr *MockSPInfoMockRecorder) GetSpByAddress(address, addressType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSpByAddress", reflect.TypeOf((*MockSPInfo)(nil).GetSpByAddress), address, addressType)
}

// GetSpByEndpoint mocks base method.
func (m *MockSPInfo) GetSpByEndpoint(endpoint string) (*types0.StorageProvider, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSpByEndpoint", endpoint)
	ret0, _ := ret[0].(*types0.StorageProvider)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSpByEndpoint indicates an expected call of GetSpByEndpoint.
func (mr *MockSPInfoMockRecorder) GetSpByEndpoint(endpoint interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSpByEndpoint", reflect.TypeOf((*MockSPInfo)(nil).GetSpByEndpoint), endpoint)
}

// SetOwnSpInfo mocks base method.
func (m *MockSPInfo) SetOwnSpInfo(sp *types0.StorageProvider) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetOwnSpInfo", sp)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetOwnSpInfo indicates an expected call of SetOwnSpInfo.
func (mr *MockSPInfoMockRecorder) SetOwnSpInfo(sp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetOwnSpInfo", reflect.TypeOf((*MockSPInfo)(nil).SetOwnSpInfo), sp)
}

// UpdateAllSp mocks base method.
func (m *MockSPInfo) UpdateAllSp(spList []*types0.StorageProvider) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAllSp", spList)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateAllSp indicates an expected call of UpdateAllSp.
func (mr *MockSPInfoMockRecorder) UpdateAllSp(spList interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAllSp", reflect.TypeOf((*MockSPInfo)(nil).UpdateAllSp), spList)
}

// MockStorageParam is a mock of StorageParam interface.
type MockStorageParam struct {
	ctrl     *gomock.Controller
	recorder *MockStorageParamMockRecorder
}

// MockStorageParamMockRecorder is the mock recorder for MockStorageParam.
type MockStorageParamMockRecorder struct {
	mock *MockStorageParam
}

// NewMockStorageParam creates a new mock instance.
func NewMockStorageParam(ctrl *gomock.Controller) *MockStorageParam {
	mock := &MockStorageParam{ctrl: ctrl}
	mock.recorder = &MockStorageParamMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorageParam) EXPECT() *MockStorageParamMockRecorder {
	return m.recorder
}

// GetStorageParams mocks base method.
func (m *MockStorageParam) GetStorageParams() (*types1.Params, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStorageParams")
	ret0, _ := ret[0].(*types1.Params)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStorageParams indicates an expected call of GetStorageParams.
func (mr *MockStorageParamMockRecorder) GetStorageParams() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStorageParams", reflect.TypeOf((*MockStorageParam)(nil).GetStorageParams))
}

// SetStorageParams mocks base method.
func (m *MockStorageParam) SetStorageParams(params *types1.Params) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetStorageParams", params)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetStorageParams indicates an expected call of SetStorageParams.
func (mr *MockStorageParamMockRecorder) SetStorageParams(params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetStorageParams", reflect.TypeOf((*MockStorageParam)(nil).SetStorageParams), params)
}

// MockTraffic is a mock of Traffic interface.
type MockTraffic struct {
	ctrl     *gomock.Controller
	recorder *MockTrafficMockRecorder
}

// MockTrafficMockRecorder is the mock recorder for MockTraffic.
type MockTrafficMockRecorder struct {
	mock *MockTraffic
}

// NewMockTraffic creates a new mock instance.
func NewMockTraffic(ctrl *gomock.Controller) *MockTraffic {
	mock := &MockTraffic{ctrl: ctrl}
	mock.recorder = &MockTrafficMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTraffic) EXPECT() *MockTrafficMockRecorder {
	return m.recorder
}

// CheckQuotaAndAddReadRecord mocks base method.
func (m *MockTraffic) CheckQuotaAndAddReadRecord(record *ReadRecord, quota *BucketQuota) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckQuotaAndAddReadRecord", record, quota)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckQuotaAndAddReadRecord indicates an expected call of CheckQuotaAndAddReadRecord.
func (mr *MockTrafficMockRecorder) CheckQuotaAndAddReadRecord(record, quota interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckQuotaAndAddReadRecord", reflect.TypeOf((*MockTraffic)(nil).CheckQuotaAndAddReadRecord), record, quota)
}

// GetBucketReadRecord mocks base method.
func (m *MockTraffic) GetBucketReadRecord(bucketID uint64, timeRange *TrafficTimeRange) ([]*ReadRecord, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBucketReadRecord", bucketID, timeRange)
	ret0, _ := ret[0].([]*ReadRecord)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBucketReadRecord indicates an expected call of GetBucketReadRecord.
func (mr *MockTrafficMockRecorder) GetBucketReadRecord(bucketID, timeRange interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBucketReadRecord", reflect.TypeOf((*MockTraffic)(nil).GetBucketReadRecord), bucketID, timeRange)
}

// GetBucketTraffic mocks base method.
func (m *MockTraffic) GetBucketTraffic(bucketID uint64, yearMonth string) (*BucketTraffic, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBucketTraffic", bucketID, yearMonth)
	ret0, _ := ret[0].(*BucketTraffic)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBucketTraffic indicates an expected call of GetBucketTraffic.
func (mr *MockTrafficMockRecorder) GetBucketTraffic(bucketID, yearMonth interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBucketTraffic", reflect.TypeOf((*MockTraffic)(nil).GetBucketTraffic), bucketID, yearMonth)
}

// GetObjectReadRecord mocks base method.
func (m *MockTraffic) GetObjectReadRecord(objectID uint64, timeRange *TrafficTimeRange) ([]*ReadRecord, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetObjectReadRecord", objectID, timeRange)
	ret0, _ := ret[0].([]*ReadRecord)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetObjectReadRecord indicates an expected call of GetObjectReadRecord.
func (mr *MockTrafficMockRecorder) GetObjectReadRecord(objectID, timeRange interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetObjectReadRecord", reflect.TypeOf((*MockTraffic)(nil).GetObjectReadRecord), objectID, timeRange)
}

// GetReadRecord mocks base method.
func (m *MockTraffic) GetReadRecord(timeRange *TrafficTimeRange) ([]*ReadRecord, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetReadRecord", timeRange)
	ret0, _ := ret[0].([]*ReadRecord)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetReadRecord indicates an expected call of GetReadRecord.
func (mr *MockTrafficMockRecorder) GetReadRecord(timeRange interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReadRecord", reflect.TypeOf((*MockTraffic)(nil).GetReadRecord), timeRange)
}

// GetUserReadRecord mocks base method.
func (m *MockTraffic) GetUserReadRecord(userAddress string, timeRange *TrafficTimeRange) ([]*ReadRecord, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserReadRecord", userAddress, timeRange)
	ret0, _ := ret[0].([]*ReadRecord)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserReadRecord indicates an expected call of GetUserReadRecord.
func (mr *MockTrafficMockRecorder) GetUserReadRecord(userAddress, timeRange interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserReadRecord", reflect.TypeOf((*MockTraffic)(nil).GetUserReadRecord), userAddress, timeRange)
}

// MockServiceConfig is a mock of ServiceConfig interface.
type MockServiceConfig struct {
	ctrl     *gomock.Controller
	recorder *MockServiceConfigMockRecorder
}

// MockServiceConfigMockRecorder is the mock recorder for MockServiceConfig.
type MockServiceConfigMockRecorder struct {
	mock *MockServiceConfig
}

// NewMockServiceConfig creates a new mock instance.
func NewMockServiceConfig(ctrl *gomock.Controller) *MockServiceConfig {
	mock := &MockServiceConfig{ctrl: ctrl}
	mock.recorder = &MockServiceConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServiceConfig) EXPECT() *MockServiceConfigMockRecorder {
	return m.recorder
}

// GetAllServiceConfigs mocks base method.
func (m *MockServiceConfig) GetAllServiceConfigs() (string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllServiceConfigs")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetAllServiceConfigs indicates an expected call of GetAllServiceConfigs.
func (mr *MockServiceConfigMockRecorder) GetAllServiceConfigs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllServiceConfigs", reflect.TypeOf((*MockServiceConfig)(nil).GetAllServiceConfigs))
}

// SetAllServiceConfigs mocks base method.
func (m *MockServiceConfig) SetAllServiceConfigs(version, config string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetAllServiceConfigs", version, config)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetAllServiceConfigs indicates an expected call of SetAllServiceConfigs.
func (mr *MockServiceConfigMockRecorder) SetAllServiceConfigs(version, config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetAllServiceConfigs", reflect.TypeOf((*MockServiceConfig)(nil).SetAllServiceConfigs), version, config)
}

// MockSPDB is a mock of SPDB interface.
type MockSPDB struct {
	ctrl     *gomock.Controller
	recorder *MockSPDBMockRecorder
}

// MockSPDBMockRecorder is the mock recorder for MockSPDB.
type MockSPDBMockRecorder struct {
	mock *MockSPDB
}

// NewMockSPDB creates a new mock instance.
func NewMockSPDB(ctrl *gomock.Controller) *MockSPDB {
	mock := &MockSPDB{ctrl: ctrl}
	mock.recorder = &MockSPDBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSPDB) EXPECT() *MockSPDBMockRecorder {
	return m.recorder
}

// CheckQuotaAndAddReadRecord mocks base method.
func (m *MockSPDB) CheckQuotaAndAddReadRecord(record *ReadRecord, quota *BucketQuota) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckQuotaAndAddReadRecord", record, quota)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckQuotaAndAddReadRecord indicates an expected call of CheckQuotaAndAddReadRecord.
func (mr *MockSPDBMockRecorder) CheckQuotaAndAddReadRecord(record, quota interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckQuotaAndAddReadRecord", reflect.TypeOf((*MockSPDB)(nil).CheckQuotaAndAddReadRecord), record, quota)
}

// CreateUploadJob mocks base method.
func (m *MockSPDB) CreateUploadJob(objectInfo *types1.ObjectInfo) (*types.JobContext, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUploadJob", objectInfo)
	ret0, _ := ret[0].(*types.JobContext)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUploadJob indicates an expected call of CreateUploadJob.
func (mr *MockSPDBMockRecorder) CreateUploadJob(objectInfo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUploadJob", reflect.TypeOf((*MockSPDB)(nil).CreateUploadJob), objectInfo)
}

// FetchAllSp mocks base method.
func (m *MockSPDB) FetchAllSp(status ...types0.Status) ([]*types0.StorageProvider, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range status {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "FetchAllSp", varargs...)
	ret0, _ := ret[0].([]*types0.StorageProvider)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchAllSp indicates an expected call of FetchAllSp.
func (mr *MockSPDBMockRecorder) FetchAllSp(status ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchAllSp", reflect.TypeOf((*MockSPDB)(nil).FetchAllSp), status...)
}

// FetchAllSpWithoutOwnSp mocks base method.
func (m *MockSPDB) FetchAllSpWithoutOwnSp(status ...types0.Status) ([]*types0.StorageProvider, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range status {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "FetchAllSpWithoutOwnSp", varargs...)
	ret0, _ := ret[0].([]*types0.StorageProvider)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchAllSpWithoutOwnSp indicates an expected call of FetchAllSpWithoutOwnSp.
func (mr *MockSPDBMockRecorder) FetchAllSpWithoutOwnSp(status ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchAllSpWithoutOwnSp", reflect.TypeOf((*MockSPDB)(nil).FetchAllSpWithoutOwnSp), status...)
}

// GetBucketReadRecord mocks base method.
func (m *MockSPDB) GetBucketReadRecord(bucketID uint64, timeRange *TrafficTimeRange) ([]*ReadRecord, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBucketReadRecord", bucketID, timeRange)
	ret0, _ := ret[0].([]*ReadRecord)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBucketReadRecord indicates an expected call of GetBucketReadRecord.
func (mr *MockSPDBMockRecorder) GetBucketReadRecord(bucketID, timeRange interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBucketReadRecord", reflect.TypeOf((*MockSPDB)(nil).GetBucketReadRecord), bucketID, timeRange)
}

// GetBucketTraffic mocks base method.
func (m *MockSPDB) GetBucketTraffic(bucketID uint64, yearMonth string) (*BucketTraffic, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBucketTraffic", bucketID, yearMonth)
	ret0, _ := ret[0].(*BucketTraffic)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBucketTraffic indicates an expected call of GetBucketTraffic.
func (mr *MockSPDBMockRecorder) GetBucketTraffic(bucketID, yearMonth interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBucketTraffic", reflect.TypeOf((*MockSPDB)(nil).GetBucketTraffic), bucketID, yearMonth)
}

// GetJobByID mocks base method.
func (m *MockSPDB) GetJobByID(jobID uint64) (*types.JobContext, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetJobByID", jobID)
	ret0, _ := ret[0].(*types.JobContext)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetJobByID indicates an expected call of GetJobByID.
func (mr *MockSPDBMockRecorder) GetJobByID(jobID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetJobByID", reflect.TypeOf((*MockSPDB)(nil).GetJobByID), jobID)
}

// GetJobByObjectID mocks base method.
func (m *MockSPDB) GetJobByObjectID(objectID uint64) (*types.JobContext, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetJobByObjectID", objectID)
	ret0, _ := ret[0].(*types.JobContext)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetJobByObjectID indicates an expected call of GetJobByObjectID.
func (mr *MockSPDBMockRecorder) GetJobByObjectID(objectID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetJobByObjectID", reflect.TypeOf((*MockSPDB)(nil).GetJobByObjectID), objectID)
}

// GetObjectInfo mocks base method.
func (m *MockSPDB) GetObjectInfo(objectID uint64) (*types1.ObjectInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetObjectInfo", objectID)
	ret0, _ := ret[0].(*types1.ObjectInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetObjectInfo indicates an expected call of GetObjectInfo.
func (mr *MockSPDBMockRecorder) GetObjectInfo(objectID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetObjectInfo", reflect.TypeOf((*MockSPDB)(nil).GetObjectInfo), objectID)
}

// GetObjectIntegrity mocks base method.
func (m *MockSPDB) GetObjectIntegrity(objectID uint64) (*IntegrityMeta, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetObjectIntegrity", objectID)
	ret0, _ := ret[0].(*IntegrityMeta)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetObjectIntegrity indicates an expected call of GetObjectIntegrity.
func (mr *MockSPDBMockRecorder) GetObjectIntegrity(objectID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetObjectIntegrity", reflect.TypeOf((*MockSPDB)(nil).GetObjectIntegrity), objectID)
}

// GetObjectReadRecord mocks base method.
func (m *MockSPDB) GetObjectReadRecord(objectID uint64, timeRange *TrafficTimeRange) ([]*ReadRecord, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetObjectReadRecord", objectID, timeRange)
	ret0, _ := ret[0].([]*ReadRecord)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetObjectReadRecord indicates an expected call of GetObjectReadRecord.
func (mr *MockSPDBMockRecorder) GetObjectReadRecord(objectID, timeRange interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetObjectReadRecord", reflect.TypeOf((*MockSPDB)(nil).GetObjectReadRecord), objectID, timeRange)
}

// GetOwnSpInfo mocks base method.
func (m *MockSPDB) GetOwnSpInfo() (*types0.StorageProvider, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOwnSpInfo")
	ret0, _ := ret[0].(*types0.StorageProvider)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOwnSpInfo indicates an expected call of GetOwnSpInfo.
func (mr *MockSPDBMockRecorder) GetOwnSpInfo() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOwnSpInfo", reflect.TypeOf((*MockSPDB)(nil).GetOwnSpInfo))
}

// GetReadRecord mocks base method.
func (m *MockSPDB) GetReadRecord(timeRange *TrafficTimeRange) ([]*ReadRecord, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetReadRecord", timeRange)
	ret0, _ := ret[0].([]*ReadRecord)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetReadRecord indicates an expected call of GetReadRecord.
func (mr *MockSPDBMockRecorder) GetReadRecord(timeRange interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReadRecord", reflect.TypeOf((*MockSPDB)(nil).GetReadRecord), timeRange)
}

// GetSpByAddress mocks base method.
func (m *MockSPDB) GetSpByAddress(address string, addressType SpAddressType) (*types0.StorageProvider, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSpByAddress", address, addressType)
	ret0, _ := ret[0].(*types0.StorageProvider)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSpByAddress indicates an expected call of GetSpByAddress.
func (mr *MockSPDBMockRecorder) GetSpByAddress(address, addressType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSpByAddress", reflect.TypeOf((*MockSPDB)(nil).GetSpByAddress), address, addressType)
}

// GetSpByEndpoint mocks base method.
func (m *MockSPDB) GetSpByEndpoint(endpoint string) (*types0.StorageProvider, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSpByEndpoint", endpoint)
	ret0, _ := ret[0].(*types0.StorageProvider)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSpByEndpoint indicates an expected call of GetSpByEndpoint.
func (mr *MockSPDBMockRecorder) GetSpByEndpoint(endpoint interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSpByEndpoint", reflect.TypeOf((*MockSPDB)(nil).GetSpByEndpoint), endpoint)
}

// GetStorageParams mocks base method.
func (m *MockSPDB) GetStorageParams() (*types1.Params, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStorageParams")
	ret0, _ := ret[0].(*types1.Params)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStorageParams indicates an expected call of GetStorageParams.
func (mr *MockSPDBMockRecorder) GetStorageParams() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStorageParams", reflect.TypeOf((*MockSPDB)(nil).GetStorageParams))
}

// GetUserReadRecord mocks base method.
func (m *MockSPDB) GetUserReadRecord(userAddress string, timeRange *TrafficTimeRange) ([]*ReadRecord, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserReadRecord", userAddress, timeRange)
	ret0, _ := ret[0].([]*ReadRecord)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserReadRecord indicates an expected call of GetUserReadRecord.
func (mr *MockSPDBMockRecorder) GetUserReadRecord(userAddress, timeRange interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserReadRecord", reflect.TypeOf((*MockSPDB)(nil).GetUserReadRecord), userAddress, timeRange)
}

// SetObjectInfo mocks base method.
func (m *MockSPDB) SetObjectInfo(objectID uint64, objectInfo *types1.ObjectInfo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetObjectInfo", objectID, objectInfo)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetObjectInfo indicates an expected call of SetObjectInfo.
func (mr *MockSPDBMockRecorder) SetObjectInfo(objectID, objectInfo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetObjectInfo", reflect.TypeOf((*MockSPDB)(nil).SetObjectInfo), objectID, objectInfo)
}

// SetObjectIntegrity mocks base method.
func (m *MockSPDB) SetObjectIntegrity(integrity *IntegrityMeta) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetObjectIntegrity", integrity)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetObjectIntegrity indicates an expected call of SetObjectIntegrity.
func (mr *MockSPDBMockRecorder) SetObjectIntegrity(integrity interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetObjectIntegrity", reflect.TypeOf((*MockSPDB)(nil).SetObjectIntegrity), integrity)
}

// SetOwnSpInfo mocks base method.
func (m *MockSPDB) SetOwnSpInfo(sp *types0.StorageProvider) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetOwnSpInfo", sp)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetOwnSpInfo indicates an expected call of SetOwnSpInfo.
func (mr *MockSPDBMockRecorder) SetOwnSpInfo(sp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetOwnSpInfo", reflect.TypeOf((*MockSPDB)(nil).SetOwnSpInfo), sp)
}

// SetStorageParams mocks base method.
func (m *MockSPDB) SetStorageParams(params *types1.Params) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetStorageParams", params)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetStorageParams indicates an expected call of SetStorageParams.
func (mr *MockSPDBMockRecorder) SetStorageParams(params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetStorageParams", reflect.TypeOf((*MockSPDB)(nil).SetStorageParams), params)
}

// UpdateAllSp mocks base method.
func (m *MockSPDB) UpdateAllSp(spList []*types0.StorageProvider) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAllSp", spList)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateAllSp indicates an expected call of UpdateAllSp.
func (mr *MockSPDBMockRecorder) UpdateAllSp(spList interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAllSp", reflect.TypeOf((*MockSPDB)(nil).UpdateAllSp), spList)
}

// UpdateJobState mocks base method.
func (m *MockSPDB) UpdateJobState(objectID uint64, state types.JobState) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateJobState", objectID, state)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateJobState indicates an expected call of UpdateJobState.
func (mr *MockSPDBMockRecorder) UpdateJobState(objectID, state interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateJobState", reflect.TypeOf((*MockSPDB)(nil).UpdateJobState), objectID, state)
}
