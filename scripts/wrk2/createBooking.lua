wrk.method = "PUT"
wrk.headers["Content-Type"] = "application/json"

function init(args)
    math.randomseed(os.time())
end

function request()
    local body = string.format([[
    {
        "room_id": %d,
        "guest_id": %d,
        "from": "2025-01-%02d",
        "to": "2025-01-%02d"
    }
    ]],
    math.random(1, 50000),
    math.random(1, 50000),
    math.random(1, 27),
    math.random(2, 28))

    return wrk.format("PUT", "/v1/booking", nil, body)
end