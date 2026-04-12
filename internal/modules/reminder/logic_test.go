package reminder

import (
	"testing"
	"time"

	"github.com/aalexanderkevin/getstarvio-backend/internal/models"
)

func TestDeductCreditsPriority(t *testing.T) {
	nw, ns, nt := deductCredits(5, 7, 9)
	if nw != 4 || ns != 7 || nt != 9 {
		t.Fatalf("welcome priority failed: got %d,%d,%d", nw, ns, nt)
	}

	nw, ns, nt = deductCredits(0, 7, 9)
	if nw != 0 || ns != 6 || nt != 9 {
		t.Fatalf("subscription priority failed: got %d,%d,%d", nw, ns, nt)
	}

	nw, ns, nt = deductCredits(0, 0, 9)
	if nw != 0 || ns != 0 || nt != 8 {
		t.Fatalf("topup priority failed: got %d,%d,%d", nw, ns, nt)
	}
}

func TestCanSendRules(t *testing.T) {
	now := time.Now().UTC()

	ok, _ := canSend(models.Wallet{
		TrialEndsAt:        now.Add(24 * time.Hour),
		SubscriptionStatus: models.SubscriptionStatusNone,
		WelcomeCreditsLeft: 1,
	}, now)
	if !ok {
		t.Fatalf("expected trial user with credits to be allowed")
	}

	ok, reason := canSend(models.Wallet{
		TrialEndsAt:        now.Add(-24 * time.Hour),
		SubscriptionStatus: models.SubscriptionStatusNone,
		WelcomeCreditsLeft: 10,
	}, now)
	if ok || reason == "" {
		t.Fatalf("expected expired trial without subscription to be blocked")
	}

	ok, _ = canSend(models.Wallet{
		TrialEndsAt:        now.Add(-24 * time.Hour),
		SubscriptionStatus: models.SubscriptionStatusActive,
		SubCreditsLeft:     1,
	}, now)
	if !ok {
		t.Fatalf("expected active subscriber to be allowed")
	}
}
