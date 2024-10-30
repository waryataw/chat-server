// Code generated by http://github.com/gojuno/minimock (v3.4.1). DO NOT EDIT.

package mocks

//go:generate minimock -i github.com/waryataw/chat-server/internal/service/chat.AuthRepository -o mocks.go -n AuthRepositoryMock -p mocks

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
	"github.com/waryataw/chat-server/internal/models"
)

// AuthRepositoryMock implements mm_chat.AuthRepository
type AuthRepositoryMock struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcGetUser          func(ctx context.Context, name string) (up1 *models.User, err error)
	funcGetUserOrigin    string
	inspectFuncGetUser   func(ctx context.Context, name string)
	afterGetUserCounter  uint64
	beforeGetUserCounter uint64
	GetUserMock          mAuthRepositoryMockGetUser
}

// NewAuthRepositoryMock returns a mock for mm_chat.AuthRepository
func NewAuthRepositoryMock(t minimock.Tester) *AuthRepositoryMock {
	m := &AuthRepositoryMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetUserMock = mAuthRepositoryMockGetUser{mock: m}
	m.GetUserMock.callArgs = []*AuthRepositoryMockGetUserParams{}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mAuthRepositoryMockGetUser struct {
	optional           bool
	mock               *AuthRepositoryMock
	defaultExpectation *AuthRepositoryMockGetUserExpectation
	expectations       []*AuthRepositoryMockGetUserExpectation

	callArgs []*AuthRepositoryMockGetUserParams
	mutex    sync.RWMutex

	expectedInvocations       uint64
	expectedInvocationsOrigin string
}

// AuthRepositoryMockGetUserExpectation specifies expectation struct of the AuthRepository.GetUser
type AuthRepositoryMockGetUserExpectation struct {
	mock               *AuthRepositoryMock
	params             *AuthRepositoryMockGetUserParams
	paramPtrs          *AuthRepositoryMockGetUserParamPtrs
	expectationOrigins AuthRepositoryMockGetUserExpectationOrigins
	results            *AuthRepositoryMockGetUserResults
	returnOrigin       string
	Counter            uint64
}

// AuthRepositoryMockGetUserParams contains parameters of the AuthRepository.GetUser
type AuthRepositoryMockGetUserParams struct {
	ctx  context.Context
	name string
}

// AuthRepositoryMockGetUserParamPtrs contains pointers to parameters of the AuthRepository.GetUser
type AuthRepositoryMockGetUserParamPtrs struct {
	ctx  *context.Context
	name *string
}

// AuthRepositoryMockGetUserResults contains results of the AuthRepository.GetUser
type AuthRepositoryMockGetUserResults struct {
	up1 *models.User
	err error
}

// AuthRepositoryMockGetUserOrigins contains origins of expectations of the AuthRepository.GetUser
type AuthRepositoryMockGetUserExpectationOrigins struct {
	origin     string
	originCtx  string
	originName string
}

// Marks this method to be optional. The default behavior of any method with Return() is '1 or more', meaning
// the test will fail minimock's automatic final call check if the mocked method was not called at least once.
// Optional() makes method check to work in '0 or more' mode.
// It is NOT RECOMMENDED to use this option unless you really need it, as default behaviour helps to
// catch the problems when the expected method call is totally skipped during test run.
func (mmGetUser *mAuthRepositoryMockGetUser) Optional() *mAuthRepositoryMockGetUser {
	mmGetUser.optional = true
	return mmGetUser
}

// Expect sets up expected params for AuthRepository.GetUser
func (mmGetUser *mAuthRepositoryMockGetUser) Expect(ctx context.Context, name string) *mAuthRepositoryMockGetUser {
	if mmGetUser.mock.funcGetUser != nil {
		mmGetUser.mock.t.Fatalf("AuthRepositoryMock.GetUser mock is already set by Set")
	}

	if mmGetUser.defaultExpectation == nil {
		mmGetUser.defaultExpectation = &AuthRepositoryMockGetUserExpectation{}
	}

	if mmGetUser.defaultExpectation.paramPtrs != nil {
		mmGetUser.mock.t.Fatalf("AuthRepositoryMock.GetUser mock is already set by ExpectParams functions")
	}

	mmGetUser.defaultExpectation.params = &AuthRepositoryMockGetUserParams{ctx, name}
	mmGetUser.defaultExpectation.expectationOrigins.origin = minimock.CallerInfo(1)
	for _, e := range mmGetUser.expectations {
		if minimock.Equal(e.params, mmGetUser.defaultExpectation.params) {
			mmGetUser.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmGetUser.defaultExpectation.params)
		}
	}

	return mmGetUser
}

// ExpectCtxParam1 sets up expected param ctx for AuthRepository.GetUser
func (mmGetUser *mAuthRepositoryMockGetUser) ExpectCtxParam1(ctx context.Context) *mAuthRepositoryMockGetUser {
	if mmGetUser.mock.funcGetUser != nil {
		mmGetUser.mock.t.Fatalf("AuthRepositoryMock.GetUser mock is already set by Set")
	}

	if mmGetUser.defaultExpectation == nil {
		mmGetUser.defaultExpectation = &AuthRepositoryMockGetUserExpectation{}
	}

	if mmGetUser.defaultExpectation.params != nil {
		mmGetUser.mock.t.Fatalf("AuthRepositoryMock.GetUser mock is already set by Expect")
	}

	if mmGetUser.defaultExpectation.paramPtrs == nil {
		mmGetUser.defaultExpectation.paramPtrs = &AuthRepositoryMockGetUserParamPtrs{}
	}
	mmGetUser.defaultExpectation.paramPtrs.ctx = &ctx
	mmGetUser.defaultExpectation.expectationOrigins.originCtx = minimock.CallerInfo(1)

	return mmGetUser
}

// ExpectNameParam2 sets up expected param name for AuthRepository.GetUser
func (mmGetUser *mAuthRepositoryMockGetUser) ExpectNameParam2(name string) *mAuthRepositoryMockGetUser {
	if mmGetUser.mock.funcGetUser != nil {
		mmGetUser.mock.t.Fatalf("AuthRepositoryMock.GetUser mock is already set by Set")
	}

	if mmGetUser.defaultExpectation == nil {
		mmGetUser.defaultExpectation = &AuthRepositoryMockGetUserExpectation{}
	}

	if mmGetUser.defaultExpectation.params != nil {
		mmGetUser.mock.t.Fatalf("AuthRepositoryMock.GetUser mock is already set by Expect")
	}

	if mmGetUser.defaultExpectation.paramPtrs == nil {
		mmGetUser.defaultExpectation.paramPtrs = &AuthRepositoryMockGetUserParamPtrs{}
	}
	mmGetUser.defaultExpectation.paramPtrs.name = &name
	mmGetUser.defaultExpectation.expectationOrigins.originName = minimock.CallerInfo(1)

	return mmGetUser
}

// Inspect accepts an inspector function that has same arguments as the AuthRepository.GetUser
func (mmGetUser *mAuthRepositoryMockGetUser) Inspect(f func(ctx context.Context, name string)) *mAuthRepositoryMockGetUser {
	if mmGetUser.mock.inspectFuncGetUser != nil {
		mmGetUser.mock.t.Fatalf("Inspect function is already set for AuthRepositoryMock.GetUser")
	}

	mmGetUser.mock.inspectFuncGetUser = f

	return mmGetUser
}

// Return sets up results that will be returned by AuthRepository.GetUser
func (mmGetUser *mAuthRepositoryMockGetUser) Return(up1 *models.User, err error) *AuthRepositoryMock {
	if mmGetUser.mock.funcGetUser != nil {
		mmGetUser.mock.t.Fatalf("AuthRepositoryMock.GetUser mock is already set by Set")
	}

	if mmGetUser.defaultExpectation == nil {
		mmGetUser.defaultExpectation = &AuthRepositoryMockGetUserExpectation{mock: mmGetUser.mock}
	}
	mmGetUser.defaultExpectation.results = &AuthRepositoryMockGetUserResults{up1, err}
	mmGetUser.defaultExpectation.returnOrigin = minimock.CallerInfo(1)
	return mmGetUser.mock
}

// Set uses given function f to mock the AuthRepository.GetUser method
func (mmGetUser *mAuthRepositoryMockGetUser) Set(f func(ctx context.Context, name string) (up1 *models.User, err error)) *AuthRepositoryMock {
	if mmGetUser.defaultExpectation != nil {
		mmGetUser.mock.t.Fatalf("Default expectation is already set for the AuthRepository.GetUser method")
	}

	if len(mmGetUser.expectations) > 0 {
		mmGetUser.mock.t.Fatalf("Some expectations are already set for the AuthRepository.GetUser method")
	}

	mmGetUser.mock.funcGetUser = f
	mmGetUser.mock.funcGetUserOrigin = minimock.CallerInfo(1)
	return mmGetUser.mock
}

// When sets expectation for the AuthRepository.GetUser which will trigger the result defined by the following
// Then helper
func (mmGetUser *mAuthRepositoryMockGetUser) When(ctx context.Context, name string) *AuthRepositoryMockGetUserExpectation {
	if mmGetUser.mock.funcGetUser != nil {
		mmGetUser.mock.t.Fatalf("AuthRepositoryMock.GetUser mock is already set by Set")
	}

	expectation := &AuthRepositoryMockGetUserExpectation{
		mock:               mmGetUser.mock,
		params:             &AuthRepositoryMockGetUserParams{ctx, name},
		expectationOrigins: AuthRepositoryMockGetUserExpectationOrigins{origin: minimock.CallerInfo(1)},
	}
	mmGetUser.expectations = append(mmGetUser.expectations, expectation)
	return expectation
}

// Then sets up AuthRepository.GetUser return parameters for the expectation previously defined by the When method
func (e *AuthRepositoryMockGetUserExpectation) Then(up1 *models.User, err error) *AuthRepositoryMock {
	e.results = &AuthRepositoryMockGetUserResults{up1, err}
	return e.mock
}

// Times sets number of times AuthRepository.GetUser should be invoked
func (mmGetUser *mAuthRepositoryMockGetUser) Times(n uint64) *mAuthRepositoryMockGetUser {
	if n == 0 {
		mmGetUser.mock.t.Fatalf("Times of AuthRepositoryMock.GetUser mock can not be zero")
	}
	mm_atomic.StoreUint64(&mmGetUser.expectedInvocations, n)
	mmGetUser.expectedInvocationsOrigin = minimock.CallerInfo(1)
	return mmGetUser
}

func (mmGetUser *mAuthRepositoryMockGetUser) invocationsDone() bool {
	if len(mmGetUser.expectations) == 0 && mmGetUser.defaultExpectation == nil && mmGetUser.mock.funcGetUser == nil {
		return true
	}

	totalInvocations := mm_atomic.LoadUint64(&mmGetUser.mock.afterGetUserCounter)
	expectedInvocations := mm_atomic.LoadUint64(&mmGetUser.expectedInvocations)

	return totalInvocations > 0 && (expectedInvocations == 0 || expectedInvocations == totalInvocations)
}

// GetUser implements mm_chat.AuthRepository
func (mmGetUser *AuthRepositoryMock) GetUser(ctx context.Context, name string) (up1 *models.User, err error) {
	mm_atomic.AddUint64(&mmGetUser.beforeGetUserCounter, 1)
	defer mm_atomic.AddUint64(&mmGetUser.afterGetUserCounter, 1)

	mmGetUser.t.Helper()

	if mmGetUser.inspectFuncGetUser != nil {
		mmGetUser.inspectFuncGetUser(ctx, name)
	}

	mm_params := AuthRepositoryMockGetUserParams{ctx, name}

	// Record call args
	mmGetUser.GetUserMock.mutex.Lock()
	mmGetUser.GetUserMock.callArgs = append(mmGetUser.GetUserMock.callArgs, &mm_params)
	mmGetUser.GetUserMock.mutex.Unlock()

	for _, e := range mmGetUser.GetUserMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.up1, e.results.err
		}
	}

	if mmGetUser.GetUserMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmGetUser.GetUserMock.defaultExpectation.Counter, 1)
		mm_want := mmGetUser.GetUserMock.defaultExpectation.params
		mm_want_ptrs := mmGetUser.GetUserMock.defaultExpectation.paramPtrs

		mm_got := AuthRepositoryMockGetUserParams{ctx, name}

		if mm_want_ptrs != nil {

			if mm_want_ptrs.ctx != nil && !minimock.Equal(*mm_want_ptrs.ctx, mm_got.ctx) {
				mmGetUser.t.Errorf("AuthRepositoryMock.GetUser got unexpected parameter ctx, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
					mmGetUser.GetUserMock.defaultExpectation.expectationOrigins.originCtx, *mm_want_ptrs.ctx, mm_got.ctx, minimock.Diff(*mm_want_ptrs.ctx, mm_got.ctx))
			}

			if mm_want_ptrs.name != nil && !minimock.Equal(*mm_want_ptrs.name, mm_got.name) {
				mmGetUser.t.Errorf("AuthRepositoryMock.GetUser got unexpected parameter name, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
					mmGetUser.GetUserMock.defaultExpectation.expectationOrigins.originName, *mm_want_ptrs.name, mm_got.name, minimock.Diff(*mm_want_ptrs.name, mm_got.name))
			}

		} else if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmGetUser.t.Errorf("AuthRepositoryMock.GetUser got unexpected parameters, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
				mmGetUser.GetUserMock.defaultExpectation.expectationOrigins.origin, *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmGetUser.GetUserMock.defaultExpectation.results
		if mm_results == nil {
			mmGetUser.t.Fatal("No results are set for the AuthRepositoryMock.GetUser")
		}
		return (*mm_results).up1, (*mm_results).err
	}
	if mmGetUser.funcGetUser != nil {
		return mmGetUser.funcGetUser(ctx, name)
	}
	mmGetUser.t.Fatalf("Unexpected call to AuthRepositoryMock.GetUser. %v %v", ctx, name)
	return
}

// GetUserAfterCounter returns a count of finished AuthRepositoryMock.GetUser invocations
func (mmGetUser *AuthRepositoryMock) GetUserAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetUser.afterGetUserCounter)
}

// GetUserBeforeCounter returns a count of AuthRepositoryMock.GetUser invocations
func (mmGetUser *AuthRepositoryMock) GetUserBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetUser.beforeGetUserCounter)
}

// Calls returns a list of arguments used in each call to AuthRepositoryMock.GetUser.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmGetUser *mAuthRepositoryMockGetUser) Calls() []*AuthRepositoryMockGetUserParams {
	mmGetUser.mutex.RLock()

	argCopy := make([]*AuthRepositoryMockGetUserParams, len(mmGetUser.callArgs))
	copy(argCopy, mmGetUser.callArgs)

	mmGetUser.mutex.RUnlock()

	return argCopy
}

// MinimockGetUserDone returns true if the count of the GetUser invocations corresponds
// the number of defined expectations
func (m *AuthRepositoryMock) MinimockGetUserDone() bool {
	if m.GetUserMock.optional {
		// Optional methods provide '0 or more' call count restriction.
		return true
	}

	for _, e := range m.GetUserMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	return m.GetUserMock.invocationsDone()
}

// MinimockGetUserInspect logs each unmet expectation
func (m *AuthRepositoryMock) MinimockGetUserInspect() {
	for _, e := range m.GetUserMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to AuthRepositoryMock.GetUser at\n%s with params: %#v", e.expectationOrigins.origin, *e.params)
		}
	}

	afterGetUserCounter := mm_atomic.LoadUint64(&m.afterGetUserCounter)
	// if default expectation was set then invocations count should be greater than zero
	if m.GetUserMock.defaultExpectation != nil && afterGetUserCounter < 1 {
		if m.GetUserMock.defaultExpectation.params == nil {
			m.t.Errorf("Expected call to AuthRepositoryMock.GetUser at\n%s", m.GetUserMock.defaultExpectation.returnOrigin)
		} else {
			m.t.Errorf("Expected call to AuthRepositoryMock.GetUser at\n%s with params: %#v", m.GetUserMock.defaultExpectation.expectationOrigins.origin, *m.GetUserMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetUser != nil && afterGetUserCounter < 1 {
		m.t.Errorf("Expected call to AuthRepositoryMock.GetUser at\n%s", m.funcGetUserOrigin)
	}

	if !m.GetUserMock.invocationsDone() && afterGetUserCounter > 0 {
		m.t.Errorf("Expected %d calls to AuthRepositoryMock.GetUser at\n%s but found %d calls",
			mm_atomic.LoadUint64(&m.GetUserMock.expectedInvocations), m.GetUserMock.expectedInvocationsOrigin, afterGetUserCounter)
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *AuthRepositoryMock) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockGetUserInspect()
		}
	})
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *AuthRepositoryMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *AuthRepositoryMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockGetUserDone()
}
