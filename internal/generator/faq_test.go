package generator

import (
	"fmt"
	"testing"
)

func TestBePredicate(t *testing.T) {
	f := Fact{Subject: "AI", Predicate: "是", Object: "生产力"}
	q := templateQuestion(f)
	if q != "AI 是什么？" {
		t.Errorf("expected 'AI 是什么？', got %s", q)
	}
	a := templateAnswer(f)
	if a != "AI 是 生产力。" {
		t.Errorf("expected 'AI 是 生产力。', got %s", a)
	}
}

func TestNormalPredicate(t *testing.T) {
	f := Fact{Subject: "订单", Predicate: "包含", Object: "商品"}
	q := templateQuestion(f)
	if q != "订单 的 包含 是什么？" {
		t.Errorf("expected '订单 的 包含 是什么？', got %s", q)
	}
	a := templateAnswer(f)
	if !contains(a, "订单") || !contains(a, "商品") {
		t.Errorf("answer should contain subject and object: %s", a)
	}
}

func TestFAQNoDuplicates(t *testing.T) {
	facts := []Fact{
		{Subject: "A", Predicate: "B", Object: "C"},
		{Subject: "A", Predicate: "B", Object: "C"},
		{Subject: "X", Predicate: "Y", Object: "Z"},
	}
	g := &FAQGenerator{}
	result, err := g.Generate(facts, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Count "Q1:" and "Q2:" headers
	qCount := 0
	for i := 1; i <= 10; i++ {
		header := fmt.Sprintf("Q%d:", i)
		if contains(result.Content, header) {
			qCount++
		}
	}
	if qCount != 2 {
		t.Errorf("expected 2 unique Qs, got %d", qCount)
	}
}
