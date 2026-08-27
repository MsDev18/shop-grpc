api -
|____ health
|        |_____ GET /health-check ✅
|
|____ auth
|        |_____ POST /auth/send-otp ✅
|        |_____ POST /auth/check-otp ✅
|        |_____ GET /auth/me ✅
|        |_____ POST /auth/logout ✅
|        |_____ POST /auth/refresh-token ✅
|
|____ user
|        |_____ GET /user/profile ✅
|        |_____ PATCH /user/update-profile ✅
|        |_____ PATCH /user/change-password ✅
|
|____ province
|        |_____ GET /province ✅
|        |_____ GET /province/:id ✅
|
|____ address
|        |_____ POST /address ✅
|        |_____ GET /address ✅
|        |_____ GET /address/:id ✅
|        |_____ PATCH /address/:id ✅
|        |_____ DELETE /address/:id ✅
|
|____ category
|        |_____ POST /category ✅ (admin)
|        |_____ GET /category ✅
|        |_____ GET /category/:slug ✅
|        |_____ PATCH /category/:slug ✅ (admin)
|        |_____ DELETE /category/:slug ✅ (admin)
|
|____ product
|        |_____ POST /product (admin) ✅
|        |_____ GET /product  -> list, pagination + filter(category,price) + search + sort
|        |_____ GET /product/:slug ✅
|        |_____ PATCH /product/:slug (admin)
|        |_____ DELETE /product/:slug (admin)
|        |_____ GET /category/:slug/products
|
|____ cart
|        |_____ GET /cart
|        |_____ POST /cart/items
|        |_____ PATCH /cart/items/:id
|        |_____ DELETE /cart/items/:id
|        |_____ DELETE /cart
|
|____ order
|        |_____ POST /order          -> checkout, build order from current cart + address_id
|        |_____ GET /order           -> current user's order history
|        |_____ GET /order/:id
|        |_____ PATCH /order/:id/cancel
|        |_____ GET /admin/order (admin)          -> list all orders
|        |_____ PATCH /order/:id/status (admin)   -> processing/shipped/delivered/canceled
|
|____ payment
|        |_____ POST /order/:id/pay        -> init payment / redirect to gateway
|        |_____ POST /payment/callback     -> gateway webhook, confirm payment & update order
|
