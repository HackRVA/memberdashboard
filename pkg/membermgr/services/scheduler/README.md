# Scheduler
This is a place to schedule things to run.
Tasks are scheduled in go routines.  However, the two mentioned below execute at startup as well.

At the moment we schedule a payment refresh from paypal everyday.
This is controled with the `checkPaymentsInterval`

We also evaluate member status everyday with the `evaluateMemberStatusInterval`

