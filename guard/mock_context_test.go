package guard

import (
	"context"
	"testing"
	"time"

	"gopkg.in/workanator/go-floc.v2"
)

func TestMockContext_Done(t *testing.T) {
	oCtx := floc.NewContext()
	oCancelCtx, oCancel := context.WithCancel(oCtx.Ctx())
	oCtx.UpdateCtx(oCancelCtx)

	mCtx := floc.NewContext()
	mCancelCtx, _ := context.WithCancel(mCtx.Ctx())
	mCtx.UpdateCtx(mCancelCtx)

	mock := mockContext{Context: oCtx, mock: mCtx}

	go func() {
		time.Sleep(time.Millisecond)
		oCancel()
	}()

	select {
	case <-oCtx.Done():
		// Ok
	case <-mock.Done():
		// Not Ok
		t.Fatalf("%s expects original context to be canceled", t.Name())
	}
}

func TestMockContext_Done2(t *testing.T) {
	oCtx := floc.NewContext()
	oCancelCtx, _ := context.WithCancel(oCtx.Ctx())
	oCtx.UpdateCtx(oCancelCtx)

	mCtx := floc.NewContext()
	mCancelCtx, mCancel := context.WithCancel(mCtx.Ctx())
	mCtx.UpdateCtx(mCancelCtx)

	mock := mockContext{Context: oCtx, mock: mCtx}

	go func() {
		time.Sleep(time.Millisecond)
		mCancel()
	}()

	select {
	case <-oCtx.Done():
		// Not Ok
		t.Fatalf("%s expects mock context to be canceled", t.Name())
	case <-mock.Done():
		// Ok
	}
}

func TestMockContext_Done3(t *testing.T) {
	oCtx := floc.NewContext()
	oCancelCtx, oCancel := context.WithCancel(oCtx.Ctx())
	oCtx.UpdateCtx(oCancelCtx)

	mCtx := floc.NewContext()
	mCancelCtx, mCancel := context.WithCancel(mCtx.Ctx())
	mCtx.UpdateCtx(mCancelCtx)

	mock := mockContext{Context: oCtx, mock: mCtx}

	go func() {
		time.Sleep(time.Millisecond)
		oCancel()
	}()

	timer := time.NewTimer(5 * time.Millisecond)
	select {
	case <-oCtx.Done():
		timer.Stop()
	case <-timer.C:
		t.Fatalf("%s expects original context to be canceled", t.Name())
	}

	go func() {
		time.Sleep(time.Millisecond)
		mCancel()
	}()

	timer = time.NewTimer(5 * time.Millisecond)
	select {
	case <-mock.Done():
		timer.Stop()
	case <-timer.C:
		t.Fatalf("%s expects mock context to be canceled", t.Name())
	}
}