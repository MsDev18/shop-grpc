## gRPC migration (strangler pattern) — done

Converted this API from REST-only to gRPC, module by module, without a rewrite: a gRPC
server ran on `:50051` in a goroutine alongside the Gin server on `:3000` until every module
below had a gRPC equivalent, then the REST/Gin layer (`handler`, `middleware`, `router`,
`server` packages) was deleted entirely. The project now runs as a single gRPC-only process
on `:50051`, single database, single binary.

- [x] health    — `health.HealthService` (Check)
- [x] auth      — `auth.AuthService` (SendOtp, CheckOtp, Me, RefreshToken, Logout)
- [x] user      — `user.UserService` (Profile, UpdateProfile, ChangePassword)
- [x] category  — `category.CategoryService` (Create, Update, Delete, GetAll, GetOne)
- [x] province  — `province.ProvinceService` (GetAll, GetOne)
- [x] address   — `address.AddressService` (Create, GetAll, GetOne, Update, Delete)
- [x] product   — `product.ProductService` (Create [client-streaming], GetAll, GetOneBySlug)
- [x] remove Gin/REST entirely
- [ ] cart / order / payment — not built in REST or gRPC yet
- [ ] decide how uploaded images get served now that Gin's static `/uploads` route is gone
      (removed along with REST; gRPC has no built-in equivalent)

---

## Microservices migration (next phase)

Now that the whole API is gRPC-only, the next phase is splitting this single monolith
process into independently deployable services. Because this is a learning project, the
end goal is deliberately fine-grained: **every** module becomes its own independent service
(own binary, own database) — not just one or two extracted while the rest stay in the
monolith.

- Gradual, one service at a time — not a big-bang split of the whole system at once.
- Each extracted service gets its own database from day one (real data ownership per
  service, not a shared-DB "distributed monolith" stepping stone).
- Every other module's direct in-process call into an extracted service (e.g. product's
  `categoryService category.Service` dependency) gets rewired from a plain Go function call
  into a real gRPC client call over the network — the actual cost of each extraction, and
  the reason it's done one module at a time instead of all at once.
- `auth` is deliberately extracted last: it's currently a cross-cutting concern (the shared
  interceptor validates every request for every module), so splitting it means deciding how
  every already-extracted service authenticates a request once that check is no longer an
  in-process call. Better to have the extraction mechanics down from simpler cases first.

Extraction order (simplest/most isolated first, hardest/most foundational last):

- [ ] 1. `health` — no database, no dependencies; pure exercise in "separate binary,
      separate port, separate deployment" with zero data or networking complexity
- [ ] 2. `province` — first real data-ownership split (own table -> own database); still no
      dependency on any other module
- [ ] 3. `category` — same shape as province (no outgoing dependency), but with real admin
      role checks and image upload
- [ ] 4. `address` — first real inter-service gRPC client: depends on `province`
- [ ] 5. `product` — second inter-service gRPC client, more complex domain: depends on
      `category`, plus client-streaming
- [ ] 6. `user` — mostly self-contained, but still needs a working answer for "how does a
      request get authenticated" once auth is external
- [ ] 7. `auth` — last: requires a real decision on how every other service validates
      tokens/sessions without an in-process call (e.g. per-request RPC to auth vs. local JWT
      signature verification + occasional revocation check)

- [ ] decide on service discovery / addressing between services
- [ ] decide on deployment shape (separate binaries? separate repos? multiple `cmd/` entries
      in this repo?)

---

## api -
|____ health
|        |_____ health.HealthService/Check ✅
|
|____ auth
|        |_____ auth.AuthService/SendOtp ✅
|        |_____ auth.AuthService/CheckOtp ✅
|        |_____ auth.AuthService/Me ✅
|        |_____ auth.AuthService/Logout ✅
|        |_____ auth.AuthService/RefreshToken ✅
|
|____ user
|        |_____ user.UserService/Profile ✅
|        |_____ user.UserService/UpdateProfile ✅
|        |_____ user.UserService/ChangePassword ✅
|
|____ province
|        |_____ province.ProvinceService/GetAll ✅
|        |_____ province.ProvinceService/GetOne ✅
|
|____ address
|        |_____ address.AddressService/Create ✅
|        |_____ address.AddressService/GetAll ✅
|        |_____ address.AddressService/GetOne ✅
|        |_____ address.AddressService/Update ✅
|        |_____ address.AddressService/Delete ✅
|
|____ category
|        |_____ category.CategoryService/Create ✅ (admin)
|        |_____ category.CategoryService/GetAll ✅
|        |_____ category.CategoryService/GetOne ✅
|        |_____ category.CategoryService/Update ✅ (admin)
|        |_____ category.CategoryService/Delete ✅ (admin)
|
|____ product
|        |_____ product.ProductService/Create ✅ (admin, client-streaming — metadata + gallery images)
|        |_____ product.ProductService/GetAll ✅ (pagination + optional category_slug filter; no price filter/search/sort yet)
|        |_____ product.ProductService/GetOneBySlug ✅
|        |_____ product.ProductService/Update -> not built yet
|        |_____ product.ProductService/Delete -> not built yet
|
|____ cart
|        |_____ get current cart
|        |_____ add item
|        |_____ update item quantity
|        |_____ remove item
|        |_____ clear cart
|
|____ order
|        |_____ checkout -> build order from current cart + address
|        |_____ current user's order history
|        |_____ get one order
|        |_____ cancel order
|        |_____ (admin) list all orders
|        |_____ (admin) update order status -> processing/shipped/delivered/canceled
|
|____ payment
|        |_____ init payment / redirect to gateway
|        |_____ gateway webhook -> confirm payment & update order
|
