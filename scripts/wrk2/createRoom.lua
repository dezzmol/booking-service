wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"

function init(args)
    math.randomseed(os.time())
end

function request()
    local body = string.format([[
    {
        "hotel_id": %d,
        "number": %d,
        "capacity": %d
    }
    ]], math.random(1, 50000), math.random(1, 999), math.random(1, 6))

    return wrk.format("POST", "/v1/room", nil, body)
end