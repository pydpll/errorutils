package errorutils

import (
	"bufio"
	"errors"
	"io"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestCompareOptions(t *testing.T) {
	//set logger level to debug
	logrus.SetLevel(logrus.DebugLevel)
	tests := []struct {
		target   Option
		template Option
		expected bool
		name     string
	}{
		{WithMsg("this is the target, it has a message"), WithMsg("test"), true, "same"},
		{WithMsg("this is the target, it has a message"), WithExitCode(1), false, "different"},
		{WithMsg("this is the target, it has a message"), WithAltPrint("test"), false, "different"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := compareOptions(tc.target, tc.template)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func BenchmarkConsole(b *testing.B) {
	tests := mockMessages()
	//logrus: set logger to /dev/null
	logrus.SetOutput(io.Discard)
	//data colection file append create
	f, err := os.OpenFile("timing.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	//bufio from f
	w := bufio.NewWriter(f)
	if err != nil {
		panic("failed to create file for data: " + err.Error())
	}

	defer f.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, test := range tests {
			start := time.Now()
			test.fn()
			w.WriteString(time.Since(start).String() + "," + test.name + "\n")
		}
	}
}

func mockMessages() []struct {
	name string
	fn   func()
} {
	return []struct {
		name string
		fn   func()
	}{
		{"LogFailures", func() {
			err := errors.New("log error: bacterial culture contamination")
			LogFailures(err, WithInner(errors.New("oops, wrong Petri dish")))
		}},
		{"Info", func() { logrus.Info("The mitochondria is the powerhouse of the cell.") }},
		{"WarnOnFail", func() {
			err := errors.New("warning: DNA sequence mismatch detected")
			WarnOnFail(err, WithMsg("check for mutations"))
		}},
		{"Error", func() { logrus.Error("Error: PCR machine is out of Taq polymerase!") }},
		{"Trace", func() { logrus.Trace("Tracing lineage: following the genetic footsteps of your ancestors...") }},
		{"WarnOnFailWithMsg", func() {
			err := errors.New("warning: gene editing unsuccessful")
			WarnOnFail(err, WithMsg("CRISPR needs a coffee break"))
		}},
		{"Debug", func() { logrus.Debug("Debug: RNA transcription in progress...") }},
		{"Info", func() { logrus.Info("Panic: lab rats have escaped the facility!") }},
		{"Error", func() { logrus.Error("Fatal: autoclave malfunction, sterilization incomplete!") }},
		{"LogFailuresWithAltPrint", func() {
			err := errors.New("log error: unexpected results in gel electrophoresis")
			LogFailures(err, WithAltPrint("was that supposed to be a smiley face?"))
		}},
		{"WarnWithLineRef", func() {
			err := errors.New("warning: plasmid vector not found")
			WarnOnFail(err, WithLineRef("GATTACA"))
		}},
		{"Info", func() { logrus.Info("Info: gene splicing complete, no superheroes detected.") }},
		{"WarnOnFail", func() {
			err := errors.New("warning: cell culture overgrowth")
			WarnOnFail(err, WithMsg("time to divide and conquer"))
		}},
		{"Error", func() { logrus.Error("Error: microscope lens needs cleaning, again.") }},
		{"LogFailures", func() {
			err := errors.New("log error: enzyme reaction failed")
			LogFailures(err, WithInner(errors.New("maybe it's enzyme envy")))
		}},
		{"Trace", func() { logrus.Trace("Trace: following the path of a rogue ribosome...") }},
		{"WarnOnFailWithMsg", func() {
			err := errors.New("warning: centrifuge imbalance")
			WarnOnFail(err, WithMsg("redistribute your samples, please"))
		}},
		{"Debug", func() { logrus.Debug("Debug: gene expression levels are off the charts!") }},
		{"Info", func() { logrus.Info("Panic: the lab cat is on the loose!") }},
		{"Error", func() { logrus.Error("Fatal: CRISPR cut too deep, genome in disarray!") }},
		{"LogFailuresWithAltPrint", func() {
			err := errors.New("log error: unexpected mutation")
			LogFailures(err, WithAltPrint("wasn't expecting the X-Men"))
		}},
		{"WarnWithLineRef", func() {
			err := errors.New("warning: biohazard spill")
			WarnOnFail(err, WithLineRef("BIO123"))
		}},
		{"Info", func() { logrus.Info("Info: genome sequencing completed, no alien DNA found.") }},
		{"WarnOnFail", func() {
			err := errors.New("warning: PCR contamination")
			WarnOnFail(err, WithMsg("someone didn't change their gloves"))
		}},
		{"Error", func() { logrus.Error("Error: failed to isolate the plasmid.") }},
		{"LogFailures", func() {
			err := errors.New("log error: cloning experiment went wrong")
			LogFailures(err, WithInner(errors.New("now we have two problems")))
		}},
		{"Trace", func() { logrus.Trace("Trace: tracking a rogue CRISPR edit...") }},
		{"WarnOnFailWithMsg", func() {
			err := errors.New("warning: gel electrophoresis running backward")
			WarnOnFail(err, WithMsg("check the connections"))
		}},
		{"Debug", func() { logrus.Debug("Debug: analyzing gene expression data...") }},
		{"Info", func() { logrus.Info("Panic: the lab fridge is empty!") }},
		{"Error", func() { logrus.Error("Fatal: lost connection to the sequencing server!") }},
		{"LogFailuresWithAltPrint", func() {
			err := errors.New("log error: protein folding simulation crashed")
			LogFailures(err, WithAltPrint("try turning it off and on again"))
		}},
		{"WarnWithLineRef", func() {
			err := errors.New("warning: CRISPR edit reverted")
			WarnOnFail(err, WithLineRef("AGCTTAGC"))
		}},
		{"Info", func() { logrus.Info("Info: all samples successfully labeled. No mix-ups today!") }},
		{"WarnOnFail", func() {
			err := errors.New("warning: contamination in the sample")
			WarnOnFail(err, WithMsg("who didn't sterilize their pipette?"))
		}},
		{"Error", func() { logrus.Error("Error: sequencing run failed midway.") }},
		{"LogFailures", func() {
			err := errors.New("log error: unexpected phenotype observed")
			LogFailures(err, WithInner(errors.New("that's not supposed to glow")))
		}},
		{"Trace", func() { logrus.Trace("Trace: mapping the human genome...") }},
		{"WarnOnFailWithMsg", func() {
			err := errors.New("warning: buffer solution expired")
			WarnOnFail(err, WithMsg("time for a new batch"))
		}},
		{"Debug", func() { logrus.Debug("Debug: analyzing CRISPR off-target effects...") }},
		{"Info", func() { logrus.Info("Panic: the incubator is overheating!") }},
		{"Error", func() { logrus.Error("Fatal: lost all data from the last experiment!") }},
		{"LogFailuresWithAltPrint", func() {
			err := errors.New("log error: protein assay failure")
			LogFailures(err, WithAltPrint("maybe it needs a pep talk"))
		}},
		{"WarnWithLineRef", func() {
			err := errors.New("warning: plasmid vector lost")
			WarnOnFail(err, WithLineRef("VECTORX"))
		}},
		{"Info", func() { logrus.Info("Info: all PCR reactions completed successfully.") }},
		{"WarnOnFail", func() {
			err := errors.New("warning: gene knockout experiment failed")
			WarnOnFail(err, WithMsg("try a different approach"))
		}},
		{"Error", func() { logrus.Error("Error: contamination detected in the cell culture.") }},
		{"LogFailures", func() {
			err := errors.New("log error: unexpected gene expression pattern")
			LogFailures(err, WithInner(errors.New("time for another round of qPCR")))
		}},
		{"Trace", func() { logrus.Trace("Trace: tracking the evolution of antibiotic resistance...") }},
		{"WarnOnFailWithMsg", func() {
			err := errors.New("warning: reagent shortage")
			WarnOnFail(err, WithMsg("order more supplies"))
		}},
		{"Debug", func() { logrus.Debug("Debug: analyzing next-gen sequencing data...") }},
		{"Info", func() { logrus.Info("Panic: the lab freezer is defrosting!") }},
		{"Error", func() { logrus.Error("Fatal: lab notebook lost!") }},
	}
}
