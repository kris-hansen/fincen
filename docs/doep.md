---
layout: page
title: Designation of Exempt Person (DOEP) - Form 110 | Moov FinCEN
hide_hero: true
show_sidebar: false
menubar: docs-menu
---

# Overview

Designation of Exempt Person (DOEP) - Form 110

# Create a report

Designation of exempt person can create using fincen go library

1. Create a [EFilingBatchXML](https://godoc.org/github.com/moov-io/fincen/pkg/batch#EFilingBatchXML) with `batch.NewReport("SARX")`.
2. Create available [ActivityType](https://godoc.org/github.com/moov-io/pkg/exempt_designation#ActivityType) records with `exempt_designation.NewActivity()`.
3. Append created activities into Batch XML report with `batch.AppendActivity(activity)`.
4. Validate Batch XML report with `Validate()` and figure out report problems.
5. Generate Batch XML report attributes with `GenerateAttrs()`
6. Getting xml contents from Batch XML report.

# Create an acknowledgement

FinCEN SAR XML batch acknowledgement can create using fincen go library

1. Create a [EFilingBatchXML](https://godoc.org/github.com/moov-io/fincen/pkg/batch#EFilingBatchXML) with `batch.NewReport("SARX")`.
2. Create a [EFilingSubmissionXML](https://godoc.org/github.com/moov-io/pkg/batch#EFilingSubmissionXML).
3. Validate Batch XML report with `Validate()` and figure out report problems.
4. Generate Batch XML report attributes with `GenerateAttrs()`
5. Getting xml contents from Batch XML report.
