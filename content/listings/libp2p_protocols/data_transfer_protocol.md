---
title: Data Transfer Protocol
dashboardWeight: 0.2
dashboardState: wip
dashboardAudit: n/a
---

# Data Transfer Protocol
{{< hint info >}}
- **Name**: Data Transfer Protocol
- **Protocol ID**: `/fil/data-transfer/0.0.1`
{{< /hint >}}

Message Protobuf

```go

message DataTransferMessage {
  message Request {
    int32 transferID = 1
    bool isPull = 2
    bytes voucher = 3
    bytes pieceID = 4
    bytes selector = 5
    bool isPartial = 6
    bool isCancel = 7
  }

  message Response {
    int32 transferID = 1
    boolean accepted = 2
  }

  bool isResponse = 1
  Request request = 2
  Response response = 3
}

```
