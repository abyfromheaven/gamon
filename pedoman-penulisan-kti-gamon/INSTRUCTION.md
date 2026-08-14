# Skill: Writing an SMK RPL PKL Report and Scientific Paper (KTI) for Project GAMON

## 1. Purpose of This Skill

This skill provides a complete writing and reasoning framework for producing, revising, and expanding a **Praktik Kerja Lapangan (PKL) Report** presented as a **Karya Tulis Ilmiah (KTI)** for an Indonesian vocational high school (SMK), specifically within the **Rekayasa Perangkat Lunak (RPL / Software Engineering)** competency.

The goal is to ensure that all discussion about the GAMON project is:

- grounded in the actual project implementation;
- consistent with the school's PKL-report structure;
- formal and scientific, but appropriate for SMK level;
- clear and concise;
- focused on the purpose of each chapter and subsection;
- free from unnecessary repetition;
- free from fabricated technical facts, features, technologies, requirements, or project decisions.

This is **not** a thesis-writing guide for undergraduate, master's, or doctoral-level academic work.

The expected level is:

> **Technically credible and academically structured, but still accessible and appropriate for an SMK RPL PKL report.**

---

# 2. Report Context

The document is a:

**Praktik Kerja Lapangan (PKL) Report**

containing a Karya Tulis Ilmiah with the title:

> **“Perancangan Sistem Monitoring Infrastruktur TI Berbasis Website pada Universitas Pertahanan RI”**

The main software project discussed in the report is:

> **GAMON (Garda Monitoring)**

GAMON is a **prototype web-based information technology infrastructure monitoring system**.

The report should explain the relationship between:

**Observed Condition → Problem → Analysis → Impact → Proposed Solution → Requirements → System Design → Implementation → Result**

The report is **not merely a daily PKL journal** and is **not merely source-code documentation**.

---

# 3. Main Problem Behind GAMON

The project originates from observations made during the PKL in an IT infrastructure management environment, including monitoring activities in a Network Operation Center (NOC).

The core problems that form the basis of the project are:

- a large number of devices need to be monitored;
- devices come from different vendors;
- different vendors may provide different management interfaces;
- monitoring is performed through multiple interfaces or media;
- monitoring is not centralized;
- device failures are not always detected quickly;
- there is no centralized notification mechanism within the scope of the prototype.

These problems must remain the foundation of the GAMON discussion throughout the report.

## 3.1 The Role of the NOC

The NOC is the **context in which the problem was observed**, not necessarily the exclusive target user of GAMON.

The report may state that the problem was identified through monitoring activities in the NOC.

However, do not state that GAMON is exclusively a “NOC monitoring system” if the actual prototype is designed generically.

The correct conceptual relationship is:

> **The problem was observed in the PKL environment, while the resulting GAMON prototype is designed generically.**

## 3.2 The Role of Universitas Pertahanan RI

Universitas Pertahanan RI is the institution where the PKL was conducted and where the problem context was observed.

The report should distinguish between:

- the **case/context**: Universitas Pertahanan RI;
- the **prototype**: GAMON;
- the **general applicability** of GAMON: potentially broader than the original observation environment.

Do not describe GAMON as a production replacement for the institution's existing monitoring infrastructure unless that is explicitly proven.

Always use the term **prototype** where appropriate.

---

# 4. Truthfulness and Anti-Hallucination Rules

These rules are mandatory.

## 4.1 Source Code Is the Primary Source for Technical Facts

When discussing:

- technologies actually used;
- libraries;
- frameworks;
- databases;
- APIs;
- endpoints;
- database tables;
- functions;
- algorithms;
- monitoring mechanisms;
- notification mechanisms;
- WebSocket behavior;
- system flows;
- modules;
- configuration;

use the actual project source code, repository structure, configuration, database, and project documentation as the primary evidence.

## 4.2 Never Invent Technical Facts

If something cannot be confirmed from the implementation, explicitly state:

> “This cannot be directly inferred from the project implementation.”

Do not invent project decisions.

## 4.3 Distinguish Fact from Inference

Example:

**Fact:**
> The project uses Go and goroutines in the monitoring engine.

**Reasonable inference:**
> The use of goroutines is consistent with the requirement to perform monitoring processes concurrently.

**Unsupported claim:**
> Go was chosen because the developer intended to make the system twice as fast as Node.js.

Never turn an inference into a historical fact about why the developer made a decision.

## 4.4 Do Not Treat Planned Features as Implemented Features

If something exists only in:

- a roadmap;
- future work;
- an idea;
- a PRD;
- a design document;
- an unfinished branch;
- comments;

do not describe it as an active implemented feature.

Use terms such as:

- planned;
- proposed;
- future development;
- not yet implemented.

---

# 5. Writing Style for an SMK-Level KTI

Use a style that is:

- formal;
- objective;
- clear;
- concise;
- technically credible;
- easy for a teacher or examiner to understand.

Avoid marketing-style language such as:

- revolutionary;
- extremely powerful;
- best;
- perfect;
- cutting-edge;
- very advanced;

unless objectively supported.

Avoid casual language.

Prefer:

- “the author”;
- “the user”;
- “the system”;
- “the application”;
- “the prototype”;
- “the implementation”;
- “the system is designed to…”;
- “the application uses…”.

Do not write like an undergraduate thesis unless the school's requirements explicitly demand it.

---

# 6. Scientific Writing Principles

Each paragraph should have a clear purpose.

A useful pattern is:

> **Statement → Explanation → GAMON Context**

Example:

> React is a JavaScript library for building component-based user interfaces. In GAMON, React is used to construct the frontend pages and reusable interface components. This approach allows the interface to be organized into reusable components and updated dynamically based on monitoring data.

Do not place all of the following into one paragraph:

- definition;
- history;
- selection rationale;
- internal mechanism;
- source code;
- test result.

Separate them according to the purpose of the section.

Use tables when information is naturally tabular.

Use diagrams when relationships or flows are more clearly represented visually.

---

# 7. Citations and References

## 7.1 Claims That Need References

External references should support theoretical claims such as:

- definition of monitoring;
- availability;
- functional requirements;
- non-functional requirements;
- Go;
- React;
- TypeScript;
- SQLite;
- REST API;
- WebSocket;
- ICMP;
- UML;
- ERD;
- PKL methodology;
- educational concepts;
- other technical concepts.

## 7.2 Reference Priority

Prefer sources in this order:

1. official Indonesian government documents;
2. Indonesian scientific journals;
3. academic books;
4. official technology documentation;
5. standards/RFCs;
6. international academic sources;
7. reputable technical articles.

Indonesian-language sources are preferred when a suitable and credible source exists.

## 7.3 Specific Source Identification

Do not write vague references such as:

> “Source: Google”

Instead identify:

- author or organization;
- publication year;
- article/document title;
- journal/publisher;
- page number or section when available;
- URL or DOI.

For PDFs, identify the exact page when possible.

## 7.4 Observation Does Not Need a Fake External Source

Statements describing:

- what was directly observed during PKL;
- what the administrator actually did;
- what monitoring conditions were actually encountered;

may be presented as observations, interviews, or project evidence.

Do not invent an external citation for something that is specifically the author's observation.

---

# 8. Main Report Structure

The report uses the following general structure:

```text
BAB I PENDAHULUAN

A. Latar Belakang Pemilihan Judul
B. Metode Penyusunan dan Pengumpulan Data
C. Sistematika Penyusunan Karya Tulis

BAB II TRANSFORMASI PKL

A. Tujuan Praktik Kerja Lapangan (PKL)
B. Manfaat Praktik Kerja Lapangan (PKL)

BAB III TINJAUAN UMUM PERUSAHAAN/INSTANSI

A. Sejarah Singkat Perusahaan/Instansi
B. Visi, Misi, dan Tujuan
C. Struktur Organisasi
... sesuai dengan pedoman sekolah

BAB IV PEMBAHASAN MASALAH

A. Pembahasan Masalah yang Mengacu pada Judul
B. Implementasi Kegiatan PKL

BAB V PENUTUP
```

The exact first-level headings must remain consistent with the school's official PKL guideline.

---

# 9. BAB I — PENDAHULUAN

## 9.1 A. Latar Belakang Pemilihan Judul

Main question:

> **Why was this topic selected?**

Recommended logical flow:

```text
Information Technology Development
↓
Importance of IT Infrastructure Monitoring
↓
PKL Observation Results
↓
Problem Identification
↓
Proposed Solution
↓
Reason for Choosing the Title
```

### Information Technology Development

Discuss:

- development of information technology;
- organizational dependence on IT;
- importance of IT infrastructure;
- service availability.

Do not introduce PKL or GAMON too early.

### Importance of IT Infrastructure Monitoring

Discuss:

- monitoring function;
- monitoring devices/services;
- early detection;
- importance of monitoring infrastructure.

### PKL Observation Results

Discuss actual field findings.

No forced journal citation is required for direct observations.

### Problem Identification

Summarize the observed problems.

### Proposed Solution

Introduce the solution conceptually.

Do not explain implementation details.

### Reason for Choosing the Title

Keep this short.

Its purpose is only to connect:

> **Problem → Solution → Title**

Do not repeat the entire background.

## 9.2 B. Metode Penyusunan dan Pengumpulan Data

Focus on:

- introduction to data collection;
- observation;
- interview;
- literature study.

Explain what the author actually did.

Do not turn this section into a long research-methodology chapter.

## 9.3 C. Sistematika Penyusunan Karya Tulis

This should be written after the entire report is structurally stable.

Its purpose is simply to summarize the contents of Chapters I–V.

---

# 10. BAB II — TRANSFORMASI PKL

## Opening Paragraph

The opening paragraph should introduce:

- what PKL is;
- its role in SMK education;
- the connection between school learning and workplace practice;
- transition to PKL objectives and benefits.

Keep it short.

## A. Tujuan PKL

Discuss the objectives of PKL based on relevant official guidance.

## B. Manfaat PKL

Discuss the benefits of PKL in the context of vocational education.

Do not turn BAB II into a broad essay on the history of vocational education.

---

# 11. BAB III — TINJAUAN UMUM PERUSAHAAN/INSTANSI

The main focus is the institution, not GAMON.

Possible content:

- history;
- profile/identity;
- vision;
- mission;
- organizational structure;
- relevant institutional activities.

Do not place:

- GAMON algorithms;
- monitoring implementation;
- source code;
- project technology discussion;

in this chapter.

---

# 12. BAB IV — PEMBAHASAN MASALAH

BAB IV is the core of the report.

The overall chain should feel like:

> **Problem → Analysis → Impact → Solution → Requirements → Design → Implementation → Result**

---

# 13. BAB IV — A. Pembahasan Masalah yang Mengacu pada Judul

## 13.1 4.1 Kondisi Permasalahan

Main question:

> **What condition was found?**

Focus on observed facts:

- monitoring conditions;
- many devices;
- multiple vendors;
- multiple management interfaces;
- manual monitoring;
- lack of centralization.

Do not jump into solution implementation.

## 13.2 4.2 Analisis Permasalahan

Main question:

> **Why does that condition become a problem?**

Discuss:

- monitoring complexity;
- switching between interfaces;
- dependency on manual checks;
- delayed awareness of failures.

## 13.3 4.3 Dampak Permasalahan

Main question:

> **What is the impact?**

Discuss:

- delayed detection;
- slower response;
- reduced monitoring efficiency;
- operational complexity.

Avoid repeating 4.2.

## 13.4 4.4 Solusi yang Diusulkan

Main question:

> **What solution is proposed?**

Introduce:

> **GAMON (Garda Monitoring)**

Discuss conceptually:

- centralized monitoring;
- device management;
- monitoring;
- alert notification;
- history.

Do not introduce implementation-level source code or technical configuration.

---

# 14. BAB IV — B. Implementasi Kegiatan PKL

This section transitions from the proposed solution to the actual prototype implementation.

---

# 15. 4.5 Gambaran Umum Aplikasi

Only answer four questions:

1. What is GAMON?
2. What is the purpose of GAMON?
3. What is its scope?
4. Who are the intended users?

Do not put detailed algorithms, architecture, APIs, source code, or internal implementation here.

---

# 16. 4.6 Teknologi yang Digunakan

The key questions are:

1. **What is this technology/concept?**
2. **What is it used for in GAMON?**
3. **How does it work in the context of GAMON?**
4. **Why is it used, if this can be supported by the project?**

This section should not be a long tutorial.

## 16.1 Important rule: Introduce the concept first

Technical explanations should **not immediately jump into GAMON implementation**.

For every important technology, technique, protocol, algorithm, or supporting concept, introduce:

> **What is it? → What is it used for in GAMON? → How does it work in GAMON?**

Then, if relevant:

> **Why is it used? → How is it implemented?**

This is important because the reader may not understand a project-specific explanation without first understanding the underlying concept.

## 16.2 Example: React

The explanation should first establish:

- what React is;
- what component-based UI means at a basic level.

Then explain:

- React is used to build GAMON's frontend;
- the frontend contains Dashboard, Monitoring, Device Management, and other UI components;
- state and components allow interface data to update based on backend/real-time data.

Detailed source code belongs in 4.9.

## 16.3 Example: CORS

If CORS is actually used:

1. explain what Cross-Origin Resource Sharing is;
2. explain the browser-origin problem it addresses;
3. explain where GAMON uses it;
4. explain briefly how the request/response relationship works in GAMON;
5. explain the middleware/configuration in 4.9 if needed.

Do not assume CORS exists just because a web application usually uses it. Verify it in the actual project.

## 16.4 Example: Middleware

If middleware is actually present:

1. explain what middleware means in web applications;
2. explain its role before the main request handler;
3. explain which middleware is used in GAMON;
4. explain the request flow;
5. show code in 4.9 when necessary.

Do not discuss middleware types that are not present in the project.

## 16.5 Example: Go and Concurrency

Introduce:

- what Go is;
- what concurrency means;
- what a goroutine is in simple terms.

Then connect the concept to GAMON:

- Go serves as backend;
- monitoring jobs can run concurrently;
- goroutines are used in the monitoring process;
- this matches the need to monitor multiple devices.

Do not turn the section into an internal Go runtime tutorial.

## 16.6 Example: REST API

Introduce:

- REST/API;
- request-response;
- HTTP methods at a basic level.

Then explain:

- GAMON uses REST API between frontend and backend;
- it handles CRUD/data operations;
- JSON is used for data exchange if supported by the project.

Detailed endpoint implementation belongs in 4.9.

## 16.7 Example: WebSocket

Introduce:

- what WebSocket is;
- how it differs conceptually from request-response polling.

Then explain:

- GAMON uses WebSocket for real-time updates;
- backend can push monitoring changes;
- frontend receives changes and updates the interface.

Detailed Hub/client/broadcast code belongs in 4.9.

## 16.8 Example: ICMP Ping

Introduce:

- what ICMP is;
- Echo Request/Echo Reply;
- basic meaning of latency.

Then explain:

- GAMON uses ICMP Ping to check device reachability;
- results are converted into monitoring status and latency;
- the result contributes to the alert logic.

Detailed timeout/parsing/source code belongs in 4.9.

## 16.9 Other Important Concepts

If actually used by GAMON, the same introductory pattern may apply to:

- CORS;
- middleware;
- concurrency;
- goroutines;
- request-response;
- push communication;
- polling;
- JSON;
- client-server architecture;
- failure threshold;
- pairing token;
- database migration;
- state management;
- routing;
- HTTP methods;
- other relevant concepts.

Do not introduce a concept merely because it is common in the technology ecosystem. Confirm that GAMON actually uses it.

## 16.10 Tools and Development Environment

Tools such as:

- Node.js;
- NPM;
- Vite;
- Git;
- GitHub;
- Zed Code Editor;

may be included if actually used and if the school's report format treats development tools as part of “technologies used”.

However, these should receive less space than technologies that directly form the system architecture.

## 16.11 Recommended Level of Detail

For core technologies:

- usually 1–2 paragraphs;
- simple technologies may use only 1 paragraph.

The goal is:

> **Introduce → contextualize → briefly explain**

not:

> **Teach the entire technology.**

---

# 17. 4.7 Analisis Kebutuhan Sistem

Structure:

```text
4.7 Analisis Kebutuhan Sistem
Opening paragraph

4.7.1 Kebutuhan Fungsional
Short opening paragraph
Functional Requirement table

4.7.2 Kebutuhan Non-Fungsional
Short opening paragraph
Non-Functional Requirement table
```

## 17.1 Functional Requirements

Answer:

> **What must the system do?**

Use a table.

The currently defined GAMON functional areas include:

- device registration;
- device listing;
- device editing;
- device deletion;
- enabling/disabling monitoring;
- periodic monitoring;
- device status determination;
- failure threshold;
- history storage;
- automatic alert creation;
- automatic alert resolution;
- alert management;
- real-time application notifications;
- Telegram notification;
- dashboard summary;
- real-time monitoring display.

Verify the final list against the project before treating it as final.

## 17.2 Non-Functional Requirements

Answer:

> **What quality characteristics must the system satisfy?**

Current categories:

- performance efficiency;
- reliability;
- availability;
- security;
- usability;
- maintainability;
- portability.

Do not invent measured performance values or SLA targets unless they were actually tested.

---

# 18. 4.8 Perancangan Sistem

Main question:

> **How was the system designed before implementation?**

Important distinction:

> **Design = blueprint.**
>
> **Implementation = realization.**

## 18.1 4.8.0 Opening

Explain that the design is derived from the functional and non-functional requirements defined in 4.7.

## 18.2 4.8.1 System Architecture

Answer:

> **What components make up the system and how are they connected?**

Possible components, subject to verification:

- frontend;
- backend;
- database;
- monitoring engine;
- REST API;
- WebSocket;
- ICMP;
- Telegram.

Focus on component relationships, not source code.

## 18.3 4.8.2 Use Case Diagram

Answer:

> **What can users do with the system?**

Use cases must map to functional requirements.

## 18.4 4.8.3 Activity Diagram

Answer:

> **How does a major activity proceed?**

Use only major activities, such as:

- monitoring;
- device management;
- Telegram pairing.

Do not create separate diagrams for every basic button.

## 18.5 4.8.4 Flowchart

Answer:

> **How does the internal logic of an important process work?**

Priority:

**Monitoring Engine Flowchart**

Possible logic:

```text
Select device
↓
Ping
↓
Response?
↓
Update failure count
↓
Compare with threshold
↓
Status change?
↓
Create/resolve alert
↓
Save history
↓
Broadcast update
↓
Wait interval
↓
Repeat
```

A Telegram pairing flowchart may be added if the mechanism is genuinely significant and not already covered adequately by an activity diagram.

Do not make broad flowcharts for the entire frontend or REST API.

## 18.6 4.8.5 ERD

Answer:

> **How are the data entities related?**

Focus on:

- entities;
- relationships;
- cardinality;
- relevant foreign keys.

## 18.7 4.8.6 Database Structure

Answer:

> **What tables exist and what does each table store?**

Focus on:

- table list;
- purpose of each table;
- important fields if necessary.

Do not document every column if it makes the report unnecessarily long.

## 18.8 UI Design

A separate UI-design section is optional.

If the UI was created directly through AI-assisted implementation and no formal wireframe/design process was documented, do not fabricate a design process.

---

# 19. 4.9 Implementasi Sistem

Main question:

> **How was the design actually implemented?**

Possible structure:

- Dashboard;
- Monitoring;
- Device Management;
- Alert Notification;
- Monitoring History;
- User Management, only if actually implemented.

Each section may contain:

1. feature purpose;
2. implementation explanation;
3. screenshot;
4. important source code;
5. explanation of the source code.

Source code is supporting evidence, not the main content.

---

# 20. 4.10 Hasil Implementasi

Main question:

> **What was the result after implementation?**

Discuss:

- implemented features;
- monitoring results;
- alert behavior;
- notification behavior;
- history;
- dashboard;
- limitations.

If testing exists, use a table such as:

| Scenario | Expected Result | Actual Result | Status |
|---|---|---|---|

Do not claim perfect success without evidence.

---

# 21. Diagram Selection Rules

| Diagram | Main Question | Focus |
|---|---|---|
| System Architecture | What components make up the system? | System structure |
| Use Case Diagram | What can users do? | Interaction |
| Activity Diagram | How does an activity proceed? | Activity flow |
| Flowchart | How does an internal process work? | Algorithm/process |
| ERD | How is data related? | Database relationships |

Do not use multiple diagrams to explain the same thing.

---

# 22. Design vs Implementation

## Design

Discuss:

- structure;
- relationships;
- process models;
- architecture;
- data models.

Examples:
- architecture diagram;
- use case;
- activity diagram;
- flowchart;
- ERD.

## Implementation

Discuss:

- screenshots;
- actual source code;
- actual endpoints;
- actual configurations;
- implemented feature behavior;
- actual UI.

---

# 23. Source Code Rules

Source code should:

- come from the actual project;
- be selected for relevance;
- not be shown as entire files;
- be explained after the snippet.

Do not insert code only to increase page count.

If the school requires more complete code, additional source code may be placed in an appendix.

---

# 24. Report-Length Rules

The standard is:

> **Concise, not shallow.**

Use tables when appropriate.

Use diagrams when appropriate.

Use one paragraph when one paragraph is enough.

Do not force every subsection into multiple paragraphs.

Do not repeat theory.

Do not repeat a technology explanation in 4.9 if it was already introduced in 4.6.

Do not repeat requirements in 4.8.

Do not repeat the full design in 4.9.

---

# 25. Terminology Consistency

Use terminology consistently:

- prototype;
- GAMON (Garda Monitoring);
- monitoring;
- device/perangkat;
- user/pengguna;
- administrator;
- dashboard;
- alert;
- monitoring history;
- failure threshold.

Do not randomly replace technical terms between chapters.

---

# 26. Important Scope Limitations

Do not state that:

- GAMON is production-ready;
- GAMON has replaced the institution's monitoring platform;
- every institutional device is already integrated;
- the system is guaranteed to be free from false positives/false negatives;
- a specific uptime or response-time target has been guaranteed;

unless there is clear evidence or actual testing.

Use proportional wording:

- prototype;
- designed to;
- intended to;
- helps;
- can;
- within the project scope;
- based on the implementation.

---

# 27. Final AI Rules

When the AI has direct access to the GAMON project:

1. inspect the source code before making technical claims;
2. use source code as the main technical evidence;
3. use official documentation and academic references for theory;
4. never invent a feature;
5. never invent a technology-selection rationale;
6. distinguish fact, inference, and theory;
7. do not turn the prototype into a production system in the wording;
8. do not expand the scope beyond the actual implementation;
9. do not repeat material unnecessarily;
10. do not make a section longer merely to make it appear scientific;
11. use tables and diagrams when they communicate better than prose;
12. use SMK-level KTI language;
13. introduce technical concepts before discussing their GAMON application;
14. if a fact cannot be confirmed, explicitly state that it cannot be concluded from the implementation;
15. prioritize accuracy, consistency, traceability, and clarity over length.

# 28. Final Writing Logic

Every section should answer one primary question:

- **BAB I:** Why was this topic selected?
- **BAB II:** What is the PKL context?
- **BAB III:** What is the institution?
- **4.1:** What is the observed problem?
- **4.2:** Why is it a problem?
- **4.3:** What is the impact?
- **4.4:** What solution is proposed?
- **4.5:** What is GAMON?
- **4.6:** What technologies/concepts are used, what are they, how are they used in GAMON, and why are they used?
- **4.7:** What must the system do and what quality characteristics must it satisfy?
- **4.8:** How is the system designed?
- **4.9:** How is the design implemented?
- **4.10:** What is the implementation result?

If a piece of information does not help answer the primary question of the current section, consider moving it to a more appropriate section or removing it.

The final target is not the longest possible report.

The target is a report that is:

**accurate, structured, understandable, technically grounded, scientifically written, appropriate for SMK RPL, and defensible during the PKL examination.**
