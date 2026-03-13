# Working Days Chart & Calendar Analysis

## Project Timeline Overview

### Milestone Dates
- **M1 - Beta Launch:** March 1, 2026 (National Insurance Day)
- **M2 - Grand Launch:** April 14, 2026 (Pohela Boishakh)
- **M3 - Complete Platform:** August 1, 2026

---

## Working Days Calculation

### M1 Period: December 20, 2025 - March 1, 2026

**Calendar Breakdown:**
- Start Date: December 20, 2025
- End Date: March 1, 2026
- Total Calendar Days: 72 days

**Holidays in M1 Period:**
| Date | Day | Holiday Name | Status |
|------|-----|--------------|--------|
| Dec 25, 2025 | Thursday | Christmas Day | Confirmed |
| Feb 4, 2026 | Wednesday | Shab-e-Barat | Moon Dependent |
| Feb 21, 2026 | Saturday | Intl. Mother Language Day | Confirmed |

**Total Holidays:** 3 days

**Weekends in M1 Period:**
- Working Days: Saturday to Thursday (6 days/week)
- Friday is weekend
- Approximate Fridays: ~10 days

**Working Days Calculation:**
```
Calendar Days:        72 days
Minus Holidays:       -3 days
Minus Weekends:       -10 days
= Net Working Days:   59 days
```

**Effective Working Days for M1:** 59 days

---

### M2 Period: March 2, 2026 - April 14, 2026

**Calendar Breakdown:**
- Start Date: March 2, 2026 (Day after M1 launch)
- End Date: April 14, 2026 (M2 Launch Day)
- Total Calendar Days: 44 days

**Holidays in M2 Period:**
| Date | Day | Holiday Name | Status |
|------|-----|--------------|--------|
| Mar 18, 2026 | Wednesday | Shab-e-Qadr | Moon Dependent |
| Mar 20, 2026 | Friday | Jumatul Bidah | Moon Dependent |
| Mar 21, 2026 | Saturday | Eid-ul-Fitr Day 1 | Moon Dependent |
| Mar 22, 2026 | Sunday | Eid-ul-Fitr Day 2 | Moon Dependent |
| Mar 23, 2026 | Monday | Eid-ul-Fitr Day 3 | Moon Dependent |
| Mar 24, 2026 | Tuesday | Eid-ul-Fitr Day 4 (Extended) | Moon Dependent |
| Mar 25, 2026 | Wednesday | Eid-ul-Fitr Day 5 (Extended) | Moon Dependent |
| Mar 26, 2026 | Thursday | Independence Day | Confirmed |
| Apr 14, 2026 | Tuesday | Pohela Boishakh (Launch Day) | Confirmed |

**Total Holidays:** 9 days (including launch day)

**Weekends in M2 Period:**
- Working Days: Saturday to Thursday (6 days/week)
- Friday is weekend
- Approximate Fridays: ~6 days

**Working Days Calculation:**
```
Calendar Days:        44 days
Minus Holidays:       -9 days
Minus Weekends:       -6 days
= Net Working Days:   29 days
```

**Effective Working Days for M2:** 29 days

---

### M3 Period: April 15, 2026 - August 1, 2026

**Calendar Breakdown:**
- Start Date: April 15, 2026 (Day after M2 launch)
- End Date: August 1, 2026
- Total Calendar Days: 109 days

**Holidays in M3 Period:**
| Date | Day | Holiday Name | Status |
|------|-----|--------------|--------|
| May 1, 2026 | Friday | May Day | Confirmed |
| May 27, 2026 | Wednesday | Eid-ul-Adha Day 1 | Moon Dependent |
| May 28, 2026 | Thursday | Eid-ul-Adha Day 2 | Moon Dependent |
| May 29, 2026 | Friday | Eid-ul-Adha Day 3 | Moon Dependent |
| May 30, 2026 | Saturday | Eid-ul-Adha Day 4 (Extended) | Moon Dependent |
| Aug 15, 2026 | Saturday | National Mourning Day | Confirmed |

**Total Holidays:** 5 days (May 1 is Friday, already weekend)

**Weekends in M3 Period:**
- Working Days: Saturday to Thursday (6 days/week)
- Friday is weekend
- Approximate Fridays: ~15 days

**Working Days Calculation:**
```
Calendar Days:        109 days
Minus Holidays:       -5 days
Minus Weekends:       -15 days
= Net Working Days:   89 days
```

**Effective Working Days for M3:** 89 days

---

## Summary: Total Working Days

| Phase | Period | Calendar Days | Holidays | Weekends | Working Days |
|-------|--------|---------------|----------|----------|--------------|
| **M1** | Dec 20, 2025 - Mar 1, 2026 | 72 | 3 | 10 | **59 days** |
| **M2** | Mar 2, 2026 - Apr 14, 2026 | 44 | 9 | 6 | **29 days** |
| **M3** | Apr 15, 2026 - Aug 1, 2026 | 109 | 5 | 15 | **89 days** |
| **TOTAL** | Dec 20, 2025 - Aug 1, 2026 | 225 | 17 | 31 | **177 days** |

---

## Work Hours Calculation

**Assumptions:**
- Work Week: Saturday to Thursday (6 days/week)
- Work Hours per Day: 8 hours
- Work Hours per Week: 48 hours (6 days × 8 hours)

### Team Size by Phase

**M1 Phase:**
- December 2025: 9 members available
- January 2026 onwards: 14 members available

**M2 & M3 Phase:**
- Full team: 14 members available

### Available Team Hours (Before Buffer)

**M1 Calculation:**
- December portion: ~2 weeks × 9 members × 48 hrs/week = 864 hours
- January-March portion: ~8 weeks × 14 members × 48 hrs/week = 5,376 hours
- **Total M1 (Gross):** 6,240 hours

**M2 Calculation:**
- 29 working days = ~5.8 weeks
- 14 members × 48 hrs/week × 5.8 weeks = 3,916 hours
- **Total M2 (Gross):** 3,916 hours

**M3 Calculation:**
- 89 working days = ~17.8 weeks
- 14 members × 48 hrs/week × 17.8 weeks = 11,980 hours
- **Total M3 (Gross):** 11,980 hours

---

## Buffer Strategy: 10% Strict Reserve

**Buffer Policy:**
- Reserve 10% of all available hours for:
  - Unexpected bugs and rework
  - Production incidents
  - Emergency changes
  - Technical debt
  - Testing extensions

### Net Available Hours (After 10% Buffer)

| Phase | Gross Hours | 10% Buffer | Net Available Hours |
|-------|-------------|------------|---------------------|
| **M1** | 6,240 hrs | 624 hrs | **5,616 hrs** |
| **M2** | 3,916 hrs | 392 hrs | **3,524 hrs** |
| **M3** | 11,980 hrs | 1,198 hrs | **10,782 hrs** |
| **TOTAL** | 22,136 hrs | 2,214 hrs | **19,922 hrs** |

---

## Critical Notes

1. **M1 is the most critical phase** with foundational services and infrastructure
2. **M2 has limited time** due to heavy Eid holidays (5 days) + Independence Day + Launch Day
3. **M3 has comfortable time** for desirable features and enhancements
4. **Buffer is non-negotiable** - must maintain 10% reserve for each phase
5. **Moon-dependent holidays** may shift by 1-2 days; plan conservatively
6. **Launch days (M1, M2)** are counted as holidays since full development work is not feasible

---

**Last Updated:** [Current Date]
**Version:** 1.0
