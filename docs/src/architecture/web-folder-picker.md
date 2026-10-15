# Web FolderPicker decision

**Decision, 2026-10-15: do not implement browser-folder RPC in the Web
executor at present.** Keep its `FolderPicker` disabled. This completes the
design investigation in series 3/3 of TG-110, depending on TG-111; it does not
schedule an implementation or change the protocol. Desktop support shipped in
[#257](https://github.com/voilelab/toolgui/pull/257), and Wasm support in
[#259](https://github.com/voilelab/toolgui/pull/259).

## Demand and alternatives

The likely first audience is a local tool whose Go server and browser run on
the same machine: batch conversion, inspecting a source tree, or exporting
results beside inputs. That is a hypothesis, not demonstrated demand. The
maintainer confirmed during this investigation that there is currently no
concrete app requiring a Web server to operate on a browser-granted folder.
The repository's upload/download examples and Local folder demo demonstrate
APIs, not an independent app that needs this transport.

Running on the same machine does not turn a browser directory handle into a
server path. The handle can be cloned to another browser thread, not sent as
a capability to Go over a WebSocket. Nor should the browser claim an absolute
path that the server then opens. Even loopback would need a browser agent for
each filesystem operation. On a shared remote server, the folder would still
be on each visitor's machine; this is a different, more sensitive deployment.

| Need | Current choice | Limitation |
| --- | --- | --- |
| Select inputs, run server processing, save outputs | `MultiFileUpload` plus `DownloadFile`; package multiple outputs into one archive | Copies inputs to the server; not directory traversal or in-place editing. Multi-file selection is currently capped at 1000 files per pick. |
| Read/write a local tree with Go-specific dependencies | Desktop `FolderPicker` | Requires installing a desktop app. |
| Offline browser processing of a tree | Wasm `FolderPicker` | Requires compatible browser APIs and an app that can run in Wasm. |
| Upload an entire directory while preserving relative names | A possible future `webkitdirectory` input | Read-only selection of file snapshots; still uploads bytes. It neither grants write access nor implements the existing live Folder interface. Not added here. |

Uploads/downloads cover the demonstrated import/process/export model. They
do not cover every hypothetical tree-processing app, especially incremental
random reads or writing results in place. That gap alone does not justify a
new bidirectional filesystem protocol without a named consumer.

## What implementing it would require

The following is a candidate design to expose the cost, **not a protocol
specification to implement**. A future implementation needs its own ticket
and review. No `tgframe.Folder` methods are added by this decision.

### Grants and messages

Opening the picker must remain a main-thread user gesture. The server cannot
request a native dialog during an arbitrary page run. A successful selection
would create a browser-side grant registry entry, followed by a client-initiated
`folder_granted` control message. It would identify the connection epoch,
picker component, opaque grant ID and granted read/readwrite mode. The server
would validate the component and pending selection before applying the trusted
`EventFolder` and rerunning. Cancellation or denial would leave state unchanged.
The handle itself and an OS path would never be sent to the server.

Every successful WebSocket handshake would negotiate `folder_fs_v1` and a
fresh unpredictable session epoch. Reusing an idle Web state ID on reconnect
must not reuse an epoch, outstanding requests or directory grants. Peers
without the capability keep the existing disabled behavior.

Illustrative server request and browser response:

```json
{"fs_request":{"v":1,"session_epoch":"e7","run_id":"r9","request_id":"q12","grant_id":"g3","op":"read","file_id":"f2","offset":0,"length":65536}}
{"fs_response":{"v":1,"session_epoch":"e7","run_id":"r9","request_id":"q12","result":{"data_base64":"aGVsbG8=","eof":true}}}
```

An error response has an `error` instead of `result`, for example
`{"code":"permission_denied"}`. Stable codes would map to `fs.ErrNotExist`,
`fs.ErrPermission`, `fs.ErrInvalid`, `fs.ErrClosed` and
`tgframe.ErrFolderReadOnly`; browser error text is diagnostic, not the wire
contract. Offsets, lengths, decoded payload sizes and response metadata must
be checked before allocation or use.

The operation set needs `open`, `read`, `stat`, `readdir`, `create`, `write`,
`mkdir_all`, `close` and `abort`. `open` returns metadata and a browser-side
file ID for a File snapshot; `create` returns a writer ID, and writer `close`
commits its stream. IDs are scoped to the grant, epoch and run. Relative names
follow `fs.ValidPath`; `Path()` remains empty. Directory listings need bounded
pages and a listing ID, with a defined snapshot and sorting policy. Writers
need ordered chunks and per-writer sequence validation. A duplicate or late
response must not complete another request or repeat a write.

This JSON/base64 shape is enough for a prototype, not a good large-file data
plane. A production design must compare binary frames or bounded HTTP streams
associated with these request IDs. The existing update-message limit is 1 MiB;
it should not be raised globally to accommodate file contents. Negotiated
chunk limits, backpressure and per-session in-flight byte limits are required.

### Dispatch, reruns and cancellation

The socket is already bidirectional, but `SendPackFunc` only delivers packs;
it has no response correlation or waiter lifecycle. It also carries query,
navigation and component-notification packs, beyond create/update/delete,
ready and result. Adding an FS request tag alone would not supply RPC.

The Web executor currently receives a frame and synchronously calls
`Session.HandleRawEvent`. A new user event cancels the previous run and waits
for it to finish. If that run is waiting for browser IO, and the FS response
is behind the user event on the socket, the receiver can deadlock. A future
receiver must keep reading, dispatch FS responses directly to a pending-request
registry, and queue normal UI events to a separate ordered handler. FS replies
must never go through `ParseEvent`, apply component state, or cause a rerun.
The frontend similarly needs an FS control dispatcher before render dispatch.

Each request would register a waiter before sending, then release all send,
state and registry locks before waiting on its response, deadline or run context.
The design needs these lifecycle rules:

| Trigger | Required behavior |
| --- | --- |
| User event interrupts a run | `Params.Context` cancellation wakes its IO waiters immediately with `context.Canceled`; late replies cannot change state. Best-effort `fs_cancel` asks the browser to abort pending work. |
| Normal run completion | Close readers and release run-owned IDs; abort any writer the app did not close. A committed writer must have completed before run success. |
| Same-session rerun | Keep the directory grant, allocate a new run ID and new run-bound proxies. `ReadyPack` is a rendering boundary, not grant revocation. |
| `App.RerunAll` / `RerunPage` | Preserve existing behavior: coalesce and wait for the active run to finish, without interrupting it. They do not rescue a stuck FS call; IO deadlines are still necessary. |
| Session replacement or connection loss | Detach the epoch and grants, wake all waiters before waiting for the runner, and reject old-session operations with `ErrSessionReplaced`. Session replacement takes precedence over generic cancellation when both apply. Reconnect requires a fresh user selection. |
| Timeout or permission revocation | Complete the waiter with a defined timeout/permission error; stop further use of an invalid grant. Never retry writes automatically. |

The TG-111 Folder methods have no context parameter, and `State.GetFolder`
currently returns a session-owned object. A future implementation must work
out an internal run-scoped facade: each `FolderPicker` result, reader and
writer captures that run's `Params.Context`, while only the underlying grant
survives reruns. Mutating one cached proxy to point at the newest run would
let old work act under a new context and is unacceptable. How to expose that
run scope through the existing container/state path needs separate design;
waiting with `context.Background()` is not a solution. If the public interface
proves insufficient, revisit TG-111 and all executors together.

Cancellation is not a transaction rollback. Abort may discard an uncommitted
stream, but a completed `mkdir_all`, newly created file entry, or already
committed write cannot be undone reliably. If a commit succeeds and the reply
is lost, the outcome is unknown to the server. Automatic replay on reconnect
would be unsafe; the app must reconcile before retrying.

## Security boundary

A directory grant would authorize **the connected server application**, not
just the local page UI, to request reads under that root. Readwrite would also
authorize creating and overwriting files without another dialog per operation.
A compromised server or injected same-origin script could read the entire
granted tree, send it elsewhere, or damage writable files. The grant must be
explained that way; it must not look like selecting a harmless display path.

The browser agent, not the server alone, must enforce the root and mode for
every operation. It must traverse from registered handles, reject invalid and
escaping paths, disallow unsupported operations such as arbitrary remove or
rename, reject unknown IDs and revoked grants, and bind requests to the exact
connection/epoch and active run. The service cannot silently widen a read grant
or reuse a browser's broader persisted permission. Origin checks do not replace
these capability checks or protect against a compromised trusted origin.

ToolGUI currently has no built-in user authentication; a Web state ID is a
bearer credential, not a verified user identity. On a shared server, a routing
mistake, state-ID leak, global cache or app-level cross-session reference must
not expose one visitor's grant or response to another. Keep capability maps
per authenticated connection and session, never in app caches; require deployment
authentication, exact origin checks and HTTPS/WSS for remote service. Disconnect
revokes the protocol grant even when the browser still remembers OS permission.
Loopback is a useful default for single-user tools, not a substitute for these
checks. Any future multi-user support needs an explicit tenant-isolation review.

Apply request, byte, open-file and directory-entry quotas on both peers, plus
timeouts and bounded queues. Do not log file contents, bearer/grant IDs or full
local names by default. Browser responses are also untrusted server input:
validate IDs, lengths, entry names and chunk sequences rather than trusting a
visitor not to exhaust memory or forge another session's response.

## Performance budget

These are analytical examples, not measurements. For sequential stop-and-wait
IO, latency alone is approximately `request_count × RTT`, in addition to disk,
serialization and byte-transfer time. Even loopback retains browser event-loop
and Promise scheduling costs. Each remote read must cross to Go; write output
must cross back, even when both endpoints are on the same machine.

| Workload | Round-trip latency floor at 1 ms | At 30 ms |
| --- | --- | --- |
| Read 1 GiB in 64 KiB chunks: 16,384 reads | 16.4 s | 491.5 s (8.2 min) |
| 10,000 sequential metadata or small-file requests | 10 s | 300 s (5 min) |

At 100 Mbit/s, 1 GiB of bytes alone takes about 86 seconds one way, ignoring
overheads; base64 increases the payload by about one third. Go readers can ask
for just a few bytes at a time, so sending an RPC for every `Read` would be much
worse than the chunked example. Open/stat/close, listing and output writes add
further requests. Reruns could repeat the work unless the app explicitly caches
its results with a safe invalidation policy.

Read-ahead buffers, metadata batching, bounded concurrent requests and streamed
writes could reduce that cost. They also require flow control, cancellation,
snapshot consistency and commit semantics. They must not block small UI/control
frames behind large transfers. Increasing a chunk size trades latency for
memory and cancellation responsiveness; it does not remove the need to upload
bytes. For import/process/export, uploading each input once and reading it on
the server avoids repeated browser round trips using the existing protocol.

## Outcome and reopening criteria

Do not open a Web FS implementation ticket under TG-110 now: the decision is
not to implement it. Keep the existing disabled UI and document the reason on
the [FolderPicker page](../components/input/folder_picker.md). Its current
generic unavailable message is not a promise that support is scheduled; this
documentation-only change intentionally does not alter component source.

Reopen the decision when there is a named app, an owner, and a workflow that
cannot reasonably use uploads/downloads, desktop or Wasm. Require representative
file counts/sizes, local versus remote deployment, an acceptable latency budget,
permission UX, and measured results from an isolated transport prototype. Only
then scope and attach an implementation ticket to TG-110, including the run-context,
deadlock, multi-user isolation and uncertain-write-outcome requirements above.

## Evidence checked

- [Session and pack delivery](https://github.com/voilelab/toolgui/blob/db7efc2/toolgui/tgframe/session.go): serialized sends, interrupt-and-wait, run cancellation.
- [Server reruns and Params.Context](https://github.com/voilelab/toolgui/blob/db7efc2/toolgui/tgframe/app.go): non-interrupting queued reruns.
- [Web receiver and state acquisition](https://github.com/voilelab/toolgui/blob/db7efc2/toolgui/tgexec/web.go): synchronous event handling, reconnect state and limits.
- [Session replacement](https://github.com/voilelab/toolgui/blob/db7efc2/toolgui/tgframe/host.go) and [TG-111 Folder](https://github.com/voilelab/toolgui/blob/db7efc2/toolgui/tgframe/folder.go).
- [Frontend socket](https://github.com/voilelab/toolgui/blob/db7efc2/toolgui-web/app/src/api/StatefulWebSocket.ts) and [render dispatch](https://github.com/voilelab/toolgui/blob/db7efc2/toolgui-web/lib/src/app/dispatch.ts).
- [Deployment trust model](exposure.md), [uploads](../components/input/fileupload.md), [downloads](../components/input/download_file.md), and [Local folder demo](https://github.com/voilelab/toolgui/blob/db7efc2/cmd/toolgui-demo/folder.go).
- Maintainer confirmation on 2026-10-15: no concrete app currently needs this Web filesystem bridge. This is not a claim that no outside user could need it.
