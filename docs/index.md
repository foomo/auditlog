---
layout: home

hero:
  name: auditlog
  text: Type-safe audit logging for Go
  tagline: A generic audit-log domain with MongoDB persistence and TTL-based retention. You own the payload.
  image:
    src: /logo.png
    alt: auditlog
  actions:
    - theme: brand
      text: Get started
      link: /guide/getting-started
    - theme: alt
      text: View on GitHub
      link: https://github.com/foomo/auditlog

features:
  - title: Generic entry envelope
    details: Entry[Payload] carries the standard fields — id, timestamp, user, request, service/func/action, entity, ip, user-agent, ttl — around a payload union you define.
  - title: MongoDB with TTL retention
    details: Entries expire automatically. Retention is configured once at repository construction via WithRetention; the TTL field is anchored on insert.
  - title: Composable command / query handlers
    details: CreateEntry, Search and Get mirror the foomo/redirects domain, so middleware composition — trace events, project wrappers — is identical.
  - title: No wire surface in the library
    details: Projects declare their own typed Service interfaces over Entry[Payload] and own all gotsrpc generation. The library stays free of any-typed shapes that confuse the generator.
---
