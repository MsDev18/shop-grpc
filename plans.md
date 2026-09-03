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
process into independently deployable services. Agreed approach:

- Gradual, one service at a time — not a big-bang split of the whole system at once.
- Each extracted service gets its own database from day one (real data ownership per
  service, not a shared-DB "distributed monolith" stepping stone).
- Every other module's direct in-process call into the extracted service (e.g. product's
  `categoryService category.Service` dependency) must be rewired from a plain Go function
  call into a real gRPC client call over the network — this is the actual cost of each
  extraction, and the reason it's done deliberately one module at a time instead of all at
  once.

- [ ] decide which module goes first
- [ ] design the data-ownership / migration plan for that module's tables
- [ ] introduce a gRPC client (instead of a direct Go interface) wherever another module
      currently calls it in-process
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
