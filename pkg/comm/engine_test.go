package comm

import (
	"github.com/hegemone/kore/pkg/msg"
	"testing"
)

type TestAdapter struct{}

func (a *TestAdapter) Name() string {
	return "test"
}

func (a *TestAdapter) Listen(inChan chan<- msg.RawIngress) {
	inChan <- msg.RawIngress{}
}

func (a *TestAdapter) SendMessage(e msg.Egress) {}

func TestStart(t *testing.T) {
	bufferSize := 8
	e := &Engine{
		rawIngressBuffer: make(chan msg.RawIngressBuffer, bufferSize),
		ingressBuffer:    make(chan msg.IngressBuffer, bufferSize),
		egressBuffer:     make(chan msg.EgressBuffer, bufferSize),
		plugins:          make(map[string]*Plugin),
		adapters:         make(map[string]Adapter),
	}

	e.adapters["test"] = &TestAdapter{}

	sleep := make(chan bool)
	funcDone = func() {
		sleep <- true
	}

	HRICalled, HICalled, HECalled := false, false, false
	eHandleRawIngress = func(*Engine, msg.RawIngressBuffer) {
		HRICalled = true
	}
	eHandleIngress = func(*Engine, msg.IngressBuffer) {
		HICalled = true
	}
	eHandleEgress = func(*Engine, msg.EgressBuffer) {
		HECalled = true
	}

	go e.Start()

	m := <-e.rawIngressBuffer
	t.Logf("Received message %v from ingress buffer", m)
	<-sleep

	e.rawIngressBuffer <- msg.RawIngressBuffer{}
	<-sleep
	e.ingressBuffer <- msg.IngressBuffer{}
	<-sleep
	e.egressBuffer <- msg.EgressBuffer{}
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
