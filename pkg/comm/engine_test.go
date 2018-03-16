package comm

import (
	"testing"
)

type TestAdapter struct{}

func (a *TestAdapter) Name() string {
	return "test"
}

func (a *TestAdapter) Listen(inChan chan<- RawIngressMessage) {
	inChan <- RawIngressMessage{}
}

func (a *TestAdapter) SendMessage(e EgressMessage) {}

func TestStart(t *testing.T) {
	bufferSize := 8
	e := &Engine{
		rawIngressBuffer: make(chan rawIngressBufferMsg, bufferSize),
		ingressBuffer:    make(chan ingressBufferMsg, bufferSize),
		egressBuffer:     make(chan egressBufferMsg, bufferSize),
		plugins:          make(map[string]*Plugin),
		adapters:         make(map[string]Adapter),
	}

	e.adapters["test"] = &TestAdapter{}

	sleep := make(chan bool)
	funcDone = func() {
		sleep <- true
	}

	HRICalled, HICalled, HECalled := false, false, false
	eHandleRawIngress = func(*Engine, rawIngressBufferMsg) {
		HRICalled = true
	}
	eHandleIngress = func(*Engine, ingressBufferMsg) {
		HICalled = true
	}
	eHandleEgress = func(*Engine, egressBufferMsg) {
		HECalled = true
	}

	go e.Start()

	m := <-e.rawIngressBuffer
	t.Logf("Received message %v from ingress buffer", m)
	<-sleep

	e.rawIngressBuffer <- rawIngressBufferMsg{}
	<-sleep
	e.ingressBuffer <- ingressBufferMsg{}
	<-sleep
	e.egressBuffer <- egressBufferMsg{}
	<-sleep

	if !HRICalled {
		t.Errorf("handleRawIngress was not called")
	}
	if !HICalled {
		t.Errorf("handleIngress was not called")
	}
	if !HECalled {
		t.Errorf("handleEgress was not called")
	}
	funcDone = func() {}
}

func TestParseRawContent(t *testing.T) {
	tests := []struct {
		raw    string
		parsed string
	}{
		{"!bacon Hello world", "bacon Hello world"},
		{"!die What's up, Doc?", "die What's up, Doc?"},
	}

	for _, test := range tests {
		if parsed := parseRawContent(test.raw); parsed != test.parsed {
			t.Errorf("parseRawContent(%s) = %s; wanted %s", test.raw, parsed, test.parsed)
		}
	}
}
